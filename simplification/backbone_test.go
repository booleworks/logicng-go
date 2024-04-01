package simplification

import (
	"testing"

	"github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/parser"
	"github.com/stretchr/testify/assert"
)

func TestBackboneSimplifierTrivial(t *testing.T) {
	assert := assert.New(t)
	fac := formula.NewFactory()
	p := parser.New(fac)
	assert.Equal(p.ParseUnsafe("$true"), SimplifyWithBackbone(fac, p.ParseUnsafe("$true")))
	assert.Equal(p.ParseUnsafe("$false"), SimplifyWithBackbone(fac, p.ParseUnsafe("$false")))
	assert.Equal(p.ParseUnsafe("$false"), SimplifyWithBackbone(fac, p.ParseUnsafe("A & (A => B) & ~B")))
	assert.Equal(p.ParseUnsafe("A"), SimplifyWithBackbone(fac, p.ParseUnsafe("A")))
	assert.Equal(p.ParseUnsafe("A & B"), SimplifyWithBackbone(fac, p.ParseUnsafe("A & B")))
	assert.Equal(p.ParseUnsafe("A | B | C"), SimplifyWithBackbone(fac, p.ParseUnsafe("A | B | C")))
}

func TestBackboneSimplifierReal(t *testing.T) {
	assert := assert.New(t)
	fac := formula.NewFactory()
	p := parser.New(fac)
	assert.Equal(p.ParseUnsafe("A & B"), SimplifyWithBackbone(fac, p.ParseUnsafe("A & B & (B | C)")))
	assert.Equal(p.ParseUnsafe("A & B & C"), SimplifyWithBackbone(fac, p.ParseUnsafe("A & B & (~B | C)")))
	assert.Equal(
		p.ParseUnsafe("A & B & C & F"),
		SimplifyWithBackbone(fac, p.ParseUnsafe("A & B & (~B | C) & (B | D) & (A => F)")),
	)
	assert.Equal(
		p.ParseUnsafe("X & Y & (~B | C) & (B | D) & (A => F)"),
		SimplifyWithBackbone(fac, p.ParseUnsafe("X & Y & (~B | C) & (B | D) & (A => F)")),
	)
	assert.Equal(
		p.ParseUnsafe("D & ~A & ~B"),
		SimplifyWithBackbone(fac, p.ParseUnsafe("~A & ~B & (~B | C) & (B | D) & (A => F)")),
	)
}
