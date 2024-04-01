package iter

import (
	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/handler"
	"booleworks.com/logicng/sat"
)

// A ModelIterator iterates over models on a solver.
type ModelIterator[R any] struct {
	vars           *f.VarSet
	additionalVars *f.VarSet
	handler        Handler
	strategy       Strategy
}

// New creates a new model iterator.  It will iterate over all models projected
// to the given vars.  The additionalVars will be included in each model, but
// they are not iterated over.  The config can be used to influence the
// iteration process with a split strategy an optional handler to abort the
// iteration.
func New[R any](vars, additionalVars *f.VarSet, config *Config) *ModelIterator[R] {
	strategy := config.Strategy
	if strategy == nil {
		strategy = NewNoSplitMEStrategy()
	}
	return &ModelIterator[R]{vars, additionalVars, config.Handler, strategy}
}

// Iterate is the main entry point for the model iteration.  It iterates over
// all model on the given solver.  The newCollector function is used to
// generate a new iteration collector for the given known variables on the
// solver, don't care variables, and additional variables.
func (m *ModelIterator[R]) Iterate(
	solver *sat.Solver,
	newCollector func(fac f.Factory, knownVars, dontCareVars, additionalVars *f.VarSet) Collector[R],
) (R, bool) {
	handler.Start(m.handler)
	knownVariables := solver.KnownVariables()
	additionalVarsNotOnSolver := difference(m.additionalVars, knownVariables)
	dontCareVariablesNotOnSolver := difference(m.vars, knownVariables)
	collector := newCollector(solver.Factory(), knownVariables, dontCareVariablesNotOnSolver, additionalVarsNotOnSolver)
	iterVars := m.getIterVars(solver.Factory(), knownVariables)
	initialSplitVars := f.NewVarSet()
	if vars := m.strategy.SplitVarsForRecursionDepth(iterVars, solver, 0); vars.Size() > 0 {
		initialSplitVars.AddAll(vars)
	}
	m.iterRecursive(collector, solver, []f.Literal{}, iterVars, initialSplitVars, 0)
	return collector.Result(), !handler.Aborted(m.handler)
}

func (m *ModelIterator[R]) iterRecursive(
	collector Collector[R],
	solver *sat.Solver,
	splitModel []f.Literal,
	iterVars *f.VarSet,
	splitVars *f.VarSet,
	recursionDepth int,
) {
	maxModelsForIter := m.strategy.MaxModelsForIter(recursionDepth)
	state := solver.SaveState()
	solver.Add(f.LiteralsAsFormulas(splitModel)...)
	iterFinished := iterate(collector, solver, iterVars, m.additionalVars, maxModelsForIter, m.handler)
	if !iterFinished {
		if !collector.Rollback(m.handler) {
			err := solver.LoadState(state)
			if err != nil {
				panic(err)
			}
			return
		}
		newSplitVars := f.NewVariableSetCopy(splitVars)
		maxModelsForSplitAssignments := m.strategy.MaxModelsForSplitAssignments(recursionDepth)
		for !iterate(collector, solver, newSplitVars, nil, maxModelsForSplitAssignments, m.handler) {
			if !collector.Rollback(m.handler) {
				err := solver.LoadState(state)
				if err != nil {
					panic(err)
				}
				return
			}
			newSplitVars = m.strategy.ReduceSplitVars(newSplitVars, recursionDepth)
		}
		if handler.Aborted(m.handler) {
			collector.Rollback(m.handler)
			return
		}

		remainingVars := f.NewVariableSetCopy(iterVars)
		remainingVars.RemoveAll(newSplitVars)
		for _, literal := range splitModel {
			remainingVars.Remove(literal.Variable())
		}

		newSplitAssignments := collector.RollbackAndReturnModels(solver, m.handler)
		recursiveSplitVars := m.strategy.SplitVarsForRecursionDepth(remainingVars, solver, recursionDepth+1)
		for _, newSplitAssignment := range newSplitAssignments {
			recursiveSplitAssignment := make([]f.Literal, newSplitAssignment.Size())
			copy(recursiveSplitAssignment, newSplitAssignment.Literals)
			recursiveSplitAssignment = append(recursiveSplitAssignment, splitModel...)
			m.iterRecursive(collector, solver, recursiveSplitAssignment, iterVars, recursiveSplitVars, recursionDepth+1)
			if !collector.Commit(m.handler) {
				err := solver.LoadState(state)
				if err != nil {
					panic(err)
				}
				return
			}
		}
	} else {
		if !collector.Commit(m.handler) {
			err := solver.LoadState(state)
			if err != nil {
				panic(err)
			}
			return
		}
	}
	err := solver.LoadState(state)
	if err != nil {
		panic(err)
	}
}

