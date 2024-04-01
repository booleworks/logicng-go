package smus

import (
	"fmt"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/sat"
)

const propositionSelector = "@PROPOSITION_SEL_"

// ComputeForFormulas computes the SMUS for the given list of formulas modulo
// some additional constraints.  Returns the SMUS as a list of formulas.
func ComputeForFormulas(
	fac f.Factory,
	formulas []f.Formula,
	additionalConstraints ...f.Formula,
) []f.Formula {
	smus, _ := ComputeForFormulasWithHandler(fac, formulas, nil, additionalConstraints...)
	return smus
}

// ComputeForFormulasWithHandler computes the SMUS for the given list of
// formulas modulo some additional constraints.  The optimization handler can
// be used to abort the computation.  Returns the SMUS as a list of formulas.
// If the handler aborted the computation, the ok flag is false.
func ComputeForFormulasWithHandler(
	fac f.Factory,
	formulas []f.Formula,
	optimizationHandler sat.OptimizationHandler,
	additionalConstraints ...f.Formula,
) (smus []f.Formula, ok bool) {
	props := make([]f.Proposition, len(formulas))
	for i, form := range formulas {
		props[i] = f.NewStandardProposition(form)
	}
	props, ok = ComputeWithHandler(fac, props, optimizationHandler, additionalConstraints...)
	if !ok {
		return nil, false
	} else if props == nil {
		return nil, true
	} else {
		forms := make([]f.Formula, len(props))
		for i, prop := range props {
			forms[i] = prop.Formula()
		}
		return forms, true
	}
}

// Compute computes the SMUS for the given list of propositions modulo some
// additional constraints.  Returns the SMUS as a list of propositions.
func Compute(
	fac f.Factory,
	propositions []f.Proposition,
	additionalConstraints ...f.Formula,
) []f.Proposition {
	smus, _ := ComputeWithHandler(fac, propositions, nil, additionalConstraints...)
	return smus
}

// ComputeWithHandler computes the SMUS for the given list of propositions
// modulo some additional constraints.  The optimization handler can be used to
// abort the computation.  Returns the SMUS as a list of propositions. If the
// handler aborted the computation, the ok flag is false.
func ComputeWithHandler(
	fac f.Factory,
	propositions []f.Proposition,
	optimizationHandler sat.OptimizationHandler,
	additionalConstraints ...f.Formula,
) (smus []f.Proposition, ok bool) {
	handler.Start(optimizationHandler)
	growSolver := sat.NewSolver(fac)
	for _, formula := range additionalConstraints {
		growSolver.Add(formula)
	}
	propositionMapping := make(map[f.Variable]f.Proposition)
	assumptions := make([]f.Variable, len(propositions))
	for i, proposition := range propositions {
		selector := fac.Var(fmt.Sprintf("%s%d", propositionSelector, len(propositionMapping)))
		assumptions[i] = selector
		propositionMapping[selector] = proposition
		growSolver.Add(fac.Equivalence(selector.AsFormula(), proposition.Formula()))
	}
	var satHandler sat.Handler
	if optimizationHandler != nil {
		satHandler = optimizationHandler.SatHandler()
	}
	satisfiable, ok := growSolver.SatWithHandler(satHandler, f.VariablesAsLiterals(assumptions)...)
	if !ok {
		return nil, false
	}
	if satisfiable {
		return nil, true
	}
	hSolver := sat.NewSolver(fac)
	for {
		h, ok := minimumHs(hSolver, assumptions, optimizationHandler)
		if !ok {
			return nil, false
		}
		c, ok := grow(growSolver, h, assumptions, optimizationHandler)
		if !ok {
			return nil, false
		}
		if c == nil {
			props := make([]f.Proposition, len(h))
			for i, sel := range h {
				props[i] = propositionMapping[sel]
			}
			return props, true
		}
		hSolver.Add(fac.Or(f.VariablesAsFormulas(c)...))
	}
}

func minimumHs(hSolver *sat.Solver, variables []f.Variable, handler sat.OptimizationHandler) ([]f.Variable, bool) {
	minimumHsModel, ok := hSolver.MinimizeWithHandler(f.VariablesAsLiterals(variables), handler)
	if !ok {
		return nil, false
	} else {
		return minimumHsModel.PosVars(), true
	}
}

func grow(growSolver *sat.Solver, h, variables []f.Variable, handler sat.OptimizationHandler) ([]f.Variable, bool) {
	solverState := growSolver.SaveState()
	growSolver.Add(f.VariablesAsFormulas(h)...)
	maxModel, ok := growSolver.MaximizeWithHandler(f.VariablesAsLiterals(variables), handler)
	if !ok {
		return nil, false
	} else if maxModel == nil {
		return nil, true
	} else {
		growSolver.LoadState(solverState)
		minimumCorrectionSet := f.NewVarSet(variables...)
		posVars := maxModel.PosVars()
		minimumCorrectionSet.RemoveAllElements(&posVars)
		return minimumCorrectionSet.Content(), true
	}
}
