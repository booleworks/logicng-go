package normalform

import (
	"testing"

	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/parser"
	"github.com/stretchr/testify/assert"
)

func TestNnfConstants(t *testing.T) {
	fac := f.NewFactory()
	d := f.NewTestData(fac)

	assert.Equal(t, d.True, NNF(fac, d.True))
	assert.Equal(t, d.False, NNF(fac, d.False))
}

func TestNnfLiterals(t *testing.T) {
	fac := f.NewFactory()
	d := f.NewTestData(fac)

	assert.Equal(t, d.A, NNF(fac, d.A))
	assert.Equal(t, d.NA, NNF(fac, d.NA))
}

func TestNnfBinaryOperators(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	p := parser.New(fac)

	assert.Equal(p.ParseUnsafe("~a | b"), NNF(fac, d.IMP1))
	assert.Equal(p.ParseUnsafe("a | ~b"), NNF(fac, d.IMP2))
	assert.Equal(p.ParseUnsafe("~a | ~b | x | y"), NNF(fac, d.IMP3))
	assert.Equal(p.ParseUnsafe("(~a | ~b) & (a | b) | (x | ~y) & (~x | y)"), NNF(fac, d.IMP4))
	assert.Equal(p.ParseUnsafe("(~a | b) & (a | ~b)"), NNF(fac, d.EQ1))
	assert.Equal(p.ParseUnsafe("(a | ~b) & (~a | b)"), NNF(fac, d.EQ2))
	assert.Equal(p.ParseUnsafe("(~a | ~b | x | y) & (a & b | ~x & ~y)"), NNF(fac, d.EQ3))
}

func TestNnfNAryOperators(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	p := parser.New(fac)

	assert.Equal(d.AND1, NNF(fac, d.AND1))
	assert.Equal(d.OR1, NNF(fac, d.OR1))
	assert.Equal(p.ParseUnsafe("~a & ~b & c & (~x | y) & (~w | z)"), NNF(fac, p.ParseUnsafe("~(a | b) & c & ~(x & ~y) & (w => z)")))
	assert.Equal(p.ParseUnsafe("~a  | ~b | c | (~x & y) | (~w | z)"), NNF(fac, p.ParseUnsafe("~(a & b) | c | ~(x | ~y) | (w => z)")))
}

func TestNnfNot(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	assert.Equal(p.ParseUnsafe("~a"), NNF(fac, p.ParseUnsafe("~a")))
	assert.Equal(p.ParseUnsafe("a"), NNF(fac, p.ParseUnsafe("~~a")))
	assert.Equal(p.ParseUnsafe("a & ~b"), NNF(fac, p.ParseUnsafe("~(a => b)")))
	assert.Equal(p.ParseUnsafe("~a & ~b & (x | y)"), NNF(fac, p.ParseUnsafe("~(~(a | b) => ~(x | y))")))
	assert.Equal(p.ParseUnsafe("(~a | b) & (a | ~b)"), NNF(fac, p.ParseUnsafe("a <=> b")))
	assert.Equal(p.ParseUnsafe("(~a | ~b) & (a | b)"), NNF(fac, p.ParseUnsafe("~(a <=> b)")))
	assert.Equal(p.ParseUnsafe("((a | b) | (x | y)) & ((~a & ~b) | (~x & ~y))"), NNF(fac, p.ParseUnsafe("~(~(a | b) <=> ~(x | y))")))
	assert.Equal(p.ParseUnsafe("~a | ~b | x | y"), NNF(fac, p.ParseUnsafe("~(a & b & ~x & ~y)")))
	assert.Equal(p.ParseUnsafe("~a & ~b & x & y"), NNF(fac, p.ParseUnsafe("~(a | b | ~x | ~y)")))
	assert.Equal(p.ParseUnsafe("~a & ~b & x & y"), NNF(fac, p.ParseUnsafe("~(a | b | ~x | ~y)")))
}

func TestNnfPredicate(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)

	assert.True(IsNNF(fac, d.F.Verum()))
	assert.True(IsNNF(fac, d.F.Falsum()))
	assert.True(IsNNF(fac, d.A))
	assert.True(IsNNF(fac, d.NA))
	assert.True(IsNNF(fac, d.OR1))
	assert.True(IsNNF(fac, d.AND1))
	assert.True(IsNNF(fac, d.AND3))
	assert.True(IsNNF(fac, d.F.And(d.OR1, d.OR2, d.A, d.NY)))
	assert.True(IsNNF(fac, d.F.And(d.OR1, d.OR2, d.AND1, d.AND2, d.AND3, d.A, d.NY)))
	assert.True(IsNNF(fac, d.OR3))
	assert.False(IsNNF(fac, d.PBC1))
	assert.False(IsNNF(fac, d.IMP1))
	assert.False(IsNNF(fac, d.EQ1))
	assert.False(IsNNF(fac, d.NOT1))
	assert.False(IsNNF(fac, d.NOT2))
	assert.False(IsNNF(fac, d.F.And(d.OR1, d.F.Not(d.OR2), d.A, d.NY)))
	assert.False(IsNNF(fac, d.F.And(d.OR1, d.EQ1)))
	assert.False(IsNNF(fac, d.F.And(d.OR1, d.IMP1, d.AND1)))
	assert.False(IsNNF(fac, d.F.And(d.OR1, d.PBC1, d.AND1)))
}
