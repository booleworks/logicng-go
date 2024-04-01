package assignment

import (
	"testing"

	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/parser"
	"github.com/stretchr/testify/assert"
)

func TestRestrictConstants(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	ass, _ := New(fac, f.Literal(d.A), f.Literal(d.NB), f.Literal(d.NX))

	assert.Equal(fac.Verum(), Restrict(fac, fac.Verum(), ass))
	assert.Equal(fac.Falsum(), Restrict(fac, fac.Falsum(), ass))
}

func TestRestrictLiterals(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	ass, _ := New(fac, f.Literal(d.A), f.Literal(d.NB), f.Literal(d.NX))

	assert.Equal(d.True, Restrict(fac, d.A, ass))
	assert.Equal(d.False, Restrict(fac, d.NA, ass))
	assert.Equal(d.False, Restrict(fac, d.X, ass))
	assert.Equal(d.True, Restrict(fac, d.NX, ass))
	assert.Equal(d.C, Restrict(fac, d.C, ass))
	assert.Equal(d.NY, Restrict(fac, d.NY, ass))
}

func TestRestrictNot(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	ass, _ := New(fac, f.Literal(d.A), f.Literal(d.NB), f.Literal(d.NX))

	assert.Equal(d.True, Restrict(fac, d.NOT1, ass))
	assert.Equal(d.NY, Restrict(fac, d.NOT2, ass))
}

func TestRestrictBinaryOperators(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	ass, _ := New(fac, f.Literal(d.A), f.Literal(d.NB), f.Literal(d.NX))

	assert.Equal(d.False, Restrict(fac, d.IMP1, ass))
	assert.Equal(d.True, Restrict(fac, d.IMP2, ass))
	assert.Equal(d.True, Restrict(fac, fac.Implication(d.NA, d.C), ass))
	assert.Equal(d.True, Restrict(fac, d.IMP3, ass))
	assert.Equal(d.C, Restrict(fac, fac.Implication(d.A, d.C), ass))
	assert.Equal(d.False, Restrict(fac, d.EQ1, ass))
	assert.Equal(d.False, Restrict(fac, d.EQ2, ass))
	assert.Equal(d.NY, Restrict(fac, d.EQ3, ass))
	assert.Equal(d.False, Restrict(fac, d.EQ4, ass))
}

func TestRestrictNaryOperators(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	p := parser.New(fac)
	ass, _ := New(fac, f.Literal(d.A), f.Literal(d.NB), f.Literal(d.NX))

	assert.Equal(d.Y, Restrict(fac, d.OR1, ass))
	assert.Equal(d.True, Restrict(fac, d.OR2, ass))
	assert.Equal(d.False, Restrict(fac, d.OR3, ass))
	assert.Equal(p.ParseUnsafe("~c | y"), Restrict(fac, p.ParseUnsafe("~a | b | ~c | x | y"), ass))
	assert.Equal(d.True, Restrict(fac, p.ParseUnsafe("~a | b | ~c | ~x | ~y"), ass))
	assert.Equal(d.False, Restrict(fac, d.AND1, ass))
	assert.Equal(d.False, Restrict(fac, d.AND2, ass))
	assert.Equal(d.Y, Restrict(fac, d.AND3, ass))
	assert.Equal(p.ParseUnsafe("c & ~y"), Restrict(fac, p.ParseUnsafe("a & ~b & c & ~x & ~y"), ass))
	assert.Equal(d.False, Restrict(fac, p.ParseUnsafe("a & b & c & ~x & y"), ass))
}

func TestRestrictPbc(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()

	lits := []f.Literal{fac.Lit("a", true), fac.Lit("b", false), fac.Lit("c", true)}
	litsA1 := []f.Literal{fac.Lit("b", false), fac.Lit("c", true)}
	litsA2 := []f.Literal{fac.Lit("c", true)}
	coeffs := []int{2, -2, 3}
	coeffsA1 := []int{-2, 3}
	coeffsA2 := []int{3}

	a1, _ := New(fac, fac.Lit("a", true))
	a2, _ := New(fac, fac.Lit("a", true), fac.Lit("b", false))
	a3, _ := New(fac, fac.Lit("a", true), fac.Lit("b", false), fac.Lit("c", true))
	a4, _ := New(fac, fac.Lit("a", false), fac.Lit("b", true), fac.Lit("c", false))
	a5, _ := New(fac, fac.Lit("a", true), fac.Lit("b", true), fac.Lit("c", false))

	pb1 := fac.PBC(f.EQ, 2, lits, coeffs)
	assert.Equal(fac.PBC(f.EQ, 0, litsA1, coeffsA1), Restrict(fac, pb1, a1))
	assert.Equal(fac.PBC(f.EQ, 2, litsA2, coeffsA2), Restrict(fac, pb1, a2))
	assert.Equal(fac.Falsum(), Restrict(fac, pb1, a3))
	assert.Equal(fac.Falsum(), Restrict(fac, pb1, a4))
	assert.Equal(fac.Verum(), Restrict(fac, pb1, a5))
}
