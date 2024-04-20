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
	assert.Equal(p.ParseUnsafe("$true"), PropagateBackbone(fac, p.ParseUnsafe("$true")))
	assert.Equal(p.ParseUnsafe("$false"), PropagateBackbone(fac, p.ParseUnsafe("$false")))
	assert.Equal(p.ParseUnsafe("$false"), PropagateBackbone(fac, p.ParseUnsafe("A & (A => B) & ~B")))
	assert.Equal(p.ParseUnsafe("A"), PropagateBackbone(fac, p.ParseUnsafe("A")))
	assert.Equal(p.ParseUnsafe("A & B"), PropagateBackbone(fac, p.ParseUnsafe("A & B")))
	assert.Equal(p.ParseUnsafe("A | B | C"), PropagateBackbone(fac, p.ParseUnsafe("A | B | C")))
}

func TestBackboneSimplifierReal(t *testing.T) {
	assert := assert.New(t)
	fac := formula.NewFactory()
	p := parser.New(fac)
	assert.Equal(p.ParseUnsafe("A & B"), PropagateBackbone(fac, p.ParseUnsafe("A & B & (B | C)")))
	assert.Equal(p.ParseUnsafe("A & B & C"), PropagateBackbone(fac, p.ParseUnsafe("A & B & (~B | C)")))
	assert.Equal(
		p.ParseUnsafe("A & B & C & F"),
		PropagateBackbone(fac, p.ParseUnsafe("A & B & (~B | C) & (B | D) & (A => F)")),
	)
	assert.Equal(
		p.ParseUnsafe("X & Y & (~B | C) & (B | D) & (A => F)"),
		PropagateBackbone(fac, p.ParseUnsafe("X & Y & (~B | C) & (B | D) & (A => F)")),
	)
	assert.Equal(
		p.ParseUnsafe("D & ~A & ~B"),
		PropagateBackbone(fac, p.ParseUnsafe("~A & ~B & (~B | C) & (B | D) & (A => F)")),
	)
}
