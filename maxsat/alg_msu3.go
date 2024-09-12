package maxsat

import (
	"math"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/sat"
)

type msu3 struct {
	*maxSatAlgorithm
	encoder             *encoder
	incrementalStrategy IncrementalStrategy
}

func newMSU3(fac f.Factory, config ...*Config) *msu3 {
	var cfg *Config
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = DefaultConfig()
	}
	return &msu3{
		maxSatAlgorithm:     newAlgorithm(fac, cfg),
		incrementalStrategy: cfg.IncrementalStrategy,
	}
}

func (m *msu3) search(maxHandler handler.Handler) (result, handler.State) {
	m.encoder = newEncoder()
	if m.problemType == weighted {
		panic(errorx.BadInput("msu3 does not support weighted MaxSAT instances"))
	}
	return m.innerSearch(maxHandler, func() (result, handler.State) {
		switch m.incrementalStrategy {
		case IncNone:
			return m.none()
		case IncIterative:
			return m.iterative()
		default:
			panic(errorx.UnknownEnumValue(m.incrementalStrategy))
		}
	})
}

func (m *msu3) none() (result, handler.State) {
	m.nbInitialVariables = m.nVars()
	var objFunction []int32
	coreMapping := make(map[int32]int)
	m.initRelaxation(&objFunction)
	solver := m.rebuildSolver()
	var assumptions []int32
	m.encoder.setIncremental(IncNone)
	activeSoft := make([]bool, m.nSoft())
	for i := 0; i < m.nSoft(); i++ {
		coreMapping[m.softClauses[i].assumptionVar] = i
	}
	for {
		res, state := searchSatSolverWithAssumptions(solver, m.hdl, assumptions)
		if !state.Success {
			return resUndef, state
		} else if res == f.TristateTrue {
			m.nbSatisfiable++
			newCost := m.computeCostModel(solver.Model(), math.MaxInt)
			m.saveModel(solver.Model())
			m.ubCost = newCost
			if m.nbSatisfiable == 1 {
				if state := m.foundUpperBound(m.ubCost); !state.Success {
					return resUndef, state
				}
				for i := 0; i < len(objFunction); i++ {
					assumptions = append(assumptions, sat.Not(objFunction[i]))
				}
			} else {
				return resOptimum, succ
			}
		} else {
			m.lbCost++
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
				activeSoft[coreMapping[solver.Conflict()[i]]] = true
			}
			var currentObjFunction []int32
			assumptions = []int32{}
			for i := 0; i < m.nSoft(); i++ {
				if activeSoft[i] {
					currentObjFunction = append(currentObjFunction, m.softClauses[i].relaxationVars[0])
				} else {
					assumptions = append(assumptions, sat.Not(m.softClauses[i].assumptionVar))
				}
			}
			solver = m.rebuildSolver()
			m.encoder.encodeCardinality(solver, currentObjFunction, m.lbCost)
		}
	}
}

func (m *msu3) iterative() (result, handler.State) {
	m.nbInitialVariables = m.nVars()
	var objFunction []int32
	coreMapping := make(map[int32]int)
	m.initRelaxation(&objFunction)
	solver := m.rebuildSolver()
	m.encoder.setIncremental(IncIterative)
	activeSoft := make([]bool, m.nSoft())
	for i := 0; i < m.nSoft(); i++ {
		coreMapping[m.softClauses[i].assumptionVar] = i
	}
	var assumptions []int32
	for {
		res, state := searchSatSolverWithAssumptions(solver, m.hdl, assumptions)
		if !state.Success {
			return resUndef, state
		} else if res == f.TristateTrue {
			m.nbSatisfiable++
			newCost := m.computeCostModel(solver.Model(), math.MaxInt)
			m.saveModel(solver.Model())
			m.ubCost = newCost
			if m.nbSatisfiable == 1 {
				if state := m.foundUpperBound(m.ubCost); !state.Success {
					return resUndef, state
				}
				for i := 0; i < len(objFunction); i++ {
					assumptions = append(assumptions, sat.Not(objFunction[i]))
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
			if len(solver.Conflict()) == 0 {
				return resUnsat, succ
			}
			if state := m.foundLowerBound(m.lbCost); !state.Success {
				return resUndef, state
			}
			var joinObjFunction []int32
			for i := 0; i < len(solver.Conflict()); i++ {
				entry, ok := coreMapping[solver.Conflict()[i]]
				if ok {
					activeSoft[coreMapping[solver.Conflict()[i]]] = true
					joinObjFunction = append(joinObjFunction, m.softClauses[entry].relaxationVars[0])
				}
			}
			var currentObjFunction []int32
			assumptions = []int32{}
			for i := 0; i < m.nSoft(); i++ {
				if activeSoft[i] {
					currentObjFunction = append(currentObjFunction, m.softClauses[i].relaxationVars[0])
				} else {
					assumptions = append(assumptions, sat.Not(m.softClauses[i].assumptionVar))
				}
			}
			var encodingAssumptions []int32
			if !m.encoder.hasCardEncoding() {
				if m.lbCost != len(currentObjFunction) {
					m.encoder.buildCardinality(solver, currentObjFunction, m.lbCost)
					joinObjFunction = []int32{}
					m.encoder.incUpdateCardinality(solver, joinObjFunction, m.lbCost, &encodingAssumptions)
				}
			} else {
				m.encoder.incUpdateCardinality(solver, joinObjFunction, m.lbCost, &encodingAssumptions)
			}
			assumptions = append(assumptions, encodingAssumptions...)
		}
	}
}

func (m *msu3) rebuildSolver() *sat.CoreSolver {
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

func (m *msu3) initRelaxation(objFunction *[]int32) {
	for i := 0; i < m.nSoft(); i++ {
		l := m.newLiteral(false)
		m.softClauses[i].relaxationVars = append(m.softClauses[i].relaxationVars, l)
		m.softClauses[i].assumptionVar = l
		*objFunction = append(*objFunction, l)
	}
}
