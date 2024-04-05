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
	result := solver.Call(Params().UpZeroIfSat())
	assert.True(result.OK())
	assert.Equal(3, len(result.UpZeroLits()))
	assert.True(slices.Contains(result.UpZeroLits(), fac.Lit("b", false)))
	assert.True(slices.Contains(result.UpZeroLits(), fac.Lit("c", true)))
	assert.True(slices.Contains(result.UpZeroLits(), fac.Lit("d", true)))
}
