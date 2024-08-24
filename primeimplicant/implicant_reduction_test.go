package primeimplicant

import (
	"slices"
	"testing"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/io"
	"github.com/booleworks/logicng-go/parser"
	"github.com/booleworks/logicng-go/randomizer"
	"github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

func TestPrimeImplicantReductionSimple(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	parser := parser.New(fac)

	pr := newPrimeReduction(fac, d.True)
	result, _ := pr.reduceImplicant([]f.Literal{d.LA, d.LB}, handler.NopHandler)
	assert.Equal(0, len(result))

	pr = newPrimeReduction(fac, parser.ParseUnsafe("a&b|c&d"))
	result, _ = pr.reduceImplicant([]f.Literal{d.LA, d.LB, d.LC, d.LD.Negate(fac)}, handler.NopHandler)
	assert.Equal(2, len(result))
	assert.True(slices.Contains(result, d.LA))
	assert.True(slices.Contains(result, d.LB))

	result, _ = pr.reduceImplicant([]f.Literal{d.LNA, d.LB, d.LC, d.LD}, handler.NopHandler)
	assert.Equal(2, len(result))
	assert.True(slices.Contains(result, d.LC))
	assert.True(slices.Contains(result, d.LD))

	pr = newPrimeReduction(fac, parser.ParseUnsafe("a|b|~a&~b"))
	result, _ = pr.reduceImplicant([]f.Literal{d.LNA, d.LB}, handler.NopHandler)
	assert.Equal(0, len(result))

	pr = newPrimeReduction(fac, parser.ParseUnsafe("(a => b) | b | c"))
	result, _ = pr.reduceImplicant([]f.Literal{d.LA, d.LB, d.LC.Negate(fac)}, handler.NopHandler)
	assert.Equal(1, len(result))
	assert.True(slices.Contains(result, d.LB))

	result, _ = pr.reduceImplicant([]f.Literal{d.LA, d.LNB, d.LC}, handler.NopHandler)
	assert.Equal(1, len(result))
	assert.True(slices.Contains(result, d.LC))
}

func TestPrimeImplicantSmallFormula(t *testing.T) {
	fac := f.NewFactory()
	formula, _ := io.ReadFormula(fac, "../test/data/formulas/small.txt")
	testImplicantFormula(t, fac, formula)
}

func TestPrimeImplicantSimplifyFormulas(t *testing.T) {
	fac := f.NewFactory()
	formula, _ := io.ReadFormula(fac, "../test/data/formulas/simplify_formulas.txt")
	testImplicantFormula(t, fac, formula)
}

func TestPrimeImplicantLargeFormulas(t *testing.T) {
	fac := f.NewFactory()
	formula, _ := io.ReadFormula(fac, "../test/data/formulas/large2.txt")
	testImplicantFormula(t, fac, formula)
}

func TestPrimeImplicantRandom(t *testing.T) {
	fac := f.NewFactory()
	for i := 0; i < 500; i++ {
		config := randomizer.DefaultConfig()
		config.NumVars = 20
		config.WeightPBC = 2
		config.Seed = int64(i * 42)
		formula := randomizer.New(fac, config).Formula(4)
		testImplicantFormula(t, fac, formula)
	}
}

func TestPrimeImplicantCornerCases(t *testing.T) {
	fac := f.NewFactory()
	for _, formula := range f.NewCornerCases(fac) {
		testImplicantFormula(t, fac, formula)
	}
}

func testImplicantFormula(t *testing.T, fac f.Factory, formula f.Formula) {
	testImplicantFormulaDetail(t, fac, formula, handler.NopHandler, false)
}

func testImplicantFormulaDetail(t *testing.T, fac f.Factory, formula f.Formula, h handler.Handler, expAborted bool) {
	assert := assert.New(t)
	solver := sat.NewSolver(fac)
	solver.Add(formula)
	result := solver.Call(sat.WithModel(f.Variables(fac, formula).Content()))
	if !result.Sat() {
		return
	}

	model := result.Model()
	pr := newPrimeReduction(fac, formula)
	primeImplicant, state := pr.reduceImplicant(model.Literals, h)
	if expAborted {
		assert.False(state.Success)
		assert.Nil(primeImplicant)
	} else {
		assert.True(state.Success)
		for _, lit := range primeImplicant {
			assert.True(slices.Contains(model.Literals, lit))
		}
		testPrimeImplicantProperty(t, fac, formula, primeImplicant)
	}
}

func testPrimeImplicantProperty(t *testing.T, fac f.Factory, formula f.Formula, primeImplicant []f.Literal) {
	assert := assert.New(t)
	solver := sat.NewSolver(fac)
	solver.Add(formula.Negate(fac))
	assert.False(solver.Call(sat.WithAssumptions(primeImplicant)).Sat())
	for _, lit := range primeImplicant {
		reducedPrimeImplicant := f.NewMutableLitSet(primeImplicant...)
		reducedPrimeImplicant.Remove(lit)
		assert.True(solver.Call(sat.WithAssumptions(reducedPrimeImplicant.Content())).Sat())
	}
}
