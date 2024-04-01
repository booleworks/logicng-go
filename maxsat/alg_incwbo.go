package maxsat

import (
	"booleworks.com/logicng/errorx"
	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/sat"
)

type incWBO struct {
	*wbo
	encoder    *encoder
	incSoft    []bool
	firstBuild bool
}

func newIncWBO(config *Config) *incWBO {
	return &incWBO{
		wbo:        newWBO(config),
		encoder:    newEncoder(),
		incSoft:    []bool{},
		firstBuild: true,
	}
}

func (m *incWBO) search(handler Handler) (result, bool) {
	m.nbInitialVariables = m.nVars()
	if m.currentWeight == 1 {
		m.problemType = unweighted
		m.weightStrategy = WeightNone
	}
	return m.innerSearch(handler, func() (result, bool) {
		if m.symmetryStrategy {
			m.initSymmetry()
		}
		if m.problemType == unweighted || m.weightStrategy == WeightNone {
			return m.normalSearchInc()
		} else if m.weightStrategy == WeightNormal || m.weightStrategy == WeightDiversify {
			return m.weightSearchInc()
		}
		panic(errorx.UnknownEnumValue(m.problemType))
	})
}

func (m *incWBO) normalSearchInc() (result, bool) {
	switch m.unsatSearch() {
	case f.TristateUndef:
		return resUndef, false
	case f.TristateFalse:
		return resUnsat, true
	}

	m.initAssumptions(&m.assumptions)
	m.solver = m.rebuildSolver()
	m.incSoft = make([]bool, m.nSoft())
	for {
		m.assumptions = []int32{}
		for i := 0; i < len(m.incSoft); i++ {
			if !m.incSoft[i] {
				m.assumptions = append(m.assumptions, sat.Not(m.softClauses[i].assumptionVar))
			}
		}
		satHandler := m.satHandler()
		res, ok := searchSatSolverWithAssumptions(m.solver, satHandler, m.assumptions)
		if !ok {
			return resUndef, false
		} else if res == f.TristateFalse {
			m.nbCores++
			coreCost := m.computeCostCore(m.solver.Conflict())
			m.lbCost += coreCost
			if m.lbCost == m.ubCost {
				return resOptimum, true
			}
			if !m.foundLowerBound(m.lbCost, nil) {
				return resUndef, false
			}
			m.relaxCoreInc(m.solver.Conflict(), coreCost)
		} else {
			m.nbSatisfiable++
			m.ubCost = m.incComputeCostModel(m.solver.Model())
			m.saveModel(m.solver.Model())
			return resOptimum, true
		}
	}
}

