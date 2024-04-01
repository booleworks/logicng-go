package enum

import (
	"math"

	"github.com/booleworks/logicng-go/bdd"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/model"
	"github.com/booleworks/logicng-go/model/iter"
	"github.com/booleworks/logicng-go/sat"
)

// ToBDDOnFormula enumerates all models of a formula over the given variables
// and tathers the result in a BDD.
func ToBDDOnFormula(
	fac f.Factory,
	formula f.Formula,
	variables []f.Variable,
) *bdd.BDD {
	models, _ := ToBDDOnFormulaWithConfig(fac, formula, variables, iter.DefaultConfig())
	return models
}

// ToBDDOnFormulaWithConfig enumerates all models of a formula over the given
// variables and gathers the result in a BDD.  The config can be used to
// influence the model iteration process by setting a handler and/or an
// iteration strategy.
func ToBDDOnFormulaWithConfig(
	fac f.Factory,
	formula f.Formula,
	variables []f.Variable,
	config *iter.Config,
) (*bdd.BDD, bool) {
	solver := sat.NewSolver(fac)
	solver.Add(formula)
	return ToBDDOnSolverWithConfig(solver, variables, config)
}

// ToBDDOnSolver enumerates all models on the given SAT solver over the given
// variables and gathers the result in a BDD.
func ToBDDOnSolver(solver *sat.Solver, variables []f.Variable) *bdd.BDD {
	models, _ := ToBDDOnSolverWithConfig(solver, variables, iter.DefaultConfig())
	return models
}

// ToBDDOnSolverWithConfig enumerates all models on the given SAT solver over
// the given variables and gathers the result in a BDD.  The config can be used
// to influence the model iteration process by setting a handler and/or an
// iteration strategy.
func ToBDDOnSolverWithConfig(
	solver *sat.Solver,
	variables []f.Variable,
	config *iter.Config,
) (*bdd.BDD, bool) {
	if config == nil {
		config = iter.DefaultConfig()
	}
	me := iter.New[*bdd.BDD](f.NewVarSet(variables...), nil, config)
	result, ok := me.Iterate(solver, generateModelCollector(variables))
	return result, ok
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
			sortedVariables = f.NewVariableSetCopy(knownVariables)
		}
		numVars := sortedVariables.Size()
		kernel := bdd.NewKernelWithOrdering(fac, sortedVariables.Content(), int32(numVars*30), int32(numVars*50))
		committedModels := bdd.BuildWithKernel(fac, fac.Falsum(), kernel)
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
	modelFromSolver []bool, solver *sat.Solver, relevantAllIndices []int32, handler iter.Handler,
) bool {
	if handler == nil || handler.FoundModels(c.dontCareFactor) {
		model := solver.CoreSolver().CreateModel(solver.Factory(), modelFromSolver, relevantAllIndices)
		c.uncommittedModels = append(c.uncommittedModels, model)
		return true
	} else {
		return false
	}
}

func (c *modelEnumBddCollector) Commit(handler iter.Handler) bool {
	for _, uncommittedModel := range c.uncommittedModels {
		modelFormula := uncommittedModel.Formula(c.kernel.Factory())
		modelBdd := bdd.BuildWithKernel(c.kernel.Factory(), modelFormula, c.kernel)
		c.committedModels = c.committedModels.Or(modelBdd)
	}
	c.uncommittedModels = make([]*model.Model, 0)
	return handler == nil || handler.Commit()
}

func (c *modelEnumBddCollector) Rollback(handler iter.Handler) bool {
	c.uncommittedModels = make([]*model.Model, 0)
	return handler == nil || handler.Rollback()
}

func (c *modelEnumBddCollector) RollbackAndReturnModels(_ *sat.Solver, handler iter.Handler) []*model.Model {
	modelsToReturn := make([]*model.Model, len(c.uncommittedModels))
	copy(modelsToReturn, c.uncommittedModels)
	c.Rollback(handler)
	return modelsToReturn
}

func (c *modelEnumBddCollector) Result() *bdd.BDD {
	return c.committedModels
}
