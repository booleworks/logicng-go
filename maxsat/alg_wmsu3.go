package maxsat

import (
	"math"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/sat"
)

type wmsu3 struct {
	*maxSatAlgorithm
	encoder             *encoder
	bmoMode             bool
	incrementalStrategy IncrementalStrategy
}

func newWMSU3(fac f.Factory, config ...*Config) *wmsu3 {
	var cfg *Config
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = DefaultConfig()
		cfg.IncrementalStrategy = IncIterative
	}
	return &wmsu3{
		maxSatAlgorithm:     newAlgorithm(fac, cfg),
		incrementalStrategy: cfg.IncrementalStrategy,
		bmoMode:             cfg.BMO,
	}
}

func (m *wmsu3) search(hdl handler.Handler) (result, handler.State) {
	if m.problemType == unweighted {
		panic(errorx.BadInput("wmsu3 does not support unweighted MaxSAT instances"))
	}
	m.encoder = newEncoder()
	isBMO := m.bmoMode && m.isBmo(true)
	if !isBMO {
		m.currentWeight = 1
	}
	return m.innerSearch(hdl, func() (result, handler.State) {
		switch m.incrementalStrategy {
		case IncNone:
			return m.none()
		case IncIterative:
			if isBMO {
				return m.iterativeBmo()
			} else {
				return m.iterative()
			}
		default:
			panic(errorx.UnknownEnumValue(m.incrementalStrategy))
		}
	})
}

func (m *wmsu3) none() (result, handler.State) {
	coreMapping := make(map[int32]int)
	m.nbInitialVariables = m.nVars()
	m.initRelaxation()
	solver := m.rebuildSolver()
	m.encoder.setIncremental(IncNone)
	activeSoft := make([]bool, m.nSoft())
	for i := 0; i < m.nSoft(); i++ {
		coreMapping[m.softClauses[i].assumptionVar] = i
	}
	var assumptions []int32
	var coeffs []int
	for {
		res, state := searchSatSolverWithAssumptions(solver, m.hdl, assumptions)
		if !state.Success {
			return resUndef, state
		} else if res == f.TristateTrue {
			m.nbSatisfiable++
			newCost := m.computeCostModel(solver.Model(), math.MaxInt)
			if newCost < m.ubCost || m.nbSatisfiable == 1 {
				m.saveModel(solver.Model())
				m.ubCost = newCost
			}
			if m.ubCost == 0 || m.lbCost == m.ubCost || (m.currentWeight == 1 && m.nbSatisfiable > 1) {
				return resOptimum, succ
			} else if state := m.foundUpperBound(m.ubCost); !state.Success {
				return resUndef, state
			}
			for i := 0; i < m.nSoft(); i++ {
				if m.softClauses[i].weight >= m.currentWeight && !activeSoft[i] {
					assumptions = append(assumptions, sat.Not(m.softClauses[i].assumptionVar))
				}
			}
		} else {
			m.nbCores++
			if m.nbSatisfiable == 0 {
				return resUnsat, succ
			} else if m.lbCost == m.ubCost {
				return resOptimum, succ
			} else if state := m.foundLowerBound(m.lbCost); !state.Success {
				return resUndef, state
			}
			m.sumSizeCores += len(solver.Conflict())
			for i := 0; i < len(solver.Conflict()); i++ {
				indexSoft := coreMapping[solver.Conflict()[i]]
				activeSoft[indexSoft] = true
			}
			var objFunction []int32
			coeffs = []int{}
			assumptions = []int32{}
			for i := 0; i < m.nSoft(); i++ {
				if activeSoft[i] {
					objFunction = append(objFunction, m.softClauses[i].relaxationVars[0])
					coeffs = append(coeffs, m.softClauses[i].weight)
				} else if m.softClauses[i].weight >= m.currentWeight {
					assumptions = append(assumptions, sat.Not(m.softClauses[i].assumptionVar))
				}
			}
			solver = m.rebuildSolver()
			m.lbCost++
			for !subsetSum(coeffs, m.lbCost) {
				m.lbCost++
			}
			m.encoder.encodePB(solver, &objFunction, &coeffs, m.lbCost)
		}
	}
}

