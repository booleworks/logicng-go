package sat

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/io"
	"github.com/booleworks/logicng-go/model"
	"github.com/booleworks/logicng-go/parser"
	"github.com/booleworks/logicng-go/randomizer"
	"github.com/stretchr/testify/assert"
)

func solvers(fac f.Factory) []*Solver {
	solvers := make([]*Solver, 4)

	mc := DefaultConfig()
	mc.InitialPhase = true
	solvers[0] = NewSolver(fac, mc)

	mc = DefaultConfig()
	mc.InitialPhase = false
	solvers[1] = NewSolver(fac, mc)

	mc = DefaultConfig()
	mc.InitialPhase = false
	mc.UseAtMostClauses = true
	solvers[2] = NewSolver(fac, mc)

	mc = DefaultConfig()
	mc.InitialPhase = true
	mc.UseAtMostClauses = true
	solvers[3] = NewSolver(fac, mc)

	return solvers
}

func TestOptimizerFunctionUnsat(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	parser := parser.New(fac)
	formula := parser.ParseUnsafe("a & b & (a => ~b)")
	vars := f.VariablesAsLiterals(f.Variables(fac, formula).Content())

	for _, solver := range solvers(fac) {
		minimumModel := optimize([]f.Formula{formula}, vars, []f.Variable{}, false, solver)
		assert.Nil(minimumModel)
		maximumModel := optimize([]f.Formula{formula}, vars, []f.Variable{}, true, solver)
		assert.Nil(maximumModel)
	}
}

func TestOptimizerFunctionSingleModel(t *testing.T) {
	fac := f.NewFactory()
	parser := parser.New(fac)
	formula := parser.ParseUnsafe("~a & ~b & ~c")
	vars := f.VariablesAsLiterals(f.Variables(fac, formula).Content())

	for _, solver := range solvers(fac) {
		minimumModel := optimize([]f.Formula{formula}, vars, []f.Variable{}, false, solver)
		testMinimumModel(t, fac, formula, minimumModel, vars)
		maximumModel := optimize([]f.Formula{formula}, vars, []f.Variable{}, true, solver)
		testMaximumModel(t, fac, formula, maximumModel, vars)
	}
}

func TestOptimizerFunctionEXOModel(t *testing.T) {
	fac := f.NewFactory()
	parser := parser.New(fac)
	formula := parser.ParseUnsafe("a + b + c = 1")
	vars := f.VariablesAsLiterals(f.Variables(fac, formula).Content())

	for _, solver := range solvers(fac) {
		minimumModel := optimize([]f.Formula{formula}, vars, []f.Variable{}, false, solver)
		testMinimumModel(t, fac, formula, minimumModel, vars)
		maximumModel := optimize([]f.Formula{formula}, vars, []f.Variable{}, true, solver)
		testMaximumModel(t, fac, formula, maximumModel, vars)
	}
}

func TestOptimizerFunctionCornerCases(t *testing.T) {
	fac := f.NewFactory()
	for _, formula := range f.NewCornerCases(fac) {
		vars := f.Variables(fac, formula).Content()
		targetLits := f.VariablesAsLiterals(vars)
		solver := NewSolver(fac)
		solver.Add(formula)
		minModel := optimize([]f.Formula{formula}, targetLits, []f.Variable{}, false, solver)
		testMinimumModel(t, fac, formula, minModel, targetLits)
		maxModel := optimize([]f.Formula{formula}, targetLits, []f.Variable{}, true, solver)
		testMaximumModel(t, fac, formula, maxModel, targetLits)
	}
}

