package sat

import (
	"slices"

	"github.com/booleworks/logicng-go/model"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
)

const (
	LitUndef = int32(-1) // constant for an undefined literal
	LitError = int32(-2) // constant for an error state literal
)

const (
	ratioRemoveClauses = 2
	lbBlockingRestart  = 10000
)

// CoreSolver represents a core SAT solver.
//
// The LogicNG solver is based on MiniSat, Glucose, and MiniCard.  Usually you
// should not need to interact with the core solver yourself but only via its
// Solver wrapper.
type CoreSolver struct {
	config   *Config
	llConfig *LowLevelConfig

	ok              bool
	qhead           int
	unitClauses     []int32
	clauses         []*clause
	learnts         []*clause
	watches         [][]*watcher
	vars            []*variable
	orderHeap       lngheap
	trail           []int32
	trailLim        []int
	model           []bool
	conflict        []int32
	assumptions     []int32
	seen            []bool
	analyzeBtLevel  int
	claInc          float64
	clausesLiterals int
	learntsLiterals int

	name2idx map[string]int32
	idx2name map[int32]string

	satHandler        Handler
	canceledByHandler bool

	pgOriginalClauses []proofInformation
	pgProof           [][]int32

	backboneCandidates  []int32
	backboneAssumptions []int32
	backboneMap         map[int32]f.Tristate
	computingBackbone   bool

	varDecay                   float64
	varInc                     float64
	learntsizeAdjustConfl      float64
	learntsizeAdjustCnt        int
	learntsizeAdjustStartConfl int
	learntsizeAdjustInc        float64
	maxLearnts                 float64

	watchesBin            [][]*watcher
	permDiff              []int
	lastDecisionLevel     []int32
	lbdQueue              boundedQueue
	trailQueue            boundedQueue
	myflag                int
	analyzeLBD            int
	nbClausesBeforeReduce int
	conflicts             int
	conflictsRestarts     int
	sumLBD                float64
	curRestart            int

	stateId     int32
	validStates []int32

	enqueueFunction func(m *CoreSolver, lit int32, reason *clause)
}

// NewCoreSolver generates a new core SAT solver with the given configuration.
// An enqueue function to insert new variable assignment decisions must be given.
// The default is the included UncheckedEnqueue function.
func NewCoreSolver(
	config *Config,
	enqueueFunction func(m *CoreSolver, lit int32, reason *clause),
) *CoreSolver {
	s := &CoreSolver{}
	initialize(s, config)
	s.enqueueFunction = enqueueFunction
	return s
}

// UncheckedEnqueue is the default function to add new variable decisions to
// the solver.
func UncheckedEnqueue(m *CoreSolver, lit int32, reason *clause) {
	vari := m.v(lit)
	vari.assignment = f.TristateFromBool(!Sign(lit))
	vari.reason = reason
	vari.level = m.decisionLevel()
	m.trail = append(m.trail, lit)
}

// MkLit generates a solver literals from a solver variable and sign.
func MkLit(v int32, sign bool) int32 {
	var s int32
	if sign {
		s = 1
	} else {
		s = 0
	}
	return v + v + s
}

// Not negates a solver literal.
func Not(lit int32) int32 {
	return lit ^ 1
}

// Sign returns if a literal on the solver is negative.
func Sign(lit int32) bool {
	return (lit & 1) == 1
}

func signAsInt(lit int32) int32 {
	if (lit & 1) == 1 {
		return 1
	} else {
		return 0
	}
}

// Vari returns the solver variable of a solver literal.
func Vari(lit int32) int32 {
	return lit >> 1
}

func initialize(m *CoreSolver, config *Config) {
	m.config = config
	m.llConfig = config.LowLevelConfig
	m.ok = true
	m.qhead = 0
	m.clauses = []*clause{}
	m.learnts = []*clause{}
	m.watches = [][]*watcher{}
	m.vars = []*variable{}
	m.orderHeap = *newLngHeap(m)
	m.trail = []int32{}
	m.trailLim = []int{}
	m.model = make([]bool, len(m.vars))
	m.conflict = []int32{}
	m.assumptions = []int32{}
	m.seen = []bool{}
	m.analyzeBtLevel = 0
	m.claInc = 1
	m.clausesLiterals = 0
	m.learntsLiterals = 0
	m.name2idx = make(map[string]int32)
	m.idx2name = make(map[int32]string)
	m.canceledByHandler = false
	if m.config.ProofGeneration {
		m.pgOriginalClauses = []proofInformation{}
		m.pgProof = [][]int32{}
	}
	m.computingBackbone = false
	m.unitClauses = []int32{}
	m.varInc = m.llConfig.VarInc
	m.varDecay = m.llConfig.VarDecay
	m.learntsizeAdjustConfl = 0
	m.learntsizeAdjustCnt = 0
	m.learntsizeAdjustStartConfl = 100
	m.learntsizeAdjustInc = 1.5
	m.maxLearnts = 0

	m.watchesBin = [][]*watcher{}
	m.permDiff = []int{}
	m.lastDecisionLevel = []int32{}
	m.lbdQueue = *newBoundedQueue()
	m.trailQueue = *newBoundedQueue()
	m.lbdQueue.initSize(m.llConfig.SizeLBDQueue)
	m.trailQueue.initSize(m.llConfig.SizeTrailQueue)
	m.myflag = 0
	m.analyzeBtLevel = 0
	m.analyzeLBD = 0
	m.nbClausesBeforeReduce = m.llConfig.FirstReduceDB
	m.conflicts = 0
	m.conflictsRestarts = 0
	m.sumLBD = 0
	m.curRestart = 1

	m.stateId = 0
	m.validStates = []int32{}
}

