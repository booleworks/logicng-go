package sat

import (
	"testing"

	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/parser"
	"github.com/stretchr/testify/assert"
)

func TestIsSatisfiable(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	f1 := p.ParseUnsafe("(a | b) & (c | ~d)")
	f2 := p.ParseUnsafe("~a & ~b & (a | b)")

	assert.False(IsSatisfiable(fac, fac.Falsum()))
	assert.True(IsSatisfiable(fac, fac.Verum()))
	assert.True(IsSatisfiable(fac, f1))
	assert.False(IsSatisfiable(fac, f2))
}

func TestIsTautology(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	f1 := p.ParseUnsafe("(a | b) & (c | ~d)")
	f2 := p.ParseUnsafe("(a & b) | (~a & b) | (a & ~b) | (~a & ~b)")

	assert.False(IsTautology(fac, fac.Falsum()))
	assert.True(IsTautology(fac, fac.Verum()))
	assert.False(IsTautology(fac, f1))
	assert.True(IsTautology(fac, f2))
}

func TestIsContradiction(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	f1 := p.ParseUnsafe("(a | b) & (c | ~d)")
	f2 := p.ParseUnsafe("~a & ~b & (a | b)")

	assert.True(IsContradiction(fac, fac.Falsum()))
	assert.False(IsContradiction(fac, fac.Verum()))
	assert.False(IsContradiction(fac, f1))
	assert.True(IsContradiction(fac, f2))
}

func TestImplies(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	f1 := p.ParseUnsafe("(a | b) & (c | ~d)")
	f2 := p.ParseUnsafe("(a | b) & (c | ~d) & (e | ~f)")
	f3 := p.ParseUnsafe("(a | b) & (c | d)")

	assert.False(Implies(fac, f1, f2))
	assert.True(Implies(fac, f2, f1))
	assert.False(Implies(fac, f1, f3))
	assert.False(Implies(fac, f2, f3))
	assert.True(Implies(fac, f2, f2))
}

func TestIsEquivalent(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	f1 := p.ParseUnsafe("(a | b) & (c | ~d)")
	f2 := p.ParseUnsafe("(a | b) & (c | ~d) & (e | ~f)")
	f3 := p.ParseUnsafe("(a & c) | (a & ~d) | (b & c) | (b & ~d)")

	assert.False(IsEquivalent(fac, f1, f2))
	assert.False(IsEquivalent(fac, f2, f1))
	assert.True(IsEquivalent(fac, f1, f3))
	assert.True(IsEquivalent(fac, f3, f1))
	assert.False(IsEquivalent(fac, f2, f3))
}
