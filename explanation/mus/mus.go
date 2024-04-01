package mus

import (
	"booleworks.com/logicng/errorx"
	e "booleworks.com/logicng/explanation"
	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/handler"
	s "booleworks.com/logicng/sat"
)

// ComputeInsertionBased computes a MUS using the insertion-based algorithm.
// The main idea of this algorithm is to start with an empty set and
// incrementally add propositions to the MUS which have been identified to be
// relevant.
//
// Returns an error if the formula is satisfiable.
func ComputeInsertionBased(fac f.Factory, propositions *[]f.Proposition) (*e.UnsatCore, error) {
	mus, _, err := ComputeInsertionBasedWithHandler(fac, propositions, nil)
	return mus, err
}

// ComputeInsertionBasedWithHandler computes a MUS using the insertion-based
// algorithm.  The given SAT handler can be used to abort the MUS computation.
// If the computation was aborted by the handler, the ok flag is false.
//
// Returns an error if the formula is satisfiable.
func ComputeInsertionBasedWithHandler(
	fac f.Factory, propositions *[]f.Proposition, satHandler s.Handler,
) (unsatCore *e.UnsatCore, ok bool, err error) {
	handler.Start(satHandler)
	currentFormula := make([]f.Proposition, len(*propositions))
	copy(currentFormula, *propositions)
	mus := make([]f.Proposition, 0, len(*propositions))

	for len(currentFormula) > 0 {
		currentSubset := make([]f.Proposition, 0, len(*propositions))
		var transitionProposition f.Proposition
		solver := s.NewSolver(fac)
		solver.AddProposition(mus...)
		count := len(currentFormula)
		for shouldProceed(solver, satHandler) {
			if count == 0 {
				return nil, true, errorx.BadInput("formula set is satisfiable")
			}
			count--
			removeProposition := currentFormula[count]
			currentSubset = append(currentSubset, removeProposition)
			transitionProposition = removeProposition
			solver.AddProposition(removeProposition)
		}
		if handler.Aborted(satHandler) {
			return nil, false, nil
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
	return e.NewUnsatCore(mus, true), true, nil
}

// ComputeDeletionBased computes a MUS using the deletion-based algorithm.
// The main idea of this algorithm is to start with all given formulas and
// iteratively test each formula for relevance. A formula is relevant for the
// conflict, if its removal yields in a satisfiable set of formulas. Only the
// relevant formulas are kept.
//
// Returns an error if the formula is satisfiable.
func ComputeDeletionBased(fac f.Factory, propositions *[]f.Proposition) (*e.UnsatCore, error) {
	mus, _, err := ComputeDeletionBasedWithHandler(fac, propositions, nil)
	return mus, err
}

// ComputeDeletionBasedWithHandler computes a MUS using the deletion-based
// algorithm. The given SAT handler can be used to abort the MUS computation.
// If the computation was aborted by the handler, the ok flag is false.
//
// Returns an error if the formula is satisfiable.
func ComputeDeletionBasedWithHandler(
	fac f.Factory, propositions *[]f.Proposition, satHandler s.Handler,
) (unsatCore *e.UnsatCore, ok bool, err error) {
	handler.Start(satHandler)
	mus := make([]f.Proposition, 0, len(*propositions))
	solverStates := make([]*s.SolverState, len(*propositions))
	solver := s.NewSolver(fac)
	for i, p := range *propositions {
		solverStates[i] = solver.SaveState()
		solver.AddProposition(p)
	}
	sat, ok := solver.SatWithHandler(satHandler)
	if !ok {
		return nil, false, nil
	}
	if sat {
		return nil, true, errorx.BadInput("formula set is satisfiable")
	}
	for i := len(solverStates) - 1; i >= 0; i-- {
		err := solver.LoadState(solverStates[i])
		if err != nil {
			return nil, false, err
		}
		for _, prop := range mus {
			solver.AddProposition(prop)
		}
		sat, ok := solver.SatWithHandler(satHandler)
		if !ok {
			return nil, false, nil
		}
		if sat {
			mus = append(mus, (*propositions)[i])
		}
	}
	return e.NewUnsatCore(mus, true), true, nil
}

func shouldProceed(solver *s.Solver, handler s.Handler) bool {
	sat, ok := solver.SatWithHandler(handler)
	return sat && ok
}
