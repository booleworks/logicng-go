package count

import (
	"math/big"

	"github.com/booleworks/logicng-go/event"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
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
) (*big.Int, handler.State) {
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
func OnSolverWithConfig(solver *sat.Solver, variables []f.Variable, config *iter.Config) (*big.Int, handler.State) {
	if config == nil {
		config = iter.DefaultConfig()
	}
	me := iter.New[*big.Int](f.NewVarSet(variables...), nil, config)
	return me.Iterate(solver, newModelCountCollector, big.NewInt(0))
}

type modelCountCollector struct {
	committedCount     *big.Int
	uncommittedModels  [][]bool
	uncommittedIndices [][]int32
	dontCareFactor     *big.Int
}

func newModelCountCollector(_ f.Factory, _, dontCaresNotOnSolver, _ *f.VarSet) iter.Collector[*big.Int] {
	dontCareFactor := big.NewInt(2).Exp(big.NewInt(2), big.NewInt(int64(dontCaresNotOnSolver.Size())), nil)
	return &modelCountCollector{
		committedCount:     big.NewInt(0),
		uncommittedModels:  make([][]bool, 0, 100),
		uncommittedIndices: make([][]int32, 0, 100),
		dontCareFactor:     dontCareFactor,
	}
}

func (c *modelCountCollector) AddModel(
	modelFromSolver []bool, _ *sat.Solver, relevantAllIndices []int32, hdl handler.Handler,
) handler.State {
	e := iter.EventIteratorFoundModels{NumberOfModels: int(c.dontCareFactor.Int64())}
	c.uncommittedModels = append(c.uncommittedModels, modelFromSolver)
	c.uncommittedIndices = append(c.uncommittedIndices, relevantAllIndices)
	if !hdl.ShouldResume(e) {
		return handler.Cancelation(e)
	}
	return succ
}

func (c *modelCountCollector) Commit(hdl handler.Handler) handler.State {
	mul := big.NewInt(1).Mul(big.NewInt(int64(len(c.uncommittedModels))), c.dontCareFactor)
	c.committedCount.Add(c.committedCount, mul)
	c.clearUncommitted()
	if e := event.ModelEnumerationCommit; !hdl.ShouldResume(e) {
		return handler.Cancelation(e)
	}
	return succ
}

func (c *modelCountCollector) Rollback(hdl handler.Handler) handler.State {
	c.clearUncommitted()
	if e := event.ModelEnumerationRollback; !hdl.ShouldResume(e) {
		return handler.Cancelation(e)
	}
	return succ
}

func (c *modelCountCollector) RollbackAndReturnModels(solver *sat.Solver, hdl handler.Handler) []*model.Model {
	modelsToReturn := make([]*model.Model, len(c.uncommittedModels))
	for i, mdl := range c.uncommittedModels {
		modelsToReturn[i] = solver.CoreSolver().CreateModel(solver.Factory(), mdl, c.uncommittedIndices[i])
	}
	c.Rollback(hdl)
	return modelsToReturn
}

func (c *modelCountCollector) Result() *big.Int {
	return c.committedCount
}

func (c *modelCountCollector) clearUncommitted() {
	c.uncommittedModels = make([][]bool, 0, 100)
	c.uncommittedIndices = make([][]int32, 0, 100)
}
