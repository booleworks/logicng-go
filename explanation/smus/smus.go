package smus

import (
	"fmt"

	"github.com/booleworks/logicng-go/event"
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
	smus, _ := ComputeForFormulasWithHandler(fac, formulas, handler.NopHandler, additionalConstraints...)
	return smus
}

// ComputeForFormulasWithHandler computes the SMUS for the given list of
// formulas modulo some additional constraints.  The optimization handler can
// be used to cancel the computation.  Returns the SMUS as a list of formulas.
func ComputeForFormulasWithHandler(
	fac f.Factory,
	formulas []f.Formula,
	hdl handler.Handler,
	additionalConstraints ...f.Formula,
) ([]f.Formula, handler.State) {
	props := make([]f.Proposition, len(formulas))
	for i, form := range formulas {
		props[i] = f.NewStandardProposition(form)
	}
	props, state := ComputeWithHandler(fac, props, hdl, additionalConstraints...)
	if !state.Success {
		return nil, state
	} else if props == nil {
		return nil, handler.Success()
	}
	forms := make([]f.Formula, len(props))
	for i, prop := range props {
		forms[i] = prop.Formula()
	}
	return forms, handler.Success()
}

// Compute computes the SMUS for the given list of propositions modulo some
// additional constraints.  Returns the SMUS as a list of propositions.
func Compute(
	fac f.Factory,
	propositions []f.Proposition,
	additionalConstraints ...f.Formula,
) []f.Proposition {
	smus, _ := ComputeWithHandler(fac, propositions, handler.NopHandler, additionalConstraints...)
	return smus
}

// ComputeWithHandler computes the SMUS for the given list of propositions
// modulo some additional constraints.  The optimization handler can be used to
// cancel the computation.  Returns the SMUS as a list of propositions.
func ComputeWithHandler(
	fac f.Factory,
	propositions []f.Proposition,
	hdl handler.Handler,
	additionalConstraints ...f.Formula,
) ([]f.Proposition, handler.State) {
	if e := event.SmusComputationStarted; !hdl.ShouldResume(e) {
		return nil, handler.Cancelation(e)
	}
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
	sResult := growSolver.Call(sat.Params().Handler(hdl).Variable(assumptions...))
	if sResult.Canceled() {
		return nil, sResult.State()
	}
	if sResult.Sat() {
		return nil, handler.Success()
	}
	hSolver := sat.NewSolver(fac)
	for {
		h, state := minimumHs(hSolver, assumptions, hdl)
		if !state.Success {
			return nil, state
		}
		c, state := grow(growSolver, h, assumptions, hdl)
		if !state.Success {
			return nil, state
		}
		if c == nil {
			props := make([]f.Proposition, len(h))
			for i, sel := range h {
				props[i] = propositionMapping[sel]
			}
			return props, handler.Success()
		}
		hSolver.Add(fac.Or(f.VariablesAsFormulas(c)...))
	}
}

func minimumHs(hSolver *sat.Solver, variables []f.Variable, hdl handler.Handler) ([]f.Variable, handler.State) {
	minimumHsModel, state := hSolver.MinimizeWithHandler(f.VariablesAsLiterals(variables), hdl)
	if !state.Success {
		return nil, state
	}
	return minimumHsModel.PosVars(), handler.Success()
}

func grow(growSolver *sat.Solver, h, variables []f.Variable, hdl handler.Handler) ([]f.Variable, handler.State) {
	solverState := growSolver.SaveState()
	growSolver.Add(f.VariablesAsFormulas(h)...)
	maxModel, state := growSolver.MaximizeWithHandler(f.VariablesAsLiterals(variables), hdl)
	if !state.Success {
		return nil, state
	} else if maxModel == nil {
		return nil, handler.Success()
	}
	err := growSolver.LoadState(solverState)
	if err != nil {
		panic(err)
	}
	minimumCorrectionSet := f.NewMutableVarSet(variables...)
	posVars := maxModel.PosVars()
	minimumCorrectionSet.RemoveAllElements(&posVars)
	return minimumCorrectionSet.Content(), handler.Success()
}
