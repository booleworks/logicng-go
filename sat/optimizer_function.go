package sat

import (
	"fmt"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/model"
)

const selPrefix = "@SEL_OPT_"

// Maximize searches for a model on the solver with the maximum of the given
// literals set to true.  The returned model will also include the additional
// variables.
func (s *Solver) Maximize(literals []f.Literal, additionalVariables ...f.Variable) *model.Model {
	opt, _ := s.optimize(true, literals, nil, additionalVariables)
	return opt
}

// MaximizeWithHandler searches for a model on the solver with the maximum of
// the given literals set to true.  The returned model will also include the
// additional variables.  The given optimizationHandler can be used to abort
// the optimization process.  The ok flag is false when the computation was
// aborted by the handler.
func (s *Solver) MaximizeWithHandler(
	literals []f.Literal, optimizationHandler OptimizationHandler, additionalVariables ...f.Variable,
) (mdl *model.Model, ok bool) {
	return s.optimize(true, literals, optimizationHandler, additionalVariables)
}

// Minimize searches for a model on the solver with the minimum of the given
// literals set to true.  The returned model will also include the additional
// variables.
func (s *Solver) Minimize(literals []f.Literal, additionalVariables ...f.Variable) *model.Model {
	opt, _ := s.optimize(false, literals, nil, additionalVariables)
	return opt
}

// MinimizeWithHandler searches for a model on the solver with the minimum of
// the given literals set to true.  The returned model will also include the
// additional variables.  The given optimizationHandler can be used to abort
// the optimization process.  The ok flag is false when the computation was
// aborted by the handler.
func (s *Solver) MinimizeWithHandler(
	literals []f.Literal, optimizationHandler OptimizationHandler, additionalVariables ...f.Variable,
) (*model.Model, bool) {
	return s.optimize(false, literals, optimizationHandler, additionalVariables)
}

func (s *Solver) optimize(
	maximize bool,
	literals []f.Literal,
	optimizationHandler OptimizationHandler,
	additionalVariables []f.Variable,
) (*model.Model, bool) {
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

	mdl, ok := s.maximize(maximize, literals, relevantIndices, optimizationHandler)
	_ = s.LoadState(initialState)
	return mdl, ok
}

func (s *Solver) maximize(
	maximize bool,
	literals []f.Literal,
	relevantIndices []int32,
	optimizationHandler OptimizationHandler,
) (*model.Model, bool) {
	handler.Start(optimizationHandler)
	fac := s.fac
	selectorMap := make(map[f.Variable]f.Literal)
	selectors := make([]f.Variable, len(literals))

	for i, lit := range literals {
		selVar := fac.Var(fmt.Sprintf("%s%d", selPrefix, len(selectorMap)))
		selectorMap[selVar] = lit
		selectors[i] = selVar
	}

	for selVar, lit := range selectorMap {
		if maximize {
			s.Add(fac.Clause(selVar.Negate(fac), lit))
			s.Add(fac.Clause(lit.Negate(fac), selVar.AsLiteral()))
		} else {
			s.Add(fac.Clause(selVar.Negate(fac), lit.Negate(fac)))
			s.Add(fac.Clause(lit, selVar.AsLiteral()))
		}
	}

	params := Params().Handler(satHandler(optimizationHandler)).WithModel(selectors)
	sResult := s.Call(params)
	if sResult.Aborted() {
		return nil, false
	}
	if !sResult.Sat() {
		return nil, true
	}
	internalModel := s.core.Model()
	currentModel := sResult.Model()
	currentBound := len(currentModel.PosVars())

	if currentBound == 0 {
		s.Add(fac.CC(f.GE, 1, selectors...))
		sResult = s.Call(params)
		if sResult.Aborted() {
			return nil, false
		} else if !sResult.Sat() {
			return s.core.CreateModel(s.fac, internalModel, relevantIndices), true
		} else {
			internalModel = s.core.Model()
			currentModel = sResult.Model()
			currentBound = len(currentModel.PosVars())
		}
	} else if currentBound == len(selectors) {
		return s.core.CreateModel(s.fac, internalModel, relevantIndices), true
	}

	cc := fac.CC(f.GE, uint32(currentBound+1), selectors...)

	incrementalData, _ := s.AddIncrementalCC(cc)
	sResult = s.Call(params)
	if sResult.Aborted() {
		optimizationHandler.SetModel(s.core.CreateModel(s.fac, internalModel, relevantIndices))
		return nil, false
	}

	for sResult.Sat() {
		internalModel = s.core.Model()
		if optimizationHandler != nil &&
			!optimizationHandler.FoundBetterBound(s.core.CreateModel(s.fac, internalModel, relevantIndices)) {
			return nil, false
		}
		currentModel = sResult.Model()
		currentBound = len(currentModel.PosVars())
		if currentBound == len(selectors) {
			return s.core.CreateModel(s.fac, internalModel, relevantIndices), true
		}
		incrementalData.NewLowerBoundForSolver(currentBound + 1)
		sResult = s.Call(params)
		if sResult.Aborted() {
			optimizationHandler.SetModel(s.core.CreateModel(s.fac, internalModel, relevantIndices))
			return nil, false
		}
	}
	return s.core.CreateModel(s.fac, internalModel, relevantIndices), true
}
