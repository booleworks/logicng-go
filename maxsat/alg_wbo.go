package maxsat

import (
	"math"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/sat"
	"github.com/emirpasic/gods/sets/treeset"
)

type wbo struct {
	*maxSatAlgorithm
	solver                    *sat.CoreSolver
	nbCurrentSoft             int
	weightStrategy            WeightStrategy
	coreMapping               map[int32]int
	assumptions               []int32
	symmetryStrategy          bool
	indexSoftCore             []int
	softMapping               [][]int
	relaxationMapping         [][]int32
	duplicatedSymmetryClauses map[intPair]present
	symmetryBreakingLimit     int
}

func newWBO(fac f.Factory, config *Config) *wbo {
	return &wbo{
		maxSatAlgorithm:           newAlgorithm(fac, config),
		solver:                    nil,
		weightStrategy:            config.WeightStrategy,
		symmetryStrategy:          config.Symmetry,
		symmetryBreakingLimit:     config.Limit,
		coreMapping:               make(map[int32]int),
		assumptions:               []int32{},
		indexSoftCore:             []int{},
		softMapping:               [][]int{},
		relaxationMapping:         [][]int32{},
		duplicatedSymmetryClauses: make(map[intPair]present),
	}
}

func (m *wbo) search(hdl handler.Handler) (result, handler.State) {
	m.nbInitialVariables = m.nVars()
	m.coreMapping = make(map[int32]int)
	m.assumptions = []int32{}
	m.indexSoftCore = []int{}
	m.softMapping = [][]int{}
	m.relaxationMapping = [][]int32{}
	m.duplicatedSymmetryClauses = make(map[intPair]present)

	if m.currentWeight == 1 {
		m.problemType = unweighted
		m.weightStrategy = WeightNone
	}
	return m.innerSearch(hdl, func() (result, handler.State) {
		if m.symmetryStrategy {
			m.initSymmetry()
		}
		if m.problemType == unweighted || m.weightStrategy == WeightNone {
			return m.normalSearch()
		} else if m.weightStrategy == WeightNormal || m.weightStrategy == WeightDiversify {
			return m.weightSearch()
		}
		panic(errorx.UnknownEnumValue(m.problemType))
	})
}

func (m *wbo) initSymmetry() {
	for i := 0; i < m.nSoft(); i++ {
		m.softMapping = append(m.softMapping, []int{})
		m.relaxationMapping = append(m.relaxationMapping, []int32{})
	}
}

func (m *wbo) normalSearch() (result, handler.State) {
	switch res, state := m.unsatSearch(); res {
	case f.TristateUndef:
		return resUndef, state
	case f.TristateFalse:
		return resUnsat, succ
	}

	m.initAssumptions()
	m.solver = m.rebuildSolver()
	for {
		res, state := searchSatSolverWithAssumptions(m.solver, m.hdl, m.assumptions)
		if !state.Success {
			return resUndef, state
		} else if res == f.TristateFalse {
			m.nbCores++
			coreCost := m.computeCostCore(m.solver.Conflict())
			m.lbCost += coreCost
			if m.lbCost == m.ubCost {
				return resOptimum, succ
			} else if state := m.foundLowerBound(m.lbCost); !state.Success {
				return resUndef, state
			}
			m.relaxCore(m.solver.Conflict(), coreCost, &m.assumptions)
			m.solver = m.rebuildSolver()
		} else {
			m.nbSatisfiable++
			m.ubCost = m.computeCostModel(m.solver.Model(), math.MaxInt)
			m.saveModel(m.solver.Model())
			return resOptimum, succ
		}
	}
}

func (m *wbo) unsatSearch() (f.Tristate, handler.State) {
	m.solver = m.rebuildHardSolver()
	res, state := searchSatSolverWithAssumptions(m.solver, m.hdl, m.assumptions)
	if !state.Success {
		return f.TristateUndef, state
	} else if res == f.TristateFalse {
		m.nbCores++
	} else if res == f.TristateTrue {
		m.nbSatisfiable++
		cost := m.computeCostModel(m.solver.Model(), math.MaxInt)
		m.ubCost = cost
		m.saveModel(m.solver.Model())
	}
	m.solver = nil
	return res, state
}

