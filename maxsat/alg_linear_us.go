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
	objFunction         []int32
	solver              *sat.CoreSolver
}

func newLinearUS(config ...*Config) *linearUS {
	var cfg *Config
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = DefaultConfig()
	}
	return &linearUS{
		maxSatAlgorithm:     newAlgorithm(),
		solver:              nil,
		encoder:             newEncoder(),
		incrementalStrategy: cfg.IncrementalStrategy,
		objFunction:         []int32{},
	}
}

func (m *linearUS) search(hdl handler.Handler) (result, handler.State) {
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
	m.initRelaxation()
	m.solver = m.rebuildSolver()
	m.encoder.setIncremental(IncNone)
	for {
		res, state := searchSatSolver(m.solver, m.hdl)
		if !state.Success {
			return resUndef, state
		} else if res == f.TristateTrue {
			m.nbSatisfiable++
			newCost := m.computeCostModel(m.solver.Model(), math.MaxInt)
			m.saveModel(m.solver.Model())
			m.ubCost = newCost
			if m.nbSatisfiable == 1 {
				if state := m.foundUpperBound(m.ubCost); !state.Success {
					return resUndef, state
				}
				m.encoder.encodeCardinality(m.solver, m.objFunction, 0)
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
			m.solver = m.rebuildSolver()
			m.encoder.encodeCardinality(m.solver, m.objFunction, m.lbCost)
		}
	}
}

func (m *linearUS) iterative() (result, handler.State) {
	m.nbInitialVariables = m.nVars()
	m.initRelaxation()
	m.solver = m.rebuildSolver()
	var assumptions []int32
	m.encoder.setIncremental(IncIterative)
	for {
		res, state := searchSatSolverWithAssumptions(m.solver, m.hdl, assumptions)
		if !state.Success {
			return resUndef, state
		} else if res == f.TristateTrue {
			m.nbSatisfiable++
			newCost := m.computeCostModel(m.solver.Model(), math.MaxInt)
			m.saveModel(m.solver.Model())
			m.ubCost = newCost
			if m.nbSatisfiable == 1 {
				if state := m.foundUpperBound(m.ubCost); !state.Success {
					return resUndef, state
				}
				for i := 0; i < len(m.objFunction); i++ {
					assumptions = append(assumptions, sat.Not(m.objFunction[i]))
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
				m.encoder.buildCardinality(m.solver, m.objFunction, m.lbCost)
			}
			m.encoder.incUpdateCardinality(m.solver, []int32{}, m.lbCost, &assumptions)
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

func (m *linearUS) initRelaxation() {
	for i := 0; i < m.nbSoft; i++ {
		l := m.newLiteral(false)
		m.softClauses[i].relaxationVars = append(m.softClauses[i].relaxationVars, l)
		m.softClauses[i].assumptionVar = l
		m.objFunction = append(m.objFunction, l)
	}
}