func (m *wmsu3) iterative() (result, handler.State) {
	coreMapping := make(map[int32]int)
	m.nbInitialVariables = m.nVars()
	m.initRelaxation()
	solver := m.rebuildSolver()
	m.encoder.setIncremental(IncIterative)
	activeSoft := make([]bool, m.nSoft())
	for i := 0; i < m.nSoft(); i++ {
		coreMapping[m.softClauses[i].assumptionVar] = i
	}
	var assumptions []int32
	var coeffs []int
	var fullCoeffsFunction []int
	for {
		res, state := searchSatSolverWithAssumptions(solver, m.hdl, assumptions)
		if !state.Success {
			return resUndef, state
		} else if res == f.TristateTrue {
			m.nbSatisfiable++
			newCost := m.computeCostModel(solver.Model(), math.MaxInt)
			if newCost < m.ubCost || m.nbSatisfiable == 1 {
				m.saveModel(solver.Model())
				m.ubCost = newCost
			}
			if m.ubCost == 0 || m.lbCost == m.ubCost || (m.currentWeight == 1 && m.nbSatisfiable > 1) {
				return resOptimum, succ
			} else if state := m.foundUpperBound(m.ubCost); !state.Success {
				return resUndef, state
			}
			for i := 0; i < m.nSoft(); i++ {
				if m.softClauses[i].weight >= m.currentWeight && !activeSoft[i] {
					assumptions = append(assumptions, sat.Not(m.softClauses[i].assumptionVar))
				}
			}
		} else {
			m.nbCores++
			if m.nbSatisfiable == 0 {
				return resUnsat, succ
			} else if m.lbCost == m.ubCost {
				return resOptimum, succ
			} else if state := m.foundLowerBound(m.lbCost); !state.Success {
				return resUndef, state
			}
			m.sumSizeCores += len(solver.Conflict())
			var objFunction []int32
			coeffs = []int{}
			assumptions = []int32{}
			for i := 0; i < len(solver.Conflict()); i++ {
				indexSoft, ok := coreMapping[solver.Conflict()[i]]
				if !ok {
					continue
				}
				if !activeSoft[indexSoft] {
					activeSoft[indexSoft] = true
					objFunction = append(objFunction, m.softClauses[indexSoft].relaxationVars[0])
					coeffs = append(coeffs, m.softClauses[indexSoft].weight)
				}
			}
			for i := 0; i < m.nSoft(); i++ {
				if !activeSoft[i] && m.softClauses[i].weight >= m.currentWeight {
					assumptions = append(assumptions, sat.Not(m.softClauses[i].assumptionVar))
				}
			}
			for i := 0; i < len(coeffs); i++ {
				fullCoeffsFunction = append(fullCoeffsFunction, coeffs[i])
			}
			m.lbCost++
			for !subsetSum(fullCoeffsFunction, m.lbCost) {
				m.lbCost++
			}
			if !m.encoder.hasPBEncoding() {
				m.encoder.incEncodePB(solver, &objFunction, &coeffs, m.lbCost, &assumptions, m.nSoft())
			} else {
				m.encoder.incUpdatePB(solver, objFunction, coeffs, m.lbCost)
				m.encoder.incUpdatePBAssumptions(&assumptions)
			}
		}
	}
}