func (m *CoreSolver) v(lit int32) *variable {
	return m.vars[lit>>1]
}

func (m *CoreSolver) value(lit int32) f.Tristate {
	if Sign(lit) {
		return m.v(lit).assignment.Negate()
	} else {
		return m.v(lit).assignment
	}
}

func (m *CoreSolver) lt(x, y int32) bool {
	return m.vars[x].activity > m.vars[y].activity
}

func (m *CoreSolver) addName(name string, id int32) {
	m.name2idx[name] = id
	m.idx2name[id] = name
}

// NewVar generates a new variable on the solver.
func (m *CoreSolver) NewVar(sign, dvar bool) int32 {
	v := m.NVars()
	newVar := newVariable(sign)
	m.vars = append(m.vars, newVar)
	m.watches = append(m.watches, []*watcher{}, []*watcher{})
	m.seen = append(m.seen, false)
	m.watchesBin = append(m.watchesBin, []*watcher{}, []*watcher{})
	m.permDiff = append(m.permDiff, 0)
	newVar.decision = dvar
	m.insertVarOrder(v)
	return v
}

func (m *CoreSolver) addUnitClause(lit int32, proposition f.Proposition) bool {
	return m.AddClause([]int32{lit}, proposition)
}

// AddClause adds a new clause to the solver.
func (m *CoreSolver) AddClause(ps []int32, proposition f.Proposition) bool {
	if m.config.ProofGeneration {
		slice := make([]int32, len(ps))
		for i := 0; i < len(ps); i++ {
			slice[i] = (Vari(ps[i]) + 1) * (-2*signAsInt(ps[i]) + 1)
		}
		m.pgOriginalClauses = append(m.pgOriginalClauses, proofInformation{slice, proposition})
	}
	if !m.ok {
		return false
	}
	slices.Sort(ps)

	flag := false
	var oc []int32
	if m.config.ProofGeneration {
		oc = make([]int32, len(ps))
		p := LitUndef
		for i := 0; i < len(ps); i++ {
			oc[i] = ps[i]
			if m.value(ps[i]) == f.TristateTrue || ps[i] == Not(p) || m.value(ps[i]) == f.TristateFalse {
				flag = true
			}
		}
	}

	p := LitUndef
	i, j := 0, 0
	for ; i < len(ps); i++ {
		if m.value(ps[i]) == f.TristateTrue || ps[i] == Not(p) {
			return true
		} else if m.value(ps[i]) != f.TristateFalse && ps[i] != p {
			p = ps[i]
			ps[j] = p
			j++
		}
	}
	if i-j > 0 {
		ps = ps[:len(ps)-(i-j)]
	}

	if flag {
		slice := make([]int32, len(ps)+1)
		slice[0] = 1
		for i := 0; i < len(ps); i++ {
			slice[i+1] = (Vari(ps[i]) + 1) * (-2*signAsInt(ps[i]) + 1)
		}
		m.pgProof = append(m.pgProof, slice)

		slice = make([]int32, len(oc)+1)
		slice[0] = -1
		for i := 0; i < len(oc); i++ {
			slice[i+1] = (Vari(oc[i]) + 1) * (-2*signAsInt(oc[i]) + 1)
		}
		m.pgProof = append(m.pgProof, slice)
	}

	if len(ps) == 0 {
		m.ok = false
		if m.config.ProofGeneration {
			m.pgProof = append(m.pgProof, []int32{0})
		}
		return false
	} else if len(ps) == 1 {
		m.enqueueFunction(m, ps[0], nil)
		m.ok = m.propagate() == nil
		m.unitClauses = append(m.unitClauses, ps[0])
		if !m.ok && m.config.ProofGeneration {
			m.pgProof = append(m.pgProof, []int32{0})
		}
		return m.ok
	} else {
		c := newClause(ps, -1)
		m.clauses = append(m.clauses, c)
		m.attachClause(c)
	}
	return true
}

func (m *CoreSolver) addAtMost(ps []int32, rhs int) {
	k := rhs
	if !m.ok {
		return
	}
	slices.Sort(ps)
	i, j := 0, 0
	p := LitUndef
	for ; i < len(ps); i++ {
		if m.value(ps[i]) == f.TristateTrue {
			k--
		} else if ps[i] == Not(p) {
			p = ps[i]
			j--
			k--
		} else if m.value(ps[i]) != f.TristateFalse {
			p = ps[i]
			ps[j] = p
			j++
		}
	}
	if i-j > 0 {
		ps = ps[:len(ps)-(i-j)]
	}
	if k >= len(ps) {
		return
	}
	if k < 0 {
		m.ok = false
		return
	}
	if k == 0 {
		for i = 0; i < len(ps); i++ {
			m.enqueueFunction(m, Not(ps[i]), nil)
			m.unitClauses = append(m.unitClauses, Not(ps[i]))
		}
		m.ok = m.propagate() == nil
		return
	}
	cr := newAtMostClause(ps, -1)
	cr.atMostWatchers = len(ps) - k + 1
	m.clauses = append(m.clauses, cr)
	m.attachClause(cr)
}