func TestOptimizerFunctionRandomSmall(t *testing.T) {
	fac := f.NewFactory()
	config := randomizer.DefaultConfig()
	config.NumVars = 6
	config.WeightPBC = 2
	config.Seed = 42
	randomizer := randomizer.New(fac, config)

	for i := 0; i < 1000; i++ {
		formula := randomizer.Formula(2)
		variables := f.Variables(fac, formula).Content()
		literals := f.VariablesAsLiterals(f.Variables(fac, formula).Content())
		targetLiterals := randomTargetLiterals(fac, randomSubset(literals, min(len(literals), 5)))
		additionalVariables := randomSubset(variables, min(len(variables), 3))

		for _, solver := range solvers(fac) {
			minimumModel := optimize([]f.Formula{formula}, targetLiterals, additionalVariables, false, solver)
			testMinimumModel(t, fac, formula, minimumModel, targetLiterals)
			maximumModel := optimize([]f.Formula{formula}, targetLiterals, additionalVariables, true, solver)
			testMaximumModel(t, fac, formula, maximumModel, targetLiterals)
		}
	}
}

func TestOptimizerFunctionIncMinMax(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	parser := parser.New(fac)
	for _, solver := range solvers(fac) {
		formula := parser.ParseUnsafe("(a|b|c|d|e) & (p|q) & (x|y|z)")
		variables := f.NewMutableVarSetCopy(f.Variables(fac, formula))
		vars := f.VariablesAsLiterals(variables.Content())
		solver.Add(formula)

		minimumModel := solver.Minimize(vars)
		maximumModel := solver.Maximize(vars)
		assert.Equal(3, len(minimumModel.PosVars()))
		assert.Equal(10, len(maximumModel.PosVars()))

		formula = parser.ParseUnsafe("~p")
		solver.Add(formula)
		minimumModel = solver.Minimize(vars)
		maximumModel = solver.Maximize(vars)
		assert.Equal(3, len(minimumModel.PosVars()))
		assert.Equal(9, len(maximumModel.PosVars()))

		formula = parser.ParseUnsafe("(x => n) & (y => m) & (a => ~b & ~c)")
		variables.AddAll(f.Variables(fac, formula))
		vars = f.VariablesAsLiterals(variables.Content())
		solver.Add(formula)
		minimumModel = solver.Minimize(vars)
		maximumModel = solver.Maximize(vars)
		assert.Equal(3, len(minimumModel.PosVars()))
		assert.True(slices.Contains(minimumModel.PosVars(), fac.Var("q")))
		assert.True(slices.Contains(minimumModel.PosVars(), fac.Var("z")))
		assert.Equal(10, len(maximumModel.PosVars()))
		assert.True(slices.Contains(maximumModel.PosVars(), fac.Var("z")))
		assert.False(slices.Contains(maximumModel.PosVars(), fac.Var("a")))

		formula = parser.ParseUnsafe("(z => v & w) & (m => v) & (b => ~c & ~d & ~e)")
		variables.AddAll(f.Variables(fac, formula))
		vars = f.VariablesAsLiterals(variables.Content())
		solver.Add(formula)
		minimumModel = solver.Minimize(vars)
		maximumModel = solver.Maximize(vars)
		assert.Equal(4, len(minimumModel.PosVars()))
		assert.True(slices.Contains(minimumModel.PosVars(), fac.Var("q")))
		assert.True(slices.Contains(minimumModel.PosVars(), fac.Var("x")))
		assert.True(slices.Contains(minimumModel.PosVars(), fac.Var("n")))
		assert.Equal(11, len(maximumModel.PosVars()))
		assert.True(slices.Contains(maximumModel.PosVars(), fac.Var("q")))
		assert.True(slices.Contains(maximumModel.PosVars(), fac.Var("x")))
		assert.True(slices.Contains(maximumModel.PosVars(), fac.Var("n")))
		assert.True(slices.Contains(maximumModel.PosVars(), fac.Var("v")))
		assert.True(slices.Contains(maximumModel.PosVars(), fac.Var("w")))
		assert.False(slices.Contains(maximumModel.PosVars(), fac.Var("b")))

		formula = parser.ParseUnsafe("~q")
		solver.Add(formula)
		minimumModel = solver.Minimize(vars)
		maximumModel = solver.Maximize(vars)
		assert.Nil(minimumModel)
		assert.Nil(maximumModel)
	}
}