func (m *incWBO) relaxCoreInc(conflict []int32, weightCore int) {
	var lits []int32
	for i := 0; i < len(conflict); i++ {
		indexSoft := m.coreMapping[conflict[i]]
		if m.softClauses[indexSoft].weight == weightCore {
			clause := make([]int32, len(m.softClauses[indexSoft].clause))
			copy(clause, m.softClauses[indexSoft].clause)
			vars := make([]int32, len(m.softClauses[indexSoft].relaxationVars))
			copy(vars, m.softClauses[indexSoft].relaxationVars)
			p := m.newLiteral(false)
			newSatVariable(m.solver)
			vars = append(vars, p)
			lits = append(lits, p)
			m.addSoftClauseWithAssumptions(weightCore, clause, vars)
			l := m.newLiteral(false)
			newSatVariable(m.solver)
			m.softClauses[m.nSoft()-1].assumptionVar = l
			m.coreMapping[l] = m.nSoft() - 1
			m.incSoft[indexSoft] = true
			m.incSoft = append(m.incSoft, false)
			for j := 0; j < len(vars); j++ {
				clause = append(clause, vars[j])
			}
			clause = append(clause, l)
			m.solver.AddClause(clause, nil)
			clause = []int32{m.softClauses[indexSoft].assumptionVar}
			m.solver.AddClause(clause, nil)
			if m.symmetryStrategy {
				cpy := make([]int, len(m.softMapping[indexSoft]))
				copy(cpy, m.softMapping[indexSoft])
				m.softMapping = append(m.softMapping, cpy)
				m.softMapping[indexSoft] = []int{}
				cpy32 := make([]int32, len(m.relaxationMapping[indexSoft]))
				copy(cpy32, m.relaxationMapping[indexSoft])
				m.relaxationMapping = append(m.relaxationMapping, cpy32)
				m.relaxationMapping[indexSoft] = []int32{}
				m.symmetryLog(m.nSoft() - 1)
			}
		} else {
			m.softClauses[indexSoft].weight = m.softClauses[indexSoft].weight - weightCore
			clause := make([]int32, len(m.softClauses[indexSoft].clause))
			copy(clause, m.softClauses[indexSoft].clause)
			vars := make([]int32, len(m.softClauses[indexSoft].relaxationVars))
			copy(vars, m.softClauses[indexSoft].relaxationVars)
			m.addSoftClauseWithAssumptions(m.softClauses[indexSoft].weight, clause, vars)
			if m.symmetryStrategy {
				cpy := make([]int, len(m.softMapping[indexSoft]))
				copy(cpy, m.softMapping[indexSoft])
				m.softMapping = append(m.softMapping, cpy)
				m.softMapping[indexSoft] = []int{}
				cpy32 := make([]int32, len(m.relaxationMapping[indexSoft]))
				copy(cpy32, m.relaxationMapping[indexSoft])
				m.relaxationMapping = append(m.relaxationMapping, cpy32)
				m.relaxationMapping[indexSoft] = []int32{}
			}
			m.incSoft[indexSoft] = true
			l := m.newLiteral(false)
			newSatVariable(m.solver)
			m.softClauses[m.nSoft()-1].assumptionVar = l
			m.coreMapping[l] = m.nSoft() - 1
			m.incSoft = append(m.incSoft, false)
			for j := 0; j < len(vars); j++ {
				clause = append(clause, vars[j])
			}
			clause = append(clause, l)
			m.solver.AddClause(clause, nil)
			clause = make([]int32, len(m.softClauses[indexSoft].clause))
			copy(clause, m.softClauses[indexSoft].clause)
			vars = make([]int32, len(m.softClauses[indexSoft].relaxationVars))
			copy(vars, m.softClauses[indexSoft].relaxationVars)
			l = m.newLiteral(false)
			newSatVariable(m.solver)
			vars = append(vars, l)
			lits = append(lits, l)
			m.addSoftClauseWithAssumptions(weightCore, clause, vars)
			l = m.newLiteral(false)
			newSatVariable(m.solver)
			m.softClauses[m.nSoft()-1].assumptionVar = l
			m.coreMapping[l] = m.nSoft() - 1
			m.incSoft = append(m.incSoft, false)
			for j := 0; j < len(vars); j++ {
				clause = append(clause, vars[j])
			}
			clause = append(clause, l)
			m.solver.AddClause(clause, nil)
			clause = []int32{m.softClauses[indexSoft].assumptionVar}
			m.solver.AddClause(clause, nil)
			if m.symmetryStrategy {
				m.softMapping = append(m.softMapping, []int{})
				m.relaxationMapping = append(m.relaxationMapping, []int32{})
				m.symmetryLog(m.nSoft() - 1)
			}
		}
	}
	m.encoder.encodeAMO(m.solver, lits)
	m.nbVars = int(m.solver.NVars())
	if m.symmetryStrategy {
		m.symmetryBreakingInc()
	}
	m.sumSizeCores += len(conflict)
}