// Solve solves the formula on the solver with the given handler.  Returns the
// result as tristate and an ok flag which is false when the computation was
// aborted by the handler.
func (m *CoreSolver) Solve(satHandler Handler) (res f.Tristate, ok bool) {
	m.satHandler = satHandler
	handler.Start(m.satHandler)
	m.model = []bool{}
	m.conflict = []int32{}
	if !m.ok {
		return f.TristateFalse, true
	}
	status := f.TristateUndef
	for status == f.TristateUndef && !m.canceledByHandler {
		status, _ = m.search()
	}

	if m.config.ProofGeneration && len(m.assumptions) == 0 {
		if status == f.TristateFalse {
			m.pgProof = append(m.pgProof, []int32{0})
		}
	}

	if status == f.TristateTrue {
		m.model = make([]bool, len(m.vars))
		for i, v := range m.vars {
			m.model[i] = v.assignment == f.TristateTrue
		}
	} else if status == f.TristateFalse && len(m.conflict) == 0 {
		m.ok = false
	}
	handlerFinishSolving(m.satHandler)
	m.cancelUntil(0)
	m.satHandler = nil
	if m.canceledByHandler {
		m.canceledByHandler = false
		return f.TristateFalse, false
	} else {
		return status, true
	}
}

// SolveWithAssumptions is used to the the formulas on the solver with the
// given assumptions.  Returns the result as tristate and an ok flag which is
// false when the computation was aborted by the handler.
func (m *CoreSolver) SolveWithAssumptions(handler Handler, assumptions []int32) (res f.Tristate, ok bool) {
	m.assumptions = assumptions
	res, ok = m.Solve(handler)
	m.assumptions = []int32{}
	return
}

func (m *CoreSolver) search() (f.Tristate, bool) {
	if !m.ok {
		return f.TristateFalse, true
	}
	for {
		confl := m.propagate()
		if confl != nil {
			if m.satHandler != nil && !m.satHandler.DetectedConflict() {
				m.canceledByHandler = true
				return f.TristateUndef, false
			}
			m.conflicts++
			m.conflictsRestarts++
			if m.conflicts%5000 == 0 && m.varDecay < m.llConfig.MaxVarDecay {
				m.varDecay += 0.01
			}
			if m.decisionLevel() == 0 {
				return f.TristateFalse, true
			}
			m.trailQueue.push(len(m.trail))
			if m.conflictsRestarts > lbBlockingRestart && m.lbdQueue.valid() &&
				len(m.trail) > int(m.llConfig.FactorR*float64(m.trailQueue.avg())) {
				m.lbdQueue.fastClear()
			}
			var learntClause []int32
			m.analyze(confl, &learntClause)
			m.lbdQueue.push(m.analyzeLBD)
			m.sumLBD += float64(m.analyzeLBD)
			m.cancelUntil(m.analyzeBtLevel)

			if m.config.ProofGeneration {
				slice := make([]int32, len(learntClause)+1)
				slice[0] = 1
				for i := 0; i < len(learntClause); i++ {
					slice[i+1] = (Vari(learntClause[i]) + 1) * (-2*signAsInt(learntClause[i]) + 1)
				}
				m.pgProof = append(m.pgProof, slice)
			}

			if len(learntClause) == 1 {
				m.enqueueFunction(m, learntClause[0], nil)
				m.unitClauses = append(m.unitClauses, learntClause[0])
			} else {
				cr := newClause(learntClause, m.stateId)
				cr.lbd = m.analyzeLBD
				cr.oneWatched = false
				m.learnts = append(m.learnts, cr)
				m.attachClause(cr)
				m.claBumpActivity(cr)
				m.enqueueFunction(m, learntClause[0], cr)
			}
			m.varDecayActivity()
			m.claDecayActivity()
		} else {
			if m.lbdQueue.valid() &&
				(float64(m.lbdQueue.avg())*m.llConfig.FactorK) > (m.sumLBD/float64(m.conflictsRestarts)) {
				m.lbdQueue.fastClear()
				m.cancelUntil(0)
				return f.TristateUndef, true
			}
			if m.conflicts >= (m.curRestart*m.nbClausesBeforeReduce) && len(m.learnts) > 0 {
				m.curRestart = (m.conflicts / m.nbClausesBeforeReduce) + 1
				m.reduceDB()
				m.nbClausesBeforeReduce += m.llConfig.IncReduceDB
			}
			next := LitUndef
			for m.decisionLevel() < len(m.assumptions) {
				p := m.assumptions[m.decisionLevel()]
				if m.value(p) == f.TristateTrue {
					m.trailLim = append(m.trailLim, len(m.trail))
				} else if m.value(p) == f.TristateFalse {
					m.analyzeFinal(Not(p))
					return f.TristateFalse, true
				} else {
					next = p
					break
				}
			}
			if next == LitUndef {
				next = m.pickBranchLit()
				if next == LitUndef {
					return f.TristateTrue, true
				}
			}
			m.trailLim = append(m.trailLim, len(m.trail))
			m.enqueueFunction(m, next, nil)
		}
	}
}

func (m *CoreSolver) reset() {
	initialize(m, m.config)
}

