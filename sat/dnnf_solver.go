package sat

import (
	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/normalform"
)

// DnnfSatSolver is a special solver used for DNNF generation.
type DnnfSatSolver interface {
	Factory() f.Factory
	Add(formula f.Formula)
	Start() bool
	Decide(variable int32, phase bool) bool
	UndoDecide(variable int32)
	AtAssertionLevel() bool
	AssertCdLiteral() bool
	NewlyImplied(knownVariables []bool) f.Formula
	VariableIndex(lit f.Literal) int32
	LitForIdx(variable int32) f.Formula
	ValueOf(lit int32) f.Tristate
}

type DnnfSolver struct {
	*CoreSolver
	fac               f.Factory
	newlyImpliedDirty bool
	assertionLevel    int
	lastLearnt        []int32
	assignment        []f.Tristate
	impliedOperands   []f.Formula
}

func NewDnnfSolver(fac f.Factory, numberOfVariables int) *DnnfSolver {
	solver := DnnfSolver{
		fac:             fac,
		assertionLevel:  -1,
		assignment:      make([]f.Tristate, 2*numberOfVariables),
		impliedOperands: []f.Formula{},
	}
	solver.CoreSolver = NewCoreSolver(DefaultConfig(), solver.dnnfUncheckedEnqueue)
	for i := 0; i < len(solver.assignment); i++ {
		solver.assignment[i] = f.TristateUndef
	}
	return &solver
}

func (m *DnnfSolver) Start() bool {
	m.newlyImpliedDirty = true
	return m.propagate() == nil
}

func (m *DnnfSolver) ValueOf(lit int32) f.Tristate {
	return m.assignment[lit]
}

func (m *DnnfSolver) VariableIndex(lit f.Literal) int32 {
	name, _, _ := m.fac.LitNamePhase(lit)
	return m.IdxForName(name)
}

func (m *DnnfSolver) LitForIdx(variable int32) f.Formula {
	return m.fac.Variable(m.idx2name[variable])
}

func (m *DnnfSolver) Factory() f.Factory {
	return m.fac
}

func (m *DnnfSolver) Add(formula f.Formula) {
	cnf := normalform.CNF(m.fac, formula)
	switch cnf.Sort() {
	case f.SortTrue:
		break
	case f.SortFalse, f.SortLiteral, f.SortOr:
		m.AddClause(m.generateClauseVector(cnf), nil)
	case f.SortAnd:
		for _, op := range m.fac.Operands(cnf) {
			m.AddClause(m.generateClauseVector(op), nil)
		}
	}
}

func (m *DnnfSolver) Decide(variable int32, phase bool) bool {
	m.newlyImpliedDirty = true
	lit := MkLit(variable, !phase)
	m.trailLim = append(m.trailLim, len(m.trail))
	m.enqueueFunction(m.CoreSolver, lit, nil)
	return m.propagateAfterDecide()
}

func (m *DnnfSolver) UndoDecide(variable int32) {
	m.newlyImpliedDirty = false
	m.cancelUntilDnnf(m.vars[variable].level - 1)
}

func (m *DnnfSolver) AtAssertionLevel() bool {
	return m.decisionLevel() == m.assertionLevel
}

func (m *DnnfSolver) AssertCdLiteral() bool {
	m.newlyImpliedDirty = true
	if !m.AtAssertionLevel() {
		panic(errorx.IllegalState("assertCdLiteral called although not at assertion level"))
	}
	if len(m.lastLearnt) == 1 {
		m.enqueueFunction(m.CoreSolver, m.lastLearnt[0], nil)
		m.unitClauses = append(m.unitClauses, m.lastLearnt[0])
	} else {
		cr := newClause(m.lastLearnt, m.stateId)
		m.learnts = append(m.learnts, cr)
		m.attachClause(cr)
		m.claBumpActivity(cr)
		m.enqueueFunction(m.CoreSolver, m.lastLearnt[0], cr)
	}
	return m.propagateAfterDecide()
}

func (m *DnnfSolver) NewlyImplied(knownVariables []bool) f.Formula {
	m.impliedOperands = []f.Formula{}
	if m.newlyImpliedDirty {
		var limit int
		if len(m.trailLim) == 0 {
			limit = -1
		} else {
			limit = m.trailLim[len(m.trailLim)-1]
		}
		for i := len(m.trail) - 1; i > limit; i-- {
			lit := m.trail[i]
			index := Vari(lit)
			if int(index) < len(knownVariables) && knownVariables[index] {
				m.impliedOperands = append(m.impliedOperands, m.intToLiteral(lit))
			}
		}
	}
	m.newlyImpliedDirty = false
	return m.fac.And(m.impliedOperands...)
}

func (m *DnnfSolver) dnnfUncheckedEnqueue(solver *CoreSolver, lit int32, reason *clause) {
	m.assignment[lit] = f.TristateTrue
	m.assignment[lit^1] = f.TristateFalse
	UncheckedEnqueue(solver, lit, reason)
}

func (m *DnnfSolver) cancelUntilDnnf(level int) {
	if m.decisionLevel() > level {
		for c := len(m.trail) - 1; c >= m.trailLim[level]; c-- {
			l := m.trail[c]
			m.assignment[l] = f.TristateUndef
			m.assignment[l^1] = f.TristateUndef
			x := Vari(l)
			v := m.vars[x]
			v.assignment = f.TristateUndef
			v.polarity = Sign(m.trail[c])
			m.insertVarOrder(x)
		}
		m.qhead = m.trailLim[level]
		removeElements(&m.trail, len(m.trail)-m.trailLim[level])
		removeElements(&m.trailLim, len(m.trailLim)-level)
	}
}

func (m *DnnfSolver) generateClauseVector(clause f.Formula) []int32 {
	lits := f.Literals(m.fac, clause)
	clauseVec := make([]int32, lits.Size())
	for i, lit := range lits.Content() {
		name, _, _ := m.fac.LitNamePhase(lit)
		index := m.IdxForName(name)
		if index == -1 {
			index = m.NewVar(false, true)
			m.addName(name, index)
		}
		var litNum int32
		if lit.IsPos() {
			litNum = index * 2
		} else {
			litNum = (index * 2) ^ 1
		}
		clauseVec[i] = litNum
	}
	return clauseVec
}

func (m *DnnfSolver) propagateAfterDecide() bool {
	conflict := m.propagate()
	if conflict != nil {
		m.handleConflict(conflict)
		return false
	}
	return true
}

func (m *DnnfSolver) handleConflict(conflict *clause) {
	if m.decisionLevel() > 0 {
		m.lastLearnt = []int32{}
		m.analyze(conflict, &m.lastLearnt)
		m.assertionLevel = m.analyzeBtLevel
	} else {
		m.cancelUntilDnnf(0)
		m.lastLearnt = nil
		m.assertionLevel = -1
	}
}

func (m *DnnfSolver) intToLiteral(lit int32) f.Formula {
	name := m.idx2name[Vari(lit)]
	return m.fac.Literal(name, !Sign(lit))
}

func removeElements[T any](slice *[]T, num int) {
	if num > 0 {
		*slice = (*slice)[:len(*slice)-num]
	}
}
