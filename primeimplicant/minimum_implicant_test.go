package primeimplicant

import (
	"testing"

	"booleworks.com/logicng/encoding"
	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/io"
	"booleworks.com/logicng/parser"
	"booleworks.com/logicng/sat"
	"github.com/stretchr/testify/assert"
)

func TestMinimumPrimeImplicantDoc(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	f1 := p.ParseUnsafe("(A | B) & (A | C ) & (C | D) & (B | ~D)")
	implicant, err := Minimum(fac, f1)
	assert.Nil(t, err)
	assert.Equal(t, []f.Literal{fac.Lit("B", true), fac.Lit("C", true)}, implicant)
}

func TestMinimumPrimeImplicantSimpleCases(t *testing.T) {
	assert := assert.New(t)
	fac := fac()
	p := parser.New(fac)
	formula := p.ParseUnsafe("a")
	pi, err := Minimum(fac, formula)
	assert.Nil(err)
	assert.Equal(1, len(pi))
	isPrimeImplicant(t, fac, formula, pi)

	formula = p.ParseUnsafe("a | b | c")
	pi, err = Minimum(fac, formula)
	assert.Nil(err)
	assert.Equal(1, len(pi))
	isPrimeImplicant(t, fac, formula, pi)

	formula = p.ParseUnsafe("a & b & (~a|~b)")
	pi, err = Minimum(fac, formula)
	assert.NotNil(err)
	assert.Nil(pi)

	formula = p.ParseUnsafe("a & b & c")
	pi, err = Minimum(fac, formula)
	assert.Nil(err)
	assert.Equal(3, len(pi))
	isPrimeImplicant(t, fac, formula, pi)

	formula = p.ParseUnsafe("a | b | ~c => e & d & f")
	pi, err = Minimum(fac, formula)
	assert.Nil(err)
	assert.Equal(3, len(pi))
	isPrimeImplicant(t, fac, formula, pi)

	formula = p.ParseUnsafe("a | b | ~c <=> e & d & f")
	pi, err = Minimum(fac, formula)
	assert.Nil(err)
	assert.Equal(4, len(pi))
	isPrimeImplicant(t, fac, formula, pi)

	formula = p.ParseUnsafe("(a | b | ~c <=> e & d & f) | (a | b | ~c => e & d & f)")
	pi, err = Minimum(fac, formula)
	assert.Nil(err)
	assert.Equal(3, len(pi))
	isPrimeImplicant(t, fac, formula, pi)

	formula = p.ParseUnsafe("(a | b | ~c <=> e & d & f) | (a | b | ~c => e & d & f) | (a & b)")
	pi, err = Minimum(fac, formula)
	assert.Nil(err)
	assert.Equal(2, len(pi))
	isPrimeImplicant(t, fac, formula, pi)

	formula = p.ParseUnsafe("(a | b | ~c <=> e & d & f) | (a | b | ~c => e & d & f) | (a & b) | (f => g)")
	pi, err = Minimum(fac, formula)
	assert.Nil(err)
	assert.Equal(1, len(pi))
	isPrimeImplicant(t, fac, formula, pi)
}

func TestMinimumPrimeImplicantSmallExamples(t *testing.T) {
	assert := assert.New(t)
	fac := fac()
	p := parser.New(fac)
	formula := p.ParseUnsafe("(~(v17 | v18) | ~v1494 & (v17 | v18)) & ~v687 => v686")
	pi, err := Minimum(fac, formula)
	assert.Nil(err)
	assert.Equal(1, len(pi))
	isPrimeImplicant(t, fac, formula, pi)

	formula = p.ParseUnsafe("(~(v17 | v18) | ~v1494 & (v17 | v18)) & v687 => ~v686")
	pi, err = Minimum(fac, formula)
	assert.Nil(err)
	assert.Equal(1, len(pi))
	isPrimeImplicant(t, fac, formula, pi)

	formula = p.ParseUnsafe(
		"v173 + v174 + v451 + v258 + v317 + v259 + v452 + v453 + v175 + v176 + v177 + v178 + v179 + v180 + v181 + v182 + v183 + v102 + v103 + v104 + v105 = 1")
	pi, err = Minimum(fac, formula)
	assert.Nil(err)
	assert.Equal(21, len(pi))
	isPrimeImplicant(t, fac, formula, pi)
}

func TestMinimumPrimeImplicantMiddleExamples(t *testing.T) {
	fac := fac()
	formulas, _ := io.ReadFormulas(fac, "../test/data/formulas/small.txt")
	for _, formula := range formulas {
		min, _ := Minimum(fac, formula)
		isPrimeImplicant(t, fac, formula, min)
	}
}

func TestMinimumPrimeImplicantLargeExamples(t *testing.T) {
	fac := fac()
	formulas, _ := io.ReadFormulas(fac, "../test/data/formulas/small_formulas.txt")
	for _, formula := range formulas {
		min, _ := Minimum(fac, formula)
		isPrimeImplicant(t, fac, formula, min)
	}
}

func isPrimeImplicant(t *testing.T, fac f.Factory, formula f.Formula, pi []f.Literal) {
	assert.True(t, sat.Implies(fac, fac.And(f.LiteralsAsFormulas(pi)...), formula))
	for _, literal := range pi {
		newSet := f.NewLitSet(pi...)
		newSet.Remove(literal)
		if !newSet.Empty() {
			assert.False(t, sat.Implies(fac, fac.And(f.LiteralsAsFormulas(newSet.Content())...), formula))
		}
	}
}

func fac() f.Factory {
	encodingConfig := encoding.DefaultConfig()
	encodingConfig.AMOEncoder = encoding.AMOPure
	fac := f.NewFactory()
	fac.PutConfiguration(encodingConfig)
	return fac
}
