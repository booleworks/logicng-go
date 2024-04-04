package iter

import (
	"strings"
	"testing"

	"github.com/booleworks/logicng-go/sat"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/parser"
	"github.com/stretchr/testify/assert"
)

func TestFixedVariableProvider(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	vars := NewFixedVarProvider(varSet(fac, "a b c d e f")).Vars(nil, nil)
	assert.Equal(6, vars.Size())
	assert.True(vars.Contains(fac.Var("a")))
	assert.True(vars.Contains(fac.Var("b")))
	assert.True(vars.Contains(fac.Var("c")))
	assert.True(vars.Contains(fac.Var("d")))
	assert.True(vars.Contains(fac.Var("e")))
	assert.True(vars.Contains(fac.Var("f")))

	vars = NewFixedVarProvider(varSet(fac, "a b c d e f")).Vars(nil, f.NewVarSet())
	assert.Equal(6, vars.Size())
	assert.True(vars.Contains(fac.Var("a")))
	assert.True(vars.Contains(fac.Var("b")))
	assert.True(vars.Contains(fac.Var("c")))
	assert.True(vars.Contains(fac.Var("d")))
	assert.True(vars.Contains(fac.Var("e")))
	assert.True(vars.Contains(fac.Var("f")))

	vars = NewFixedVarProvider(varSet(fac, "a b c d e f")).Vars(nil, f.NewVarSet(fac.Var("a"), fac.Var("b")))
	assert.Equal(6, vars.Size())
	assert.True(vars.Contains(fac.Var("a")))
	assert.True(vars.Contains(fac.Var("b")))
	assert.True(vars.Contains(fac.Var("c")))
	assert.True(vars.Contains(fac.Var("d")))
	assert.True(vars.Contains(fac.Var("e")))
	assert.True(vars.Contains(fac.Var("f")))
}

func TestLeastCommonVarsProvider(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	solver := sat.NewSolver(fac)
	a := fac.Var("a")
	b := fac.Var("b")
	c := fac.Var("c")
	d := fac.Var("d")
	e := fac.Var("e")
	x := fac.Var("x")
	g := fac.Var("g")
	h := fac.Var("h")
	i := fac.Var("i")
	j := fac.Var("j")
	varSet := f.NewVarSet(a, b, c, d, e, x, g, h, i, j)

	solver.Add(p.ParseUnsafe("(a | b | c) & (~b | c) & (d | ~e) & (~a | e) & (a | d | b | g | h) & (~h | i) & (x | g | j) & (x | b | j | ~g) & (g | c)"))

	vars := NewLeastCommonVarProvider(.1, 100).Vars(solver, nil)
	assert.Equal(1, vars.Size())
	assert.True(vars.Contains(i))

	vars = NewLeastCommonVarProvider(.1, 100).Vars(solver, varSet)
	assert.Equal(1, vars.Size())
	assert.True(vars.Contains(i))

	vars = NewLeastCommonVarProvider(.0001, 100).Vars(solver, varSet)
	assert.Equal(1, vars.Size())
	assert.True(vars.Contains(i))

	vars = NewLeastCommonVarProvider(.2, 100).Vars(solver, nil)
	assert.Equal(2, vars.Size())
	assert.True(vars.Contains(i))

	vars = NewLeastCommonVarProvider(.6, 100).Vars(solver, nil)
	assert.Equal(6, vars.Size())
	assert.True(vars.Contains(e))
	assert.True(vars.Contains(d))
	assert.True(vars.Contains(x))
	assert.True(vars.Contains(i))
	assert.True(vars.Contains(h))
	assert.True(vars.Contains(j))

	vars = NewLeastCommonVarProvider(.6, 1).Vars(solver, nil)
	assert.Equal(1, vars.Size())
	assert.True(vars.Contains(i))

	vars = NewLeastCommonVarProvider(.6, 2).Vars(solver, nil)
	assert.Equal(2, vars.Size())
	assert.True(vars.Contains(i))

	vars = NewLeastCommonVarProvider(.25, 100).Vars(solver, f.NewVarSet(a, b, g))
	assert.Equal(1, vars.Size())
	assert.True(vars.Contains(a))

	vars = NewLeastCommonVarProvider(.5, 100).Vars(solver, f.NewVarSet(a, c, b, g))
	assert.Equal(2, vars.Size())
	assert.True(vars.Contains(a))
	assert.True(vars.Contains(c))

	vars = NewLeastCommonVarProvider(1, 100).Vars(solver, f.NewVarSet(a, c, b, g))
	assert.Equal(4, vars.Size())

	vars = NewLeastCommonVarProvider(1, 100).Vars(solver, nil)
	assert.Equal(10, vars.Size())

	vars = DefaultLeastCommonVarProvider().Vars(solver, f.NewVarSet(a, c, b, g))
	assert.Equal(2, vars.Size())
	assert.True(vars.Contains(a))
	assert.True(vars.Contains(c))
}