func TestOptimizerAdditionalVariables(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	parser := parser.New(fac)
	for _, solver := range solvers(fac) {
		va := fac.Var("a")
		vc := fac.Var("c")
		vy := fac.Var("y")
		a := fac.Lit("a", true)
		b := fac.Lit("b", true)
		c := fac.Lit("c", true)
		x := fac.Lit("x", true)
		y := fac.Lit("y", true)
		na := fac.Lit("a", false)
		nb := fac.Lit("b", false)
		nx := fac.Lit("x", false)

		formula := parser.ParseUnsafe("(a|b) & (~a => c) & (x|y)")

		literalsANBX := []f.Literal{a, nb, x}
		minimumModel := optimize([]f.Formula{formula}, literalsANBX, []f.Variable{}, false, solver)
		assert.True(slices.Contains(minimumModel.Literals, na))
		assert.True(slices.Contains(minimumModel.Literals, b))
		assert.True(slices.Contains(minimumModel.Literals, nx))

		minimumModelWithY := optimize([]f.Formula{formula}, literalsANBX, []f.Variable{vy}, false, solver)
		assert.True(slices.Contains(minimumModelWithY.Literals, na))
		assert.True(slices.Contains(minimumModelWithY.Literals, b))
		assert.True(slices.Contains(minimumModelWithY.Literals, nx))
		assert.True(slices.Contains(minimumModelWithY.Literals, y))

		minimumModelWithCY := optimize([]f.Formula{formula}, literalsANBX, []f.Variable{vc, vy}, false, solver)
		assert.True(slices.Contains(minimumModelWithCY.Literals, na))
		assert.True(slices.Contains(minimumModelWithCY.Literals, b))
		assert.True(slices.Contains(minimumModelWithCY.Literals, nx))
		assert.True(slices.Contains(minimumModelWithCY.Literals, y))
		assert.True(slices.Contains(minimumModelWithCY.Literals, c))

		literalsNBNX := []f.Literal{na, nx}
		maximumModel := optimize([]f.Formula{formula}, literalsNBNX, []f.Variable{}, true, solver)
		assert.True(slices.Contains(maximumModel.Literals, na))
		assert.True(slices.Contains(maximumModel.Literals, nx))
		maximumModelWithC := optimize([]f.Formula{formula}, literalsNBNX, []f.Variable{vc}, true, solver)
		assert.True(slices.Contains(maximumModelWithC.Literals, na))
		assert.True(slices.Contains(maximumModelWithC.Literals, nx))
		assert.True(slices.Contains(maximumModelWithC.Literals, c))
		maximumModelWithACY := optimize([]f.Formula{formula}, literalsNBNX, []f.Variable{va, vc, vy}, true, solver)
		assert.True(slices.Contains(maximumModelWithACY.Literals, na))
		assert.True(slices.Contains(maximumModelWithACY.Literals, c))
		assert.True(slices.Contains(maximumModelWithACY.Literals, nx))
		assert.True(slices.Contains(maximumModelWithACY.Literals, y))
	}
}

func TestOptimizerFunctionLargeFormulaMinimize(t *testing.T) {
	fac := f.NewFactory()
	formula, _ := io.ReadFormula(fac, "../test/data/formulas/large2.txt")
	vars := f.VariablesAsLiterals(f.Variables(fac, formula).Content())
	for _, solver := range solvers(fac) {
		minimumModel := optimize([]f.Formula{formula}, vars, []f.Variable{}, false, solver)
		assert.Equal(t, 25, len(minimumModel.PosVars()))
		testMinimumModel(t, fac, formula, minimumModel, vars)
	}
}

func TestOptimizerFunctionLargeFormulaMaximize(t *testing.T) {
	fac := f.NewFactory()
	formula, _ := io.ReadFormula(fac, "../test/data/formulas/large2.txt")
	vars := f.VariablesAsLiterals(f.Variables(fac, formula).Content())
	for _, solver := range solvers(fac) {
		minimumModel := optimize([]f.Formula{formula}, vars, []f.Variable{}, true, solver)
		assert.Equal(t, 162, len(minimumModel.PosVars()))
		testMaximumModel(t, fac, formula, minimumModel, vars)
	}
}

