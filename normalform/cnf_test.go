package normalform

import (
	"testing"

	"github.com/booleworks/logicng-go/event"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/parser"
	"github.com/stretchr/testify/assert"
)

func TestCNFConstants(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()

	assert.Equal(fac.Verum(), FactorizedCNF(fac, fac.Verum()))
	assert.Equal(fac.Falsum(), FactorizedCNF(fac, fac.Falsum()))
}

func TestCNFPGConstants(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()

	assert.Equal(fac.Verum(), PGCNFDefault(fac, fac.Verum()))
	assert.Equal(fac.Falsum(), PGCNFDefault(fac, fac.Falsum()))
}

func TestTseitinConstants(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()

	assert.Equal(fac.Verum(), TseitinCNFDefault(fac, fac.Verum()))
	assert.Equal(fac.Falsum(), TseitinCNFDefault(fac, fac.Falsum()))
}

func TestCNFLiterals(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)

	assert.Equal(d.A, FactorizedCNF(fac, d.A))
	assert.Equal(d.NA, FactorizedCNF(fac, d.NA))
}

func TestCNFPGLiterals(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)

	assert.Equal(d.A, PGCNFDefault(fac, d.A))
	assert.Equal(d.NA, PGCNFDefault(fac, d.NA))
}

func TestTseitinLiterals(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)

	assert.Equal(d.A, TseitinCNFDefault(fac, d.A))
	assert.Equal(d.NA, TseitinCNFDefault(fac, d.NA))
}

func TestCNFBinaryOperators(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	p := parser.New(fac)

	assert.Equal(p.ParseUnsafe("~a | b"), FactorizedCNF(fac, d.IMP1))
	assert.Equal(p.ParseUnsafe("a | ~b"), FactorizedCNF(fac, d.IMP2))
	assert.Equal(p.ParseUnsafe("~a | ~b | x | y"), FactorizedCNF(fac, d.IMP3))
	assert.Equal(p.ParseUnsafe("(~a | b) & (a | ~b)"), FactorizedCNF(fac, d.EQ1))
	assert.Equal(p.ParseUnsafe("(a | ~b) & (~a | b)"), FactorizedCNF(fac, d.EQ2))

	assert.True(IsCNF(fac, FactorizedCNF(fac, d.IMP1)))
	assert.True(IsCNF(fac, FactorizedCNF(fac, d.IMP2)))
	assert.True(IsCNF(fac, FactorizedCNF(fac, d.IMP3)))
	assert.True(IsCNF(fac, FactorizedCNF(fac, d.EQ1)))
	assert.True(IsCNF(fac, FactorizedCNF(fac, d.EQ2)))
	assert.True(IsCNF(fac, FactorizedCNF(fac, d.EQ3)))
	assert.True(IsCNF(fac, FactorizedCNF(fac, d.EQ4)))
}

func TestCNFPGBinaryOperators(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	p := parser.New(fac)

	assert.Equal(p.ParseUnsafe("~a | b"), PGCNFDefault(fac, d.IMP1))
	assert.Equal(p.ParseUnsafe("a | ~b"), PGCNFDefault(fac, d.IMP2))
	assert.Equal(p.ParseUnsafe("~a | ~b | x | y"), PGCNFDefault(fac, d.IMP3))
	assert.Equal(p.ParseUnsafe("(~a | b) & (a | ~b)"), PGCNFDefault(fac, d.EQ1))
	assert.Equal(p.ParseUnsafe("(a | ~b) & (~a | b)"), PGCNFDefault(fac, d.EQ2))

	assert.True(IsCNF(fac, PGCNFDefault(fac, d.IMP1)))
	assert.True(IsCNF(fac, PGCNFDefault(fac, d.IMP2)))
	assert.True(IsCNF(fac, PGCNFDefault(fac, d.IMP3)))
	assert.True(IsCNF(fac, PGCNFDefault(fac, d.EQ1)))
	assert.True(IsCNF(fac, PGCNFDefault(fac, d.EQ2)))
	assert.True(IsCNF(fac, PGCNFDefault(fac, d.EQ3)))
	assert.True(IsCNF(fac, PGCNFDefault(fac, d.EQ4)))
}

func TestCNFTseitinBinaryOperators(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	p := parser.New(fac)

	assert.Equal(p.ParseUnsafe("~a | b"), TseitinCNFDefault(fac, d.IMP1))
	assert.Equal(p.ParseUnsafe("a | ~b"), TseitinCNFDefault(fac, d.IMP2))
	assert.Equal(p.ParseUnsafe("~a | ~b | x | y"), TseitinCNFDefault(fac, d.IMP3))
	assert.Equal(p.ParseUnsafe("(~a | b) & (a | ~b)"), TseitinCNFDefault(fac, d.EQ1))
	assert.Equal(p.ParseUnsafe("(a | ~b) & (~a | b)"), TseitinCNFDefault(fac, d.EQ2))

	assert.True(IsCNF(fac, TseitinCNFDefault(fac, d.IMP1)))
	assert.True(IsCNF(fac, TseitinCNFDefault(fac, d.IMP2)))
	assert.True(IsCNF(fac, TseitinCNFDefault(fac, d.IMP3)))
	assert.True(IsCNF(fac, TseitinCNFDefault(fac, d.EQ1)))
	assert.True(IsCNF(fac, TseitinCNFDefault(fac, d.EQ2)))
	assert.True(IsCNF(fac, TseitinCNFDefault(fac, d.EQ3)))
	assert.True(IsCNF(fac, TseitinCNFDefault(fac, d.EQ4)))
}

