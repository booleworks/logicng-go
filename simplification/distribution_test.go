package simplification

import (
	"testing"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/function"
	"github.com/booleworks/logicng-go/parser"
	"github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

func TestDistributeSimple(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	p := parser.New(fac)

	assert.Equal(fac.Falsum(), Distribute(fac, fac.Falsum()))
	assert.Equal(fac.Verum(), Distribute(fac, fac.Verum()))
	assert.Equal(d.A, Distribute(fac, d.A))
	assert.Equal(d.NA, Distribute(fac, d.NA))
	assert.Equal(d.AND1, Distribute(fac, d.AND1))
	assert.Equal(d.AND2, Distribute(fac, d.AND2))
	assert.Equal(d.OR1, Distribute(fac, d.OR1))
	assert.Equal(d.OR2, Distribute(fac, d.OR2))
	assert.Equal(d.IMP1, Distribute(fac, d.IMP1))
	assert.Equal(d.EQ1, Distribute(fac, d.EQ1))
	assert.Equal(d.NOT1, Distribute(fac, d.NOT1))

	assert.Equal(d.AND1, Distribute(fac, fac.And(d.AND1, d.A)))
	assert.Equal(d.False, Distribute(fac, fac.And(d.AND2, d.A)))
	assert.Equal(fac.And(d.X, d.OR1), Distribute(fac, fac.And(d.OR1, d.X)))
	assert.Equal(fac.And(d.X, d.OR2), Distribute(fac, fac.And(d.OR2, d.X)))
	assert.Equal(fac.Or(d.A, d.AND1), Distribute(fac, fac.Or(d.AND1, d.A)))
	assert.Equal(fac.Or(d.A, d.AND2), Distribute(fac, fac.Or(d.AND2, d.A)))
	assert.Equal(d.OR1, Distribute(fac, fac.Or(d.OR1, d.X)))
	formula := p.ParseUnsafe("(a | b | ~c) & (~a | ~d) & (~c | d) & (~b | e | ~f | g) & (e | f | g | h) & (e | ~f | ~g | h) & f & c")
	dist := Distribute(fac, formula)

	assert.True(sat.IsEquivalent(fac, formula, dist))
	assert.Equal(19, function.NumberOfAtoms(fac, dist))
}

func TestDistributeComplex(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	cAnd := p.ParseUnsafe("(a | b | ~c) & (~a | ~d) & (~c | d | b) & (~c | ~b)")
	cAndD1 := Distribute(fac, cAnd)
	verify(t, fac, cAndD1, cAnd)
	verify(t, fac, Distribute(fac, cAndD1), cAnd)

	assert.Equal(fac.Not(cAndD1), Distribute(fac, fac.Not(cAnd)))

	cOr := p.ParseUnsafe("(x & y & z) | (x & y & ~z) | (x & ~y & z)")
	cOrD1 := Distribute(fac, cOr)
	verify(t, fac, cOrD1, cOr)
	verify(t, fac, Distribute(fac, cOrD1), cOr)

	assert.Equal(fac.Equivalence(cOrD1, cAndD1), Distribute(fac, fac.Equivalence(cOr, cAnd)))
	assert.Equal(fac.Implication(cOrD1, cAndD1), Distribute(fac, fac.Implication(cOr, cAnd)))
	assert.Equal(fac.Not(cOrD1), Distribute(fac, fac.Not(cOr)))
}

func verify(t *testing.T, fac f.Factory, f1, f2 f.Formula) {
	assert.True(t, sat.IsEquivalent(fac, f1, f2))
	assert.Less(t, len(f1.Sprint(fac)), len(f2.Sprint(fac)))
}