func (m *wbo) rebuildHardSolver() *sat.CoreSolver {
	s := m.newSatSolver()
	for i := 0; i < m.nVars(); i++ {
		newSatVariable(s)
	}
	for i := 0; i < m.nHard(); i++ {
		s.AddClause(m.hardClauses[i].clause, nil)
	}
	return s
}

func (m *wbo) initAssumptions() {
	for i := 0; i < m.nSoft(); i++ {
		l := m.newLiteral(false)
		m.softClauses[i].assumptionVar = l
		m.coreMapping[l] = i
		m.assumptions = append(m.assumptions, sat.Not(l))
	}
}

func (m *wbo) rebuildSolver() *sat.CoreSolver {
	s := m.newSatSolver()
	for i := 0; i < m.nVars(); i++ {
		newSatVariable(s)
	}
	for i := 0; i < m.nHard(); i++ {
		s.AddClause(m.hardClauses[i].clause, nil)
	}
	if m.symmetryStrategy {
		m.symmetryBreaking()
	}
	for i := 0; i < m.nSoft(); i++ {
		clause := make([]int32, len(m.softClauses[i].clause))
		copy(clause, m.softClauses[i].clause)
		for j := 0; j < len(m.softClauses[i].relaxationVars); j++ {
			clause = append(clause, m.softClauses[i].relaxationVars[j])
		}
		clause = append(clause, m.softClauses[i].assumptionVar)
		s.AddClause(clause, nil)
	}
	return s
}

func (m *wbo) symmetryBreaking() {
	if len(m.indexSoftCore) != 0 && m.nbSymmetryClauses < m.symmetryBreakingLimit {
		coreIntersection := make([][]int32, m.nbCores)
		coreIntersectionCurrent := make([][]int32, m.nbCores)
		for i := 0; i < m.nbCores; i++ {
			coreIntersection[i] = []int32{}
			coreIntersectionCurrent[i] = []int32{}
		}
		var coreList []int
		for i := 0; i < len(m.indexSoftCore); i++ {
			p := m.indexSoftCore[i]
			var addCores []int
			for j := 0; j < len(m.softMapping[p])-1; j++ {
				core := m.softMapping[p][j]
				addCores = append(addCores, core)
				if len(coreIntersection[core]) == 0 {
					coreList = append(coreList, core)
				}
				coreIntersection[core] = append(coreIntersection[core], m.relaxationMapping[p][j])
			}
			for j := 0; j < len(addCores); j++ {
				core := addCores[j]
				b := len(m.softMapping[p]) - 1
				coreIntersectionCurrent[core] = append(coreIntersectionCurrent[core], m.relaxationMapping[p][b])
			}
			for k := 0; k < len(coreList); k++ {
				for n := 0; n < len(coreIntersection[coreList[k]]); n++ {
					for j := n + 1; j < len(coreIntersectionCurrent[coreList[k]]); j++ {
						clause := make([]int32, 2)
						clause[0] = sat.Not(coreIntersection[coreList[k]][n])
						clause[1] = sat.Not(coreIntersectionCurrent[coreList[k]][j])
						symClause := intPair{
							sat.Vari(coreIntersection[coreList[k]][n]),
							sat.Vari(coreIntersectionCurrent[coreList[k]][j]),
						}
						if sat.Vari(coreIntersection[coreList[k]][n]) >
							sat.Vari(coreIntersectionCurrent[coreList[k]][j]) {
							symClause = intPair{
								sat.Vari(coreIntersectionCurrent[coreList[k]][j]),
								sat.Vari(coreIntersection[coreList[k]][n]),
							}
						}
						_, ok := m.duplicatedSymmetryClauses[symClause]
						if !ok {
							m.duplicatedSymmetryClauses[symClause] = present{}
							m.addHardClause(clause)
							m.nbSymmetryClauses++
							if m.symmetryBreakingLimit == m.nbSymmetryClauses {
								break
							}
						}
					}
					if m.symmetryBreakingLimit == m.nbSymmetryClauses {
						break
					}
				}
				if m.symmetryBreakingLimit == m.nbSymmetryClauses {
					break
				}
			}
			if m.symmetryBreakingLimit == m.nbSymmetryClauses {
				break
			}
		}
	}
	m.indexSoftCore = []int{}
}

