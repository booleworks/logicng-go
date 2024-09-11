package maxsat

import (
	"math"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/sat"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/emirpasic/gods/utils"
)

type oll struct {
	*maxSatAlgorithm
	encoder *encoder
}

func newOLL(fac f.Factory) *oll {
	return &oll{
		maxSatAlgorithm: newAlgorithm(fac, DefaultConfig()),
	}
}

func (m *oll) search(hdl handler.Handler) (result, handler.State) {
	return m.innerSearch(hdl, func() (result, handler.State) {
		m.encoder = newEncoder()
		if m.problemType == weighted {
			return m.weighted()
		} else {
			return m.unweighted()
		}
	})
}

func (m *oll) unweighted() (result, handler.State) {
	coreMapping := make(map[int32]int)
	boundMapping := make(map[int32]intTriple)
	m.nbInitialVariables = m.nVars()
	m.initRelaxation()
	solver := m.rebuildSolver()

	var assumptions []int32
	var encodingAssumptions []int32
	m.encoder.setIncremental(IncIterative)
	activeSoft := make([]bool, m.nSoft())
	for i := 0; i < m.nSoft(); i++ {
		coreMapping[m.softClauses[i].assumptionVar] = i
	}
	cardinalityAssumptions := treeset.NewWith(utils.Int32Comparator)
	var softCardinality []*encoder

	for {
		res, state := searchSatSolverWithAssumptions(solver, m.hdl, assumptions)
		if !state.Success {
			return resUndef, state
		} else if res == f.TristateTrue {
			m.nbSatisfiable++
			model := solver.Model()
			newCost := m.computeCostModel(model, math.MaxInt)
			m.saveModel(model)

			m.ubCost = newCost
			if m.nbSatisfiable == 1 {
				if newCost == 0 {
					return resOptimum, succ
				}
				for i := 0; i < m.nSoft(); i++ {
					assumptions = append(assumptions, sat.Not(m.softClauses[i].assumptionVar))
				}
			} else {
				return resOptimum, succ
			}
		} else {
			m.lbCost++
			m.nbCores++
			if m.nbSatisfiable == 0 {
				return resUnsat, succ
			}
			if m.lbCost == m.ubCost {
				return resOptimum, succ
			}

			m.sumSizeCores += len(solver.Conflict())
			var softRelax []int32
			var cardinalityRelax []int32

			for i := 0; i < len(solver.Conflict()); i++ {
				p := solver.Conflict()[i]
				if entry, ok := coreMapping[p]; ok {
					activeSoft[entry] = true
					softRelax = append(softRelax, p)
				}

				if softId, ok := boundMapping[p]; ok {
					cardinalityAssumptions.Remove(p)
					cardinalityRelax = append(cardinalityRelax, p)

					encodingAssumptions = []int32{}
					softCardinality[softId.id].incUpdateCardinality(
						solver, []int32{}, softId.bound+1, &encodingAssumptions,
					)

					// if the bound is the same as the number of literals
					// then no restriction is applied
					if softId.bound+1 < len(softCardinality[softId.id].outputs()) {
						out := softCardinality[softId.id].outputs()[softId.bound+1]
						boundMapping[out] = intTriple{softId.id, softId.bound + 1, 1}
						cardinalityAssumptions.Add(out)
					}
				}
			}

			if len(softRelax) == 1 && len(cardinalityRelax) == 0 {
				solver.AddClause([]int32{softRelax[0]}, nil)
			}
			if len(softRelax)+len(cardinalityRelax) > 1 {
				relaxHarden := make([]int32, len(softRelax))
				copy(relaxHarden, softRelax)
				for i := 0; i < len(cardinalityRelax); i++ {
					relaxHarden = append(relaxHarden, cardinalityRelax[i])
				}
				e := newEncoder()
				e.setIncremental(IncIterative)
				e.buildCardinality(solver, relaxHarden, 1)
				softCardinality = append(softCardinality, e)
				out := e.outputs()[1]
				boundMapping[out] = intTriple{int32(len(softCardinality) - 1), 1, 1}
				cardinalityAssumptions.Add(out)
			}
			// reset the assumptions
			assumptions = []int32{}
			for i := 0; i < m.nSoft(); i++ {
				if !activeSoft[i] {
					assumptions = append(assumptions, sat.Not(m.softClauses[i].assumptionVar))
				}
			}
			cardinalityAssumptions.Each(func(_ int, value interface{}) {
				assumptions = append(assumptions, sat.Not(value.(int32)))
			})
		}
	}
}

