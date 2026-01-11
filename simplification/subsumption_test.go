package simplification

import (
	"testing"
	"time"

	"github.com/booleworks/logicng-go/event"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/io"
	"github.com/booleworks/logicng-go/normalform"
	"github.com/booleworks/logicng-go/parser"
	"github.com/booleworks/logicng-go/sat"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/stretchr/testify/assert"
)

func TestUbTreeSingleSet(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	a := fac.Lit("A", true)
	b := fac.Lit("B", true)
	c := fac.Lit("C", true)

	tree := newUbtree()
	tree.addSet(f.NewLitSet(a, b, c))

	assert.Equal(1, tree.rootNodes.Size())
	_tt, ok := tree.rootNodes.Get(a)
	node := _tt.(*ubnode)
	assert.True(ok)
	assert.Equal(1, node.children.Size())
	assert.False(node.isEndOfPath())

	_tt, ok = node.children.Get(b)
	node = _tt.(*ubnode)
	assert.True(ok)
	assert.Equal(1, node.children.Size())
	assert.False(node.isEndOfPath())

	_tt, ok = node.children.Get(c)
	node = _tt.(*ubnode)
	assert.True(ok)
	assert.Equal(0, node.children.Size())
	assert.True(node.isEndOfPath())
}

func TestUbTreeFromPaper(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	e0 := fac.Lit("e0", true)
	e1 := fac.Lit("e1", true)
	e2 := fac.Lit("e2", true)
	e3 := fac.Lit("e3", true)

	tree := newUbtree()
	tree.addSet(f.NewLitSet(e0, e1, e2, e3))
	tree.addSet(f.NewLitSet(e0, e1, e3))
	tree.addSet(f.NewLitSet(e0, e1, e2))
	tree.addSet(f.NewLitSet(e2, e3))

	assert.Equal(2, tree.rootNodes.Size())
	_tt, ok := tree.rootNodes.Get(e0)
	nodeE0 := _tt.(*ubnode)
	assert.True(ok)
	_tt, ok = tree.rootNodes.Get(e2)
	nodeE2 := _tt.(*ubnode)
	assert.True(ok)

	// root nodes
	assert.Equal(1, nodeE0.children.Size())
	assert.False(nodeE0.isEndOfPath())
	assert.Equal(1, nodeE2.children.Size())
	assert.False(nodeE2.isEndOfPath())

	// first level
	_tt, ok = nodeE0.children.Get(e1)
	assert.True(ok)
	e0e1 := _tt.(*ubnode)
	_tt, ok = nodeE2.children.Get(e3)
	assert.True(ok)
	e2e3 := _tt.(*ubnode)

	assert.Equal(2, e0e1.children.Size())
	assert.False(e0e1.isEndOfPath())
	assert.Equal(0, e2e3.children.Size())
	assert.True(e2e3.isEndOfPath())

	// second level
	_tt, ok = e0e1.children.Get(e2)
	assert.True(ok)
	e0e1e2 := _tt.(*ubnode)
	_tt, ok = e0e1.children.Get(e3)
	assert.True(ok)
	e0e1e3 := _tt.(*ubnode)

	assert.True(e0e1e2.isEndOfPath())
	assert.Equal(1, e0e1e2.children.Size())
	assert.True(e0e1e3.isEndOfPath())
	assert.Equal(0, e0e1e3.children.Size())

	// third level
	_tt, ok = e0e1e2.children.Get(e3)
	assert.True(ok)
	e0e1e2e3 := _tt.(*ubnode)
	assert.True(e0e1e2e3.isEndOfPath())
	assert.Equal(0, e0e1e2e3.children.Size())
}

