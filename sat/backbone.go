package sat

import (
	"github.com/booleworks/logicng-go/event"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
)

// BackboneSort describes sort of backbone should be computed:
type BackboneSort byte

const (
	BBPos  BackboneSort = iota // only variables which occur positive in every model
	BBNeg                      // only variables which occur negative in every model
	BBBoth                     // positive, negative and optional variables
)

// A Backbone of a formula is a set of literals (positive and/or negative)
// which are present in their respective polarity in every model of the given
// formula.  Therefore, the literals must be set accordingly in order for the
// formula to evaluate to true.
type Backbone struct {
	Sat      bool         // flag whether the formula was satisfiable
	Positive []f.Variable // variables that occur positive in each model of the formula
	Negative []f.Variable // variables that occur negative in each model of the formula
	Optional []f.Variable // variables that are neither in the positive nor in the negative backbone
}

// ToFormula returns the conjunction of positive and negative literals of the
// backbone.
func (b *Backbone) ToFormula(fac f.Factory) f.Formula {
	if !b.Sat {
		return fac.Falsum()
	} else {
		return fac.Minterm(b.CompleteBackbone(fac)...)
	}
}

// CompleteBackbone returns the positive and negative literals of the backbone.
func (b *Backbone) CompleteBackbone(fac f.Factory) []f.Literal {
	if !b.Sat {
		return []f.Literal{}
	} else {
		completeBackbone := make([]f.Literal, 0, len(b.Positive)+len(b.Negative))
		for i := 0; i < len(b.Positive); i++ {
			completeBackbone = append(completeBackbone, b.Positive[i].AsLiteral())
		}
		for i := 0; i < len(b.Negative); i++ {
			completeBackbone = append(completeBackbone, b.Negative[i].Negate(fac))
		}
		return completeBackbone
	}
}

// ComputeBackbone computes the positive and negative backbone on the solver.
func (s *Solver) ComputeBackbone(
	fac f.Factory, variables []f.Variable, backboneSort ...BackboneSort,
) *Backbone {
	backbone, _ := s.ComputeBackboneWithHandler(fac, variables, handler.NopHandler, backboneSort...)
	return backbone
}

// ComputeBackboneWithHandler computes the positive and negative backbone on
// the solver.  The given handler can be used to cancel the solver used for
// the backbone computation.
func (s *Solver) ComputeBackboneWithHandler(
	fac f.Factory, variables []f.Variable, hdl handler.Handler, backboneSort ...BackboneSort,
) (*Backbone, handler.State) {
	if e := event.BackboneComputationStarted; !hdl.ShouldResume(e) {
		return nil, handler.Cancelation(e)
	}
	m := s.core
	var bbSort BackboneSort
	if len(backboneSort) == 0 {
		bbSort = BBBoth
	} else {
		bbSort = backboneSort[0]
	}
	state := m.saveState()
	sat, hdlState := m.Solve(hdl)
	var backbone *Backbone
	if !hdlState.Success {
		// do nothing
	} else if sat == f.TristateTrue {
		m.computingBackbone = true
		relevantVarIndices := m.getRelevantVarIndices(fac, variables)
		m.initBackboneDS(relevantVarIndices)
		hdlState = m.computeBackboneInternal(relevantVarIndices, bbSort, hdl)
		if hdlState.Success {
			backbone = m.buildBackbone(fac, variables, bbSort)
		}
		m.computingBackbone = false
	} else {
		backbone = &Backbone{Sat: false}
	}
	err := m.loadState(state)
	if err != nil {
		panic(err)
	}
	return backbone, hdlState
}

func (m *CoreSolver) getRelevantVarIndices(fac f.Factory, variables []f.Variable) []int32 {
	relevantVarIndices := make([]int32, 0, len(variables))
	for _, v := range variables {
		name, _ := fac.VarName(v)
		idx, ok := m.name2idx[name]
		if ok {
			relevantVarIndices = append(relevantVarIndices, idx)
		}
	}
	return relevantVarIndices
}

func (m *CoreSolver) initBackboneDS(variables []int32) {
	m.backboneCandidates = make([]int32, 0, len(variables))
	m.backboneAssumptions = make([]int32, 0, len(variables))
	m.backboneMap = make(map[int32]f.Tristate)
	for _, vari := range variables {
		m.backboneMap[vari] = f.TristateUndef
	}
}

func (m *CoreSolver) computeBackboneInternal(
	variables []int32, bbSort BackboneSort, hdl handler.Handler,
) handler.State {
	m.createInitialCandidates(variables, bbSort)
	for len(m.backboneCandidates) > 0 {
		lit := m.backboneCandidates[len(m.backboneCandidates)-1]
		m.backboneCandidates = m.backboneCandidates[:len(m.backboneCandidates)-1]
		sat, state := m.solveWithLit(lit, hdl)
		if !state.Success {
			return state
		}
		if sat {
			m.refineUpperBound()
		} else {
			m.addBackboneLiteral(lit)
		}
	}
	return succ
}

