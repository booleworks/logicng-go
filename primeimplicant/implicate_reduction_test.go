package primeimplicant

import (
	"slices"
	"testing"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/io"
	"github.com/booleworks/logicng-go/parser"
	"github.com/booleworks/logicng-go/randomizer"
	"github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

func TestPrimeImplicateReductionSimple(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	parser := parser.New(fac)

	pr := newPrimeReduction(fac, parser.ParseUnsafe("a&b"))
	result, _ := pr.reduceImplicate(fac, []f.Literal{d.LA, d.LB}, nil)
	assert.Equal(1, len(result))

	pr = newPrimeReduction(fac, parser.ParseUnsafe("(a => b) | b | c"))
	result, _ = pr.reduceImplicate(fac, []f.Literal{d.LNA, d.LB, d.LC}, nil)
	assert.Equal(3, len(result))
	assert.True(slices.Contains(result, d.LNA))
	assert.True(slices.Contains(result, d.LB))
	assert.True(slices.Contains(result, d.LC))

	pr = newPrimeReduction(fac, parser.ParseUnsafe("(a => b) & b & c"))
	result, _ = pr.reduceImplicate(fac, []f.Literal{d.LB, d.LC}, nil)
	assert.Equal(1, len(result))
}

func TestPrimeImplicateSmallFormula(t *testing.T) {
	fac := f.NewFactory()
	formula, _ := io.ReadFormula(fac, "../test/data/formulas/small.txt")
	testImplicateFormula(t, fac, formula)
}

func TestPrimeImplicateSimplifyFormulas(t *testing.T) {
	fac := f.NewFactory()
	formula, _ := io.ReadFormula(fac, "../test/data/formulas/simplify_formulas.txt")
	testImplicateFormula(t, fac, formula)
}

func TestPrimeImplicateLargeFormulas(t *testing.T) {
	fac := f.NewFactory()
	formula, _ := io.ReadFormula(fac, "../test/data/formulas/large2.txt")
	testImplicateFormula(t, fac, formula)
}

func TestPrimeImplicateRandom(t *testing.T) {
	fac := f.NewFactory()
	for i := 0; i < 500; i++ {
		config := randomizer.DefaultConfig()
		config.NumVars = 20
		config.WeightPBC = 2
		config.Seed = int64(i * 42)
		formula := randomizer.New(fac, config).Formula(4)
		testImplicateFormula(t, fac, formula)
	}
}

func TestPrimeImplicateCornerCases(t *testing.T) {
	fac := f.NewFactory()
	for _, formula := range f.NewCornerCases(fac) {
		testImplicateFormula(t, fac, formula)
	}
}

func testImplicateFormula(t *testing.T, fac f.Factory, formula f.Formula) {
	testImplicateFormulaDetail(t, fac, formula, nil, false)
}

func testImplicateFormulaDetail(t *testing.T, fac f.Factory, formula f.Formula, handler sat.Handler, expAborted bool) {
	assert := assert.New(t)
	solver := sat.NewSolver(fac)
	solver.Add(formula.Negate(fac))
	result := solver.Call(sat.WithModel(f.Variables(fac, formula).Content()))
	if !result.Sat() {
		return
	}

	model := result.Model()
	falsifyingAssignment := make([]f.Literal, model.Size())
	for i, lit := range model.Literals {
		falsifyingAssignment[i] = lit.Negate(fac)
	}
	pr := newPrimeReduction(fac, formula)
	primeImplicate, ok := pr.reduceImplicate(fac, falsifyingAssignment, handler)
	if expAborted {
		assert.False(ok)
		assert.True(handler.Aborted())
		assert.Nil(primeImplicate)
	} else {
		assert.True(ok)
		for _, lit := range primeImplicate {
			assert.True(slices.Contains(falsifyingAssignment, lit))
		}
		testPrimeImplicateProperty(t, fac, formula, primeImplicate)
	}
}

func testPrimeImplicateProperty(t *testing.T, fac f.Factory, formula f.Formula, primeImplicate []f.Literal) {
	assert := assert.New(t)
	solver := sat.NewSolver(fac)
	solver.Add(formula)
	negatedLiterals := make([]f.Literal, len(primeImplicate))
	for i, lit := range primeImplicate {
		negatedLiterals[i] = lit.Negate(fac)
	}
	assert.False(solver.Call(sat.WithAssumptions(negatedLiterals)).Sat())
	for _, lit := range negatedLiterals {
		reducedNegatedLiterals := f.NewMutableLitSet(negatedLiterals...)
		reducedNegatedLiterals.Remove(lit)
		assert.True(solver.Call(sat.WithAssumptions(reducedNegatedLiterals.Content())).Sat())
	}
}
