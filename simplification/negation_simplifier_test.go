package simplification

import (
	"testing"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/parser"
	"github.com/booleworks/logicng-go/randomizer"
	"github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

func TestNegationSimplifierSimple(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	assert.Equal(p.ParseUnsafe("$false"), SimplifyNegations(fac, p.ParseUnsafe("$false")))
	assert.Equal(p.ParseUnsafe("$true"), SimplifyNegations(fac, p.ParseUnsafe("$true")))
	assert.Equal(p.ParseUnsafe("a"), SimplifyNegations(fac, p.ParseUnsafe("a")))
	assert.Equal(p.ParseUnsafe("~a"), SimplifyNegations(fac, p.ParseUnsafe("~a")))

	assert.Equal(p.ParseUnsafe("A&~B&~C&~D"), SimplifyNegations(fac, p.ParseUnsafe("A&~B&~C&~D")))
	assert.Equal(p.ParseUnsafe("~(A|B|C|D)"), SimplifyNegations(fac, p.ParseUnsafe("~A&~B&~C&~D")))
	assert.Equal(p.ParseUnsafe("~(A|B|C|D)"), SimplifyNegations(fac, p.ParseUnsafe("~A&~B&~C&~D")))
	assert.Equal(p.ParseUnsafe("~(A&B&C&D)"), SimplifyNegations(fac, p.ParseUnsafe("~A|~B|~C|~D")))
	assert.Equal(p.ParseUnsafe("D|~(A&B&C&E&G)"), SimplifyNegations(fac, p.ParseUnsafe("~A|~B|~C|D|~E|~G")))
	assert.Equal(p.ParseUnsafe("D&~(A|B|C|E|G)"), SimplifyNegations(fac, p.ParseUnsafe("~A&~B&~C&D&~E&~G")))

	assert.Equal(p.ParseUnsafe("~E&G|~(A&B&(H|B|C)&X)"), SimplifyNegations(fac, p.ParseUnsafe("~A|~B|~E&G|~H&~B&~C|~X")))
	assert.Equal(p.ParseUnsafe("~E&G|~(A&B&(H|B|C)&X)"), SimplifyNegations(fac, p.ParseUnsafe("~(A&B&~(~E&G)&(H|B|C)&X)")))
	assert.Equal(p.ParseUnsafe("~A|B|~(E|G|H|K)"), SimplifyNegations(fac, p.ParseUnsafe("~A|B|(~E&~G&~H&~K)")))

	assert.Equal(p.ParseUnsafe("~A|~B"), SimplifyNegations(fac, p.ParseUnsafe("~A|~B")))
	assert.Equal(p.ParseUnsafe("~A|~B|~C"), SimplifyNegations(fac, p.ParseUnsafe("~A|~B|~C")))
	assert.Equal(p.ParseUnsafe("~(A&B&C&D)"), SimplifyNegations(fac, p.ParseUnsafe("~A|~B|~C|~D")))

	assert.Equal(p.ParseUnsafe("X&~(A&B)"), SimplifyNegations(fac, p.ParseUnsafe("X&(~A|~B)")))
	assert.Equal(p.ParseUnsafe("X&~(A&B&C)"), SimplifyNegations(fac, p.ParseUnsafe("X&(~A|~B|~C)")))
	assert.Equal(p.ParseUnsafe("X&~(A&B&C&D)"), SimplifyNegations(fac, p.ParseUnsafe("X&(~A|~B|~C|~D)")))

	assert.Equal(p.ParseUnsafe("~A&~B"), SimplifyNegations(fac, p.ParseUnsafe("~A&~B")))
	assert.Equal(p.ParseUnsafe("~A&~B&~C"), SimplifyNegations(fac, p.ParseUnsafe("~A&~B&~C")))
	assert.Equal(p.ParseUnsafe("~(A|B|C|D)"), SimplifyNegations(fac, p.ParseUnsafe("~A&~B&~C&~D")))

	assert.Equal(p.ParseUnsafe("X|~A&~B"), SimplifyNegations(fac, p.ParseUnsafe("X|~A&~B")))
	assert.Equal(p.ParseUnsafe("X|~A&~B&~C"), SimplifyNegations(fac, p.ParseUnsafe("X|~A&~B&~C")))
	assert.Equal(p.ParseUnsafe("X|~(A|B|C|D)"), SimplifyNegations(fac, p.ParseUnsafe("X|~A&~B&~C&~D")))

	assert.Equal(p.ParseUnsafe("A&(X|Y|H|~(B&C&D&E&G))"), SimplifyNegations(fac, p.ParseUnsafe("A&(~B|~C|~D|~E|~G|X|Y|H)")))
}

func TestNegationSimplifierCornerCases(t *testing.T) {
	fac := f.NewFactory()
	for _, formula := range f.NewCornerCases(fac) {
		computeAndVerify(t, fac, formula)
	}
}

func TestNegationSimplifierRandomized(t *testing.T) {
	fac := f.NewFactory()
	for i := 0; i < 100; i++ {
		config := randomizer.DefaultConfig()
		config.NumVars = 5
		config.WeightPBC = 1
		config.Seed = int64(i * 42)
		formula := randomizer.New(fac, config).Formula(6)
		computeAndVerifyNN(t, fac, formula)
	}
}

func computeAndVerifyNN(t *testing.T, fac f.Factory, formula f.Formula) {
	simplified := SimplifyNegations(fac, formula)
	assert.True(t, sat.IsEquivalent(fac, formula, simplified))
	assert.LessOrEqual(t, len(simplified.Sprint(fac)), len(formula.Sprint(fac)))
}