func (m *CoreSolver) saveState() *SolverState {
	state := make([]int, 6)
	if m.ok {
		state[0] = 1
	} else {
		state[0] = 0
	}
	state[1] = len(m.vars)
	state[2] = len(m.clauses)
	state[3] = len(m.unitClauses)
	if m.config.ProofGeneration {
		state[4] = len(m.pgOriginalClauses)
		state[5] = len(m.pgProof)
	}
	id := m.stateId
	m.stateId++
	m.validStates = append(m.validStates, id)
	return &SolverState{
		id:    id,
		state: state,
	}
}

func (m *CoreSolver) loadState(solverState *SolverState) error {
	index := -1
	for i := len(m.validStates) - 1; i >= 0 && index == -1; i-- {
		if m.validStates[i] == solverState.id {
			index = i
		}
	}
	if index == -1 {
		return errorx.BadInput("solver state %d is not valid any more", solverState.id)
	}
	state := solverState.state
	shrinkTo(&m.validStates, index+1)
	m.completeBacktrack()
	m.ok = state[0] == 1
	newVarsSize := min(state[1], len(m.vars))
	for i := len(m.vars) - 1; i >= newVarsSize; i-- {
		idx := int32(i)
		varName := m.idx2name[idx]
		delete(m.idx2name, idx)
		delete(m.name2idx, varName)
		m.orderHeap.remove(idx)
	}
	shrinkTo(&m.vars, newVarsSize)
	shrinkTo(&m.permDiff, newVarsSize)
	newClausesSize := min(state[2], len(m.clauses))
	for i := len(m.clauses) - 1; i >= newClausesSize; i-- {
		m.simpleRemoveClause(m.clauses[i])
	}
	shrinkTo(&m.clauses, newClausesSize)

	newLearntsSize := 0
	for i := 0; i < len(m.learnts); i++ {
		learnt := m.learnts[i]
		if learnt.learntOnState <= solverState.id {
			m.learnts[newLearntsSize] = learnt
			newLearntsSize++
		} else {
			m.simpleRemoveClause(learnt)
		}
	}
	shrinkTo(&m.learnts, newLearntsSize)

	shrinkTo(&m.watches, newVarsSize*2)
	shrinkTo(&m.watchesBin, newVarsSize*2)
	shrinkTo(&m.unitClauses, state[3])
	for i := 0; m.ok && i < len(m.unitClauses); i++ {
		m.enqueueFunction(m, m.unitClauses[i], nil)
		m.ok = m.propagate() == nil
	}
	if m.config.ProofGeneration {
		newPgOriginalSize := min(state[4], len(m.pgOriginalClauses))
		shrinkTo(&m.pgOriginalClauses, newPgOriginalSize)
		newPgProofSize := min(state[5], len(m.pgProof))
		shrinkTo(&m.pgProof, newPgProofSize)
	}
	return nil
}

func (m *CoreSolver) completeBacktrack() {
	for v := 0; v < len(m.vars); v++ {
		vari := m.vars[v]
		vari.assignment = f.TristateUndef
		vari.reason = nil
		if !m.orderHeap.inHeap(int32(v)) && vari.decision {
			m.orderHeap.insert(int32(v))
		}
	}
	m.trail = []int32{}
	m.trailLim = []int{}
	m.qhead = 0
}

func (m *CoreSolver) simpleRemoveClause(c *clause) {
	if c.isAtMost {
		for i := 0; i < c.atMostWatchers; i++ {
			removeWatcher(&m.watches[c.get(i)], c)
		}
	} else if c.size() == 2 {
		removeWatcher(&m.watchesBin[Not(c.get(0))], c)
		removeWatcher(&m.watchesBin[Not(c.get(1))], c)
	} else {
		removeWatcher(&m.watches[Not(c.get(0))], c)
		removeWatcher(&m.watches[Not(c.get(1))], c)
	}
}

// NVars returns the number of vars on the solver.
func (m *CoreSolver) NVars() int32 {
	return int32(len(m.vars))
}

func (m *CoreSolver) nAssigns() int {
	return len(m.trail)
}

func (m *CoreSolver) decisionLevel() int {
	return len(m.trailLim)
}

func (m *CoreSolver) abstractLevel(x int32) int {
	return 1 << (m.vars[x].level & 31)
}

func (m *CoreSolver) insertVarOrder(x int32) {
	if !m.orderHeap.inHeap(x) && m.vars[x].decision {
		m.orderHeap.insert(x)
	}
}

func (m *CoreSolver) pickBranchLit() int32 {
	next := int32(-1)
	for next == -1 || m.vars[next].assignment != f.TristateUndef || !m.vars[next].decision {
		if m.orderHeap.isEmpty() {
			return -1
		} else {
			next = m.orderHeap.removeMin()
		}
	}
	return MkLit(next, m.vars[next].polarity)
}

func (m *CoreSolver) varDecayActivity() {
	m.varInc *= 1 / m.varDecay
}

func (m *CoreSolver) varBumpActivity(v int32) {
	m.varBumpActivityWithInc(v, m.varInc)
}

func (m *CoreSolver) varBumpActivityWithInc(v int32, inc float64) {
	variable := m.vars[v]
	variable.incrementActivity(inc)
	if variable.activity > 1e100 {
		for _, vari := range m.vars {
			vari.rescaleActivity()
		}
		m.varInc *= 1e-100
	}
	if m.orderHeap.inHeap(v) {
		m.orderHeap.decrease(v)
	}
}

func (m *CoreSolver) rebuildOrderHeap() {
	var vs []int32
	for v := int32(0); v < m.NVars(); v++ {
		if m.vars[v].decision && m.vars[v].assignment == f.TristateUndef {
			vs = append(vs, v)
		}
	}
	m.orderHeap.build(vs)
}

