package normalform

import (
	"testing"

	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/parser"
	"github.com/stretchr/testify/assert"
)

func TestAIGConstants(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()

	assert.Equal(fac.Falsum(), AIG(fac, fac.Falsum()))
	assert.Equal(fac.Verum(), AIG(fac, fac.Verum()))
	assert.True(IsAIG(fac, AIG(fac, fac.Falsum())))
	assert.True(IsAIG(fac, AIG(fac, fac.Verum())))
}

func TestAIGLiterals(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)

	assert.Equal(d.A, AIG(fac, d.A))
	assert.Equal(d.NA, AIG(fac, d.NA))
	assert.True(IsAIG(fac, AIG(fac, d.A)))
	assert.True(IsAIG(fac, AIG(fac, d.NA)))
}

func TestAIGBinaryOperators(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	p := parser.New(fac)

	assert.Equal(p.ParseUnsafe("~(a & ~b)"), AIG(fac, d.IMP1))
	assert.Equal(p.ParseUnsafe("~(~a & b)"), AIG(fac, d.IMP2))
	assert.Equal(p.ParseUnsafe("~((a & b) & (~x & ~y))"), AIG(fac, d.IMP3))
	assert.Equal(p.ParseUnsafe("~(a & ~b) & ~(~a & b)"), AIG(fac, d.EQ1))
	assert.Equal(p.ParseUnsafe("~(~a & b) & ~(a & ~b)"), AIG(fac, d.EQ2))

	assert.True(IsAIG(fac, AIG(fac, d.IMP1)))
	assert.True(IsAIG(fac, AIG(fac, d.IMP2)))
	assert.True(IsAIG(fac, AIG(fac, d.IMP3)))
	assert.True(IsAIG(fac, AIG(fac, d.EQ1)))
	assert.True(IsAIG(fac, AIG(fac, d.EQ2)))

	assert.False(IsAIG(fac, d.IMP1))
	assert.False(IsAIG(fac, d.IMP2))
	assert.False(IsAIG(fac, d.IMP3))
	assert.False(IsAIG(fac, d.EQ1))
	assert.False(IsAIG(fac, d.EQ2))
}

func TestAIGNaryOperators(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	p := parser.New(fac)

	assert.Equal(p.ParseUnsafe("a & b"), AIG(fac, d.AND1))
	assert.Equal(
		p.ParseUnsafe("(~a & ~b) & c & ~(x & ~y) & ~(w & ~z)"),
		AIG(fac, p.ParseUnsafe("~(a | b) & c & ~(x & ~y) & (w => z)")),
	)
	assert.Equal(
		p.ParseUnsafe("~(a & b & ~c & ~(~x & y))"),
		AIG(fac, p.ParseUnsafe("~(a & b) | c | ~(x | ~y)")),
	)
	assert.Equal(
		p.ParseUnsafe("~(~a & ~b & ~(~x & ~y))"),
		AIG(fac, p.ParseUnsafe("a | b | (~x & ~y)")),
	)
	assert.Equal(
		p.ParseUnsafe("~(~a & ~b & ~(~x & ~y))"),
		AIG(fac, p.ParseUnsafe("a | b | (~x & ~y)")),
	)

	assert.True(IsAIG(fac, AIG(fac, p.ParseUnsafe("~(a | b) & c & ~(x & ~y) & (w => z)"))))
	assert.True(IsAIG(fac, AIG(fac, p.ParseUnsafe("~(a & b) | c | ~(x | ~y)"))))
	assert.True(IsAIG(fac, AIG(fac, p.ParseUnsafe("a | b | (~x & ~y)"))))
	assert.True(IsAIG(fac, AIG(fac, d.AND1)))
	assert.True(IsAIG(fac, AIG(fac, d.AND2)))
	assert.True(IsAIG(fac, AIG(fac, d.AND3)))
	assert.True(IsAIG(fac, AIG(fac, d.OR1)))
	assert.True(IsAIG(fac, AIG(fac, d.OR2)))
	assert.True(IsAIG(fac, AIG(fac, d.OR3)))

	assert.False(IsAIG(fac, p.ParseUnsafe("~(a | b) & c & ~(x & ~y) & (w => z)")))
	assert.False(IsAIG(fac, p.ParseUnsafe("~(a & b) | c | ~(x | ~y)")))
	assert.False(IsAIG(fac, p.ParseUnsafe("a | b | (~x & ~y)")))
	assert.False(IsAIG(fac, d.OR2))
	assert.False(IsAIG(fac, d.OR3))
}

func TestAIGNot(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	assert.Equal(p.ParseUnsafe("a & ~b"), AIG(fac, p.ParseUnsafe("~(a => b)")))
	assert.Equal(p.ParseUnsafe("(~a & ~b) & ~(~x & ~y)"), AIG(fac, p.ParseUnsafe("~(~(a | b) => ~(x | y))")))
	assert.Equal(p.ParseUnsafe("~(~(a & ~b) & ~(~a & b))"), AIG(fac, p.ParseUnsafe("~(a <=> b)")))
	assert.Equal(p.ParseUnsafe("~a & ~b & x & y"), AIG(fac, p.ParseUnsafe("~(a | b | ~x | ~y)")))
	assert.Equal(p.ParseUnsafe("~(a & b & ~x & ~y)"), AIG(fac, p.ParseUnsafe("~(a & b & ~x & ~y)")))
	assert.Equal(
		p.ParseUnsafe("~(~(~a & ~b & ~(~x & ~y)) & ~((a | b) & ~(x | y)))"),
		AIG(fac, p.ParseUnsafe("~(~(a | b) <=> ~(x | y))")),
	)
	assert.Equal(p.ParseUnsafe("~a & ~b & x & y"), AIG(fac, p.ParseUnsafe("~(a | b | ~x | ~y)")))
}

func TestAIGPbc(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)

	assert.True(IsAIG(fac, AIG(fac, d.PBC1)))
	assert.True(IsAIG(fac, AIG(fac, d.PBC2)))
	assert.True(IsAIG(fac, AIG(fac, d.PBC3)))
	assert.True(IsAIG(fac, AIG(fac, d.PBC4)))
	assert.True(IsAIG(fac, AIG(fac, d.PBC5)))
	assert.True(IsAIG(fac, AIG(fac, d.CC1)))
	assert.True(IsAIG(fac, AIG(fac, d.CC2)))

	assert.False(IsAIG(fac, d.PBC1))
	assert.False(IsAIG(fac, d.PBC2))
	assert.False(IsAIG(fac, d.PBC3))
	assert.False(IsAIG(fac, d.PBC4))
	assert.False(IsAIG(fac, d.PBC5))
	assert.False(IsAIG(fac, d.CC1))
	assert.False(IsAIG(fac, d.CC2))
}
