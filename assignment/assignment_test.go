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

func TestSprintDeterministic(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()

	ass1, _ := New(fac,
		fac.Lit("z", true),
		fac.Lit("a", false),
		fac.Lit("m", true),
		fac.Lit("b", false),
		fac.Lit("x", false),
	)

	ass2, _ := New(fac,
		fac.Lit("b", false),
		fac.Lit("x", false),
		fac.Lit("a", false),
		fac.Lit("z", true),
		fac.Lit("m", true),
	)

	expected := "[~a, ~b, m, ~x, z]"
	assert.Equal(expected, ass1.Sprint(fac))
	assert.Equal(expected, ass2.Sprint(fac))

	for range 10 {
		assert.Equal(expected, ass1.Sprint(fac))
		assert.Equal(expected, ass2.Sprint(fac))
	}

	emptyAss := Empty()
	assert.Equal("[]", emptyAss.Sprint(fac))

	singlePos, _ := New(fac, fac.Lit("var", true))
	assert.Equal("[var]", singlePos.Sprint(fac))

	singleNeg, _ := New(fac, fac.Lit("var", false))
	assert.Equal("[~var]", singleNeg.Sprint(fac))

	ass3, _ := New(fac,
		fac.Lit("var10", true),
		fac.Lit("var2", false),
		fac.Lit("var1", true),
		fac.Lit("a_variable", false),
	)
	assert.Equal("[~a_variable, var1, var10, ~var2]", ass3.Sprint(fac))
}
