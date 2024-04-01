package maxsat

import (
	"math"

	"booleworks.com/logicng/errorx"
	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/sat"
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

func (m *linearUS) search(handler Handler) (result, bool) {
	if m.problemType == weighted {
		panic(errorx.BadInput("linearUS does not support weighted MaxSAT instances"))
	}
	return m.innerSearch(handler, func() (result, bool) {
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

func (m *linearUS) none() (result, bool) {
	m.nbInitialVariables = m.nVars()
	m.initRelaxation()
	m.solver = m.rebuildSolver()
	m.encoder.setIncremental(IncNone)
	for {
		satHandler := m.satHandler()
		res, ok := searchSatSolver(m.solver, satHandler)
		if !ok {
			return resUndef, false
		} else if res == f.TristateTrue {
			m.nbSatisfiable++
			newCost := m.computeCostModel(m.solver.Model(), math.MaxInt)
			m.saveModel(m.solver.Model())
			m.ubCost = newCost
			if m.nbSatisfiable == 1 {
				if !m.foundUpperBound(m.ubCost, nil) {
					return resUndef, false
				}
				m.encoder.encodeCardinality(m.solver, m.objFunction, 0)
			} else {
				return resOptimum, true
			}
		} else {
			m.lbCost++
			if m.nbSatisfiable == 0 {
				return resUnsat, true
			} else if m.lbCost == m.ubCost {
				if m.nbSatisfiable > 0 {
					return resOptimum, true
				} else {
					return resUnsat, true
				}
			} else if !m.foundLowerBound(m.lbCost, nil) {
				return resUndef, false
			}
			m.solver = m.rebuildSolver()
			m.encoder.encodeCardinality(m.solver, m.objFunction, m.lbCost)
		}
	}
}

func (m *linearUS) iterative() (result, bool) {
	m.nbInitialVariables = m.nVars()
	m.initRelaxation()
	m.solver = m.rebuildSolver()
	var assumptions []int32
	m.encoder.setIncremental(IncIterative)
	for {
		satHandler := m.satHandler()
		res, ok := searchSatSolverWithAssumptions(m.solver, satHandler, assumptions)
		if !ok {
			return resUndef, false
		} else if res == f.TristateTrue {
			m.nbSatisfiable++
			newCost := m.computeCostModel(m.solver.Model(), math.MaxInt)
			m.saveModel(m.solver.Model())
			m.ubCost = newCost
			if m.nbSatisfiable == 1 {
				if !m.foundUpperBound(m.ubCost, nil) {
					return resUndef, false
				}
				for i := 0; i < len(m.objFunction); i++ {
					assumptions = append(assumptions, sat.Not(m.objFunction[i]))
				}
			} else {
				return resOptimum, true
			}
		} else {
			m.nbCores++
			m.lbCost++
			if m.nbSatisfiable == 0 {
				return resUnsat, true
			}
			if m.lbCost == m.ubCost {
				if m.nbSatisfiable > 0 {
					return resOptimum, true
				} else {
					return resUnsat, true
				}
			}
			if !m.foundLowerBound(m.lbCost, nil) {
				return resUndef, false
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