func (m *CoreSolver) locked(c *clause) bool {
	return m.value(c.get(0)) == f.TristateTrue && m.v(c.get(0)).reason != nil && m.v(c.get(0)).reason == c
}

func (m *CoreSolver) claDecayActivity() {
	m.claInc *= 1 / m.llConfig.ClauseDecay
}

func (m *CoreSolver) claBumpActivity(c *clause) {
	c.incrementActivity(m.claInc)
	if c.activity > 1e20 {
		for _, clause := range m.learnts {
			clause.rescaleActivity()
		}
		m.claInc *= 1e-20
	}
}

func (m *CoreSolver) attachClause(c *clause) {
	if c.isAtMost {
		for i := 0; i < c.atMostWatchers; i++ {
			l := c.get(i)
			m.watches[l] = append(m.watches[l], newWatcher(c, LitUndef))
		}
		m.clausesLiterals += c.size()
	} else {
		if c.size() == 2 {
			m.watchesBin[Not(c.get(0))] = append(m.watchesBin[Not(c.get(0))], newWatcher(c, c.get(1)))
			m.watchesBin[Not(c.get(1))] = append(m.watchesBin[Not(c.get(1))], newWatcher(c, c.get(0)))
		} else {
			m.watches[Not(c.get(0))] = append(m.watches[Not(c.get(0))], newWatcher(c, c.get(1)))
			m.watches[Not(c.get(1))] = append(m.watches[Not(c.get(1))], newWatcher(c, c.get(0)))
		}
		if c.learnt() {
			m.learntsLiterals += c.size()
		} else {
			m.clausesLiterals += c.size()
		}
	}
}

func (m *CoreSolver) detachClause(c *clause) {
	m.simpleRemoveClause(c)
	if c.learnt() {
		m.learntsLiterals -= c.size()
	} else {
		m.clausesLiterals -= c.size()
	}
}

func (m *CoreSolver) detachAtMost(c *clause) {
	for i := 0; i < c.atMostWatchers; i++ {
		removeWatcher(&m.watches[c.get(i)], c)
	}
	m.clausesLiterals -= c.size()
}

func removeWatcher(watcher *[]*watcher, clause *clause) {
	for i, w := range *watcher {
		if w.clause == clause {
			*watcher = append((*watcher)[:i], (*watcher)[i+1:]...)
		}
	}
}

func (m *CoreSolver) removeClause(c *clause) {
	if c.isAtMost {
		m.detachAtMost(c)
		for i := 0; i < c.atMostWatchers; i++ {
			if m.value(c.get(i)) == f.TristateFalse && m.v(c.get(i)).reason != nil && m.v(c.get(i)).reason == c {
				m.v(c.get(i)).reason = nil
			}
		}
	} else {
		if m.config.ProofGeneration {
			slice := make([]int32, c.size()+1)
			slice[0] = -1
			for i := 0; i < c.size(); i++ {
				slice[i+1] = (Vari(c.get(i)) + 1) * (-2*signAsInt(c.get(i)) + 1)
			}
			m.pgProof = append(m.pgProof, slice)
		}
		m.detachClause(c)
		if m.locked(c) {
			m.v(c.get(0)).reason = nil
		}
	}
}

func (m *CoreSolver) propagate() *clause {
	var confl *clause
	for m.qhead < len(m.trail) {
		p := m.trail[m.qhead]
		m.qhead++
		ws := &m.watches[p]
		var iInd, jInd int
		wbin := m.watchesBin[p]
		for k := 0; k < len(wbin); k++ {
			imp := wbin[k].blocker
			if m.value(imp) == f.TristateFalse {
				return wbin[k].clause
			}
			if m.value(imp) == f.TristateUndef {
				m.enqueueFunction(m, imp, wbin[k].clause)
			}
		}
		for iInd < len(*ws) {
			i := (*ws)[iInd]
			blocker := i.blocker
			if blocker != LitUndef && m.value(blocker) == f.TristateTrue {
				(*ws)[jInd] = i
				jInd++
				iInd++
				continue
			}
			c := i.clause

			if c.isAtMost {
				switch newWatch := m.findNewWatchForAtMostClause(c, p); newWatch {
				case LitUndef:
					for k := 0; k < c.atMostWatchers; k++ {
						if c.get(k) != p && m.value(c.get(k)) != f.TristateFalse {
							m.enqueueFunction(m, Not(c.get(k)), c)
						}
					}
					(*ws)[jInd] = (*ws)[iInd]
					jInd++
					iInd++
				case LitError:
					confl = c
					m.qhead = len(m.trail)
					for iInd < len(*ws) {
						(*ws)[jInd] = (*ws)[iInd]
						jInd++
						iInd++
					}
				case p:
					(*ws)[jInd] = (*ws)[iInd]
					jInd++
					iInd++
				default:
					iInd++
					w := newWatcher(c, LitUndef)
					m.watches[newWatch] = append(m.watches[newWatch], w)
				}
			} else {
				falseLit := Not(p)
				if c.get(0) == falseLit {
					c.set(0, c.get(1))
					c.set(1, falseLit)
				}
				iInd++
				first := c.get(0)
				w := newWatcher(c, first)
				if first != blocker && m.value(first) == f.TristateTrue {
					(*ws)[jInd] = w
					jInd++
					continue
				}
				foundWatch := false
				for k := 2; k < c.size() && !foundWatch; k++ {
					if m.value(c.get(k)) != f.TristateFalse {
						c.set(1, c.get(k))
						c.set(k, falseLit)
						m.watches[Not(c.get(1))] = append(m.watches[Not(c.get(1))], w)
						foundWatch = true
					}
				}
				if !foundWatch {
					(*ws)[jInd] = w
					jInd++
					if m.value(first) == f.TristateFalse {
						confl = c
						m.qhead = len(m.trail)
						for iInd < len(*ws) {
							(*ws)[jInd] = (*ws)[iInd]
							jInd++
							iInd++
						}
					} else {
						m.enqueueFunction(m, first, c)
					}
				}
			}
		}
		if del := iInd - jInd; del > 0 {
			*ws = (*ws)[:len(*ws)-del]
		}
	}
	return confl
}

