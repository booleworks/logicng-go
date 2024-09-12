package maxsat

import (
	"math"
	"slices"

	"github.com/booleworks/logicng-go/errorx"
	"github.com/booleworks/logicng-go/event"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/sat"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets/treeset"
)

var succ = handler.Success()

type algorithm interface {
	search(hdl handler.Handler) (result, handler.State)
	result() int
	newLiteral(bool) int32
	newVar() int32
	addHardClause(lits []int32)
	addSoftClause(weight int, lits []int32)
	setCurrentWeight(weight int)
	updateSumWeights(weight int)
	getCurrentWeight() int
	setProblemType(problemType problemType)
	getModel() []bool
	addClause(formula f.Formula, weight int)
	saveState() *SolverState
	loadState(state *SolverState) error
	literal(lit f.Literal) int32
	addClauseVec(clauseVec []int32, weight int)
	varForIndex(index int) (f.Variable, bool)
}

type maxSatAlgorithm struct {
	fac                f.Factory
	cfg                *Config
	model              []bool
	var2index          map[f.Variable]int32
	index2var          map[int32]f.Variable
	softClauses        []*softClause
	hardClauses        []*hardClause
	orderWeights       []int
	hdl                handler.Handler
	hardWeight         int
	problemType        problemType
	nbVars             int
	nbInitialVariables int
	nbCores            int
	nbSymmetryClauses  int
	sumSizeCores       int
	nbSatisfiable      int
	ubCost             int
	lbCost             int
	currentWeight      int

	stateId     int32
	validStates []int32
}

func newAlgorithm(fac f.Factory, cfg *Config) *maxSatAlgorithm {
	return &maxSatAlgorithm{
		fac:           fac,
		cfg:           cfg,
		var2index:     make(map[f.Variable]int32),
		index2var:     make(map[int32]f.Variable),
		hardClauses:   []*hardClause{},
		softClauses:   []*softClause{},
		hardWeight:    math.MaxInt,
		problemType:   unweighted,
		currentWeight: 1,
		model:         []bool{},
		orderWeights:  []int{},
		validStates:   []int32{},
	}
}

func newSatVariable(s *sat.CoreSolver) {
	s.NewVar(true, true)
}

func searchSatSolver(s *sat.CoreSolver, hdl handler.Handler) (f.Tristate, handler.State) {
	return s.Solve(hdl)
}

func searchSatSolverWithAssumptions(
	s *sat.CoreSolver, hdl handler.Handler, assumptions []int32,
) (f.Tristate, handler.State) {
	return s.SolveWithAssumptions(hdl, assumptions)
}

func (m *maxSatAlgorithm) innerSearch(
	hdl handler.Handler,
	search func() (result, handler.State),
) (result, handler.State) {
	m.hdl = hdl
	if e := event.MaxSATCallStarted; !hdl.ShouldResume(e) {
		return resUndef, handler.Cancelation(e)
	}
	stateBeforeSolving := m.saveState()
	result, state := search()
	if e := event.MaxSatCallFinished; !hdl.ShouldResume(e) {
		return resUndef, handler.Cancelation(e)
	}
	_ = m.loadState(stateBeforeSolving)
	m.hdl = nil
	return result, state
}

func (m *maxSatAlgorithm) nVars() int {
	return m.nbVars
}

func (m *maxSatAlgorithm) nSoft() int {
	return len(m.softClauses)
}

func (m *maxSatAlgorithm) nHard() int {
	return len(m.hardClauses)
}

func (m *maxSatAlgorithm) newVar() int32 {
	n := m.nbVars
	m.nbVars++
	return int32(n)
}

func (m *maxSatAlgorithm) addHardClause(lits []int32) {
	m.hardClauses = append(m.hardClauses, newHardClause(lits))
}

func (m *maxSatAlgorithm) addSoftClause(weight int, lits []int32) {
	m.softClauses = append(m.softClauses, newSoftClause(lits, []int32{}, weight, sat.LitUndef))
}

