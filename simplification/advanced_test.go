package simplification

import (
	"testing"
	"time"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/io"
	"github.com/booleworks/logicng-go/parser"
	"github.com/booleworks/logicng-go/randomizer"
	"github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

func TestAdvancedSimplifierSimple(t *testing.T) {
	fac := f.NewFactory()
	computeAndVerifyAS(t, fac, fac.Falsum())
	computeAndVerifyAS(t, fac, fac.Verum())
}

func TestAdvancedSimplifierCornerCases(t *testing.T) {
	fac := f.NewFactory()
	for _, formula := range f.NewCornerCases(fac) {
		computeAndVerifyAS(t, fac, formula)
	}
}

func TestAdvancedSimplifierRandomized(t *testing.T) {
	fac := f.NewFactory()
	for i := 0; i < 100; i++ {
		config := randomizer.DefaultConfig()
		config.NumVars = 8
		config.WeightPBC = 2
		config.Seed = int64(i * 42)
		formula := randomizer.New(fac, config).Formula(5)
		computeAndVerify(t, fac, formula)
	}
}

func TestAdvancedSimplifierTimeoutHandlerSmall(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	duration, _ := time.ParseDuration("1s")
	handler := sat.OptimizationHandlerWithTimeout(*handler.NewTimeoutWithDuration(duration))
	formula := p.ParseUnsafe("a & b | ~c & a")
	testHandler(t, fac, handler, formula, false)
}

func TestAdvancedSimplifierTimeoutHandlerLarge(t *testing.T) {
	fac := f.NewFactory()
	duration, _ := time.ParseDuration("500ms")
	handler := sat.OptimizationHandlerWithTimeout(*handler.NewTimeoutWithDuration(duration))
	formula, _ := io.ReadFormula(fac, "../test/data/formulas/large.txt")
	testHandler(t, fac, handler, formula, true)
}

func TestAdvancedSimplifierConfigs(t *testing.T) {
	configs := []*Config{
		{FactorOut: true, RestrictBackbone: true, SimplifyNegations: true, RatingFunction: DefaultRatingFunction},
		{FactorOut: true, RestrictBackbone: true, SimplifyNegations: false, RatingFunction: DefaultRatingFunction},
		{FactorOut: true, RestrictBackbone: false, SimplifyNegations: true, RatingFunction: DefaultRatingFunction},
		{FactorOut: true, RestrictBackbone: false, SimplifyNegations: false, RatingFunction: DefaultRatingFunction},
		{FactorOut: false, RestrictBackbone: true, SimplifyNegations: true, RatingFunction: DefaultRatingFunction},
		{FactorOut: false, RestrictBackbone: true, SimplifyNegations: false, RatingFunction: DefaultRatingFunction},
		{FactorOut: false, RestrictBackbone: false, SimplifyNegations: true, RatingFunction: DefaultRatingFunction},
		{FactorOut: false, RestrictBackbone: false, SimplifyNegations: false, RatingFunction: DefaultRatingFunction},
	}

	fac := f.NewFactory()
	for _, config := range configs {
		for i := 0; i < 10; i++ {
			formula := randomizer.NewWithSeed(fac, int64(i)).Formula(3)
			simplified := Advanced(fac, formula, config)
			assert.True(t, sat.IsEquivalent(fac, formula, simplified))
		}
	}
}

func computeAndVerifyAS(t *testing.T, fac f.Factory, formula f.Formula) {
	simplified := Advanced(fac, formula)
	assert.True(t, sat.IsEquivalent(fac, formula, simplified))
}

func testHandler(t *testing.T, fac f.Factory, handler sat.OptimizationHandler, formula f.Formula, exp bool) {
	assert := assert.New(t)
	result, ok := AdvancedWithHandler(fac, formula, handler)
	if exp {
		assert.True(handler.Aborted())
		assert.False(ok)
		assert.Equal(fac.Falsum(), result)
	} else {
		assert.False(handler.Aborted())
		assert.True(ok)
		assert.NotEqual(fac.Falsum(), result)
	}
}