func (m *wmsu3) iterativeBmo() (result, handler.State) {
	coreMapping := make(map[int32]int)
	m.nbInitialVariables = m.nVars()
	m.initRelaxation()
	solver := m.rebuildSolver()
	m.encoder.setIncremental(IncIterative)
	var encodingAssumptions []int32
	activeSoft := make([]bool, m.nSoft())
	for i := 0; i < m.nSoft(); i++ {
		coreMapping[m.softClauses[i].assumptionVar] = i
	}
	var minWeight, posWeight, localCost int
	functions := [][]int32{{}}
	weights := []int{0}
	e := newEncoder()
	e.setIncremental(IncIterative)
	bmoEncodings := []*encoder{e}
	firstEncoding := []bool{true}
	var objFunction []int32
	var assumptions []int32
	var coeffs []int
	for {
		res, state := searchSatSolverWithAssumptions(solver, m.hdl, assumptions)
		if !state.Success {
			return resUndef, state
		} else if res == f.TristateTrue {
			m.nbSatisfiable++
			newCost := m.computeCostModel(solver.Model(), math.MaxInt)
			if newCost < m.ubCost || m.nbSatisfiable == 1 {
				m.saveModel(solver.Model())
				m.ubCost = newCost
			}
			if m.nbSatisfiable == 1 {
				if m.ubCost == 0 {
					return resOptimum, succ
				} else if state := m.foundUpperBound(m.ubCost); !state.Success {
					return resUndef, state
				}
				minWeight = m.orderWeights[len(m.orderWeights)-1]
				m.currentWeight = m.orderWeights[0]
				for i := 0; i < m.nSoft(); i++ {
					if m.softClauses[i].weight >= m.currentWeight {
						assumptions = append(assumptions, sat.Not(m.softClauses[i].assumptionVar))
					}
				}
			} else {
				if m.currentWeight == 1 || m.currentWeight == minWeight {
					return resOptimum, succ
				} else {
					if state := m.foundUpperBound(m.ubCost); !state.Success {
						return resUndef, state
					}
					assumptions = []int32{}
					previousWeight := m.currentWeight
					posWeight++
					m.currentWeight = m.orderWeights[posWeight]
					if len(objFunction) > 0 {
						cpy := make([]int32, len(objFunction))
						copy(cpy, objFunction)
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
						solver.AddClause([]int32{encodingAssumptions[i]}, nil)
					}
					encodingAssumptions = []int32{}
					for i := 0; i < m.nSoft(); i++ {
						if !activeSoft[i] && previousWeight == m.softClauses[i].weight {
							solver.AddClause([]int32{sat.Not(m.softClauses[i].assumptionVar)}, nil)
						}
						if m.currentWeight == m.softClauses[i].weight {
							assumptions = append(assumptions, sat.Not(m.softClauses[i].assumptionVar))
						}
						if activeSoft[i] {
							activeSoft[i] = false
						}
					}
				}
			}
		} else {
			localCost++
			m.lbCost += m.currentWeight
			m.nbCores++
			if m.nbSatisfiable == 0 {
				return resUnsat, succ
			} else if m.lbCost == m.ubCost {
				return resOptimum, succ
			} else if state := m.foundLowerBound(m.lbCost); !state.Success {
				return resUndef, state
			}
			m.sumSizeCores += len(solver.Conflict())
			var joinObjFunction []int32
			for i := 0; i < len(solver.Conflict()); i++ {
				entry, ok := coreMapping[solver.Conflict()[i]]
				if ok {
					if activeSoft[entry] {
						continue
					}
					activeSoft[entry] = true
					joinObjFunction = append(joinObjFunction, m.softClauses[entry].relaxationVars[0])
				}
			}
			objFunction = []int32{}
			coeffs = []int{}
			assumptions = []int32{}
			for i := 0; i < m.nSoft(); i++ {
				if activeSoft[i] {
					objFunction = append(objFunction, m.softClauses[i].relaxationVars[0])
					coeffs = append(coeffs, m.softClauses[i].weight)
				} else if m.currentWeight == m.softClauses[i].weight {
					assumptions = append(assumptions, sat.Not(m.softClauses[i].assumptionVar))
				}
			}
			cpy := make([]int32, len(objFunction))
			copy(cpy, objFunction)
			functions[posWeight] = cpy
			weights[posWeight] = localCost
			if firstEncoding[posWeight] {
				if weights[posWeight] != len(objFunction) {
					bmoEncodings[posWeight].buildCardinality(solver, objFunction, weights[posWeight])
					joinObjFunction = []int32{}
					bmoEncodings[posWeight].incUpdateCardinality(solver, joinObjFunction,
						weights[posWeight], &encodingAssumptions)
					firstEncoding[posWeight] = false
				}
			} else {
				bmoEncodings[posWeight].incUpdateCardinality(solver, joinObjFunction,
					weights[posWeight], &encodingAssumptions)
			}
			for i := 0; i < len(encodingAssumptions); i++ {
				assumptions = append(assumptions, encodingAssumptions[i])
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
	for i := 0; i < m.nSoft(); i++ {
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