func (m *oll) weighted() (result, handler.State) {
	coreMapping := make(map[int32]int)
	boundMapping := make(map[int32]intTriple)
	m.nbInitialVariables = m.nVars()
	m.initRelaxation()
	solver := m.rebuildSolver()

	var assumptions []int32
	var encodingAssumptions []int32
	m.encoder.setIncremental(IncIterative)

	activeSoft := make([]bool, m.nSoft())
	for i := 0; i < m.nSoft(); i++ {
		coreMapping[m.softClauses[i].assumptionVar] = i
	}

	cardinalityAssumptions := treeset.NewWith(utils.Int32Comparator)
	var softCardinality []*encoder
	minWeight := m.currentWeight

	for {
		res, state := searchSatSolverWithAssumptions(solver, m.hdl, assumptions)
		if !state.Success {
			return resUndef, state
		} else if res == f.TristateTrue {
			m.nbSatisfiable++
			model := solver.Model()
			newCost := m.computeCostModel(model, math.MaxInt)
			if newCost < m.ubCost || m.nbSatisfiable == 1 {
				m.saveModel(model)
				m.ubCost = newCost
			}
			if m.nbSatisfiable == 1 {
				minWeight := m.findNextWeightDiversity(minWeight, cardinalityAssumptions, boundMapping)
				for i := 0; i < m.nSoft(); i++ {
					if m.softClauses[i].weight >= minWeight {
						assumptions = append(assumptions, sat.Not(m.softClauses[i].assumptionVar))
					}
				}
			} else {
				// compute min weight in soft
				notConsidered := 0
				for i := 0; i < m.nSoft(); i++ {
					if m.softClauses[i].weight < minWeight {
						notConsidered++
					}
				}
				cardinalityAssumptions.Each(func(_ int, value interface{}) {
					softId := boundMapping[value.(int32)]
					if softId.weight < minWeight {
						notConsidered++
					}
				})
				if notConsidered != 0 {
					minWeight = m.findNextWeightDiversity(minWeight, cardinalityAssumptions, boundMapping)
					assumptions = []int32{}
					for i := 0; i < m.nSoft(); i++ {
						if !activeSoft[i] && m.softClauses[i].weight >= minWeight {
							assumptions = append(assumptions, sat.Not(m.softClauses[i].assumptionVar))
						}
					}
					cardinalityAssumptions.Each(func(_ int, value interface{}) {
						softId := boundMapping[value.(int32)]
						if softId.weight >= minWeight {
							assumptions = append(assumptions, sat.Not(value.(int32)))
						}
					})
				} else {
					return resOptimum, succ
				}
			}
		} else if res == f.TristateFalse {
			// reduce the weighted to the unweighted case
			minCore := math.MaxInt
			for i := 0; i < len(solver.Conflict()); i++ {
				p := solver.Conflict()[i]
				if entry, ok := coreMapping[p]; ok {
					if m.softClauses[entry].weight < minCore {
						minCore = m.softClauses[entry].weight
					}
				}
				if softId, ok := boundMapping[p]; ok {
					if softId.weight < minCore {
						minCore = softId.weight
					}
				}
			}
			m.lbCost += minCore
			m.nbCores++
			if m.nbSatisfiable == 0 {
				return resUnsat, succ
			}
			if m.lbCost == m.ubCost {
				return resOptimum, succ
			}
			m.sumSizeCores += len(solver.Conflict())
			var softRelax []int32
			var cardinalityRelax []int32

			for i := 0; i < len(solver.Conflict()); i++ {
				p := solver.Conflict()[i]
				if entry, ok := coreMapping[p]; ok {
					if m.softClauses[entry].weight > minCore {
						// Split the clause
						indexSoft := coreMapping[p]

						// Update the weight of the soft clause.
						m.softClauses[indexSoft].weight = m.softClauses[indexSoft].weight - minCore
						clause := make([]int32, len(m.softClauses[indexSoft].clause))
						copy(clause, m.softClauses[indexSoft].clause)
						var vars []int32

						// Since cardinality constraints are added the variables are not in sync
						for m.nVars() < int(solver.NVars()) {
							m.newLiteral(false)
						}
						l := m.newLiteral(false)
						vars = append(vars, l)

						// Add a new soft clause with the weight of the core
						m.addSoftClauseWithAssumptions(minCore, clause, vars)
						activeSoft = append(activeSoft, true)

						// Add information to the SAT solver
						newSatVariable(solver)
						clause = append(clause, l)
						solver.AddClause(clause, nil)

						// Create a new assumption literal.
						m.softClauses[m.nSoft()-1].assumptionVar = l
						// Map the new soft clause to its assumption literal
						coreMapping[l] = m.nSoft() - 1
						softRelax = append(softRelax, l)
					} else {
						softRelax = append(softRelax, p)
						activeSoft[coreMapping[p]] = true
					}
				}
				if softId, ok := boundMapping[p]; ok {
					// this is a soft cardinality -- bound must be increased

					// increase the bound
					if softId.weight == minCore {
						cardinalityAssumptions.Remove(p)
						cardinalityRelax = append(cardinalityRelax, p)
						encodingAssumptions = []int32{}
						softCardinality[softId.id].incUpdateCardinality(solver, []int32{},
							softId.bound+1, &encodingAssumptions)

						// if the bound is the same as the number of literals then no restriction is applied
						if softId.bound+1 < len(softCardinality[softId.id].outputs()) {
							out := softCardinality[softId.id].outputs()[softId.bound+1]
							boundMapping[out] = intTriple{softId.id, softId.bound + 1, minCore}
							cardinalityAssumptions.Add(out)
						}
					} else {
						// Duplicate cardinality constraint
						e := newEncoder()
						e.setIncremental(IncIterative)
						e.buildCardinality(solver, softCardinality[softId.id].lits(), softId.bound)
						out := e.outputs()[softId.bound]
						softCardinality = append(softCardinality, e)
						boundMapping[out] = intTriple{int32(len(softCardinality) - 1), softId.bound, minCore}
						cardinalityRelax = append(cardinalityRelax, out)

						// Update value of the previous cardinality constraint
						boundMapping[p] = intTriple{softId.id, softId.bound, softId.weight - minCore}

						// Update bound as usual...
						softCoreId := boundMapping[out]
						encodingAssumptions = []int32{}
						softCardinality[softCoreId.id].incUpdateCardinality(solver, []int32{},
							softCoreId.bound+1, &encodingAssumptions)

						// if the bound is the same as the number of literals then no restriction is applied
						if softCoreId.bound+1 < len(softCardinality[softCoreId.id].outputs()) {
							out2 := softCardinality[softCoreId.id].outputs()[softCoreId.bound+1]
							boundMapping[out2] = intTriple{softCoreId.id, softCoreId.bound + 1, minCore}
							cardinalityAssumptions.Add(out2)
						}
					}
				}
			}
			if len(softRelax) == 1 && len(cardinalityRelax) == 0 {
				solver.AddClause([]int32{softRelax[0]}, nil)
			}
			if len(softRelax)+len(cardinalityRelax) > 1 {
				relaxHarden := make([]int32, len(softRelax))
				copy(relaxHarden, softRelax)
				for i := 0; i < len(cardinalityRelax); i++ {
					relaxHarden = append(relaxHarden, cardinalityRelax[i])
				}
				e := newEncoder()
				e.setIncremental(IncIterative)
				e.buildCardinality(solver, relaxHarden, 1)
				softCardinality = append(softCardinality, e)
				out := e.outputs()[1]
				boundMapping[out] = intTriple{int32(len(softCardinality) - 1), 1, minCore}
				cardinalityAssumptions.Add(out)
			}
			assumptions = []int32{}
			for i := 0; i < m.nSoft(); i++ {
				if !activeSoft[i] && m.softClauses[i].weight >= minWeight {
					assumptions = append(assumptions, sat.Not(m.softClauses[i].assumptionVar))
				}
			}
			cardinalityAssumptions.Each(func(_ int, value interface{}) {
				softId := boundMapping[value.(int32)]
				if softId.weight >= minWeight {
					assumptions = append(assumptions, sat.Not(value.(int32)))
				}
			})
		}
	}
}