func (m *CoreSolver) createInitialCandidates(variables []int32, bbSort BackboneSort) {
	for _, vari := range variables {
		if m.isUpZeroLit(vari) {
			backboneLit := MkLit(vari, !m.model[vari])
			m.addBackboneLiteral(backboneLit)
		} else {
			modelPhase := m.model[vari]
			if isBothOrNegative(bbSort) && !modelPhase || isBothOrPositive(bbSort) && modelPhase {
				lit := MkLit(vari, !modelPhase)
				if !m.isRotatable(lit) {
					m.backboneCandidates = append(m.backboneCandidates, lit)
				}
			}
		}
	}
}

func (m *CoreSolver) refineUpperBound() {
	candidates := make([]int32, len(m.backboneCandidates))
	copy(candidates, m.backboneCandidates)
	for _, lit := range candidates {
		vari := Vari(lit)
		if m.isUpZeroLit(vari) {
			removeFromSlice(&m.backboneCandidates, lit)
			m.addBackboneLiteral(lit)
		} else if m.model[vari] == Sign(lit) {
			removeFromSlice(&m.backboneCandidates, lit)
		} else if m.isRotatable(lit) {
			removeFromSlice(&m.backboneCandidates, lit)
		}
	}
}

func (m *CoreSolver) solveWithLit(lit int32, hdl handler.Handler) (bool, handler.State) {
	m.backboneAssumptions = append(m.backboneAssumptions, Not(lit))
	sat, state := m.SolveWithAssumptions(hdl, m.backboneAssumptions)
	m.backboneAssumptions = m.backboneAssumptions[:len(m.backboneAssumptions)-1]
	return sat == f.TristateTrue, state
}

func (m *CoreSolver) buildBackbone(fac f.Factory, variables []f.Variable, bbSort BackboneSort) *Backbone {
	var posBackboneVars, negBackboneVars, optionalVars []f.Variable
	if isBothOrPositive(bbSort) {
		posBackboneVars = make([]f.Variable, 0)
	}
	if isBothOrNegative(bbSort) {
		negBackboneVars = make([]f.Variable, 0)
	}
	if isBoth(bbSort) {
		optionalVars = make([]f.Variable, 0)
	}

	for _, v := range variables {
		name, _ := fac.VarName(v)
		idx, ok := m.name2idx[name]
		if !ok {
			if isBoth(bbSort) {
				optionalVars = append(optionalVars, v)
			}
		} else {
			switch m.backboneMap[idx] {
			case f.TristateTrue:
				if isBothOrPositive(bbSort) {
					posBackboneVars = append(posBackboneVars, v)
				}
			case f.TristateFalse:
				if isBothOrNegative(bbSort) {
					negBackboneVars = append(negBackboneVars, v)
				}
			case f.TristateUndef:
				if isBoth(bbSort) {
					optionalVars = append(optionalVars, v)
				}
			}
		}
	}
	return &Backbone{true, posBackboneVars, negBackboneVars, optionalVars}
}

func (m *CoreSolver) isUpZeroLit(variable int32) bool {
	return m.vars[variable].level == 0
}

func (m *CoreSolver) isRotatable(lit int32) bool {
	if m.v(lit).reason != nil {
		return false
	}
	for _, watcher := range m.watches[Not(lit)] {
		if m.isUnit(lit, watcher.clause) {
			return false
		}
	}
	for _, watcher := range m.watchesBin[Not(lit)] {
		if m.isUnit(lit, watcher.clause) {
			return false
		}
	}
	return true
}

func (m *CoreSolver) isUnit(lit int32, clause *clause) bool {
	if !clause.isAtMost {
		for i := 0; i < clause.size(); i++ {
			clauseLit := clause.get(i)
			if lit != clauseLit && m.model[Vari(clauseLit)] != Sign(clauseLit) {
				return false
			}
		}
		return true
	} else {
		countPos := 0
		cardinality := clause.cardinality()
		for i := 0; i < clause.size(); i++ {
			variable := Vari(clause.get(i))
			if Vari(lit) != variable && m.model[variable] {
				countPos++
				if countPos == cardinality {
					return true
				}
			}
		}
		return false
	}
}

func (m *CoreSolver) addBackboneLiteral(lit int32) {
	if Sign(lit) {
		m.backboneMap[Vari(lit)] = f.TristateFalse
	} else {
		m.backboneMap[Vari(lit)] = f.TristateTrue
	}
	m.backboneAssumptions = append(m.backboneAssumptions, lit)
}

func isBothOrPositive(bbSort BackboneSort) bool {
	return bbSort == BBBoth || bbSort == BBPos
}

func isBothOrNegative(bbSort BackboneSort) bool {
	return bbSort == BBBoth || bbSort == BBNeg
}

func isBoth(bbSort BackboneSort) bool {
	return bbSort == BBBoth
}

func removeFromSlice(elems *[]int32, elem int32) {
	for i, e := range *elems {
		if e == elem {
			*elems = append((*elems)[:i], (*elems)[i+1:]...)
			return
		}
	}
}