func TestUbTreeContainsSubset(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	e0 := fac.Lit("e0", true)
	e1 := fac.Lit("e1", true)
	e2 := fac.Lit("e2", true)
	e3 := fac.Lit("e3", true)
	e4 := fac.Lit("e4", true)

	e0123 := f.NewLitSet(e0, e1, e2, e3)
	e013 := f.NewLitSet(e0, e1, e3)
	e012 := f.NewLitSet(e0, e1, e2)
	e23 := f.NewLitSet(e2, e3)

	tree := newUbtree()
	tree.addSet(e0123)
	tree.addSet(e013)
	tree.addSet(e012)
	tree.addSet(e23)

	assert.Nil(tree.firstSubset(set(e0)))
	assert.Nil(tree.firstSubset(set(e1)))
	assert.Nil(tree.firstSubset(set(e2)))
	assert.Nil(tree.firstSubset(set(e3)))

	assert.Nil(tree.firstSubset(set(e0, e1)))
	assert.Nil(tree.firstSubset(set(e0, e2)))
	assert.Nil(tree.firstSubset(set(e0, e3)))
	assert.Nil(tree.firstSubset(set(e1, e2)))
	assert.Nil(tree.firstSubset(set(e1, e3)))

	equalSubset(t, e23, tree.firstSubset(set(e2, e3)))
	equalSubset(t, e012, tree.firstSubset(set(e0, e1, e2)))
	equalSubset(t, e013, tree.firstSubset(set(e0, e1, e3)))
	equalSubset(t, e23, tree.firstSubset(set(e0, e2, e3)))
	equalSubset(t, e23, tree.firstSubset(set(e1, e2, e3)))
	assert.NotNil(tree.firstSubset(set(e0, e1, e2, e3)))

	assert.Nil(tree.firstSubset(set(e0, e4)))
	assert.Nil(tree.firstSubset(set(e1, e4)))
	assert.Nil(tree.firstSubset(set(e2, e4)))
	assert.Nil(tree.firstSubset(set(e3, e4)))

	assert.Nil(tree.firstSubset(set(e0, e1, e4)))
	assert.Nil(tree.firstSubset(set(e0, e2, e4)))
	assert.Nil(tree.firstSubset(set(e0, e3, e4)))
	assert.Nil(tree.firstSubset(set(e1, e2, e4)))
	assert.Nil(tree.firstSubset(set(e1, e3, e4)))

	equalSubset(t, e23, tree.firstSubset(set(e2, e3, e4)))
	equalSubset(t, e23, tree.firstSubset(set(e2, e3, e4)))
	equalSubset(t, e012, tree.firstSubset(set(e0, e1, e2, e4)))
	equalSubset(t, e013, tree.firstSubset(set(e0, e1, e3, e4)))
	equalSubset(t, e23, tree.firstSubset(set(e0, e2, e3, e4)))
	equalSubset(t, e23, tree.firstSubset(set(e1, e2, e3, e4)))
	assert.NotNil(tree.firstSubset(set(e0, e1, e2, e3, e4)))
}

func TestUbTreeAllSets(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	e0 := fac.Lit("e0", true)
	e1 := fac.Lit("e1", true)
	e2 := fac.Lit("e2", true)
	e3 := fac.Lit("e3", true)

	e0123 := f.NewLitSet(e0, e1, e2, e3)
	e013 := f.NewLitSet(e0, e1, e3)
	e012 := f.NewLitSet(e0, e1, e2)
	e23 := f.NewLitSet(e2, e3)

	tree := newUbtree()
	tree.addSet(e0123)
	assert.Equal(1, tree.allSets().Size())
	tree.addSet(e013)
	assert.Equal(2, tree.allSets().Size())
	tree.addSet(e012)
	assert.Equal(3, tree.allSets().Size())
	tree.addSet(e23)
	assert.Equal(4, tree.allSets().Size())
}

