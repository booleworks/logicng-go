package simplification

import (
	"testing"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/normalform"
	"github.com/booleworks/logicng-go/parser"
	"github.com/booleworks/logicng-go/randomizer"
	"github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

func TestFactorOutSimple(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	assert.Equal(p.ParseUnsafe("$false"), FactorOut(fac, p.ParseUnsafe("$false")))
	assert.Equal(p.ParseUnsafe("$true"), FactorOut(fac, p.ParseUnsafe("$true")))
	assert.Equal(p.ParseUnsafe("a"), FactorOut(fac, p.ParseUnsafe("a")))
	assert.Equal(p.ParseUnsafe("~a"), FactorOut(fac, p.ParseUnsafe("~a")))

	assert.Equal(p.ParseUnsafe("A&~B&~C&~D"), FactorOut(fac, p.ParseUnsafe("A&~B&~C&~D")))
	assert.Equal(p.ParseUnsafe("~A&~B&~C&~D"), FactorOut(fac, p.ParseUnsafe("~A&~B&~C&~D")))

	assert.Equal(p.ParseUnsafe("A"), FactorOut(fac, p.ParseUnsafe("A|A&B")))
	assert.Equal(p.ParseUnsafe("C&D|A"), FactorOut(fac, p.ParseUnsafe("A|A&B|C&D")))
	assert.Equal(p.ParseUnsafe("~A"), FactorOut(fac, p.ParseUnsafe("~(A&(A|B))")))
	assert.Equal(p.ParseUnsafe("C|A"), FactorOut(fac, p.ParseUnsafe("A|A&B|C")))

	assert.Equal(p.ParseUnsafe("A"), FactorOut(fac, p.ParseUnsafe("A&(A|B)")))
	assert.Equal(p.ParseUnsafe("(C|D)&A"), FactorOut(fac, p.ParseUnsafe("A&(A|B)&(C|D)")))
	assert.Equal(p.ParseUnsafe("~A"), FactorOut(fac, p.ParseUnsafe("~(A|A&B)")))
	assert.Equal(p.ParseUnsafe("C&A"), FactorOut(fac, p.ParseUnsafe("A&(A|B)&C")))

	assert.Equal(
		p.ParseUnsafe("B&C&D|A&(X&Y|B&C|Z)"),
		FactorOut(fac, p.ParseUnsafe("A&X&Y|A&B&C|B&C&D|A&Z")),
	)
	assert.Equal(
		p.ParseUnsafe("G&(B&C&D|A&(X&Y|B&C|Z))"),
		FactorOut(fac, p.ParseUnsafe("G&(A&X&Y|A&B&C|B&C&D|A&Z)")),
	)

	assert.Equal(
		p.ParseUnsafe("G&(~(A&X&Y)|~(A&B&C))"),
		FactorOut(fac, p.ParseUnsafe("G&(~(A&X&Y)|~(A&B&C))")),
	)
}

func TestFactorOutSimplifierCornerCases(t *testing.T) {
	fac := f.NewFactory()
	for _, formula := range f.NewCornerCases(fac) {
		computeAndVerify(t, fac, formula)
	}
}

func TestFactorOutRandomized(t *testing.T) {
	fac := f.NewFactory()
	for i := 0; i < 100; i++ {
		config := randomizer.DefaultConfig()
		config.NumVars = 5
		config.WeightPBC = 2
		config.Seed = int64(i * 42)
		formula := randomizer.New(fac, config).Formula(6)
		computeAndVerify(t, fac, formula)
		computeAndVerify(t, fac, normalform.NNF(fac, formula))
	}
}

func computeAndVerify(t *testing.T, fac f.Factory, formula f.Formula) {
	simplified := FactorOut(fac, formula)
	assert.True(t, sat.IsEquivalent(fac, formula, simplified))
	assert.LessOrEqual(t, len(simplified.Sprint(fac)), len(formula.Sprint(fac)))
}