func (m *incWBO) incComputeCostModel(currentModel []bool) int {
	currentCost := 0
	for i := 0; i < m.nSoft(); i++ {
		unsatisfied := true
		for j := 0; j < len(m.softClauses[i].clause); j++ {
			if m.incSoft[i] {
				unsatisfied = false
				continue
			}
			if sat.Sign(m.softClauses[i].clause[j]) && !currentModel[sat.Vari(m.softClauses[i].clause[j])] ||
				!sat.Sign(m.softClauses[i].clause[j]) && currentModel[sat.Vari(m.softClauses[i].clause[j])] {
				unsatisfied = false
				break
			}
		}
		if unsatisfied {
			currentCost += m.softClauses[i].weight
		}
	}
	return currentCost
}

func (m *incWBO) symmetryBreakingInc() {
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
						clause := []int32{
							sat.Not(coreIntersection[coreList[k]][n]),
							sat.Not(coreIntersectionCurrent[coreList[k]][j]),
						}
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
							m.solver.AddClause(clause, nil)
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

func (m *incWBO) weightSearchInc() (result, bool) {
	switch m.unsatSearch() {
	case f.TristateUndef:
		return resUndef, false
	case f.TristateFalse:
		return resUnsat, true
	}

	m.initAssumptions(&m.assumptions)
	m.updateCurrentWeight(m.weightStrategy)
	m.incrementalBuildWeightSolver()
	m.incSoft = make([]bool, m.nSoft())
	for {
		m.assumptions = []int32{}
		for i := 0; i < len(m.incSoft); i++ {
			if !m.incSoft[i] {
				m.assumptions = append(m.assumptions, sat.Not(m.softClauses[i].assumptionVar))
			}
		}
		satHandler := m.satHandler()
		res, ok := searchSatSolverWithAssumptions(m.solver, satHandler, m.assumptions)
		if !ok {
			return resUndef, false
		} else if res == f.TristateFalse {
			m.nbCores++
			coreCost := m.computeCostCore(m.solver.Conflict())
			m.lbCost += coreCost
			if !m.foundLowerBound(m.lbCost, nil) {
				return resUndef, false
			}
			m.relaxCoreInc(m.solver.Conflict(), coreCost)
			m.incrementalBuildWeightSolver()
		} else {
			m.nbSatisfiable++
			if m.nbCurrentSoft == m.nSoft() {
				if m.lbCost < m.ubCost {
					m.ubCost = m.lbCost
					m.saveModel(m.solver.Model())
				}
				return resOptimum, true
			} else {
				m.updateCurrentWeight(m.weightStrategy)
				cost := m.incComputeCostModel(m.solver.Model())
				if cost < m.ubCost {
					m.ubCost = cost
					m.saveModel(m.solver.Model())
				}
				if m.lbCost == m.ubCost {
					return resOptimum, true
				} else if !m.foundUpperBound(m.ubCost, nil) {
					return resUndef, false
				}
				m.incrementalBuildWeightSolver()
			}
		}
	}
}

func (m *incWBO) incrementalBuildWeightSolver() {
	if m.firstBuild {
		m.solver = m.newSatSolver()
		for i := 0; i < m.nVars(); i++ {
			newSatVariable(m.solver)
		}
		for i := 0; i < m.nHard(); i++ {
			m.solver.AddClause(m.hardClauses[i].clause, nil)
		}
		if m.symmetryStrategy {
			m.symmetryBreakingInc()
		}
		m.firstBuild = false
	}
	m.nbCurrentSoft = 0
	for i := 0; i < m.nSoft(); i++ {
		if m.softClauses[i].weight >= m.currentWeight && m.softClauses[i].weight != 0 {
			m.nbCurrentSoft++
			clause := make([]int32, len(m.softClauses[i].clause))
			copy(clause, m.softClauses[i].clause)
			for j := 0; j < len(m.softClauses[i].relaxationVars); j++ {
				clause = append(clause, m.softClauses[i].relaxationVars[j])
			}
			clause = append(clause, m.softClauses[i].assumptionVar)
			m.solver.AddClause(clause, nil)
		}
	}
}