func (m *CoreSolver) findNewWatchForAtMostClause(c *clause, p int32) int32 {
	numFalse, numTrue := 0, 0
	maxTrue := c.size() - c.atMostWatchers + 1
	for q := 0; q < c.atMostWatchers; q++ {
		val := m.value(c.get(q))
		if val == f.TristateUndef {
			continue
		} else if val == f.TristateFalse {
			numFalse++
			if numFalse >= c.atMostWatchers-1 {
				return p
			}
			continue
		}
		numTrue++
		if numTrue > maxTrue {
			return LitError
		}
		if c.get(q) == p {
			for next := c.atMostWatchers; next < c.size(); next++ {
				if m.value(c.get(next)) != f.TristateTrue {
					newWatch := c.get(next)
					c.set(next, c.get(q))
					c.set(q, newWatch)
					return newWatch
				}
			}
		}
	}
	if numTrue > 1 {
		return LitError
	} else {
		return LitUndef
	}
}

func (m *CoreSolver) analyze(conflictClause *clause, outLearnt *[]int32) {
	c := conflictClause
	pathC := 0
	p := LitUndef
	*outLearnt = append(*outLearnt, -1)
	index := len(m.trail) - 1
	for ok := true; ok; ok = pathC > 0 {
		if c.isAtMost {
			for j := 0; j < c.size(); j++ {
				if m.value(c.get(j)) != f.TristateTrue {
					continue
				}
				q := Not(c.get(j))
				if !m.seen[Vari(q)] && m.v(q).level > 0 {
					m.varBumpActivity(Vari(q))
					m.seen[Vari(q)] = true
					if m.v(q).level >= m.decisionLevel() {
						pathC++
					} else {
						*outLearnt = append(*outLearnt, q)
					}
				}
			}
		} else {
			if p != LitUndef && c.size() == 2 && m.value(c.get(0)) == f.TristateFalse {
				tmp := c.get(0)
				c.set(0, c.get(1))
				c.set(1, tmp)
			}
			if c.learnt() {
				m.claBumpActivity(c)
			} else {
				if !c.seen {
					c.seen = true
				}
			}
			if c.learnt() && c.lbd > 2 {
				nblevels := m.computeLBDUnit(c)
				if nblevels+1 < c.lbd {
					if c.lbd <= m.llConfig.LBLBDFrozenClause {
						c.canBeDel = false
					}
					c.lbd = nblevels
				}
			}
			var j int
			if p == LitUndef {
				j = 0
			} else {
				j = 1
			}
			for ; j < c.size(); j++ {
				q := c.get(j)
				if !m.seen[Vari(q)] && m.v(q).level != 0 {
					m.varBumpActivity(Vari(q))
					m.seen[Vari(q)] = true
					if m.v(q).level >= m.decisionLevel() {
						pathC++
						if (m.v(q).reason != nil) && m.v(q).reason.learnt() {
							m.lastDecisionLevel = append(m.lastDecisionLevel, q)
						}
					} else {
						*outLearnt = append(*outLearnt, q)
					}
				}
			}
		}

		for !m.seen[Vari(m.trail[index])] {
			index--
		}

		p = m.trail[index]
		c = m.v(p).reason
		m.seen[Vari(p)] = false
		pathC--
	}
	(*outLearnt)[0] = Not(p)
	m.simplifyClause(outLearnt)
}

func (m *CoreSolver) computeLBD(lits *[]int32) int {
	nbLevels := 0
	m.myflag++
	for i := 0; i < len(*lits); i++ {
		l := m.v((*lits)[i]).level
		if m.permDiff[l] != m.myflag {
			m.permDiff[l] = m.myflag
			nbLevels++
		}
	}
	if !m.llConfig.ReduceOnSize {
		return nbLevels
	}
	if len(*lits) < m.llConfig.ReduceOnSizeSize {
		return len(*lits)
	}
	return len(*lits) + nbLevels
}

func (m *CoreSolver) computeLBDUnit(c *clause) int {
	return m.computeLBD(&c.data)
}

