package maxsat

import (
	"math"

	"booleworks.com/logicng/errorx"
	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/sat"
)

type wmsu3 struct {
	*maxSatAlgorithm
	encoder             *encoder
	bmoMode             bool
	incrementalStrategy IncrementalStrategy
	assumptions         []int32
	objFunction         []int32
	coeffs              []int
	coreMapping         map[int32]int
	activeSoft          []bool
	solver              *sat.CoreSolver
	bmo                 bool
}

func newWMSU3(config ...*Config) *wmsu3 {
	var cfg *Config
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = DefaultConfig()
		cfg.IncrementalStrategy = IncIterative
	}
	return &wmsu3{
		maxSatAlgorithm:     newAlgorithm(),
		solver:              nil,
		incrementalStrategy: cfg.IncrementalStrategy,
		encoder:             newEncoder(),
		bmoMode:             cfg.BMO,
		bmo:                 false,
		assumptions:         []int32{},
		objFunction:         []int32{},
		coeffs:              []int{},
		coreMapping:         make(map[int32]int),
		activeSoft:          []bool{},
	}
}

func (m *wmsu3) search(handler Handler) (result, bool) {
	if m.problemType == unweighted {
		panic(errorx.BadInput("wmsu3 does not support unweighted MaxSAT instances"))
	}
	if m.bmoMode {
		m.bmo = m.isBmo(true)
	}
	if !m.bmo {
		m.currentWeight = 1
	}
	return m.innerSearch(handler, func() (result, bool) {
		switch m.incrementalStrategy {
		case IncNone:
			return m.none()
		case IncIterative:
			if m.bmo {
				return m.iterativeBmo()
			} else {
				return m.iterative()
			}
		default:
			panic(errorx.UnknownEnumValue(m.incrementalStrategy))
		}
	})
}

func (m *wmsu3) none() (result, bool) {
	m.nbInitialVariables = m.nVars()
	m.initRelaxation()
	m.solver = m.rebuildSolver()
	m.encoder.setIncremental(IncNone)
	m.activeSoft = make([]bool, m.nSoft())
	for i := 0; i < m.nSoft(); i++ {
		m.coreMapping[m.softClauses[i].assumptionVar] = i
	}
	m.assumptions = []int32{}
	for {
		satHandler := m.satHandler()
		res, ok := searchSatSolverWithAssumptions(m.solver, satHandler, m.assumptions)
		if !ok {
			return resUndef, false
		} else if res == f.TristateTrue {
			m.nbSatisfiable++
			newCost := m.computeCostModel(m.solver.Model(), math.MaxInt)
			if newCost < m.ubCost || m.nbSatisfiable == 1 {
				m.saveModel(m.solver.Model())
				m.ubCost = newCost
			}
			if m.ubCost == 0 || m.lbCost == m.ubCost || (m.currentWeight == 1 && m.nbSatisfiable > 1) {
				return resOptimum, true
			} else if !m.foundUpperBound(m.ubCost, nil) {
				return resUndef, true
			}
			for i := 0; i < m.nSoft(); i++ {
				if m.softClauses[i].weight >= m.currentWeight && !m.activeSoft[i] {
					m.assumptions = append(m.assumptions, sat.Not(m.softClauses[i].assumptionVar))
				}
			}
		} else {
			m.nbCores++
			if m.nbSatisfiable == 0 {
				return resUnsat, true
			} else if m.lbCost == m.ubCost {
				return resOptimum, true
			} else if !m.foundLowerBound(m.lbCost, nil) {
				return resUndef, false
			}
			m.sumSizeCores += len(m.solver.Conflict())
			for i := 0; i < len(m.solver.Conflict()); i++ {
				indexSoft := m.coreMapping[m.solver.Conflict()[i]]
				m.activeSoft[indexSoft] = true
			}
			m.objFunction = []int32{}
			m.coeffs = []int{}
			m.assumptions = []int32{}
			for i := 0; i < m.nSoft(); i++ {
				if m.activeSoft[i] {
					m.objFunction = append(m.objFunction, m.softClauses[i].relaxationVars[0])
					m.coeffs = append(m.coeffs, m.softClauses[i].weight)
				} else if m.softClauses[i].weight >= m.currentWeight {
					m.assumptions = append(m.assumptions, sat.Not(m.softClauses[i].assumptionVar))
				}
			}
			m.solver = m.rebuildSolver()
			m.lbCost++
			for !subsetSum(m.coeffs, m.lbCost) {
				m.lbCost++
			}
			m.encoder.encodePB(m.solver, &m.objFunction, &m.coeffs, m.lbCost)
		}
	}
}

