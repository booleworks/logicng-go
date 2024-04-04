package enum

import (
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/model"
	"github.com/booleworks/logicng-go/model/iter"
	"github.com/booleworks/logicng-go/sat"
)

// OnFormula enumerates all models of a formula over the given variables.  The
// additionalVariables will be included in each model, but are not iterated
// over.
func OnFormula(
	fac f.Factory,
	formula f.Formula,
	variables []f.Variable,
	additionalVariables ...f.Variable,
) []*model.Model {
	models, _ := OnFormulaWithConfig(fac, formula, variables, iter.DefaultConfig(), additionalVariables...)
	return models
}

// OnFormulaWithConfig enumerates all models of a formula over the given
// variables.  The additionalVariables will be included in each model, but are
// not iterated over.  The config can be used to influence the model iteration
// process by setting a handler and/or an iteration strategy.
func OnFormulaWithConfig(
	fac f.Factory,
	formula f.Formula,
	variables []f.Variable,
	config *iter.Config,
	additionalVariables ...f.Variable,
) ([]*model.Model, bool) {
	solver := sat.NewSolver(fac)
	solver.Add(formula)
	return OnSolverWithConfig(solver, variables, config, additionalVariables...)
}

// OnSolver enumerates all models on the given SAT solver over the given
// variables.  The additionalVariables will be included in each model, but are
// not iterated over.
func OnSolver(solver *sat.Solver, variables []f.Variable, additionalVariables ...f.Variable) []*model.Model {
	models, _ := OnSolverWithConfig(solver, variables, iter.DefaultConfig(), additionalVariables...)
	return models
}

// OnSolverWithConfig enumerates all models on the given SAT solver over
// the given variables.  The additionalVariables will be included in each
// model, but are not iterated over.  The config can be used to influence the
// model iteration process by setting a handler and/or an iteration strategy.
func OnSolverWithConfig(
	solver *sat.Solver,
	variables []f.Variable,
	config *iter.Config,
	additionalVariables ...f.Variable,
) ([]*model.Model, bool) {
	var add *f.VarSet
	if additionalVariables != nil {
		add = f.NewVarSet(additionalVariables...)
	}
	if config == nil {
		config = iter.DefaultConfig()
	}
	me := iter.New[[]*model.Model](f.NewVarSet(variables...), add, config)
	result, ok := me.Iterate(solver, newModelEnumCollector)
	return result, ok
}

type modelEnumCollector struct {
	committedModels                []*model.Model
	uncommittedModels              [][]f.Literal
	baseModels                     [][]f.Literal
	additionalVariablesNotOnSolver *f.LitSet
}

func newModelEnumCollector(
	fac f.Factory, _, dontCaresNotOnSolver, additionalVarsNotOnSolver *f.VarSet,
) iter.Collector[[]*model.Model] {
	baseModels := getCartesianProduct(fac, dontCaresNotOnSolver)
	addVars := f.NewMutableLitSet()
	for _, v := range additionalVarsNotOnSolver.Content() {
		addVars.Add(v.Negate(fac))
	}
	return &modelEnumCollector{
		committedModels:                []*model.Model{},
		uncommittedModels:              [][]f.Literal{},
		baseModels:                     baseModels,
		additionalVariablesNotOnSolver: addVars.AsImmutable(),
	}
}

func (c *modelEnumCollector) AddModel(
	modelFromSolver []bool, solver *sat.Solver, relevantAllIndices []int32, handler iter.Handler,
) bool {
	if handler == nil || handler.FoundModels(len(c.baseModels)) {
		model := solver.CoreSolver().CreateModel(solver.Factory(), modelFromSolver, relevantAllIndices)
		modelLiterals := c.additionalVariablesNotOnSolver.Content()
		modelLiterals = append(modelLiterals, model.Literals...)
		c.uncommittedModels = append(c.uncommittedModels, modelLiterals)
		return true
	} else {
		return false
	}
}

func (c *modelEnumCollector) Commit(handler iter.Handler) bool {
	c.committedModels = append(c.committedModels, c.expandUncommittedModels()...)
	c.uncommittedModels = make([][]f.Literal, 0, 4)
	return handler == nil || handler.Commit()
}

func (c *modelEnumCollector) Rollback(handler iter.Handler) bool {
	c.uncommittedModels = make([][]f.Literal, 0, 4)
	return handler == nil || handler.Rollback()
}

func (c *modelEnumCollector) RollbackAndReturnModels(_ *sat.Solver, handler iter.Handler) []*model.Model {
	modelsToReturn := make([]*model.Model, len(c.uncommittedModels))
	for i, lits := range c.uncommittedModels {
		modelsToReturn[i] = model.New(lits...)
	}
	c.Rollback(handler)
	return modelsToReturn
}

func (c *modelEnumCollector) Result() []*model.Model {
	return c.committedModels
}

func (c *modelEnumCollector) expandUncommittedModels() []*model.Model {
	expanded := make([]*model.Model, 0, len(c.baseModels))
	for _, baseModel := range c.baseModels {
		for _, uncommittedModel := range c.uncommittedModels {
			completeModel := make([]f.Literal, 0, len(baseModel)+len(uncommittedModel))
			completeModel = append(completeModel, baseModel...)
			completeModel = append(completeModel, uncommittedModel...)
			expanded = append(expanded, model.New(completeModel...))
		}
	}
	return expanded
}

func getCartesianProduct(fac f.Factory, variables *f.VarSet) [][]f.Literal {
	result := make([][]f.Literal, 1)
	result[0] = make([]f.Literal, 0)
	for _, v := range variables.Content() {
		extended := make([][]f.Literal, 0, len(result)*2)
		for _, literals := range result {
			extended = append(extended, extendedByLiteral(literals, v.AsLiteral()))
			extended = append(extended, extendedByLiteral(literals, v.Negate(fac)))
		}
		result = extended
	}
	return result
}

func extendedByLiteral(literals []f.Literal, lit f.Literal) []f.Literal {
	extended := make([]f.Literal, len(literals), len(literals)+1)
	copy(extended, literals)
	extended = append(extended, lit)
	return extended
}
