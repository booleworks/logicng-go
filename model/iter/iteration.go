package iter

import (
	"github.com/booleworks/logicng-go/event"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/sat"
)

var succ = handler.Success()

type EventIteratorFoundModels struct {
	NumberOfModels int
}

func (EventIteratorFoundModels) EventType() string {
	return "Model Iterator Found Models"
}

// A ModelIterator iterates over models on a solver.
type ModelIterator[R any] struct {
	vars           *f.VarSet
	additionalVars *f.VarSet
	hdl            handler.Handler
	strategy       Strategy
}

// New creates a new model iterator.  It will iterate over all models projected
// to the given vars.  The additionalVars will be included in each model, but
// they are not iterated over.  The config can be used to influence the
// iteration process with a split strategy an optional handler to cancel the
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
	emptyElement R,
) (R, handler.State) {
	if e := event.ModelEnumerationStarted; !m.hdl.ShouldResume(e) {
		return emptyElement, handler.Cancelation(e)
	}
	knownVariables := solver.CoreSolver().KnownVariables(solver.Factory())
	additionalVarsNotOnSolver := difference(m.additionalVars, knownVariables)
	dontCareVariablesNotOnSolver := difference(m.vars, knownVariables)
	collector := newCollector(solver.Factory(), knownVariables, dontCareVariablesNotOnSolver, additionalVarsNotOnSolver)
	initialSplitVars := f.NewMutableVarSet()
	if vars := m.strategy.SplitVarsForRecursionDepth(m.vars, solver, 0); vars.Size() > 0 {
		initialSplitVars.AddAll(vars)
	}
	state := m.iterRecursive(collector, solver, []f.Literal{}, m.vars, initialSplitVars.AsImmutable(), 0)
	return collector.Result(), state
}

func (m *ModelIterator[R]) iterRecursive(
	collector Collector[R],
	solver *sat.Solver,
	splitModel []f.Literal,
	iterVars *f.VarSet,
	splitVars *f.VarSet,
	recursionDepth int,
) handler.State {
	maxModelsForIter := m.strategy.MaxModelsForIter(recursionDepth)
	state := solver.SaveState()
	solver.Add(f.LiteralsAsFormulas(splitModel)...)
	iterFinished, iterState := iterate(collector, solver, iterVars, m.additionalVars, maxModelsForIter, m.hdl)
	if !iterState.Success {
		collector.Commit(m.hdl)
		return iterState
	}
	if !iterFinished {
		if s := collector.Rollback(m.hdl); !s.Success {
			err := solver.LoadState(state)
			if err != nil {
				panic(err)
			}
			return s
		}
		newSplitVars := f.NewVarSetCopy(splitVars)
		maxModelsForSplitAssignments := m.strategy.MaxModelsForSplitAssignments(recursionDepth)
		for {
			itForSplit, itState := iterate(collector, solver, newSplitVars, nil, maxModelsForSplitAssignments, m.hdl)
			if !itState.Success {
				err := solver.LoadState(state)
				if err != nil {
					panic(err)
				}
				collector.Rollback(m.hdl)
				return itState
			} else if itForSplit {
				break
			} else {
				if s := collector.Rollback(m.hdl); !s.Success {
					err := solver.LoadState(state)
					if err != nil {
						panic(err)
					}
					return s
				}
				newSplitVars = m.strategy.ReduceSplitVars(newSplitVars, recursionDepth)
			}
		}

		remainingVars := f.NewMutableVarSetCopy(iterVars)
		remainingVars.RemoveAll(newSplitVars)
		for _, literal := range splitModel {
			remainingVars.Remove(literal.Variable())
		}

		newSplitAssignments, st := collector.RollbackAndReturnModels(solver, m.hdl)
		if !st.Success {
			solver.LoadState(state)
			return st
		}
		recursiveSplitVars := m.strategy.SplitVarsForRecursionDepth(remainingVars.AsImmutable(), solver, recursionDepth+1)
		for _, newSplitAssignment := range newSplitAssignments {
			recursiveSplitAssignment := make([]f.Literal, newSplitAssignment.Size())
			copy(recursiveSplitAssignment, newSplitAssignment.Literals)
			recursiveSplitAssignment = append(recursiveSplitAssignment, splitModel...)
			m.iterRecursive(collector, solver, recursiveSplitAssignment, iterVars, recursiveSplitVars, recursionDepth+1)
			if s := collector.Commit(m.hdl); !s.Success {
				err := solver.LoadState(state)
				if err != nil {
					panic(err)
				}
				return s
			}
		}
	} else {
		if s := collector.Commit(m.hdl); !s.Success {
			err := solver.LoadState(state)
			if err != nil {
				panic(err)
			}
			return s
		}
	}
	err := solver.LoadState(state)
	if err != nil {
		panic(err)
	}
	return succ
}

func iterate[R any](
	collector Collector[R],
	solver *sat.Solver,
	variables *f.VarSet,
	additionalVariables *f.VarSet,
	maxModels int,
	hdl handler.Handler,
) (bool, handler.State) {
	stateBeforeIter := solver.SaveState()
	relevantIndices := relevantIndicesFromSolver(variables, solver)
	relevantAllIndices := relevantAllIndicesFromSolver(variables, additionalVariables, relevantIndices, solver)

	foundModels := 0
	state := handler.Success()
	for iterSATCall(solver, hdl) {
		modelFromSolver := solver.CoreSolver().Model()
		foundModels++
		if foundModels >= maxModels {
			err := solver.LoadState(stateBeforeIter)
			if err != nil {
				panic(err)
			}
			return false, succ
		}
		state = collector.AddModel(modelFromSolver, solver, relevantAllIndices, hdl)
		if state.Success && len(modelFromSolver) > 0 {
			blockingClause := generateBlockingClause(modelFromSolver, relevantIndices)
			solver.CoreSolver().AddClause(blockingClause, nil)
		} else {
			break
		}
	}
	err := solver.LoadState(stateBeforeIter)
	if err != nil {
		panic(err)
	}
	if !state.Success {
		return false, state
	}
	return true, succ
}

func iterSATCall(solver *sat.Solver, hdl handler.Handler) bool {
	sResult := solver.Call(sat.Params().Handler(hdl))
	return sResult.OK() && sResult.Sat()
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
	uniqueAdditionalVariables := f.NewMutableVarSet()
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
	result := f.NewMutableVarSet()
	if col1 != nil {
		result.AddAll(col1)
	}
	if col2 != nil {
		result.RemoveAll(col2)
	}
	return result.AsImmutable()
}