func (m *oll) initRelaxation() {
	for i := 0; i < m.nSoft(); i++ {
		l := m.newLiteral(false)
		m.softClauses[i].relaxationVars = append(m.softClauses[i].relaxationVars, l)
		m.softClauses[i].assumptionVar = l
	}
}

func (m *oll) rebuildSolver() *sat.CoreSolver {
	s := m.newSatSolver()
	for i := 0; i < m.nVars(); i++ {
		newSatVariable(s)
	}
	for i := 0; i < m.nHard(); i++ {
		s.AddClause(m.hardClauses[i].clause, nil)
	}
	for i := 0; i < m.nSoft(); i++ {
		clause := make([]int32, len(m.softClauses[i].clause))
		copy(clause, m.softClauses[i].clause)
		for j := 0; j < len(m.softClauses[i].relaxationVars); j++ {
			clause = append(clause, m.softClauses[i].relaxationVars[j])
		}
		s.AddClause(clause, nil)
	}
	return s
}

func (m *oll) findNextWeightDiversity(
	weight int, cardinalityAssumptions *treeset.Set, boundMapping map[int32]intTriple,
) int {
	nextWeight := weight
	var nbClauses int
	nbWeights := treeset.NewWithIntComparator()
	alpha := 1.25
	findNext := false
	for {
		if m.nbSatisfiable > 1 || findNext {
			nextWeight = m.findNextWeight(nextWeight, cardinalityAssumptions, boundMapping)
		}
		nbClauses = 0
		nbWeights.Clear()
		for i := 0; i < m.nSoft(); i++ {
			if m.softClauses[i].weight >= nextWeight {
				nbClauses++
				nbWeights.Add(m.softClauses[i].weight)
			}
		}
		cardinalityAssumptions.Each(func(_ int, value interface{}) {
			softId := boundMapping[value.(int32)]
			if softId.weight >= nextWeight {
				nbClauses++
				nbWeights.Add(softId.weight)
			}
		})

		if float64(nbClauses)/float64(nbWeights.Size()) > alpha ||
			nbClauses == m.nSoft()+cardinalityAssumptions.Size() {
			break
		}
		if m.nbSatisfiable == 1 && !findNext {
			findNext = true
		}
	}
	return nextWeight
}

func (m *oll) findNextWeight(
	weight int, cardinalityAssumptions *treeset.Set, boundMapping map[int32]intTriple,
) int {
	nextWeight := 1
	for i := 0; i < m.nSoft(); i++ {
		if m.softClauses[i].weight > nextWeight && m.softClauses[i].weight < weight {
			nextWeight = m.softClauses[i].weight
		}
	}
	cardinalityAssumptions.Each(func(_ int, value interface{}) {
		softId := boundMapping[value.(int32)]
		if softId.weight > nextWeight && softId.weight < weight {
			nextWeight = softId.weight
		}
	})
	return nextWeight
}

type intTriple struct {
	id     int32
	bound  int
	weight int
}
