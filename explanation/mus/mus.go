package mus

import (
	"github.com/booleworks/logicng-go/errorx"
	"github.com/booleworks/logicng-go/event"
	e "github.com/booleworks/logicng-go/explanation"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	s "github.com/booleworks/logicng-go/sat"
)

// ComputeInsertionBased computes a MUS using the insertion-based algorithm.
// The main idea of this algorithm is to start with an empty set and
// incrementally add propositions to the MUS which have been identified to be
// relevant.
//
// Returns an error if the formula is satisfiable.
func ComputeInsertionBased(fac f.Factory, propositions *[]f.Proposition) (*e.UnsatCore, error) {
	mus, _, err := ComputeInsertionBasedWithHandler(fac, propositions, handler.NopHandler)
	return mus, err
}

// ComputeInsertionBasedWithHandler computes a MUS using the insertion-based
// algorithm.  The given handler can be used to cancel the MUS computation.
//
// Returns an error if the formula is satisfiable.
func ComputeInsertionBasedWithHandler(
	fac f.Factory, propositions *[]f.Proposition, hdl handler.Handler,
) (*e.UnsatCore, handler.State, error) {
	if e := event.MusComputationStarted; !hdl.ShouldResume(e) {
		return nil, handler.Cancelation(e), nil
	}
	currentFormula := make([]f.Proposition, len(*propositions))
	copy(currentFormula, *propositions)
	mus := make([]f.Proposition, 0, len(*propositions))

	for len(currentFormula) > 0 {
		currentSubset := make([]f.Proposition, 0, len(*propositions))
		var transitionProposition f.Proposition
		solver := s.NewSolver(fac)
		solver.AddProposition(mus...)
		count := len(currentFormula)
		for {
			sat := solver.Call(s.Params().Handler(hdl))
			if sat.Canceled() {
				return nil, sat.State(), nil
			}
			if !sat.Sat() {
				break
			}
			if count == 0 {
				return nil, handler.Success(), errorx.BadInput("formula set is satisfiable")
			}
			count--
			removeProposition := currentFormula[count]
			currentSubset = append(currentSubset, removeProposition)
			transitionProposition = removeProposition
			solver.AddProposition(removeProposition)
		}
		currentFormula = make([]f.Proposition, len(currentSubset))
		copy(currentFormula, currentSubset)
		if transitionProposition != nil {
			for i, p := range currentFormula {
				if p == transitionProposition {
					currentFormula = append(currentFormula[:i], currentFormula[i+1:]...)
				}
			}
			mus = append(mus, transitionProposition)
		}
	}
	return e.NewUnsatCore(mus, true), handler.Success(), nil
}

// ComputeDeletionBased computes a MUS using the deletion-based algorithm.
// The main idea of this algorithm is to start with all given formulas and
// iteratively test each formula for relevance. A formula is relevant for the
// conflict, if its removal yields in a satisfiable set of formulas. Only the
// relevant formulas are kept.
//
// Returns an error if the formula is satisfiable.
func ComputeDeletionBased(fac f.Factory, propositions *[]f.Proposition) (*e.UnsatCore, error) {
	mus, _, err := ComputeDeletionBasedWithHandler(fac, propositions, handler.NopHandler)
	return mus, err
}

// ComputeDeletionBasedWithHandler computes a MUS using the deletion-based
// algorithm. The given handler can be used to cancel the MUS computation.
//
// Returns an error if the formula is satisfiable.
func ComputeDeletionBasedWithHandler(
	fac f.Factory, propositions *[]f.Proposition, hdl handler.Handler,
) (*e.UnsatCore, handler.State, error) {
	if e := event.MusComputationStarted; !hdl.ShouldResume(e) {
		return nil, handler.Cancelation(e), nil
	}
	mus := make([]f.Proposition, 0, len(*propositions))
	solverStates := make([]*s.SolverState, len(*propositions))
	solver := s.NewSolver(fac)
	for i, p := range *propositions {
		solverStates[i] = solver.SaveState()
		solver.AddProposition(p)
	}
	sResult := solver.Call(s.Params().Handler(hdl))
	if sResult.Canceled() {
		return nil, sResult.State(), nil
	}
	if sResult.Sat() {
		return nil, handler.Success(), errorx.BadInput("formula set is satisfiable")
	}
	for i := len(solverStates) - 1; i >= 0; i-- {
		err := solver.LoadState(solverStates[i])
		if err != nil {
			return nil, handler.Success(), err
		}
		for _, prop := range mus {
			solver.AddProposition(prop)
		}
		sResult := solver.Call(s.Params().Handler(hdl))
		if sResult.Canceled() {
			return nil, sResult.State(), nil
		}
		if sResult.Sat() {
			mus = append(mus, (*propositions)[i])
		}
	}
	return e.NewUnsatCore(mus, true), handler.Success(), nil
}
