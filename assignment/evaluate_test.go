package assignment

import (
	"testing"

	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/parser"
	"github.com/stretchr/testify/assert"
)

func TestEvalConstants(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	ass, _ := New(fac, f.Literal(d.A), f.Literal(d.B), f.Literal(d.C), f.Literal(d.NX), f.Literal(d.NY))

	assert.True(Evaluate(fac, fac.Verum(), ass))
	assert.False(Evaluate(fac, fac.Falsum(), ass))
}

func TestEvalLiterals(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	ass, _ := New(fac, f.Literal(d.A), f.Literal(d.B), f.Literal(d.C), f.Literal(d.NX), f.Literal(d.NY))

	assert.True(Evaluate(fac, d.A, ass))
	assert.False(Evaluate(fac, d.NA, ass))
	assert.False(Evaluate(fac, d.X, ass))
	assert.True(Evaluate(fac, d.NX, ass))
}

func TestEvalNot(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	ass, _ := New(fac, f.Literal(d.A), f.Literal(d.B), f.Literal(d.C), f.Literal(d.NX), f.Literal(d.NY))

	assert.False(Evaluate(fac, d.NOT1, ass))
	assert.True(Evaluate(fac, d.NOT2, ass))
}

func TestEvalBinaryOperators(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	ass, _ := New(fac, f.Literal(d.A), f.Literal(d.B), f.Literal(d.C), f.Literal(d.NX), f.Literal(d.NY))

	assert.True(Evaluate(fac, d.IMP1, ass))
	assert.True(Evaluate(fac, d.IMP2, ass))
	assert.False(Evaluate(fac, d.IMP3, ass))
	assert.True(Evaluate(fac, d.IMP4, ass))
	assert.True(Evaluate(fac, d.EQ1, ass))
	assert.True(Evaluate(fac, d.EQ2, ass))
	assert.False(Evaluate(fac, d.EQ3, ass))
	assert.True(Evaluate(fac, d.EQ4, ass))
}

func TestEvalNaryOperators(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	p := parser.New(fac)
	ass, _ := New(fac, f.Literal(d.A), f.Literal(d.B), f.Literal(d.C), f.Literal(d.NX), f.Literal(d.NY))

	assert.False(Evaluate(fac, d.OR1, ass))
	assert.True(Evaluate(fac, d.OR2, ass))
	assert.True(Evaluate(fac, d.OR3, ass))
	assert.False(Evaluate(fac, p.ParseUnsafe("~a | ~b | ~c | x | y"), ass))
	assert.True(Evaluate(fac, p.ParseUnsafe("~a | ~b | ~c | x | ~y"), ass))
	assert.True(Evaluate(fac, d.AND1, ass))
	assert.False(Evaluate(fac, d.AND2, ass))
	assert.False(Evaluate(fac, d.AND3, ass))
	assert.True(Evaluate(fac, p.ParseUnsafe("a & b & c & ~x & ~y"), ass))
	assert.False(Evaluate(fac, p.ParseUnsafe("a & b & c & ~x & y"), ass))
}

func TestEvalPbc(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	lits := []f.Literal{fac.Lit("a", true), fac.Lit("b", false), fac.Lit("c", true)}
	coeffs := []int{2, -2, 3}
	a1, _ := New(fac, fac.Lit("a", true), fac.Lit("b", true), fac.Lit("c", false))
	a2, _ := New(fac, fac.Lit("a", true), fac.Lit("b", false), fac.Lit("c", false))

	pb1 := fac.PBC(f.EQ, 2, lits, coeffs)
	pb3 := fac.PBC(f.GE, 1, lits, coeffs)
	pb4 := fac.PBC(f.GT, 0, lits, coeffs)
	pb5 := fac.PBC(f.LE, 1, lits, coeffs)
	pb6 := fac.PBC(f.LT, 2, lits, coeffs)

	assert.True(Evaluate(fac, pb1, a1))
	assert.False(Evaluate(fac, pb1, a2))
	assert.True(Evaluate(fac, pb3, a1))
	assert.False(Evaluate(fac, pb3, a2))
	assert.True(Evaluate(fac, pb4, a1))
	assert.False(Evaluate(fac, pb4, a2))
	assert.False(Evaluate(fac, pb5, a1))
	assert.True(Evaluate(fac, pb5, a2))
	assert.False(Evaluate(fac, pb6, a1))
	assert.True(Evaluate(fac, pb6, a2))
}