func (m *maxSatAlgorithm) addSoftClauseWithAssumptions(weight int, lits, vars []int32) {
	m.softClauses = append(m.softClauses, newSoftClause(lits, vars, weight, sat.LitUndef))
}

func (m *maxSatAlgorithm) newLiteral(sign bool) int32 {
	return sat.MkLit(m.newVar(), sign)
}

func (m *maxSatAlgorithm) updateSumWeights(weight int) {
	if weight != m.hardWeight {
		m.ubCost += weight
	}
}

func (m *maxSatAlgorithm) setCurrentWeight(weight int) {
	if weight > m.currentWeight && weight != m.hardWeight {
		m.currentWeight = weight
	}
}

func (m *maxSatAlgorithm) newSatSolver() *sat.CoreSolver {
	return sat.NewCoreSolver(sat.DefaultConfig(), sat.UncheckedEnqueue)
}

func (m *maxSatAlgorithm) saveModel(currentModel []bool) {
	m.model = make([]bool, m.nbInitialVariables)
	for i := 0; i < m.nbInitialVariables; i++ {
		m.model[i] = currentModel[i]
	}
}

func (m *maxSatAlgorithm) computeCostModel(currentModel []bool, weight int) int {
	currentCost := 0
	for i := 0; i < m.nSoft(); i++ {
		unsatisfied := true
		for j := 0; j < len(m.softClauses[i].clause); j++ {
			if weight != math.MaxInt && m.softClauses[i].weight != weight {
				unsatisfied = false
				continue
			}
			if (sat.Sign(m.softClauses[i].clause[j]) &&
				!currentModel[sat.Vari(m.softClauses[i].clause[j])]) ||
				(!sat.Sign(m.softClauses[i].clause[j]) &&
					currentModel[sat.Vari(m.softClauses[i].clause[j])]) {
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

func (m *maxSatAlgorithm) isBmo(cache bool) bool {
	bmo := true
	partitionWeights := treeset.NewWithIntComparator()
	nbPartitionWeights := treemap.NewWithIntComparator()
	for i := 0; i < m.nSoft(); i++ {
		weight := m.softClauses[i].weight
		partitionWeights.Add(weight)
		val, ok := nbPartitionWeights.Get(weight)
		if !ok {
			nbPartitionWeights.Put(weight, 1)
		} else {
			nbPartitionWeights.Put(weight, val.(int)+1)
		}
	}
	partitionWeights.Each(func(_ int, value interface{}) {
		m.orderWeights = append(m.orderWeights, value.(int))
	})
	slices.Sort(m.orderWeights)
	slices.Reverse(m.orderWeights)

	totalWeights := 0
	for i := 0; i < len(m.orderWeights); i++ {
		val, _ := nbPartitionWeights.Get(m.orderWeights[i])
		totalWeights += m.orderWeights[i] * val.(int)
	}
	for i := 0; i < len(m.orderWeights); i++ {
		val, _ := nbPartitionWeights.Get(m.orderWeights[i])
		totalWeights -= m.orderWeights[i] * val.(int)
		if m.orderWeights[i] < totalWeights {
			bmo = false
			break
		}
	}
	if !cache {
		m.orderWeights = []int{}
	}
	return bmo
}

func (m *maxSatAlgorithm) result() int {
	return m.ubCost
}

func (m *maxSatAlgorithm) foundLowerBound(lowerBound int) handler.State {
	e := EventMaxSatNewLowerBound{lowerBound}
	if m.hdl.ShouldResume(e) {
		return succ
	} else {
		return handler.Cancelation(e)
	}
}

func (m *maxSatAlgorithm) foundUpperBound(upperBound int) handler.State {
	e := EventMaxSatNewUpperBound{upperBound}
	if m.hdl.ShouldResume(e) {
		return succ
	} else {
		return handler.Cancelation(e)
	}
}

func (m *maxSatAlgorithm) getCurrentWeight() int {
	return m.currentWeight
}

func (m *maxSatAlgorithm) setProblemType(problemType problemType) {
	m.problemType = problemType
}

func (m *maxSatAlgorithm) getModel() []bool {
	return m.model
}

func (m *maxSatAlgorithm) addClause(formula f.Formula, weight int) {
	clauseVec := make([]int32, f.NumberOfAtoms(m.fac, formula))
	for i, lit := range f.Literals(m.fac, formula).Content() {
		variable := lit.Variable()
		index, ok := m.var2index[variable]
		if !ok {
			index = m.newLiteral(false) >> 1
			m.var2index[variable] = index
			m.index2var[index] = variable
		}
		var litNum int32
		if lit.IsPos() {
			litNum = index * 2
		} else {
			litNum = (index * 2) ^ 1
		}
		clauseVec[i] = litNum
	}
	m.addClauseVec(clauseVec, weight)
}

func (m *maxSatAlgorithm) addClauseVec(clauseVec []int32, weight int) {
	if weight == -1 {
		m.addHardClause(clauseVec)
	} else {
		m.setCurrentWeight(weight)
		m.updateSumWeights(weight)
		m.addSoftClause(weight, clauseVec)
	}
}

func (m *maxSatAlgorithm) literal(lit f.Literal) int32 {
	variable := lit.Variable()
	index, ok := m.var2index[variable]
	if !ok {
		index = m.newLiteral(false) >> 1
		m.var2index[variable] = index
		m.index2var[index] = variable
	}
	if lit.IsPos() {
		return index * 2
	} else {
		return (index * 2) ^ 1
	}
}

func (m *maxSatAlgorithm) saveState() *SolverState {
	softWeights := make([]int, len(m.softClauses))
	for i, c := range m.softClauses {
		softWeights[i] = c.weight
	}
	id := m.stateId
	m.stateId++
	m.validStates = append(m.validStates, id)
	return &SolverState{
		id:            id,
		nbVars:        m.nbVars,
		nbHard:        len(m.hardClauses),
		nbSoft:        len(m.softClauses),
		ubCost:        m.ubCost,
		currentWeight: m.currentWeight,
		softWeights:   softWeights,
	}
}

func (m *maxSatAlgorithm) loadState(state *SolverState) error {
	index := -1
	for i := len(m.validStates) - 1; i >= 0 && index == -1; i-- {
		if m.validStates[i] == state.id {
			index = i
		}
	}
	if index == -1 {
		return errorx.BadInput("solver state %d is not valid any more", state.id)
	}
	shrinkTo(&m.validStates, index+1)

	shrinkTo(&m.hardClauses, state.nbHard)
	shrinkTo(&m.softClauses, state.nbSoft)
	m.orderWeights = []int{}
	for i := int32(state.nbVars); i < int32(m.nbVars); i++ {
		if v, ok := m.index2var[i]; ok {
			delete(m.index2var, i)
			delete(m.var2index, v)
		}
	}
	m.nbVars = state.nbVars
	m.nbCores = 0
	m.nbSymmetryClauses = 0
	m.sumSizeCores = 0
	m.nbSatisfiable = 0
	m.ubCost = state.ubCost
	m.lbCost = 0
	m.currentWeight = state.currentWeight
	for i := 0; i < len(m.softClauses); i++ {
		clause := m.softClauses[i]
		clause.relaxationVars = []int32{}
		clause.weight = state.softWeights[i]
		clause.assumptionVar = sat.LitUndef
	}
	return nil
}

func (m *maxSatAlgorithm) varForIndex(index int) (f.Variable, bool) {
	v, ok := m.index2var[int32(index)]
	return v, ok
}

func shrinkTo[T any](slice *[]T, newSize int) {
	if newSize < len(*slice) {
		*slice = (*slice)[:newSize]
	}
}
