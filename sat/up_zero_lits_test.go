package sat

import (
	"slices"
	"testing"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/parser"
	"github.com/stretchr/testify/assert"
)

func TestUpZerLiterals(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	solver := NewSolver(fac)
	formula := p.ParseUnsafe("(b | c) & ~b & (~c | d)")
	solver.Add(formula)
	sat := solver.Sat()
	assert.True(sat)
	upZeroLits, err := solver.UpZeroLits()
	assert.Nil(err)
	assert.Equal(3, len(upZeroLits))
	assert.True(slices.Contains(upZeroLits, fac.Lit("b", false)))
	assert.True(slices.Contains(upZeroLits, fac.Lit("c", true)))
	assert.True(slices.Contains(upZeroLits, fac.Lit("d", true)))
}
