package normalform

import (
	"testing"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/parser"
	"github.com/stretchr/testify/assert"
)

func TestDNFPredicate(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)

	assert.True(IsDNF(fac, fac.Verum()))
	assert.True(IsDNF(fac, fac.Falsum()))
	assert.True(IsDNF(fac, d.A))
	assert.True(IsDNF(fac, d.NA))
	assert.True(IsDNF(fac, d.AND1))
	assert.True(IsDNF(fac, d.OR1))
	assert.True(IsDNF(fac, d.OR3))
	assert.True(IsDNF(fac, fac.Or(d.AND1, d.AND2, d.A, d.NY)))
	assert.False(IsDNF(fac, d.PBC1))
	assert.False(IsDNF(fac, d.AND3))
	assert.False(IsDNF(fac, d.IMP1))
	assert.False(IsDNF(fac, d.EQ1))
	assert.False(IsDNF(fac, d.NOT1))
	assert.False(IsDNF(fac, d.NOT2))
	assert.False(IsDNF(fac, fac.Or(d.AND1, d.EQ1)))
}

func TestDNFConstants(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()

	assert.Equal(fac.Verum(), FactorizedDNF(fac, fac.Verum()))
	assert.Equal(fac.Falsum(), FactorizedDNF(fac, fac.Falsum()))
}

func TestDNFLiterals(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)

	assert.Equal(d.A, FactorizedDNF(fac, d.A))
	assert.Equal(d.NA, FactorizedDNF(fac, d.NA))
}

func TestDNFBinaryOperators(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	p := parser.New(fac)

	assert.Equal(p.ParseUnsafe("~a | b"), FactorizedDNF(fac, d.IMP1))
	assert.Equal(p.ParseUnsafe("a | ~b"), FactorizedDNF(fac, d.IMP2))
	assert.Equal(p.ParseUnsafe("~a | ~b | x | y"), FactorizedDNF(fac, d.IMP3))
	assert.Equal(p.ParseUnsafe("~b & ~a | a & b"), FactorizedDNF(fac, d.EQ1))
	assert.Equal(p.ParseUnsafe("b & a | ~a & ~b"), FactorizedDNF(fac, d.EQ2))

	assert.True(IsDNF(fac, FactorizedDNF(fac, d.IMP1)))
	assert.True(IsDNF(fac, FactorizedDNF(fac, d.IMP2)))
	assert.True(IsDNF(fac, FactorizedDNF(fac, d.IMP3)))
	assert.True(IsDNF(fac, FactorizedDNF(fac, d.EQ1)))
	assert.True(IsDNF(fac, FactorizedDNF(fac, d.EQ2)))
	assert.True(IsDNF(fac, FactorizedDNF(fac, d.EQ3)))
	assert.True(IsDNF(fac, FactorizedDNF(fac, d.EQ4)))
}

func TestDNFNaryOperators(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	p := parser.New(fac)

	assert.Equal(d.AND1, FactorizedDNF(fac, d.AND1))
	assert.Equal(d.OR1, FactorizedDNF(fac, d.OR1))

	assert.Equal(p.ParseUnsafe("~a | ~b | c | (~x & y)"), FactorizedDNF(fac, p.ParseUnsafe("~(a & b) | c | ~(x | ~y)")))
	assert.Equal(p.ParseUnsafe("~x & ~a & ~b & c | y & ~a & ~b & c"), FactorizedDNF(fac, p.ParseUnsafe("~(a | b) & c & ~(x & ~y)")))
	assert.Equal(p.ParseUnsafe("~x & a & b | ~y & a & b"), FactorizedDNF(fac, p.ParseUnsafe("a & b & (~x | ~y)")))

	assert.True(IsDNF(fac, FactorizedDNF(fac, d.AND1)))
	assert.True(IsDNF(fac, FactorizedDNF(fac, d.AND2)))
	assert.True(IsDNF(fac, FactorizedDNF(fac, d.AND3)))
	assert.True(IsDNF(fac, FactorizedDNF(fac, d.OR1)))
	assert.True(IsDNF(fac, FactorizedDNF(fac, d.OR2)))
	assert.True(IsDNF(fac, FactorizedDNF(fac, d.OR3)))
	assert.True(IsDNF(fac, FactorizedDNF(fac, p.ParseUnsafe("~(a | b) & c & ~(x & ~y) & (w => z)"))))
	assert.True(IsDNF(fac, FactorizedDNF(fac, p.ParseUnsafe("~(a & b) | c | ~(x | ~y)"))))
	assert.True(IsDNF(fac, FactorizedDNF(fac, p.ParseUnsafe("a | b | (~x & ~y)"))))
}

func TestDNFNot(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	assert.Equal(p.ParseUnsafe("~a"), FactorizedDNF(fac, p.ParseUnsafe("~a")))
	assert.Equal(p.ParseUnsafe("a & ~b"), FactorizedCNF(fac, p.ParseUnsafe("~(a => b)")))
	assert.Equal(p.ParseUnsafe("b & ~a | a & ~b"), FactorizedDNF(fac, p.ParseUnsafe("~(a <=> b)")))
	assert.Equal(p.ParseUnsafe("x & ~a & ~b | y & ~a & ~b"), FactorizedDNF(fac, p.ParseUnsafe("~(~(a | b) => ~(x | y))")))
	assert.Equal(p.ParseUnsafe("~a2 | ~b2 | x2 | y2"), FactorizedDNF(fac, p.ParseUnsafe("~(a2 & b2 & ~x2 & ~y2)")))
	assert.Equal(p.ParseUnsafe("~a2 & ~b2 & x2 & y2"), FactorizedDNF(fac, p.ParseUnsafe("~(a2 | b2 | ~x2 | ~y2)")))
	assert.Equal(p.ParseUnsafe("~a2 & ~b2 & x2 & y2"), FactorizedDNF(fac, p.ParseUnsafe("~(a2 | b2 | ~x2 | ~y2)")))
}
