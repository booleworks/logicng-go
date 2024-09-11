package maxsat

import (
	"math"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/sat"
)

type linearUS struct {
	*maxSatAlgorithm
	encoder             *encoder
	incrementalStrategy IncrementalStrategy
}

func newLinearUS(fac f.Factory, config ...*Config) *linearUS {
	var cfg *Config
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = DefaultConfig()
	}
	return &linearUS{
		maxSatAlgorithm:     newAlgorithm(fac, cfg),
		incrementalStrategy: cfg.IncrementalStrategy,
	}
}

func (m *linearUS) search(hdl handler.Handler) (result, handler.State) {
	m.encoder = newEncoder()
	if m.problemType == weighted {
		panic(errorx.BadInput("linearUS does not support weighted MaxSAT instances"))
	}
	return m.innerSearch(hdl, func() (result, handler.State) {
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

func (m *linearUS) none() (result, handler.State) {
	m.nbInitialVariables = m.nVars()
	objFunction := []int32{}
	m.initRelaxation(&objFunction)
	solver := m.rebuildSolver()
	m.encoder.setIncremental(IncNone)
	for {
		res, state := searchSatSolver(solver, m.hdl)
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
				m.encoder.encodeCardinality(solver, objFunction, 0)
			} else {
				return resOptimum, succ
			}
		} else {
			m.lbCost++
			if m.nbSatisfiable == 0 {
				return resUnsat, succ
			} else if m.lbCost == m.ubCost {
				if m.nbSatisfiable > 0 {
					return resOptimum, succ
				} else {
					return resUnsat, succ
				}
			} else if state := m.foundLowerBound(m.lbCost); !state.Success {
				return resUndef, state
			}
			solver = m.rebuildSolver()
			m.encoder.encodeCardinality(solver, objFunction, m.lbCost)
		}
	}
}

func (m *linearUS) iterative() (result, handler.State) {
	objFunction := []int32{}
	m.nbInitialVariables = m.nVars()
	m.initRelaxation(&objFunction)
	solver := m.rebuildSolver()
	var assumptions []int32
	m.encoder.setIncremental(IncIterative)
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
			m.nbCores++
			m.lbCost++
			if m.nbSatisfiable == 0 {
				return resUnsat, succ
			}
			if m.lbCost == m.ubCost {
				if m.nbSatisfiable > 0 {
					return resOptimum, succ
				} else {
					return resUnsat, succ
				}
			}
			if state := m.foundLowerBound(m.lbCost); !state.Success {
				return resUndef, state
			}
			if !m.encoder.hasCardEncoding() {
				m.encoder.buildCardinality(solver, objFunction, m.lbCost)
			}
			m.encoder.incUpdateCardinality(solver, []int32{}, m.lbCost, &assumptions)
		}
	}
}

func (m *linearUS) rebuildSolver() *sat.CoreSolver {
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
		clause = append(clause, m.softClauses[i].relaxationVars...)
		s.AddClause(clause, nil)
	}
	return s
}

func (m *linearUS) initRelaxation(objFunction *[]int32) {
	for i := 0; i < m.nSoft(); i++ {
		l := m.newLiteral(false)
		m.softClauses[i].relaxationVars = append(m.softClauses[i].relaxationVars, l)
		m.softClauses[i].assumptionVar = l
		*objFunction = append(*objFunction, l)
	}
}
