package iter

import (
	"math"

	"github.com/booleworks/logicng-go/sat"

	f "github.com/booleworks/logicng-go/formula"
)

const minModels = 3

// Strategy represents a strategy for fine-tuning the SAT solver based model
// iteration.
type Strategy interface {
	// MaxModelsForIter returns the maximum number of models to be iterated
	// on the given recursion depth.
	//
	// If this number of models is exceeded, the algorithm will compute new
	// split assignments and proceed to the next recursion step.
	//
	// This number refers to actual iterations on the solver, not to the
	// expanded number of models in the presence of don't care variables.
	MaxModelsForIter(recursionDepth int) int

	// MaxModelsForSplitAssignments returns the maximum number of models to be
	// iterated for split variables on the given recursion depth.
	//
	// This method is used to determine how many split assignments should at
	// most be computed and used for the next recursion step. If this limit is
	// exceeded, the algorithm will reduce the number of split variables using
	// ReduceSplitVars and then try again.
	MaxModelsForSplitAssignments(recursionDepth int) int

	// SplitVarsForRecursionDepth selects the split variables for the given
	// recursion depth from the given variables.
	//
	// This method is called before the algorithm makes another recursive call
	// to determine the initial split variables for this call.
	SplitVarsForRecursionDepth(variables *f.VarSet, solver *sat.Solver, recursionDepth int) *f.VarSet

	// ReduceSplitVars reduces the split variables for the given recursion
	// depth in case of MaxModelsForSplitAssignments was exceeded.
	ReduceSplitVars(variables *f.VarSet, recursionDepth int) *f.VarSet
}

// NoSplitStrategy is strategy for the model iteration where there are no splits.
type NoSplitStrategy struct{}

// NewNoSplitMEStrategy returns a new model iteration strategy without
// splits.
func NewNoSplitMEStrategy() NoSplitStrategy {
	return NoSplitStrategy{}
}

func (s NoSplitStrategy) MaxModelsForIter(_ int) int {
	return math.MaxInt
}

func (s NoSplitStrategy) MaxModelsForSplitAssignments(_ int) int {
	return math.MaxInt
}

func (s NoSplitStrategy) SplitVarsForRecursionDepth(
	_ *f.VarSet, _ *sat.Solver, _ int,
) *f.VarSet {
	return f.NewVarSet()
}

func (s NoSplitStrategy) ReduceSplitVars(_ *f.VarSet, _ int) *f.VarSet {
	return f.NewVarSet()
}

// BasicStrategy represents the default strategy for the model iteration.
//
// It takes a SplitVariableProvider and a maximum number of models.
//
// The split variable provider is used to compute the initial split variables
// on recursion depth 0. Afterward (including the ReduceSplitVars the reduction
// of split variables), always the first half of the variables is returned.
//
// MaxModels is always returned for both MaxModelsForIter and
// MaxModelsForSplitAssignments, ignoring the recursion depth.
//
// This struct can potentially be embedded if you want to fine-tune some
// methods, e.g. to change the maximum number of models depending on the
// recursion depth or whether the models are required for iteration or for
// split assignments.
type BasicStrategy struct {
	SplitProvider SplitProvider
	MaxModels     int
}

// NewBasicStrategy constructs a new default model iteration strategy
// with a given split provider and a maximum number of models before a split is
// performed.  In order to guarantee termination of the iteration algorithm,
// this number must be > 2. If a smaller number is provided, it is
// automatically set to 3.
func NewBasicStrategy(splitProvider SplitProvider, maxModels int) BasicStrategy {
	return BasicStrategy{splitProvider, maxModels}
}

// DefaultStrategy constructs a new default model iteration strategy with
// the most common split provider and a maximum number of 500 models.
func DefaultStrategy() BasicStrategy {
	return BasicStrategy{DefaultMostCommonVarProvider(), 500}
}

func (s BasicStrategy) MaxModelsForIter(_ int) int {
	return max(s.MaxModels, minModels)
}

func (s BasicStrategy) MaxModelsForSplitAssignments(_ int) int {
	return max(s.MaxModels, minModels)
}

func (s BasicStrategy) SplitVarsForRecursionDepth(
	variables *f.VarSet, solver *sat.Solver, recursionDepth int,
) *f.VarSet {
	if recursionDepth == 0 {
		return s.SplitProvider.Vars(solver, variables)
	}
	return s.ReduceSplitVars(variables, recursionDepth)
}

func (s BasicStrategy) ReduceSplitVars(variables *f.VarSet, _ int) *f.VarSet {
	return f.NewVarSet(variables.Content()[:variables.Size()/2]...)
}
