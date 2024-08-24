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
	handler := handler.NewTimeoutWithDuration(duration)
	formula := p.ParseUnsafe("a & b | ~c & a")
	testHandler(t, fac, handler, formula, false)
}

func TestAdvancedSimplifierTimeoutHandlerLarge(t *testing.T) {
	fac := f.NewFactory()
	duration, _ := time.ParseDuration("500ms")
	handler := handler.NewTimeoutWithDuration(duration)
	formula, _ := io.ReadFormula(fac, "../test/data/formulas/large.txt")
	testHandler(t, fac, handler, formula, true)
}

func TestAdvancedDeterministic(t *testing.T) {
	for i := 0; i < 100; i++ {
		fac := f.NewFactory()
		p := parser.New(fac)
		input := p.ParseUnsafe("A&B&F&D&E&~H&~C&~G&~I&~J|A&H&D&C&G&~B&~F&~E&~I&~J|A&D&C&~B&~F&~H&~E&~G&~I&~J|" +
			"A&F&H&D&C&G&~B&~E&~I&~J|A&B&H&~F&~D&~C&~E&~G&~I&~J|A&B&H&D&G&~F&~C&~E&~I&~J|A&H&C&~B&~F&~D&~E&~G&~I&~J|" +
			"A&B&G&~F&~H&~D&~C&~E&~I&~J|A&H&C&E&G&~B&~F&~D&~I&~J|A&C&G&~B&~F&~H&~D&~E&~I&~J|A&B&H&G&~F&~D&~C&~E&~I&~J|" +
			"A&C&~B&~F&~H&~D&~E&~G&~I&~J|A&D&C&G&~B&~F&~H&~E&~I&~J|A&B&D&G&~F&~H&~C&~E&~I&~J|A&H&D&C&~B&~F&~E&~G&~I&~J|" +
			"A&H&D&C&E&G&~B&~F&~I&~J|A&B&D&E&~F&~H&~C&~G&~I&~J|A&C&E&G&~B&~F&~H&~D&~I&~J")
		simpl := Advanced(fac, input)
		assert.Equal(t, "A & ~I & ~J & (~F & (~B & C & (~(E | G & H) | G & (~D & E | D & H)) | "+
			"B & ~E & ~C & (G | ~D & H)) | D & (~B & ~E & H & C & G | B & E & ~H & ~C & ~G))", simpl.Sprint(fac))

	}
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

func testHandler(t *testing.T, fac f.Factory, h handler.Handler, formula f.Formula, exp bool) {
	assert := assert.New(t)
	result, state := AdvancedWithHandler(fac, formula, h)
	if exp {
		assert.False(state.Success)
		assert.Equal(fac.Falsum(), result)
	} else {
		assert.True(state.Success)
		assert.NotEqual(fac.Falsum(), result)
	}
}