func (m *CoreSolver) simplifyClause(outLearnt *[]int32) {
	var i, j int
	analyzeToClear := make([]int32, len(*outLearnt))
	copy(analyzeToClear, *outLearnt)
	if m.config.ClauseMinimization == ClauseMinDeep {
		abstractLevel := 0
		for i = 1; i < len(*outLearnt); i++ {
			abstractLevel |= m.abstractLevel(Vari((*outLearnt)[i]))
		}
		i, j = 1, 1
		for ; i < len(*outLearnt); i++ {
			if m.v((*outLearnt)[i]).reason == nil || !m.litRedundant((*outLearnt)[i], abstractLevel, &analyzeToClear) {
				(*outLearnt)[j] = (*outLearnt)[i]
				j++
			}
		}
	} else if m.config.ClauseMinimization == ClauseMinBasic {
		i, j = 1, 1
		for ; i < len(*outLearnt); i++ {
			c := m.v((*outLearnt)[i]).reason
			if c == nil {
				(*outLearnt)[j] = (*outLearnt)[i]
				j++
			} else {
				var k int
				if c.size() == 2 {
					k = 0
				} else {
					k = 1
				}
				for ; k < c.size(); k++ {
					if !m.seen[Vari(c.get(k))] && m.v(c.get(k)).level > 0 {
						(*outLearnt)[j] = (*outLearnt)[i]
						j++
						break
					}
				}
			}
		}
	} else {
		i = len(*outLearnt)
		j = i
	}
	if i-j > 0 {
		*outLearnt = (*outLearnt)[:len(*outLearnt)-(i-j)]
	}
	if len(*outLearnt) <= m.llConfig.LBSizeMinimizingClause {
		m.minimisationWithBinaryResolution(outLearnt)
	}
	m.analyzeBtLevel = 0
	if len(*outLearnt) > 1 {
		maxLit := 1
		for k := 2; k < len(*outLearnt); k++ {
			if m.v((*outLearnt)[k]).level > m.v((*outLearnt)[maxLit]).level {
				maxLit = k
			}
		}
		p := (*outLearnt)[maxLit]
		(*outLearnt)[maxLit] = (*outLearnt)[1]
		(*outLearnt)[1] = p
		m.analyzeBtLevel = m.v(p).level
	}
	m.analyzeLBD = m.computeLBD(outLearnt)
	if len(m.lastDecisionLevel) > 0 {
		for k := 0; k < len(m.lastDecisionLevel); k++ {
			if (m.v(m.lastDecisionLevel[k]).reason).lbd < m.analyzeLBD {
				m.varBumpActivity(Vari(m.lastDecisionLevel[k]))
			}
		}
		m.lastDecisionLevel = []int32{}
	}
	for l := 0; l < len(analyzeToClear); l++ {
		m.seen[Vari(analyzeToClear[l])] = false
	}
}

func (m *CoreSolver) minimisationWithBinaryResolution(outLearnt *[]int32) {
	lbd := m.computeLBD(outLearnt)
	p := Not((*outLearnt)[0])
	if lbd <= m.llConfig.LBLBDMinimizingClause {
		m.myflag++
		for i := 1; i < len(*outLearnt); i++ {
			m.permDiff[Vari((*outLearnt)[i])] = m.myflag
		}
		nb := 0
		for _, wbin := range m.watchesBin[p] {
			imp := wbin.blocker
			if m.permDiff[Vari(imp)] == m.myflag && m.value(imp) == f.TristateTrue {
				nb++
				m.permDiff[Vari(imp)] = m.myflag - 1
			}
		}
		l := len(*outLearnt) - 1
		if nb > 0 {
			for i := 1; i < len(*outLearnt)-nb; i++ {
				if m.permDiff[Vari((*outLearnt)[i])] != m.myflag {
					p = (*outLearnt)[l]
					(*outLearnt)[l] = (*outLearnt)[i]
					(*outLearnt)[i] = p
					l--
					i--
				}
			}
			*outLearnt = (*outLearnt)[:len(*outLearnt)-nb]
		}
	}
}

func (m *CoreSolver) litRedundant(p int32, abstractLevel int, analyzeToClear *[]int32) bool {
	analyzeStack := []int32{p}
	top := len(*analyzeToClear)
	for len(analyzeStack) > 0 {
		c := m.v(analyzeStack[len(analyzeStack)-1]).reason
		pop(&analyzeStack)
		if c.isAtMost {
			for i := 0; i < c.size(); i++ {
				if m.value(c.get(i)) != f.TristateTrue {
					continue
				}
				q := Not(c.get(i))
				if !m.seen[Vari(q)] && m.v(q).level > 0 {
					if m.v(q).reason != nil && (m.abstractLevel(Vari(q))&abstractLevel) != 0 {
						m.seen[Vari(q)] = true
						analyzeStack = append(analyzeStack, q)
						*analyzeToClear = append(*analyzeToClear, q)
					} else {
						for j := top; j < len(*analyzeToClear); j++ {
							m.seen[Vari((*analyzeToClear)[j])] = false
						}
						if del := len(*analyzeToClear) - top; del > 0 {
							*analyzeToClear = (*analyzeToClear)[:len(*analyzeToClear)-del]
						}
						return false
					}
				}
			}
		} else {
			if c.size() == 2 && m.value(c.get(0)) == f.TristateFalse {
				tmp := c.get(0)
				c.set(0, c.get(1))
				c.set(1, tmp)
			}
			for i := 1; i < c.size(); i++ {
				q := c.get(i)
				if !m.seen[Vari(q)] && m.v(q).level > 0 {
					if m.v(q).reason != nil && (m.abstractLevel(Vari(q))&abstractLevel) != 0 {
						m.seen[Vari(q)] = true
						analyzeStack = append(analyzeStack, q)
						*analyzeToClear = append(*analyzeToClear, q)
					} else {
						for j := top; j < len(*analyzeToClear); j++ {
							m.seen[Vari((*analyzeToClear)[j])] = false
						}
						if del := len(*analyzeToClear) - top; del > 0 {
							*analyzeToClear = (*analyzeToClear)[:len(*analyzeToClear)-del]
						}
						return false
					}
				}
			}
		}
	}
	return true
}