func (m *wbo) computeCostCore(conflict []int32) int {
	if m.problemType == unweighted {
		return 1
	}
	coreCost := math.MaxInt
	for i := 0; i < len(conflict); i++ {
		indexSoft := m.coreMapping[conflict[i]]
		if m.softClauses[indexSoft].weight < coreCost {
			coreCost = m.softClauses[indexSoft].weight
		}
	}
	return coreCost
}

func (m *wbo) relaxCore(conflict []int32, weightCore int, assumps *[]int32) {
	var lits []int32
	for i := 0; i < len(conflict); i++ {
		indexSoft := m.coreMapping[conflict[i]]
		if m.softClauses[indexSoft].weight == weightCore {
			p := m.newLiteral(false)
			m.softClauses[indexSoft].relaxationVars = append(m.softClauses[indexSoft].relaxationVars, p)
			lits = append(lits, p)
			if m.symmetryStrategy {
				m.symmetryLog(indexSoft)
			}
		} else {
			m.softClauses[indexSoft].weight = m.softClauses[indexSoft].weight - weightCore
			clause := make([]int32, len(m.softClauses[indexSoft].clause))
			copy(clause, m.softClauses[indexSoft].clause)
			vars := make([]int32, len(m.softClauses[indexSoft].relaxationVars))
			copy(vars, m.softClauses[indexSoft].relaxationVars)
			p := m.newLiteral(false)
			vars = append(vars, p)
			lits = append(lits, p)
			m.addSoftClauseWithAssumptions(weightCore, clause, vars)
			l := m.newLiteral(false)
			m.softClauses[m.nSoft()-1].assumptionVar = l
			// Map the new soft clause to its assumption literal
			m.coreMapping[l] = m.nSoft() - 1
			// Update the assumption vector
			*assumps = append(*assumps, sat.Not(l))
			if m.symmetryStrategy {
				m.symmetryLog(m.nSoft() - 1)
			}
		}
	}
	m.encodeEO(lits)
	m.sumSizeCores += len(conflict)
}

func (m *wbo) symmetryLog(p int) {
	if m.nbSymmetryClauses < m.symmetryBreakingLimit {
		for len(m.softMapping) <= p {
			m.softMapping = append(m.softMapping, []int{})
			m.relaxationMapping = append(m.relaxationMapping, []int32{})
		}
		m.softMapping[p] = append(m.softMapping[p], m.nbCores)
		back := m.softClauses[p].relaxationVars[len(m.softClauses[p].relaxationVars)-1]
		m.relaxationMapping[p] = append(m.relaxationMapping[p], back)
		if len(m.softMapping[p]) > 1 {
			m.indexSoftCore = append(m.indexSoftCore, p)
		}
	}
}

func (m *wbo) encodeEO(lits []int32) {
	if len(lits) == 1 {
		m.addHardClause([]int32{lits[0]})
	} else {
		auxVariables := make([]int32, len(lits)-1)
		for i := 0; i < len(lits)-1; i++ {
			auxVariables[i] = m.newLiteral(false)
		}
		for i := 0; i < len(lits); i++ {
			if i == 0 {
				m.addHardClause([]int32{lits[i], sat.Not(auxVariables[i])})
				m.addHardClause([]int32{sat.Not(lits[i]), auxVariables[i]})
			} else if i == len(lits)-1 {
				m.addHardClause([]int32{lits[i], auxVariables[i-1]})
				m.addHardClause([]int32{sat.Not(lits[i]), sat.Not(auxVariables[i-1])})
			} else {
				m.addHardClause([]int32{sat.Not(auxVariables[i-1]), auxVariables[i]})
				m.addHardClause([]int32{lits[i], sat.Not(auxVariables[i]), auxVariables[i-1]})
				m.addHardClause([]int32{sat.Not(lits[i]), auxVariables[i]})
				m.addHardClause([]int32{sat.Not(lits[i]), sat.Not(auxVariables[i-1])})
			}
		}
	}
}

