package enum

import (
	"math"

	"github.com/booleworks/logicng-go/bdd"
	"github.com/booleworks/logicng-go/event"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/model"
	"github.com/booleworks/logicng-go/model/iter"
	"github.com/booleworks/logicng-go/sat"
)

// ToBddOnFormula enumerates all models of a formula over the given variables
// and gathers the result in a BDD.
func ToBddOnFormula(
	fac f.Factory,
	formula f.Formula,
	variables []f.Variable,
) *bdd.BDD {
	models, _ := ToBddOnFormulaWithConfig(fac, formula, variables, iter.DefaultConfig())
	return models
}

// ToBddOnFormulaWithConfig enumerates all models of a formula over the given
// variables and gathers the result in a BDD.  The config can be used to
// influence the model iteration process by setting a handler and/or an
// iteration strategy.
func ToBddOnFormulaWithConfig(
	fac f.Factory,
	formula f.Formula,
	variables []f.Variable,
	config *iter.Config,
) (*bdd.BDD, handler.State) {
	solver := sat.NewSolver(fac)
	solver.Add(formula)
	return ToBddOnSolverWithConfig(solver, variables, config)
}

// ToBddOnSolver enumerates all models on the given SAT solver over the given
// variables and gathers the result in a BDD.
func ToBddOnSolver(solver *sat.Solver, variables []f.Variable) *bdd.BDD {
	models, _ := ToBddOnSolverWithConfig(solver, variables, iter.DefaultConfig())
	return models
}

// ToBddOnSolverWithConfig enumerates all models on the given SAT solver over
// the given variables and gathers the result in a BDD.  The config can be used
// to influence the model iteration process by setting a handler and/or an
// iteration strategy.
func ToBddOnSolverWithConfig(
	solver *sat.Solver,
	variables []f.Variable,
	config *iter.Config,
) (*bdd.BDD, handler.State) {
	if config == nil {
		config = iter.DefaultConfig()
	}
	me := iter.New[*bdd.BDD](f.NewVarSet(variables...), nil, config)
	falseBdd := bdd.Compile(solver.Factory(), solver.Factory().Falsum())
	return me.Iterate(solver, generateModelCollector(variables), falseBdd)
}

type modelEnumBddCollector struct {
	kernel            *bdd.Kernel
	committedModels   *bdd.BDD
	uncommittedModels []*model.Model
	dontCareFactor    int
}

func generateModelCollector(
	variables []f.Variable,
) func(fac f.Factory, knownVars, dontCareVars, additionalVars *f.VarSet) iter.Collector[*bdd.BDD] {
	return func(fac f.Factory, knownVariables, dontCaresNotOnSolver, additionalVarsNotOnSolver *f.VarSet) iter.Collector[*bdd.BDD] {
		var sortedVariables *f.VarSet
		if variables != nil {
			sortedVariables = f.NewVarSet(variables...)
		} else {
			sortedVariables = f.NewVarSetCopy(knownVariables)
		}
		numVars := sortedVariables.Size()
		kernel := bdd.NewKernelWithOrdering(fac, sortedVariables.Content(), int32(numVars*30), int32(numVars*50))
		committedModels := bdd.CompileWithKernel(fac, fac.Falsum(), kernel)
		dontCareFactor := int(math.Pow(2, float64(dontCaresNotOnSolver.Size())))
		return &modelEnumBddCollector{
			kernel:            kernel,
			committedModels:   committedModels,
			uncommittedModels: make([]*model.Model, 0, 4),
			dontCareFactor:    dontCareFactor,
		}
	}
}

func (c *modelEnumBddCollector) AddModel(
	modelFromSolver []bool, solver *sat.Solver, relevantAllIndices []int32, hdl handler.Handler,
) handler.State {
	e := iter.EventIteratorFoundModels{NumberOfModels: c.dontCareFactor}
	mdl := solver.CoreSolver().CreateModel(solver.Factory(), modelFromSolver, relevantAllIndices)
	c.uncommittedModels = append(c.uncommittedModels, mdl)
	if !hdl.ShouldResume(e) {
		return handler.Cancelation(e)
	}
	return succ
}

func (c *modelEnumBddCollector) Commit(hdl handler.Handler) handler.State {
	for _, uncommittedModel := range c.uncommittedModels {
		modelFormula := uncommittedModel.Formula(c.kernel.Factory())
		modelBdd := bdd.CompileWithKernel(c.kernel.Factory(), modelFormula, c.kernel)
		c.committedModels = c.committedModels.Or(modelBdd)
	}
	c.uncommittedModels = make([]*model.Model, 0)
	if e := event.ModelEnumerationCommit; !hdl.ShouldResume(e) {
		return handler.Cancelation(e)
	}
	return succ
}

func (c *modelEnumBddCollector) Rollback(hdl handler.Handler) handler.State {
	c.uncommittedModels = make([]*model.Model, 0)
	if e := event.ModelEnumerationRollback; !hdl.ShouldResume(e) {
		return handler.Cancelation(e)
	}
	return succ
}

func (c *modelEnumBddCollector) RollbackAndReturnModels(_ *sat.Solver, hdl handler.Handler) ([]*model.Model, handler.State) {
	modelsToReturn := make([]*model.Model, len(c.uncommittedModels))
	copy(modelsToReturn, c.uncommittedModels)
	state := c.Rollback(hdl)
	return modelsToReturn, state
}

func (c *modelEnumBddCollector) Result() *bdd.BDD {
	return c.committedModels
}
