package smus

import (
	"slices"
	"testing"
	"time"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/io"
	"github.com/booleworks/logicng-go/parser"
	"github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

func TestSMUSFromPaper(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	parser := parser.New(fac)

	input := []f.Formula{
		parser.ParseUnsafe("~s"),
		parser.ParseUnsafe("s|~p"),
		parser.ParseUnsafe("p"),
		parser.ParseUnsafe("~p|m"),
		parser.ParseUnsafe("~m|n"),
		parser.ParseUnsafe("~n"),
		parser.ParseUnsafe("~m|l"),
		parser.ParseUnsafe("~l"),
	}
	smus := ComputeForFormulas(fac, input)
	assert.Equal(3, len(smus))
	assert.True(slices.Contains(smus, parser.ParseUnsafe("~s")))
	assert.True(slices.Contains(smus, parser.ParseUnsafe("s|~p")))
	assert.True(slices.Contains(smus, parser.ParseUnsafe("p")))
}

func TestSMUSWithAdditionalConstraint(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	parser := parser.New(fac)

	input := []f.Formula{
		parser.ParseUnsafe("~s"),
		parser.ParseUnsafe("s|~p"),
		parser.ParseUnsafe("p"),
		parser.ParseUnsafe("~p|m"),
		parser.ParseUnsafe("~m|n"),
		parser.ParseUnsafe("~n"),
		parser.ParseUnsafe("~m|l"),
		parser.ParseUnsafe("~l"),
	}
	smus := ComputeForFormulas(fac, input, parser.ParseUnsafe("n|l"))
	assert.Equal(2, len(smus))
	assert.True(slices.Contains(smus, parser.ParseUnsafe("~n")))
	assert.True(slices.Contains(smus, parser.ParseUnsafe("~l")))
}

func TestSMUSSatisfiable(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	parser := parser.New(fac)

	input := []f.Formula{
		parser.ParseUnsafe("~s"),
		parser.ParseUnsafe("s|~p"),
		parser.ParseUnsafe("~p|m"),
		parser.ParseUnsafe("~m|n"),
		parser.ParseUnsafe("~n"),
		parser.ParseUnsafe("~m|l"),
	}
	smus := ComputeForFormulas(fac, input)
	assert.Nil(smus)
}

func TestSMUSUnsatisfiableAdditionalConstraints(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	parser := parser.New(fac)

	input := []f.Formula{
		parser.ParseUnsafe("~s"),
		parser.ParseUnsafe("s|~p"),
		parser.ParseUnsafe("~p|m"),
		parser.ParseUnsafe("~m|n"),
		parser.ParseUnsafe("~n|s"),
	}
	smus := ComputeForFormulas(fac, input, parser.ParseUnsafe("~a&b"), parser.ParseUnsafe("a|~b"))
	assert.Equal(0, len(smus))
}

func TestSMUsTrivialUnsatFormula(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	parser := parser.New(fac)

	input := []f.Formula{
		parser.ParseUnsafe("~s"),
		parser.ParseUnsafe("s|~p"),
		parser.ParseUnsafe("p"),
		parser.ParseUnsafe("~p|m"),
		parser.ParseUnsafe("~m|n"),
		parser.ParseUnsafe("~n"),
		parser.ParseUnsafe("~m|l"),
		parser.ParseUnsafe("~l"),
		parser.ParseUnsafe("a&~a"),
	}
	smus := ComputeForFormulas(fac, input, parser.ParseUnsafe("n|l"))
	assert.Equal(1, len(smus))
	assert.True(slices.Contains(smus, parser.ParseUnsafe("$false")))
}

func TestSMUSUnsatFormula(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	parser := parser.New(fac)

	input := []f.Formula{
		parser.ParseUnsafe("~s"),
		parser.ParseUnsafe("s|~p"),
		parser.ParseUnsafe("p"),
		parser.ParseUnsafe("~p|m"),
		parser.ParseUnsafe("~m|n"),
		parser.ParseUnsafe("~n"),
		parser.ParseUnsafe("~m|l"),
		parser.ParseUnsafe("~l"),
		parser.ParseUnsafe("(a<=>b)&(~a<=>b)"),
	}
	smus := ComputeForFormulas(fac, input, parser.ParseUnsafe("n|l"))
	assert.Equal(1, len(smus))
	assert.True(slices.Contains(smus, parser.ParseUnsafe("(a<=>b)&(~a<=>b)")))
}

func TestSMUSShorterConflict(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	parser := parser.New(fac)

	input := []f.Formula{
		parser.ParseUnsafe("~s"),
		parser.ParseUnsafe("s|~p"),
		parser.ParseUnsafe("p"),
		parser.ParseUnsafe("p&~s"),
		parser.ParseUnsafe("~p|m"),
		parser.ParseUnsafe("~m|n"),
		parser.ParseUnsafe("~n"),
		parser.ParseUnsafe("~m|l"),
		parser.ParseUnsafe("~l"),
	}
	smus := ComputeForFormulas(fac, input)
	assert.Equal(2, len(smus))
	assert.True(slices.Contains(smus, parser.ParseUnsafe("s|~p")))
	assert.True(slices.Contains(smus, parser.ParseUnsafe("p&~s")))
}

func TestSMUSCompleteConflict(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	parser := parser.New(fac)

	input := []f.Formula{
		parser.ParseUnsafe("~s"),
		parser.ParseUnsafe("s|~p"),
		parser.ParseUnsafe("p|~m"),
		parser.ParseUnsafe("m|~n"),
		parser.ParseUnsafe("n|~l"),
		parser.ParseUnsafe("l|s"),
	}
	smus := ComputeForFormulas(fac, input)
	assert.Equal(6, len(smus))
}

func TestSMUSLongConflictWithShortcut(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	parser := parser.New(fac)

	input := []f.Formula{
		parser.ParseUnsafe("~s"),
		parser.ParseUnsafe("s|~p"),
		parser.ParseUnsafe("p|~m"),
		parser.ParseUnsafe("m|~n"),
		parser.ParseUnsafe("n|~l"),
		parser.ParseUnsafe("l|s"),
		parser.ParseUnsafe("n|s"),
	}
	smus := ComputeForFormulas(fac, input)
	assert.Equal(5, len(smus))
	assert.True(slices.Contains(smus, parser.ParseUnsafe("~s")))
	assert.True(slices.Contains(smus, parser.ParseUnsafe("s|~p")))
	assert.True(slices.Contains(smus, parser.ParseUnsafe("p|~m")))
	assert.True(slices.Contains(smus, parser.ParseUnsafe("m|~n")))
	assert.True(slices.Contains(smus, parser.ParseUnsafe("n|s")))
}

func TestSMUSManyConflicts(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	parser := parser.New(fac)

	input := []f.Formula{
		parser.ParseUnsafe("a"),
		parser.ParseUnsafe("~a|b"),
		parser.ParseUnsafe("~b|c"),
		parser.ParseUnsafe("~c|~a"),
		parser.ParseUnsafe("a1"),
		parser.ParseUnsafe("~a1|b1"),
		parser.ParseUnsafe("~b1|c1"),
		parser.ParseUnsafe("~c1|~a1"),
		parser.ParseUnsafe("a2"),
		parser.ParseUnsafe("~a2|b2"),
		parser.ParseUnsafe("~b2|c2"),
		parser.ParseUnsafe("~c2|~a2"),
		parser.ParseUnsafe("a3"),
		parser.ParseUnsafe("~a3|b3"),
		parser.ParseUnsafe("~b3|c3"),
		parser.ParseUnsafe("~c3|~a3"),
		parser.ParseUnsafe("a1|a2|a3|a4|b1|x|y"),
		parser.ParseUnsafe("x&~y"),
		parser.ParseUnsafe("x=>y"),
	}
	smus := ComputeForFormulas(fac, input)
	assert.Equal(2, len(smus))
	assert.True(slices.Contains(smus, parser.ParseUnsafe("x&~y")))
	assert.True(slices.Contains(smus, parser.ParseUnsafe("x=>y")))
}

func TestSMUSTimeoutHandlerSmall(t *testing.T) {
	duration, _ := time.ParseDuration("5s")
	handler := sat.OptimizationHandlerWithTimeout(*handler.NewTimeoutWithDuration(duration))
	fac := f.NewFactory()
	formulas := []f.Formula{fac.Variable("a"), fac.Literal("a", false)}
	testHandler(t, fac, handler, formulas, false)
}

func TestSMUSTimeoutHandlerLarge(t *testing.T) {
	duration, _ := time.ParseDuration("1s")
	handler := sat.OptimizationHandlerWithTimeout(*handler.NewTimeoutWithDuration(duration))
	fac := f.NewFactory()
	formulas, _ := io.ReadFormulas(fac, "../../test/data/formulas/large.txt")
	testHandler(t, fac, handler, formulas, true)
}

func testHandler(t *testing.T, fac f.Factory, handler sat.OptimizationHandler, formulas []f.Formula, exp bool) {
	assert := assert.New(t)
	result, ok := ComputeForFormulasWithHandler(fac, formulas, handler)
	if exp {
		assert.True(handler.Aborted())
		assert.False(ok)
		assert.Nil(result)
	} else {
		assert.False(handler.Aborted())
		assert.True(ok)
		assert.NotNil(result)
	}
}