func TestMostCommonVarsProvider(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	solver := sat.NewSolver(fac)
	a := fac.Var("a")
	b := fac.Var("b")
	c := fac.Var("c")
	d := fac.Var("d")
	e := fac.Var("e")
	x := fac.Var("x")
	g := fac.Var("g")
	h := fac.Var("h")
	i := fac.Var("i")
	j := fac.Var("j")
	varSet := f.NewVarSet(a, b, c, d, e, x, g, h, i, j)

	solver.Add(p.ParseUnsafe("(a | b | c) & (~b | c) & (d | ~e) & (~a | e) & (a | d | b | g | h) & (~h | i) & (x | g | j) & (x | b | j | ~g) & (g | c)"))

	vars := NewMostCommonVarProvider(.1, 100).Vars(solver, nil)
	assert.Equal(1, vars.Size())
	assert.True(vars.Contains(b) || vars.Contains(g))

	vars = NewMostCommonVarProvider(.1, 100).Vars(solver, varSet)
	assert.Equal(1, vars.Size())
	assert.True(vars.Contains(b) || vars.Contains(g))

	vars = NewMostCommonVarProvider(.0001, 100).Vars(solver, varSet)
	assert.Equal(1, vars.Size())
	assert.True(vars.Contains(b) || vars.Contains(g))

	vars = NewMostCommonVarProvider(.2, 100).Vars(solver, nil)
	assert.Equal(2, vars.Size())
	assert.True(vars.Contains(b))
	assert.True(vars.Contains(g))

	vars = NewMostCommonVarProvider(.4, 100).Vars(solver, nil)
	assert.Equal(4, vars.Size())
	assert.True(vars.Contains(b))
	assert.True(vars.Contains(g))
	assert.True(vars.Contains(a))
	assert.True(vars.Contains(c))

	vars = NewMostCommonVarProvider(.9, 2).Vars(solver, nil)
	assert.Equal(2, vars.Size())
	assert.True(vars.Contains(b))
	assert.True(vars.Contains(g))

	vars = NewMostCommonVarProvider(.25, 100).Vars(solver, f.NewVarSet(x, i, c))
	assert.Equal(1, vars.Size())
	assert.True(vars.Contains(c))

	vars = NewMostCommonVarProvider(.5, 100).Vars(solver, f.NewVarSet(c, b, x, h))
	assert.Equal(2, vars.Size())
	assert.True(vars.Contains(b))
	assert.True(vars.Contains(c))

	vars = NewMostCommonVarProvider(1, 100).Vars(solver, f.NewVarSet(a, c, b, g))
	assert.Equal(4, vars.Size())

	vars = NewMostCommonVarProvider(1, 100).Vars(solver, nil)
	assert.Equal(10, vars.Size())

	vars = DefaultMostCommonVarProvider().Vars(solver, f.NewVarSet(c, b, x, h))
	assert.Equal(2, vars.Size())
	assert.True(vars.Contains(b))
	assert.True(vars.Contains(c))
}

func varSet(fac f.Factory, varString string) *f.VarSet {
	tokens := strings.Split(varString, " ")
	vars := f.NewMutableVarSet()
	for _, t := range tokens {
		vars.Add(fac.Var(t))
	}
	return vars.AsImmutable()
}