func TestOptimizerFunctionLargerFormulaMinimize(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	fac := f.NewFactory()
	formula, _ := io.ReadFormula(fac, "../test/data/formulas/small_formulas.txt")
	vars := f.VariablesAsLiterals(f.Variables(fac, formula).Content())
	for _, solver := range solvers(fac) {
		minimumModel := optimize([]f.Formula{formula}, vars, []f.Variable{}, false, solver)
		assert.Equal(t, 50, len(minimumModel.PosVars()))
		testMinimumModel(t, fac, formula, minimumModel, vars)
	}
}

func TestOptimizerFunctionLargerFormulaMaximize(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	fac := f.NewFactory()
	formula, _ := io.ReadFormula(fac, "../test/data/formulas/small_formulas.txt")
	vars := f.VariablesAsLiterals(f.Variables(fac, formula).Content())
	for _, solver := range solvers(fac) {
		minimumModel := optimize([]f.Formula{formula}, vars, []f.Variable{}, true, solver)
		assert.Equal(t, 270, len(minimumModel.PosVars()))
		testMaximumModel(t, fac, formula, minimumModel, vars)
	}
}

func optimize(
	formulas []f.Formula,
	literals []f.Literal,
	additionalVariables []f.Variable,
	maximize bool,
	solver *Solver,
) *model.Model {
	solver.Reset()
	solver.Add(formulas...)
	if maximize {
		return solver.Maximize(literals, additionalVariables...)
	} else {
		return solver.Minimize(literals, additionalVariables...)
	}
}

func testMinimumModel(
	t *testing.T,
	fac f.Factory,
	formula f.Formula,
	mdl *model.Model,
	literals []f.Literal,
) {
	testOptimumModel(t, fac, formula, mdl, literals, false)
}

func testMaximumModel(
	t *testing.T,
	fac f.Factory,
	formula f.Formula,
	mdl *model.Model,
	literals []f.Literal,
) {
	testOptimumModel(t, fac, formula, mdl, literals, true)
}

func testOptimumModel(
	t *testing.T,
	fac f.Factory,
	formula f.Formula,
	mdl *model.Model,
	literals []f.Literal,
	maximize bool,
) {
	assert := assert.New(t)
	if IsSatisfiable(fac, formula) {
		assert.True(IsSatisfiable(fac, fac.And(formula, fac.Minterm(mdl.Literals...))))
		numSatisfiedLiterals := len(satisfiedLiterals(mdl, literals))
		var selVars []f.Variable
		solver := NewSolver(fac)
		solver.Add(formula)
		for _, lit := range literals {
			selVar := fac.Variable(fmt.Sprintf("SEL_VAR_%d", len(selVars)))
			if maximize {
				solver.Add(fac.Equivalence(selVar.Negate(fac), lit.AsFormula()))
			} else {
				solver.Add(fac.Equivalence(selVar.Negate(fac), lit.Negate(fac).AsFormula()))
			}
		}
		solver.Add(fac.CC(f.GT, uint32(numSatisfiedLiterals+1), selVars...))
		assert.False(solver.Sat())
	} else {
		assert.Nil(mdl)
	}
}

func satisfiedLiterals(mdl *model.Model, literals []f.Literal) []f.Literal {
	modelLiterals := f.NewLitSet(mdl.Literals...)
	result := make([]f.Literal, 0, modelLiterals.Size())
	for _, lit := range literals {
		if modelLiterals.Contains(lit) {
			result = append(result, lit)
		}
	}
	return result
}

func randomTargetLiterals(fac f.Factory, literals []f.Literal) []f.Literal {
	result := make([]f.Literal, len(literals))
	for i, l := range literals {
		name, _, _ := fac.LitNamePhase(l)
		result[i] = fac.Lit(name, rand.Intn(2) == 0)
	}
	return result
}

func randomSubset[T any](elements []T, subsetSize int) []T {
	if subsetSize > len(elements) {
		panic(errorx.IllegalState("not good"))
	}
	rand.Shuffle(len(elements), func(i, j int) {
		elements[i], elements[j] = elements[j], elements[i]
	})
	return elements[:subsetSize]
}
