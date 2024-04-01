package assignment

import (
	"errors"
	"slices"
	"testing"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/stretchr/testify/assert"
)

func TestAddLiteral(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	a := fac.Lit("a", true)
	b := fac.Lit("b", false)
	ass, _ := New(fac, a)
	assert.True(slices.Contains(ass.PosVars(), f.Variable(a)))

	ass.AddLit(fac, b)
	assert.True(slices.Contains(ass.NegVars(), f.Variable(b)))

	err := ass.AddLit(fac, a.Negate(fac))
	assert.NotNil(err)
	assert.True(errors.Is(err, errorx.ErrBadInput))
	assert.Equal("bad input: ~a (opposite phase present)", err.Error())

	err = ass.AddLit(fac, b.Negate(fac))
	assert.NotNil(err)
	assert.True(errors.Is(err, errorx.ErrBadInput))
	assert.Equal("bad input: b (opposite phase present)", err.Error())
}
