package maxsat

import (
	"math"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/sat"
)

type msu3 struct {
	*maxSatAlgorithm
	encoder             *encoder
	incrementalStrategy IncrementalStrategy
	objFunction         []int32
	coreMapping         map[int32]int
	activeSoft          []bool
	solver              *sat.CoreSolver
}

func newMSU3(config ...*Config) *msu3 {
	var cfg *Config
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = DefaultConfig()
	}
	return &msu3{
		maxSatAlgorithm:     newAlgorithm(),
		solver:              nil,
		incrementalStrategy: cfg.IncrementalStrategy,
		encoder:             newEncoder(),
		objFunction:         []int32{},
		coreMapping:         make(map[int32]int),
		activeSoft:          []bool{},
	}
}

func (m *msu3) search(handler Handler) (result, bool) {
	if m.problemType == weighted {
		panic(errorx.BadInput("msu3 does not support weighted MaxSAT instances"))
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

func (m *msu3) none() (result, bool) {
	m.nbInitialVariables = m.nVars()
	m.initRelaxation()
	m.solver = m.rebuildSolver()
	var assumptions []int32
	m.encoder.setIncremental(IncNone)
	m.activeSoft = make([]bool, m.nSoft())
	for i := 0; i < m.nSoft(); i++ {
		m.coreMapping[m.softClauses[i].assumptionVar] = i
	}
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
			m.lbCost++
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
				m.activeSoft[m.coreMapping[m.solver.Conflict()[i]]] = true
			}
			var currentObjFunction []int32
			assumptions = []int32{}
			for i := 0; i < m.nSoft(); i++ {
				if m.activeSoft[i] {
					currentObjFunction = append(currentObjFunction, m.softClauses[i].relaxationVars[0])
				} else {
					assumptions = append(assumptions, sat.Not(m.softClauses[i].assumptionVar))
				}
			}
			m.solver = m.rebuildSolver()
			m.encoder.encodeCardinality(m.solver, currentObjFunction, m.lbCost)
		}
	}
}

func (m *msu3) iterative() (result, bool) {
	m.nbInitialVariables = m.nVars()
	m.initRelaxation()
	m.solver = m.rebuildSolver()
	m.encoder.setIncremental(IncIterative)
	m.activeSoft = make([]bool, m.nSoft())
	for i := 0; i < m.nSoft(); i++ {
		m.coreMapping[m.softClauses[i].assumptionVar] = i
	}
	var assumptions []int32
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
			m.lbCost++
			m.nbCores++
			if m.nbSatisfiable == 0 {
				return resUnsat, true
			}
			if m.lbCost == m.ubCost {
				return resOptimum, true
			}
			m.sumSizeCores += len(m.solver.Conflict())
			if len(m.solver.Conflict()) == 0 {
				return resUnsat, true
			}
			if !m.foundLowerBound(m.lbCost, nil) {
				return resUndef, false
			}
			var joinObjFunction []int32
			for i := 0; i < len(m.solver.Conflict()); i++ {
				entry, ok := m.coreMapping[m.solver.Conflict()[i]]
				if ok {
					m.activeSoft[m.coreMapping[m.solver.Conflict()[i]]] = true
					joinObjFunction = append(joinObjFunction, m.softClauses[entry].relaxationVars[0])
				}
			}
			var currentObjFunction []int32
			assumptions = []int32{}
			for i := 0; i < m.nSoft(); i++ {
				if m.activeSoft[i] {
					currentObjFunction = append(currentObjFunction, m.softClauses[i].relaxationVars[0])
				} else {
					assumptions = append(assumptions, sat.Not(m.softClauses[i].assumptionVar))
				}
			}
			var encodingAssumptions []int32
			if !m.encoder.hasCardEncoding() {
				if m.lbCost != len(currentObjFunction) {
					m.encoder.buildCardinality(m.solver, currentObjFunction, m.lbCost)
					joinObjFunction = []int32{}
					m.encoder.incUpdateCardinality(m.solver, joinObjFunction, m.lbCost, &encodingAssumptions)
				}
			} else {
				m.encoder.incUpdateCardinality(m.solver, joinObjFunction, m.lbCost, &encodingAssumptions)
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

func (m *msu3) initRelaxation() {
	for i := 0; i < m.nbSoft; i++ {
		l := m.newLiteral(false)
		m.softClauses[i].relaxationVars = append(m.softClauses[i].relaxationVars, l)
		m.softClauses[i].assumptionVar = l
		m.objFunction = append(m.objFunction, l)
	}
}