func TestCnfSubsumptionSimple(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	assert.Equal(p.ParseUnsafe("$false"), cnfs(t, fac, p.ParseUnsafe("$false")))
	assert.Equal(p.ParseUnsafe("$true"), cnfs(t, fac, p.ParseUnsafe("$true")))
	assert.Equal(p.ParseUnsafe("a"), cnfs(t, fac, p.ParseUnsafe("a")))
	assert.Equal(p.ParseUnsafe("~a"), cnfs(t, fac, p.ParseUnsafe("~a")))
	assert.Equal(p.ParseUnsafe("a | b | c"), cnfs(t, fac, p.ParseUnsafe("a | b | c")))
	assert.Equal(p.ParseUnsafe("a & b & c"), cnfs(t, fac, p.ParseUnsafe("a & b & c")))
	assert.Equal(p.ParseUnsafe("a"), cnfs(t, fac, p.ParseUnsafe("a & (a | b)")))
	assert.Equal(p.ParseUnsafe("a | b"), cnfs(t, fac, p.ParseUnsafe("(a | b) & (a | b | c)")))
	assert.Equal(p.ParseUnsafe("a"), cnfs(t, fac, p.ParseUnsafe("a & (a | b) & (a | b | c)")))
	assert.Equal(p.ParseUnsafe("a & b"), cnfs(t, fac, p.ParseUnsafe("a & (a | b) & b")))
	assert.Equal(p.ParseUnsafe("a & c"), cnfs(t, fac, p.ParseUnsafe("a & (a | b) & c & (c | b)")))
	assert.Equal(p.ParseUnsafe("(a | b) & (a | c)"), cnfs(t, fac, p.ParseUnsafe("(a | b) & (a | c) & (a | b | c)")))
}

func TestDnfSubsumptionSimple(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	assert.Equal(p.ParseUnsafe("$false"), dnfs(t, fac, p.ParseUnsafe("$false")))
	assert.Equal(p.ParseUnsafe("$true"), dnfs(t, fac, p.ParseUnsafe("$true")))
	assert.Equal(p.ParseUnsafe("a"), dnfs(t, fac, p.ParseUnsafe("a")))
	assert.Equal(p.ParseUnsafe("~a"), dnfs(t, fac, p.ParseUnsafe("~a")))
	assert.Equal(p.ParseUnsafe("a | b | c"), dnfs(t, fac, p.ParseUnsafe("a | b | c")))
	assert.Equal(p.ParseUnsafe("a & b & c"), dnfs(t, fac, p.ParseUnsafe("a & b & c")))
	assert.Equal(p.ParseUnsafe("a"), dnfs(t, fac, p.ParseUnsafe("a | (a & b)")))
	assert.Equal(p.ParseUnsafe("a & b"), dnfs(t, fac, p.ParseUnsafe("(a & b) | (a & b & c)")))
	assert.Equal(p.ParseUnsafe("a"), dnfs(t, fac, p.ParseUnsafe("a | (a & b) | (a & b & c)")))
	assert.Equal(p.ParseUnsafe("a | b"), dnfs(t, fac, p.ParseUnsafe("a | (a & b) | b")))
	assert.Equal(p.ParseUnsafe("a | c"), dnfs(t, fac, p.ParseUnsafe("a | (a & b) | c | (c & b)")))
	assert.Equal(p.ParseUnsafe("(a & b) | (a & c)"), dnfs(t, fac, p.ParseUnsafe("(a & b) | (a & c) | (a & b & c)")))
}

func TestCnfSubsumptionMid(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	assert.Equal(
		p.ParseUnsafe("(a | b | c)"),
		cnfs(t, fac, p.ParseUnsafe("(a | b | c | d) & (a | b | c | e) & (a | b | c)")),
	)
	assert.Equal(
		p.ParseUnsafe("(a | b) & (a | c) & (b | c)"),
		cnfs(t, fac, p.ParseUnsafe("(a | b) & (a | c) & (a | b | c) & (a | ~b | c) & (a | b | ~c) & (b | c)")),
	)
	assert.Equal(
		p.ParseUnsafe("(a | b) & (a | c) & (b | c)"),
		cnfs(t, fac, p.ParseUnsafe("(a | b) & (a | c) & (a | b | c) & (a | ~b | c) & (a | b | ~c) & (b | c)")),
	)
	assert.Equal(
		p.ParseUnsafe("a & ~b & (c | d)"),
		cnfs(t, fac, p.ParseUnsafe("a & ~b & (c | d) & (~a | ~b | ~c) & (b | c | d) & (a | b | c | d)")),
	)
	assert.Equal(
		p.ParseUnsafe("(a | c | e | g) & (b | d | f)"),
		cnfs(t, fac, p.ParseUnsafe("(a | b | c | d | e | f | g) & (b | d | f) & (a | c | e | g)")),
	)
}

