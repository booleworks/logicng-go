package sat

import (
	"testing"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/parser"
	"github.com/stretchr/testify/assert"
)

func TestVarOccsOnEmptySolver(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	a := fac.Var("a")
	b := fac.Var("b")
	c := fac.Var("c")

	solver := NewSolver(fac)
	assert.Equal(0, len(solver.VarOccurrences(nil)))

	count := solver.VarOccurrences(f.NewVarSet(a, b, c))
	assert.Equal(3, len(count))
	cnt, ok := count[a]
	assert.True(ok)
	assert.Equal(0, cnt)
	cnt, ok = count[b]
	assert.True(ok)
	assert.Equal(0, cnt)
	cnt, ok = count[c]
	assert.True(ok)
	assert.Equal(0, cnt)
}

func TestVarOccsOnSolver(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	a := fac.Var("a")
	b := fac.Var("b")
	c := fac.Var("c")
	d := fac.Var("d")
	e := fac.Var("e")
	g := fac.Var("g")
	h := fac.Var("h")
	i := fac.Var("i")
	x := fac.Var("x")
	y := fac.Var("y")

	solver := NewSolver(fac)
	solver.Add(p.ParseUnsafe("(a | b | c) & (~b | c) & (d | ~e) & x & (~a | e) & (a | d | b | g | h) & (~h | i) & y"))
	count := solver.VarOccurrences(nil)
	assert.Equal(10, len(count))
	assert.Equal(3, count[a])
	assert.Equal(3, count[b])
	assert.Equal(2, count[c])
	assert.Equal(2, count[d])
	assert.Equal(2, count[e])
	assert.Equal(1, count[g])
	assert.Equal(2, count[h])
	assert.Equal(1, count[i])
	assert.Equal(1, count[x])
	assert.Equal(1, count[y])
}

func TestVarOccsOnSolverWithRelevant(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	a := fac.Var("a")
	c := fac.Var("c")
	x := fac.Var("x")
	j := fac.Var("j")

	solver := NewSolver(fac)
	solver.Add(p.ParseUnsafe("(a | b | c) & (~b | c) & (d | ~e) & x & (~a | e) & (a | d | b | g | h) & (~h | i) & y"))
	count := solver.VarOccurrences(f.NewVarSet(a, c, x, j))
	assert.Equal(4, len(count))
	assert.Equal(3, count[a])
	assert.Equal(2, count[c])
	assert.Equal(1, count[x])
	assert.Equal(0, count[j])
}
