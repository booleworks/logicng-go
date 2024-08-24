package sat

import (
	"fmt"

	"github.com/booleworks/logicng-go/event"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/model"
)

type EventFoundBetterBound struct {
	Model func() *model.Model
}

func (e EventFoundBetterBound) EventType() string {
	return "Found Better Bound"
}

const selPrefix = "@SEL_OPT_"

// Maximize searches for a model on the solver with the maximum of the given
// literals set to true.  The returned model will also include the additional
// variables.
func (s *Solver) Maximize(literals []f.Literal, additionalVariables ...f.Variable) *model.Model {
	opt, _ := s.optimize(true, literals, handler.NopHandler, additionalVariables)
	return opt
}

// MaximizeWithHandler searches for a model on the solver with the maximum of
// the given literals set to true.  The returned model will also include the
// additional variables.  The given optimizationHandler can be used to cancel
// the optimization process.
func (s *Solver) MaximizeWithHandler(
	literals []f.Literal, hdl handler.Handler, additionalVariables ...f.Variable,
) (*model.Model, handler.State) {
	return s.optimize(true, literals, hdl, additionalVariables)
}

// Minimize searches for a model on the solver with the minimum of the given
// literals set to true.  The returned model will also include the additional
// variables.
func (s *Solver) Minimize(literals []f.Literal, additionalVariables ...f.Variable) *model.Model {
	opt, _ := s.optimize(false, literals, handler.NopHandler, additionalVariables)
	return opt
}

// MinimizeWithHandler searches for a model on the solver with the minimum of
// the given literals set to true.  The returned model will also include the
// additional variables.  The given optimizationHandler can be used to cancel
// the optimization process.
func (s *Solver) MinimizeWithHandler(
	literals []f.Literal, hdl handler.Handler, additionalVariables ...f.Variable,
) (*model.Model, handler.State) {
	return s.optimize(false, literals, hdl, additionalVariables)
}

func (s *Solver) optimize(
	maximize bool,
	literals []f.Literal,
	hdl handler.Handler,
	additionalVariables []f.Variable,
) (*model.Model, handler.State) {
	initialState := s.SaveState()
	resultModelVariables := f.NewMutableVarSet(additionalVariables...)
	for _, lit := range literals {
		variable := lit.Variable()
		resultModelVariables.Add(variable)
	}
	relevantIndices := make([]int32, 0, resultModelVariables.Size())
	for _, variable := range resultModelVariables.Content() {
		name, _ := s.fac.VarName(variable)
		idx := s.core.IdxForName(name)
		if idx != -1 {
			relevantIndices = append(relevantIndices, idx)
		}
	}
	mdl, state := s.maximize(maximize, literals, relevantIndices, hdl)
	_ = s.LoadState(initialState)
	return mdl, state
}

func (s *Solver) maximize(
	maximize bool,
	literals []f.Literal,
	relevantIndices []int32,
	hdl handler.Handler,
) (*model.Model, handler.State) {
	if !hdl.ShouldResume(event.OptimizationFunctionStarted) {
		return nil, handler.Cancellation(event.OptimizationFunctionStarted)
	}
	fac := s.fac
	selectorMap := make(map[f.Variable]f.Literal)
	selectors := make([]f.Variable, len(literals))

	for i, lit := range literals {
		selVar := fac.Var(fmt.Sprintf("%s%d", selPrefix, len(selectorMap)))
		selectorMap[selVar] = lit
		selectors[i] = selVar
	}

	for _, selVar := range selectors {
		lit := selectorMap[selVar]
		if maximize {
			s.Add(fac.Clause(selVar.Negate(fac), lit))
			s.Add(fac.Clause(lit.Negate(fac), selVar.AsLiteral()))
		} else {
			s.Add(fac.Clause(selVar.Negate(fac), lit.Negate(fac)))
			s.Add(fac.Clause(lit, selVar.AsLiteral()))
		}
	}

	params := Params().Handler(hdl).WithModel(selectors)
	sResult := s.Call(params)
	if sResult.Cancelled() {
		return nil, sResult.state
	}
	if !sResult.Sat() {
		return nil, succ
	}
	internalModel := s.core.Model()
	currentModel := sResult.Model()
	currentBound := len(currentModel.PosVars())

	if currentBound == 0 {
		s.Add(fac.CC(f.GE, 1, selectors...))
		sResult = s.Call(params)
		if sResult.Cancelled() {
			return nil, sResult.state
		} else if !sResult.Sat() {
			return s.core.CreateModel(s.fac, internalModel, relevantIndices), succ
		} else {
			internalModel = s.core.Model()
			currentModel = sResult.Model()
			currentBound = len(currentModel.PosVars())
		}
	} else if currentBound == len(selectors) {
		return s.core.CreateModel(s.fac, internalModel, relevantIndices), succ
	}

	cc := fac.CC(f.GE, uint32(currentBound+1), selectors...)

	incrementalData, _ := s.AddIncrementalCC(cc)
	sResult = s.Call(params)
	if sResult.Cancelled() {
		return nil, sResult.state
	}

	for sResult.Sat() {
		internalModel = s.core.Model()
		betterBoundEvent := EventFoundBetterBound{func() *model.Model {
			return s.core.CreateModel(s.fac, internalModel, relevantIndices)
		}}
		if !hdl.ShouldResume(betterBoundEvent) {
			return nil, handler.Cancellation(betterBoundEvent)
		}
		currentModel = sResult.Model()
		currentBound = len(currentModel.PosVars())
		if currentBound == len(selectors) {
			return s.core.CreateModel(s.fac, internalModel, relevantIndices), succ
		}
		incrementalData.NewLowerBoundForSolver(currentBound + 1)
		sResult = s.Call(params)
		if sResult.Cancelled() {
			return nil, sResult.state
		}
	}
	return s.core.CreateModel(s.fac, internalModel, relevantIndices), succ
}