func iterate[R any](
	collector Collector[R],
	solver *sat.Solver,
	variables *f.VarSet,
	additionalVariables *f.VarSet,
	maxModels int,
	handler Handler,
) bool {
	stateBeforeIter := solver.SaveState()
	relevantIndices := relevantIndicesFromSolver(variables, solver)
	relevantAllIndices := relevantAllIndicesFromSolver(variables, additionalVariables, relevantIndices, solver)

	foundModels := 0
	proceed := true
	for proceed && iterSATCall(solver, handler) {
		modelFromSolver := solver.CoreSolver().Model()
		foundModels++
		if foundModels >= maxModels {
			err := solver.LoadState(stateBeforeIter)
			if err != nil {
				panic(err)
			}
			return false
		}
		proceed = collector.AddModel(modelFromSolver, solver, relevantAllIndices, handler)
		if len(modelFromSolver) > 0 {
			blockingClause := generateBlockingClause(modelFromSolver, relevantIndices)
			solver.CoreSolver().AddClause(blockingClause, nil)
			solver.SetResult(f.TristateUndef)
		} else {
			break
		}
	}
	err := solver.LoadState(stateBeforeIter)
	if err != nil {
		panic(err)
	}
	return true
}

func (m *ModelIterator[R]) getIterVars(fac f.Factory, knownVariables *f.VarSet) *f.VarSet {
	result := f.NewVarSet()
	for _, v := range knownVariables.Content() {
		if !fac.IsAuxVar(v) && (m.vars == nil || m.vars.Contains(v)) {
			result.Add(v)
		}
	}
	return result
}

func iterSATCall(solver *sat.Solver, handler Handler) bool {
	var satHandler sat.Handler
	if handler != nil {
		satHandler = handler.SatHandler()
	}
	sat, ok := solver.SatWithHandler(satHandler)
	return ok && sat
}

func generateBlockingClause(modelFromSolver []bool, relevantVars []int32) []int32 {
	var blockingClause []int32
	if relevantVars != nil {
		blockingClause = make([]int32, 0, len(relevantVars))
		for i := 0; i < len(relevantVars); i++ {
			varIndex := relevantVars[i]
			if varIndex != -1 {
				varAssignment := modelFromSolver[varIndex]
				var lit int32
				if varAssignment {
					lit = (varIndex * 2) ^ 1
				} else {
					lit = varIndex * 2
				}
				blockingClause = append(blockingClause, lit)
			}
		}
	} else {
		blockingClause = make([]int32, len(modelFromSolver))
		for i := int32(0); i < int32(len(modelFromSolver)); i++ {
			varAssignment := modelFromSolver[i]
			var lit int32
			if varAssignment {
				lit = (i * 2) ^ 1
			} else {
				lit = i * 2
			}
			blockingClause[i] = lit
		}
	}
	return blockingClause
}

func relevantIndicesFromSolver(variables *f.VarSet, solver *sat.Solver) []int32 {
	fac := solver.Factory()
	relevantIndices := make([]int32, variables.Size())
	for i, v := range variables.Content() {
		name, _ := fac.VarName(v)
		relevantIndices[i] = solver.CoreSolver().IdxForName(name)
	}
	return relevantIndices
}

func relevantAllIndicesFromSolver(
	variables *f.VarSet,
	additionalVariables *f.VarSet,
	relevantIndices []int32,
	solver *sat.Solver,
) []int32 {
	fac := solver.Factory()
	var relevantAllIndices []int32
	uniqueAdditionalVariables := f.NewVarSet()
	if additionalVariables != nil && additionalVariables.Size() > 0 {
		uniqueAdditionalVariables.AddAll(additionalVariables)
	}
	uniqueAdditionalVariables.RemoveAll(variables)
	if relevantIndices != nil {
		if uniqueAdditionalVariables.Empty() {
			relevantAllIndices = relevantIndices
		} else {
			relevantAllIndices = make([]int32, len(relevantIndices), len(relevantIndices)+uniqueAdditionalVariables.Size())
			copy(relevantAllIndices, relevantIndices)
			for _, v := range uniqueAdditionalVariables.Content() {
				name, _ := fac.VarName(v)
				idx := solver.CoreSolver().IdxForName(name)
				if idx != -1 {
					relevantAllIndices = append(relevantAllIndices, idx)
				}
			}
		}
	}
	return relevantAllIndices
}

func difference(col1, col2 *f.VarSet) *f.VarSet {
	result := f.NewVarSet()
	if col1 != nil {
		result.AddAll(col1)
	}
	if col2 != nil {
		result.RemoveAll(col2)
	}
	return result
}