func TestDnfSubsumptionMid(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	assert.Equal(
		p.ParseUnsafe("(a & b & c)"),
		dnfs(t, fac, p.ParseUnsafe("(a & b & c & d) | (a & b & c & e) | (a & b & c)")),
	)
	assert.Equal(
		p.ParseUnsafe("(a & b) | (a & c) | (b & c)"),
		dnfs(t, fac, p.ParseUnsafe("(a & b) | (a & c) | (a & b & c) | (a & ~b & c) | (a & b & ~c) | (b & c)")),
	)
	assert.Equal(
		p.ParseUnsafe("(a & b) | (a & c) | (b & c)"),
		dnfs(t, fac, p.ParseUnsafe("(a & b) | (a & c) | (a & b & c) | (a & ~b & c) | (a & b & ~c) | (b & c)")),
	)
	assert.Equal(
		p.ParseUnsafe("a | ~b | (c & d)"),
		dnfs(t, fac, p.ParseUnsafe("a | ~b | (c & d) | (~a & ~b & ~c) | (b & c & d) | (a & b & c & d)")),
	)
	assert.Equal(
		p.ParseUnsafe("(a & c & e & g) | (b & d & f)"),
		dnfs(t, fac, p.ParseUnsafe("(a & b & c & d & e & f & g) | (b & d & f) | (a & c & e & g)")),
	)
}

func TestCnfSubsumptionLarge(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	formula, _ := io.ReadFormula(fac, "./../test/data/formulas/large2.txt")
	cnf := normalform.FactorizedCNF(fac, formula)
	subsumed := cnfs(t, fac, cnf)
	assert.True(sat.IsEquivalent(fac, cnf, formula))
	assert.Greater(len(fac.Operands(cnf)), len(fac.Operands(subsumed)))
}

func TestCnfSubsumptionWithHandler(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	formula, _ := io.ReadFormula(fac, "./../test/data/formulas/large2.txt")
	cnf := normalform.FactorizedCNF(fac, formula)
	duration, _ := time.ParseDuration("5ms")
	hdl := handler.NewTimeoutWithDuration(duration)
	subsumed, err, state := CNFSubsumptionWithHandler(fac, cnf, hdl)
	assert.Nil(err)
	assert.False(state.Success)
	assert.NotEqual(event.Nothing, state.CancelCause)
	assert.Equal(fac.Falsum(), subsumed)
}

func TestDnfSubsumptionLarge(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	formulas, _ := io.ReadFormulas(fac, "./../test/data/formulas/small_formulas.txt")
	for i := range 10 {
		dnf := normalform.FactorizedDNF(fac, formulas[i])
		subsumed := dnfs(t, fac, dnf)
		assert.True(sat.IsEquivalent(fac, dnf, formulas[i]))
		assert.Greater(len(fac.Operands(dnf)), len(fac.Operands(subsumed)))
	}
}

func equalSubset(t *testing.T, expected *f.LitSet, actual *treeset.Set) {
	assert := assert.New(t)
	assert.Equal(expected.Size(), actual.Size())
	for _, element := range expected.Content() {
		assert.True(actual.Contains(element))
	}
}

func set(elements ...f.Literal) *f.LitSet {
	return f.NewLitSet(elements...)
}

func cnfs(t *testing.T, fac f.Factory, formula f.Formula) f.Formula {
	res, err := CNFSubsumption(fac, formula)
	assert.Nil(t, err)
	return res
}

func dnfs(t *testing.T, fac f.Factory, formula f.Formula) f.Formula {
	res, err := DNFSubsumption(fac, formula)
	assert.Nil(t, err)
	return res
}
