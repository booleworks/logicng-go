package maxsat

import (
	"math"
	"slices"

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
	addHardClause(lits []int32)
	addSoftClause(weight int, lits []int32)
	setCurrentWeight(weight int)
	updateSumWeights(weight int)
	getCurrentWeight() int
	setProblemType(problemType problemType)
	getModel() []bool
}

type maxSatAlgorithm struct {
	model              []bool
	softClauses        []*softClause
	hardClauses        []*hardClause
	orderWeights       []int
	hdl                handler.Handler
	hardWeight         int
	problemType        problemType
	nbVars             int
	nbSoft             int
	nbHard             int
	nbInitialVariables int
	nbCores            int
	nbSymmetryClauses  int
	sumSizeCores       int
	nbSatisfiable      int
	ubCost             int
	lbCost             int
	currentWeight      int
}

func newAlgorithm() *maxSatAlgorithm {
	return &maxSatAlgorithm{
		hardClauses:   []*hardClause{},
		softClauses:   []*softClause{},
		hardWeight:    math.MaxInt,
		problemType:   unweighted,
		currentWeight: 1,
		model:         make([]bool, 0),
		orderWeights:  []int{},
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
	if !hdl.ShouldResume(event.MaxSATCallStarted) {
		return resUndef, handler.Cancellation(event.MaxSATCallStarted)
	}
	result, state := search()
	if !hdl.ShouldResume(event.MaxSatCallFinished) {
		return resUndef, handler.Cancellation(event.MaxSatCallFinished)
	}
	m.hdl = nil
	return result, state
}

func (m *maxSatAlgorithm) nVars() int {
	return m.nbVars
}

func (m *maxSatAlgorithm) nSoft() int {
	return m.nbSoft
}

func (m *maxSatAlgorithm) nHard() int {
	return m.nbHard
}

func (m *maxSatAlgorithm) newVar() {
	m.nbVars++
}

func (m *maxSatAlgorithm) addHardClause(lits []int32) {
	m.hardClauses = append(m.hardClauses, newHardClause(lits))
	m.nbHard++
}

func (m *maxSatAlgorithm) addSoftClause(weight int, lits []int32) {
	m.softClauses = append(m.softClauses, newSoftClause(lits, []int32{}, weight, sat.LitUndef))
	m.nbSoft++
}

func (m *maxSatAlgorithm) addSoftClauseWithAssumptions(weight int, lits, vars []int32) {
	m.softClauses = append(m.softClauses, newSoftClause(lits, vars, weight, sat.LitUndef))
	m.nbSoft++
}

func (m *maxSatAlgorithm) newLiteral(sign bool) int32 {
	p := sat.MkLit(int32(m.nVars()), sign)
	m.newVar()
	return p
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
		return handler.Cancellation(e)
	}
}

func (m *maxSatAlgorithm) foundUpperBound(upperBound int) handler.State {
	e := EventMaxSatNewUpperBound{upperBound}
	if m.hdl.ShouldResume(e) {
		return succ
	} else {
		return handler.Cancellation(e)
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
