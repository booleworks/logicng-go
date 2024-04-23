package test

import (
	"testing"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/parser"
	"github.com/stretchr/testify/assert"
)

func TestProfileConstants(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()

	assert.Equal(map[f.Variable]int{}, f.VariableProfile(fac, fac.Verum()))
	assert.Equal(map[f.Variable]int{}, f.VariableProfile(fac, fac.Falsum()))
	assert.Equal(map[f.Literal]int{}, f.LiteralProfile(fac, fac.Verum()))
	assert.Equal(map[f.Literal]int{}, f.LiteralProfile(fac, fac.Falsum()))
}

func TestProfileLiterals(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	expectedVars := map[f.Variable]int{fac.Var("a"): 1}

	assert.Equal(expectedVars, f.VariableProfile(fac, p.ParseUnsafe("a")))
	assert.Equal(expectedVars, f.VariableProfile(fac, p.ParseUnsafe("~a")))
	expectedLits := map[f.Literal]int{fac.Lit("a", true): 1}
	assert.Equal(expectedLits, f.LiteralProfile(fac, p.ParseUnsafe("a")))
	expectedLits = map[f.Literal]int{fac.Lit("a", false): 1}
	assert.Equal(expectedLits, f.LiteralProfile(fac, p.ParseUnsafe("~a")))
}

func TestProfileNot(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	expectedVars := make(map[f.Variable]int, 3)
	expectedVars[fac.Var("a")] = 1
	expectedVars[fac.Var("b")] = 2
	expectedVars[fac.Var("c")] = 3

	expectedLits := make(map[f.Literal]int, 3)
	expectedLits[fac.Lit("a", true)] = 1
	expectedLits[fac.Lit("b", true)] = 1
	expectedLits[fac.Lit("c", true)] = 2
	expectedLits[fac.Lit("b", false)] = 1
	expectedLits[fac.Lit("c", false)] = 1

	formula := p.ParseUnsafe("~(a & (b | c) & ((~b | ~c) => c))")
	assert.Equal(expectedVars, f.VariableProfile(fac, formula))
	assert.Equal(expectedLits, f.LiteralProfile(fac, formula))
}

func TestProfileBinaryOperator(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	expectedVars := make(map[f.Variable]int, 3)
	expectedVars[fac.Var("a")] = 1
	expectedVars[fac.Var("b")] = 2
	expectedVars[fac.Var("c")] = 3

	expectedLits := make(map[f.Literal]int, 3)
	expectedLits[fac.Lit("a", true)] = 1
	expectedLits[fac.Lit("b", true)] = 1
	expectedLits[fac.Lit("c", true)] = 2
	expectedLits[fac.Lit("b", false)] = 1
	expectedLits[fac.Lit("c", false)] = 1

	impl := p.ParseUnsafe("(a & (b | c) & (~b | ~c)) => c")
	equiv := p.ParseUnsafe("(a & (b | c) & (~b | ~c)) <=> c")

	assert.Equal(expectedVars, f.VariableProfile(fac, impl))
	assert.Equal(expectedVars, f.VariableProfile(fac, equiv))
	assert.Equal(expectedLits, f.LiteralProfile(fac, impl))
	assert.Equal(expectedLits, f.LiteralProfile(fac, equiv))
}

func TestProfileNAryOperator(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	expectedVars := make(map[f.Variable]int, 3)
	expectedVars[fac.Var("a")] = 1
	expectedVars[fac.Var("b")] = 2
	expectedVars[fac.Var("c")] = 3

	expectedLits := make(map[f.Literal]int, 3)
	expectedLits[fac.Lit("a", true)] = 1
	expectedLits[fac.Lit("b", true)] = 1
	expectedLits[fac.Lit("c", true)] = 2
	expectedLits[fac.Lit("b", false)] = 1
	expectedLits[fac.Lit("c", false)] = 1

	formula := p.ParseUnsafe("a & (b | c) & (~b | ~c) & c")
	assert.Equal(expectedVars, f.VariableProfile(fac, formula))
	assert.Equal(expectedLits, f.LiteralProfile(fac, formula))
}

func TestProfilePbc(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	expectedVars1 := make(map[f.Variable]int, 3)
	expectedVars1[fac.Var("a")] = 1
	expectedVars2 := make(map[f.Variable]int, 3)
	expectedVars2[fac.Var("a")] = 1
	expectedVars2[fac.Var("b")] = 1
	expectedVars2[fac.Var("c")] = 1

	expectedVars1L := make(map[f.Literal]int, 3)
	expectedVars1L[fac.Lit("a", true)] = 1
	expectedVars2L := make(map[f.Literal]int, 3)
	expectedVars2L[fac.Lit("a", true)] = 1
	expectedVars2L[fac.Lit("b", true)] = 1
	expectedVars2L[fac.Lit("c", true)] = 1

	expectedLits1 := make(map[f.Literal]int, 3)
	expectedLits1[fac.Lit("a", false)] = 1

	expectedLits2 := make(map[f.Literal]int, 3)
	expectedLits2[fac.Lit("a", true)] = 1
	expectedLits2[fac.Lit("c", true)] = 1
	expectedLits2[fac.Lit("b", false)] = 1

	pb1 := p.ParseUnsafe("3*~a <= 2")
	pb2 := p.ParseUnsafe("3*a + -2*b + 7*c <= 8")
	cc1 := p.ParseUnsafe("a < 1")
	cc2 := p.ParseUnsafe("a + b + c > 2")
	amo := p.ParseUnsafe("a + ~b + c <= 1")

	assert.Equal(expectedVars1, f.VariableProfile(fac, pb1))
	assert.Equal(expectedLits1, f.LiteralProfile(fac, pb1))
	assert.Equal(expectedVars2, f.VariableProfile(fac, pb2))
	assert.Equal(expectedVars2L, f.LiteralProfile(fac, pb2))
	assert.Equal(expectedVars1, f.VariableProfile(fac, cc1))
	assert.Equal(expectedVars1L, f.LiteralProfile(fac, cc1))
	assert.Equal(expectedVars2, f.VariableProfile(fac, cc2))
	assert.Equal(expectedVars2L, f.LiteralProfile(fac, cc2))
	assert.Equal(expectedVars2, f.VariableProfile(fac, amo))
	assert.Equal(expectedLits2, f.LiteralProfile(fac, amo))
}