func TestCNFNaryOperators(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	p := parser.New(fac)

	assert.Equal(d.AND1, FactorizedCNF(fac, d.AND1))
	assert.Equal(d.OR1, FactorizedCNF(fac, d.OR1))
	assert.Equal(p.ParseUnsafe("~a & ~b & c & (~x | y) & (~w | z)"), FactorizedCNF(fac, p.ParseUnsafe("~(a | b) & c & ~(x & ~y) & (w => z)")))
	assert.Equal(p.ParseUnsafe("(~x | ~a | ~b | c) & (y | ~a | ~b | c)"), FactorizedCNF(fac, p.ParseUnsafe("~(a & b) | c | ~(x | ~y)")))
	assert.Equal(p.ParseUnsafe("(~x | a | b) & (~y | a | b)"), FactorizedCNF(fac, p.ParseUnsafe("a | b | (~x & ~y)")))

	assert.True(IsCNF(fac, FactorizedCNF(fac, d.AND1)))
	assert.True(IsCNF(fac, FactorizedCNF(fac, d.AND2)))
	assert.True(IsCNF(fac, FactorizedCNF(fac, d.AND3)))
	assert.True(IsCNF(fac, FactorizedCNF(fac, d.OR1)))
	assert.True(IsCNF(fac, FactorizedCNF(fac, d.OR2)))
	assert.True(IsCNF(fac, FactorizedCNF(fac, d.OR3)))
	assert.True(IsCNF(fac, FactorizedCNF(fac, p.ParseUnsafe("~(a | b) & c & ~(x & ~y) & (w => z)"))))
	assert.True(IsCNF(fac, FactorizedCNF(fac, p.ParseUnsafe("~(a & b) | c | ~(x | ~y)"))))
	assert.True(IsCNF(fac, FactorizedCNF(fac, p.ParseUnsafe("a | b | (~x & ~y)"))))
}

func TestCNFNot(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	assert.Equal(p.ParseUnsafe("~a"), FactorizedCNF(fac, p.ParseUnsafe("~a")))
	assert.Equal(p.ParseUnsafe("a2 & ~b2"), FactorizedCNF(fac, p.ParseUnsafe("~(a2 => b2)")))
	assert.Equal(p.ParseUnsafe("(~a2 | ~b2) & (a2 | b2)"), FactorizedCNF(fac, p.ParseUnsafe("~(a2 <=> b2)")))
	assert.Equal(p.ParseUnsafe("~a2 | ~b2 | x2 | y2"), FactorizedCNF(fac, p.ParseUnsafe("~(a2 & b2 & ~x2 & ~y2)")))
	assert.Equal(p.ParseUnsafe("~a2 & ~b2 & x2 & y2"), FactorizedCNF(fac, p.ParseUnsafe("~(a2 | b2 | ~x2 | ~y2)")))
	assert.Equal(p.ParseUnsafe("~a2 & ~b2 & x2 & y2"), FactorizedCNF(fac, p.ParseUnsafe("~(a2 | b2 | ~x2 | ~y2)")))

	handler := NewFactorizationHandler(-1, -1)
	cnf, state := FactorizedCNFWithHandler(fac, p.ParseUnsafe("~(~(a2 | b2) <=> ~(x2 | y2))"), handler)
	assert.Equal(p.ParseUnsafe("(a2 | b2 | x2 | y2) & (~x2 | ~a2) & (~y2 | ~a2) & (~x2 | ~b2) & (~y2 | ~b2)"), cnf)
	assert.True(state.Success)
	assert.Equal(7, handler.currentDistributions)
	assert.Equal(4, handler.currentClauses)
}

func TestCNFWithHandler(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	formula := p.ParseUnsafe("(~(~(a | b) => ~(x | y))) & ((a | x) => ~(b | y))")
	handler := NewFactorizationHandler(-1, 2)

	cnf, state := FactorizedCNFWithHandler(fac, formula, handler)
	assert.Equal(f.Formula(0), cnf)
	assert.False(state.Success)
	assert.NotEqual(event.Nothing, state.CancelCause)

	formula = p.ParseUnsafe("~(a | b)")
	handler = NewFactorizationHandler(-1, 2)
	cnf, state = FactorizedCNFWithHandler(fac, formula, handler)
	assert.Equal(p.ParseUnsafe("~a & ~b"), cnf)
	assert.True(state.Success)
}