func (m *wmsu3) iterative() (result, bool) {
	m.nbInitialVariables = m.nVars()
	m.initRelaxation()
	m.solver = m.rebuildSolver()
	m.encoder.setIncremental(IncIterative)
	m.activeSoft = make([]bool, m.nSoft())
	for i := 0; i < m.nSoft(); i++ {
		m.coreMapping[m.softClauses[i].assumptionVar] = i
	}
	m.assumptions = []int32{}
	var fullCoeffsFunction []int
	for {
		satHandler := m.satHandler()
		res, ok := searchSatSolverWithAssumptions(m.solver, satHandler, m.assumptions)
		if !ok {
			return resUndef, false
		} else if res == f.TristateTrue {
			m.nbSatisfiable++
			newCost := m.computeCostModel(m.solver.Model(), math.MaxInt)
			if newCost < m.ubCost || m.nbSatisfiable == 1 {
				m.saveModel(m.solver.Model())
				m.ubCost = newCost
			}
			if m.ubCost == 0 || m.lbCost == m.ubCost || (m.currentWeight == 1 && m.nbSatisfiable > 1) {
				return resOptimum, true
			} else if !m.foundUpperBound(m.ubCost, nil) {
				return resUndef, true
			}
			for i := 0; i < m.nSoft(); i++ {
				if m.softClauses[i].weight >= m.currentWeight && !m.activeSoft[i] {
					m.assumptions = append(m.assumptions, sat.Not(m.softClauses[i].assumptionVar))
				}
			}
		} else {
			m.nbCores++
			if m.nbSatisfiable == 0 {
				return resUnsat, true
			} else if m.lbCost == m.ubCost {
				return resOptimum, true
			} else if !m.foundLowerBound(m.lbCost, nil) {
				return resUndef, false
			}
			m.sumSizeCores += len(m.solver.Conflict())
			m.objFunction = []int32{}
			m.coeffs = []int{}
			m.assumptions = []int32{}
			for i := 0; i < len(m.solver.Conflict()); i++ {
				indexSoft, ok := m.coreMapping[m.solver.Conflict()[i]]
				if !ok {
					continue
				}
				if !m.activeSoft[indexSoft] {
					m.activeSoft[indexSoft] = true
					m.objFunction = append(m.objFunction, m.softClauses[indexSoft].relaxationVars[0])
					m.coeffs = append(m.coeffs, m.softClauses[indexSoft].weight)
				}
			}
			for i := 0; i < m.nSoft(); i++ {
				if !m.activeSoft[i] && m.softClauses[i].weight >= m.currentWeight {
					m.assumptions = append(m.assumptions, sat.Not(m.softClauses[i].assumptionVar))
				}
			}
			for i := 0; i < len(m.coeffs); i++ {
				fullCoeffsFunction = append(fullCoeffsFunction, m.coeffs[i])
			}
			m.lbCost++
			for !subsetSum(fullCoeffsFunction, m.lbCost) {
				m.lbCost++
			}
			if !m.encoder.hasPBEncoding() {
				m.encoder.incEncodePB(m.solver, &m.objFunction, &m.coeffs, m.lbCost, &m.assumptions, m.nSoft())
			} else {
				m.encoder.incUpdatePB(m.solver, m.objFunction, m.coeffs, m.lbCost)
				m.encoder.incUpdatePBAssumptions(&m.assumptions)
			}
		}
	}
}