func (m *wbo) weightSearch() (result, handler.State) {
	switch res, state := m.unsatSearch(); res {
	case f.TristateUndef:
		return resUndef, state
	case f.TristateFalse:
		return resUnsat, succ
	}

	m.initAssumptions()
	m.updateCurrentWeight(m.weightStrategy)
	m.solver = m.rebuildWeightSolver()

	for {
		res, state := searchSatSolverWithAssumptions(m.solver, m.hdl, m.assumptions)
		if !state.Success {
			return resUndef, state
		} else if res == f.TristateFalse {
			m.nbCores++
			coreCost := m.computeCostCore(m.solver.Conflict())
			m.lbCost += coreCost
			if state := m.foundLowerBound(m.lbCost); !state.Success {
				return resUndef, state
			}
			m.relaxCore(m.solver.Conflict(), coreCost, &m.assumptions)
			m.solver = m.rebuildWeightSolver()
		} else {
			m.nbSatisfiable++
			if m.nbCurrentSoft == m.nSoft() {
				if m.lbCost < m.ubCost {
					m.ubCost = m.lbCost
					m.saveModel(m.solver.Model())
				}
				return resOptimum, succ
			} else {
				m.updateCurrentWeight(m.weightStrategy)
				cost := m.computeCostModel(m.solver.Model(), math.MaxInt)
				if cost < m.ubCost {
					m.ubCost = cost
					m.saveModel(m.solver.Model())
				}
				if m.lbCost == m.ubCost {
					return resOptimum, succ
				} else if state := m.foundUpperBound(m.ubCost); !state.Success {
					return resUndef, state
				}
				m.solver = m.rebuildWeightSolver()
			}
		}
	}
}

func (m *wbo) updateCurrentWeight(strategy WeightStrategy) {
	switch strategy {
	case WeightNormal:
		m.currentWeight = m.findNextWeight(m.currentWeight)
	case WeightDiversify:
		m.currentWeight = m.findNextWeightDiversity(m.currentWeight)
	}
}

func (m *wbo) findNextWeight(weight int) int {
	nextWeight := 1
	for i := 0; i < m.nSoft(); i++ {
		if m.softClauses[i].weight > nextWeight && m.softClauses[i].weight < weight {
			nextWeight = m.softClauses[i].weight
		}
	}
	return nextWeight
}

func (m *wbo) findNextWeightDiversity(weight int) int {
	nextWeight := weight
	nbWeights := treeset.NewWithIntComparator()
	alpha := 1.25
	findNext := false
	for {
		if m.nbSatisfiable > 1 || findNext {
			nextWeight = m.findNextWeight(nextWeight)
		}
		nbClauses := 0
		nbWeights.Clear()
		for i := 0; i < m.nSoft(); i++ {
			if m.softClauses[i].weight >= nextWeight {
				nbClauses++
				nbWeights.Add(m.softClauses[i].weight)
			}
		}
		if float64(nbClauses)/float64(nbWeights.Size()) > alpha || nbClauses == m.nSoft() {
			break
		}
		if m.nbSatisfiable == 1 && !findNext {
			findNext = true
		}
	}
	return nextWeight
}

func (m *wbo) rebuildWeightSolver() *sat.CoreSolver {
	s := m.newSatSolver()
	for i := 0; i < m.nVars(); i++ {
		newSatVariable(s)
	}
	for i := 0; i < m.nHard(); i++ {
		s.AddClause(m.hardClauses[i].clause, nil)
	}
	if m.symmetryStrategy {
		m.symmetryBreaking()
	}
	m.nbCurrentSoft = 0
	for i := 0; i < m.nSoft(); i++ {
		if m.softClauses[i].weight >= m.currentWeight {
			m.nbCurrentSoft++
			clause := make([]int32, len(m.softClauses[i].clause))
			copy(clause, m.softClauses[i].clause)
			for j := 0; j < len(m.softClauses[i].relaxationVars); j++ {
				clause = append(clause, m.softClauses[i].relaxationVars[j])
			}
			clause = append(clause, m.softClauses[i].assumptionVar)
			s.AddClause(clause, nil)
		}
	}
	return s
}

type intPair struct {
	i1 int32
	i2 int32
}

type present struct{}
