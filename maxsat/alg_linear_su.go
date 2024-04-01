package maxsat

import (
	"math"

	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/sat"
)

type linearSU struct {
	*maxSatAlgorithm
	encoder     *encoder
	bmoMode     bool
	objFunction []int32
	coeffs      []int
	solver      *sat.CoreSolver
	bmo         bool
}

func newLinearSU(config ...*Config) *linearSU {
	var cfg *Config
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = DefaultConfig()
	}
	return &linearSU{
		maxSatAlgorithm: newAlgorithm(),
		solver:          nil,
		encoder:         newEncoder(),
		bmoMode:         cfg.BMO,
		bmo:             false,
		objFunction:     []int32{},
		coeffs:          []int{},
	}
}

func (m *linearSU) search(handler Handler) (result, bool) {
	return m.innerSearch(handler, func() (result, bool) {
		m.nbInitialVariables = m.nVars()
		if m.currentWeight == 1 {
			m.problemType = unweighted
		} else {
			m.bmo = m.isBmo(true)
		}
		if m.problemType == weighted {
			if m.bmoMode && m.bmo {
				return m.bmoSearch()
			} else {
				return m.normalSearch()
			}
		} else {
			return m.normalSearch()
		}
	})
}

func (m *linearSU) bmoSearch() (result, bool) {
	m.initRelaxation()
	currentWeight := m.orderWeights[0]
	minWeight := m.orderWeights[len(m.orderWeights)-1]
	posWeight := 0
	var functions [][]int32
	var weights []int
	m.solver = m.rebuildBmo(functions, weights, currentWeight)
	localCost := 0
	m.ubCost = 0
	for {
		satHandler := m.satHandler()
		res, ok := searchSatSolver(m.solver, satHandler)
		if !ok {
			return resUndef, false
		}
		if res == f.TristateTrue {
			m.nbSatisfiable++
			newCost := m.computeCostModel(m.solver.Model(), currentWeight)
			if currentWeight == minWeight {
				m.saveModel(m.solver.Model())
				m.ubCost = newCost + m.lbCost
				if newCost > 0 && !m.foundUpperBound(m.ubCost, nil) {
					return resUndef, false
				}
			}
			if newCost == 0 && currentWeight == minWeight {
				return resOptimum, true
			} else {
				if newCost == 0 {
					obj := make([]int32, len(m.objFunction))
					copy(obj, m.objFunction)
					functions = append(functions, obj)
					localCost = newCost
					weights = append(weights, 0)
					posWeight++
					currentWeight = m.orderWeights[posWeight]
					m.solver = m.rebuildBmo(functions, weights, currentWeight)
				} else {
					if localCost == 0 {
						m.encoder.encodeCardinality(m.solver, m.objFunction, newCost/currentWeight-1)
					} else {
						m.encoder.updateCardinality(m.solver, newCost/currentWeight-1)
					}
					localCost = newCost
				}
			}
		} else {
			m.nbCores++
			if currentWeight == minWeight {
				if len(m.model) == 0 {
					return resUnsat, true
				} else {
					return resOptimum, true
				}
			} else {
				obj := make([]int32, len(m.objFunction))
				copy(obj, m.objFunction)
				functions = append(functions, obj)
				weights = append(weights, localCost/currentWeight)
				m.lbCost += localCost
				posWeight++
				currentWeight = m.orderWeights[posWeight]
				localCost = 0
				if !m.foundLowerBound(m.lbCost, nil) {
					return resUndef, false
				}
				m.solver = m.rebuildBmo(functions, weights, currentWeight)
			}
		}
	}
}

func (m *linearSU) normalSearch() (result, bool) {
	m.initRelaxation()
	m.solver = m.rebuildSolver(1)
	for {
		satHandler := m.satHandler()
		res, ok := searchSatSolver(m.solver, satHandler)
		if !ok {
			return resUndef, false
		} else if res == f.TristateTrue {
			m.nbSatisfiable++
			newCost := m.computeCostModel(m.solver.Model(), math.MaxInt)
			m.saveModel(m.solver.Model())
			if newCost == 0 {
				m.ubCost = newCost
				return resOptimum, true
			} else {
				if m.problemType == weighted {
					if !m.encoder.hasPBEncoding() {
						m.encoder.encodePB(m.solver, &m.objFunction, &m.coeffs, newCost-1)
					} else {
						m.encoder.updatePB(m.solver, newCost-1)
					}
				} else {
					if !m.encoder.hasCardEncoding() {
						m.encoder.encodeCardinality(m.solver, m.objFunction, newCost-1)
					} else {
						m.encoder.updateCardinality(m.solver, newCost-1)
					}
				}
				m.ubCost = newCost
				if !m.foundUpperBound(m.ubCost, nil) {
					return resUndef, false
				}
			}
		} else {
			m.nbCores++
			if len(m.model) == 0 {
				return resUnsat, true
			} else {
				return resOptimum, true
			}
		}
	}
}

func (m *linearSU) rebuildSolver(minWeight int) *sat.CoreSolver {
	s := m.newSatSolver()
	for i := 0; i < m.nVars(); i++ {
		newSatVariable(s)
	}
	for i := 0; i < m.nHard(); i++ {
		s.AddClause(m.hardClauses[i].clause, nil)
	}
	for i := 0; i < m.nSoft(); i++ {
		if m.softClauses[i].weight < minWeight {
			continue
		}
		clause := make([]int32, len(m.softClauses[i].clause))
		copy(clause, m.softClauses[i].clause)
		for j := 0; j < len(m.softClauses[i].relaxationVars); j++ {
			clause = append(clause, m.softClauses[i].relaxationVars[j])
		}
		s.AddClause(clause, nil)
	}
	return s
}

func (m *linearSU) rebuildBmo(functions [][]int32, rhs []int, currentWeight int) *sat.CoreSolver {
	s := m.rebuildSolver(currentWeight)
	m.objFunction = []int32{}
	m.coeffs = []int{}
	for i := 0; i < m.nSoft(); i++ {
		if m.softClauses[i].weight == currentWeight {
			m.objFunction = append(m.objFunction, m.softClauses[i].relaxationVars[0])
			m.coeffs = append(m.coeffs, m.softClauses[i].weight)
		}
	}
	for i := 0; i < len(functions); i++ {
		m.encoder.encodeCardinality(s, functions[i], rhs[i])
	}
	return s
}

func (m *linearSU) initRelaxation() {
	for _, softClause := range m.softClauses {
		l := m.newLiteral(false)
		softClause.relaxationVars = append(softClause.relaxationVars, l)
		m.objFunction = append(m.objFunction, l)
		m.coeffs = append(m.coeffs, softClause.weight)
	}
}