func (m *wmsu3) iterativeBmo() (result, bool) {
	m.nbInitialVariables = m.nVars()
	m.initRelaxation()
	m.solver = m.rebuildSolver()
	m.encoder.setIncremental(IncIterative)
	var encodingAssumptions []int32
	m.activeSoft = make([]bool, m.nSoft())
	for i := 0; i < m.nSoft(); i++ {
		m.coreMapping[m.softClauses[i].assumptionVar] = i
	}
	var minWeight, posWeight, localCost int
	functions := [][]int32{{}}
	weights := []int{0}
	e := newEncoder()
	e.setIncremental(IncIterative)
	bmoEncodings := []*encoder{e}
	firstEncoding := []bool{true}
	for {
		satHandler := m.satHandler()
		res, ok := searchSatSolverWithAssumptions(m.solver, satHandler, m.assumptions)
		if !ok {
			return resUndef, false
		} else if res == f.TristateTrue {
			m.nbSatisfiable++
			newCost := m.computeCostModel(m.solver.Model(), math.MaxInt)
			if newCost < m.ubCost || m.nbSatisfiable == 1 {
				m.saveModel(m.solver.Model())
				m.ubCost = newCost
			}
			if m.nbSatisfiable == 1 {
				if m.ubCost == 0 {
					return resOptimum, true
				} else if !m.foundUpperBound(m.ubCost, nil) {
					return resUndef, false
				}
				minWeight = m.orderWeights[len(m.orderWeights)-1]
				m.currentWeight = m.orderWeights[0]
				for i := 0; i < m.nSoft(); i++ {
					if m.softClauses[i].weight >= m.currentWeight {
						m.assumptions = append(m.assumptions, sat.Not(m.softClauses[i].assumptionVar))
					}
				}
			} else {
				if m.currentWeight == 1 || m.currentWeight == minWeight {
					return resOptimum, true
				} else {
					if !m.foundUpperBound(m.ubCost, nil) {
						return resUndef, false
					}
					m.assumptions = []int32{}
					previousWeight := m.currentWeight
					posWeight++
					m.currentWeight = m.orderWeights[posWeight]
					if len(m.objFunction) > 0 {
						cpy := make([]int32, len(m.objFunction))
						copy(cpy, m.objFunction)
						functions[len(functions)-1] = cpy
					}
					functions = append(functions, []int32{})
					weights = append(weights, 0)
					localCost = 0
					e = newEncoder()
					e.setIncremental(IncIterative)
					bmoEncodings = append(bmoEncodings, e)
					firstEncoding = append(firstEncoding, true)
					for i := 0; i < len(encodingAssumptions); i++ {
						m.solver.AddClause([]int32{encodingAssumptions[i]}, nil)
					}
					encodingAssumptions = []int32{}
					for i := 0; i < m.nSoft(); i++ {
						if !m.activeSoft[i] && previousWeight == m.softClauses[i].weight {
							m.solver.AddClause([]int32{sat.Not(m.softClauses[i].assumptionVar)}, nil)
						}
						if m.currentWeight == m.softClauses[i].weight {
							m.assumptions = append(m.assumptions, sat.Not(m.softClauses[i].assumptionVar))
						}
						if m.activeSoft[i] {
							m.activeSoft[i] = false
						}
					}
				}
			}
		} else {
			localCost++
			m.lbCost += m.currentWeight
			m.nbCores++
			if m.nbSatisfiable == 0 {
				return resUnsat, true
			} else if m.lbCost == m.ubCost {
				return resOptimum, true
			} else if !m.foundLowerBound(m.lbCost, nil) {
				return resUndef, false
			}
			m.sumSizeCores += len(m.solver.Conflict())
			var joinObjFunction []int32
			for i := 0; i < len(m.solver.Conflict()); i++ {
				entry, ok := m.coreMapping[m.solver.Conflict()[i]]
				if ok {
					if m.activeSoft[entry] {
						continue
					}
					m.activeSoft[entry] = true
					joinObjFunction = append(joinObjFunction, m.softClauses[entry].relaxationVars[0])
				}
			}
			m.objFunction = []int32{}
			m.coeffs = []int{}
			m.assumptions = []int32{}
			for i := 0; i < m.nSoft(); i++ {
				if m.activeSoft[i] {
					m.objFunction = append(m.objFunction, m.softClauses[i].relaxationVars[0])
					m.coeffs = append(m.coeffs, m.softClauses[i].weight)
				} else if m.currentWeight == m.softClauses[i].weight {
					m.assumptions = append(m.assumptions, sat.Not(m.softClauses[i].assumptionVar))
				}
			}
			cpy := make([]int32, len(m.objFunction))
			copy(cpy, m.objFunction)
			functions[posWeight] = cpy
			weights[posWeight] = localCost
			if firstEncoding[posWeight] {
				if weights[posWeight] != len(m.objFunction) {
					bmoEncodings[posWeight].buildCardinality(m.solver, m.objFunction, weights[posWeight])
					joinObjFunction = []int32{}
					bmoEncodings[posWeight].incUpdateCardinality(m.solver, joinObjFunction,
						weights[posWeight], &encodingAssumptions)
					firstEncoding[posWeight] = false
				}
			} else {
				bmoEncodings[posWeight].incUpdateCardinality(m.solver, joinObjFunction,
					weights[posWeight], &encodingAssumptions)
			}
			for i := 0; i < len(encodingAssumptions); i++ {
				m.assumptions = append(m.assumptions, encodingAssumptions[i])
			}
		}
	}
}

func (m *wmsu3) rebuildSolver() *sat.CoreSolver {
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

func (m *wmsu3) initRelaxation() {
	for i := 0; i < m.nbSoft; i++ {
		l := m.newLiteral(false)
		m.softClauses[i].relaxationVars = append(m.softClauses[i].relaxationVars, l)
		m.softClauses[i].assumptionVar = l
	}
}

func subsetSum(set []int, sum int) bool {
	n := len(set)
	subset := make([][]bool, sum+1)
	for i := 0; i <= sum; i++ {
		subset[i] = make([]bool, n+1)
	}
	for i := 0; i <= n; i++ {
		subset[0][i] = true
	}
	for i := 1; i <= sum; i++ {
		subset[i][0] = false
	}
	for i := 1; i <= sum; i++ {
		for j := 1; j <= n; j++ {
			subset[i][j] = subset[i][j-1]
			if i >= set[j-1] {
				subset[i][j] = subset[i][j] || subset[i-set[j-1]][j-1]
			}
		}
	}
	return subset[sum][n]
}