func (m *CoreSolver) analyzeFinal(p int32) {
	m.conflict = []int32{p}
	if m.decisionLevel() == 0 {
		return
	}
	m.seen[Vari(p)] = true
	for i := len(m.trail) - 1; i >= m.trailLim[0]; i-- {
		x := Vari(m.trail[i])
		if m.seen[x] {
			v := m.vars[x]
			if v.reason == nil {
				m.conflict = append(m.conflict, Not(m.trail[i]))
			} else {
				c := v.reason
				if !c.isAtMost {
					var j int
					if c.size() == 2 {
						j = 0
					} else {
						j = 1
					}
					for ; j < c.size(); j++ {
						if m.v(c.get(j)).level > 0 {
							m.seen[Vari(c.get(j))] = true
						}
					}
				} else {
					for j := 0; j < c.size(); j++ {
						if m.value(c.get(j)) == f.TristateTrue && m.v(c.get(j)).level > 0 {
							m.seen[Vari(c.get(j))] = true
						}
					}
				}
			}
			m.seen[x] = false
		}
	}
	m.seen[Vari(p)] = false
}

func (m *CoreSolver) cancelUntil(level int) {
	if m.decisionLevel() > level {
		if !m.computingBackbone {
			for c := len(m.trail) - 1; c >= m.trailLim[level]; c-- {
				x := Vari(m.trail[c])
				v := m.vars[x]
				v.assignment = f.TristateUndef
				v.polarity = Sign(m.trail[c])
				m.insertVarOrder(x)
			}
		} else {
			for c := len(m.trail) - 1; c >= m.trailLim[level]; c-- {
				x := Vari(m.trail[c])
				v := m.vars[x]
				v.assignment = f.TristateUndef
				v.polarity = !m.computingBackbone && Sign(m.trail[c])
				m.insertVarOrder(x)
			}
		}
		m.qhead = m.trailLim[level]
		if del := len(m.trail) - m.trailLim[level]; del > 0 {
			m.trail = m.trail[:len(m.trail)-del]
		}
		if del := len(m.trailLim) - level; del > 0 {
			m.trailLim = m.trailLim[:len(m.trailLim)-del]
		}
	}
}

func (m *CoreSolver) reduceDB() {
	sortClauses(&m.learnts)
	if m.learnts[len(m.learnts)/ratioRemoveClauses].lbd <= 3 {
		m.nbClausesBeforeReduce += m.llConfig.SpecialIncReduceDB
	}
	if m.learnts[len(m.learnts)-1].lbd <= 5 {
		m.nbClausesBeforeReduce += m.llConfig.SpecialIncReduceDB
	}
	limit := len(m.learnts) / 2
	var i, j int
	for ; i < len(m.learnts); i++ {
		c := m.learnts[i]
		if c.lbd > 2 && c.size() > 2 && c.canBeDel && !m.locked(c) && (i < limit) {
			m.removeClause(m.learnts[i])
		} else {
			if !c.canBeDel {
				limit++
			}
			c.canBeDel = true
			m.learnts[j] = m.learnts[i]
			j++
		}
	}

	if del := i - j; del > 0 {
		m.learnts = m.learnts[:len(m.learnts)-del]
	}
}

func (m *CoreSolver) IdxForName(name string) int32 {
	id, ok := m.name2idx[name]
	if !ok {
		return -1
	}
	return id
}

func (m *CoreSolver) upZeroLiterals() []int32 {
	var upZeroLiterals []int32
	for i := 0; i < len(m.trail); i++ {
		lit := m.trail[i]
		if m.v(lit).level > 0 {
			break
		} else {
			upZeroLiterals = append(upZeroLiterals, lit)
		}
	}
	return upZeroLiterals
}

// Model returns the current model of the solver.
func (m *CoreSolver) Model() []bool {
	model := make([]bool, len(m.model))
	copy(model, m.model)
	return model
}

// Conflict returns the current conflict of the solver.
func (m *CoreSolver) Conflict() []int32 {
	conflict := make([]int32, len(m.conflict))
	copy(conflict, m.conflict)
	return conflict
}

// CreateModel is used to create a model data-structure from the given model.
func (m *CoreSolver) CreateModel(fac f.Factory, mVec []bool, relevantIndices []int32) *model.Model {
	mdl := model.New()
	for i := 0; i < len(relevantIndices); i++ {
		index := relevantIndices[i]
		if index != -1 {
			name := m.idx2name[index]
			mdl.AddLiteral(fac.Lit(name, mVec[index]))
		}
	}
	return mdl
}

// KnownVariables returns the variables currently known by the solver.
func (m *CoreSolver) KnownVariables(fac f.Factory) *f.VarSet {
	result := f.NewVarSet()
	for name := range m.name2idx {
		result.Add(fac.Var(name))
	}
	return result
}

func shrinkTo[T any](slice *[]T, newSize int) {
	if newSize < len(*slice) {
		*slice = (*slice)[:newSize]
	}
}

func pop[T any](slice *[]T) {
	*slice = (*slice)[:len(*slice)-1]
}
