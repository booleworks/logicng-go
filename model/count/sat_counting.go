package count

import (
	"math/big"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/model"
	"github.com/booleworks/logicng-go/model/iter"
	"github.com/booleworks/logicng-go/sat"
)

// OnFormula counts all models of a formula over the given variables.
func OnFormula(fac f.Factory, formula f.Formula, variables []f.Variable) *big.Int {
	count, _ := OnFormulaWithConfig(fac, formula, variables, iter.DefaultConfig())
	return count
}

// OnFormulaWithConfig counts all models of a formula over the given variables.
// The config can be used to influence the model iteration process by setting a
// handler and/or an iteration strategy.
func OnFormulaWithConfig(
	fac f.Factory,
	formula f.Formula,
	variables []f.Variable,
	config *iter.Config,
) (*big.Int, bool) {
	solver := sat.NewSolver(fac)
	solver.Add(formula)
	return OnSolverWithConfig(solver, variables, config)
}

// OnSolver counts all models on the given SAT solver over the given variables.
func OnSolver(solver *sat.Solver, variables []f.Variable) *big.Int {
	count, _ := OnSolverWithConfig(solver, variables, iter.DefaultConfig())
	return count
}

// OnSolverWithConfig counts all models on the given SAT solver over the given
// variables.  The config can be used to influence the model iteration process
// by setting a handler and/or an iteration strategy.
func OnSolverWithConfig(solver *sat.Solver, variables []f.Variable, config *iter.Config) (*big.Int, bool) {
	if config == nil {
		config = iter.DefaultConfig()
	}
	me := iter.New[*big.Int](f.NewVarSet(variables...), nil, config)
	result, ok := me.Iterate(solver, newModelCountCollector)
	return result, ok
}

type modelCountCollector struct {
	committedCount     *big.Int
	uncommittedModels  [][]bool
	uncommittedIndices [][]int32
	dontCareFactor     *big.Int
}

func newModelCountCollector(
	_ f.Factory, _, dontCaresNotOnSolver, _ *f.VarSet,
) iter.Collector[*big.Int] {
	dontCareFactor := big.NewInt(2).Exp(big.NewInt(2), big.NewInt(int64(dontCaresNotOnSolver.Size())), nil)
	return &modelCountCollector{
		committedCount:     big.NewInt(0),
		uncommittedModels:  make([][]bool, 0, 100),
		uncommittedIndices: make([][]int32, 0, 100),
		dontCareFactor:     dontCareFactor,
	}
}

func (c *modelCountCollector) AddModel(
	modelFromSolver []bool, _ *sat.Solver, relevantAllIndices []int32, handler iter.Handler,
) bool {
	if handler == nil || handler.FoundModels(int(c.dontCareFactor.Int64())) {
		c.uncommittedModels = append(c.uncommittedModels, modelFromSolver)
		c.uncommittedIndices = append(c.uncommittedIndices, relevantAllIndices)
		return true
	} else {
		return false
	}
}

func (c *modelCountCollector) Commit(handler iter.Handler) bool {
	mul := big.NewInt(1).Mul(big.NewInt(int64(len(c.uncommittedModels))), c.dontCareFactor)
	c.committedCount.Add(c.committedCount, mul)
	c.clearUncommitted()
	return handler == nil || handler.Commit()
}

func (c *modelCountCollector) Rollback(handler iter.Handler) bool {
	c.clearUncommitted()
	return handler == nil || handler.Rollback()
}

func (c *modelCountCollector) RollbackAndReturnModels(solver *sat.Solver, handler iter.Handler) []*model.Model {
	modelsToReturn := make([]*model.Model, len(c.uncommittedModels))
	for i, model := range c.uncommittedModels {
		modelsToReturn[i] = solver.CoreSolver().CreateModel(solver.Factory(), model, c.uncommittedIndices[i])
	}
	c.Rollback(handler)
	return modelsToReturn
}

func (c *modelCountCollector) Result() *big.Int {
	return c.committedCount
}

func (c *modelCountCollector) clearUncommitted() {
	c.uncommittedModels = make([][]bool, 0, 100)
	c.uncommittedIndices = make([][]int32, 0, 100)
}
