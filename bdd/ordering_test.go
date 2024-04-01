package bdd

import (
	"testing"

	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/parser"
	"booleworks.com/logicng/sat"
	"github.com/stretchr/testify/assert"
)

func TestBFSOrder(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	a := fac.Var("A")
	b := fac.Var("B")
	c := fac.Var("C")
	d := fac.Var("D")

	assert.Equal([]f.Variable{}, BFSOrder(fac, p.ParseUnsafe("$true")))
	assert.Equal([]f.Variable{}, BFSOrder(fac, p.ParseUnsafe("$false")))
	assert.Equal([]f.Variable{a}, BFSOrder(fac, p.ParseUnsafe("A")))
	assert.Equal([]f.Variable{a, b}, BFSOrder(fac, p.ParseUnsafe("A => ~B")))
	assert.Equal([]f.Variable{a, b}, BFSOrder(fac, p.ParseUnsafe("A <=> B")))
	assert.Equal([]f.Variable{a, b}, BFSOrder(fac, p.ParseUnsafe("~(A <=> ~B)")))
	assert.Equal([]f.Variable{a, b, d, c}, BFSOrder(fac, p.ParseUnsafe("A | ~C | B | D")))
	assert.Equal([]f.Variable{a, b, d, c}, BFSOrder(fac, p.ParseUnsafe("A & ~C & B & D")))
	assert.Equal([]f.Variable{a, c, b, d}, BFSOrder(fac, p.ParseUnsafe("A + C + B + D < 2")))

	formula := p.ParseUnsafe("(A => ~B) & ((A & C) | (D & ~C)) & (A | Y | X) & (Y <=> (X | (W + A + F < 1)))")
	expected := []f.Variable{
		fac.Var("A"),
		fac.Var("Y"),
		fac.Var("X"),
		fac.Var("B"),
		fac.Var("C"),
		fac.Var("D"),
		fac.Var("W"),
		fac.Var("F"),
	}
	assert.Equal(expected, BFSOrder(fac, formula))
}

func TestDFSOrder(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	a := fac.Var("A")
	b := fac.Var("B")
	c := fac.Var("C")
	d := fac.Var("D")

	assert.Equal([]f.Variable{}, DFSOrder(fac, p.ParseUnsafe("$true")))
	assert.Equal([]f.Variable{}, DFSOrder(fac, p.ParseUnsafe("$false")))
	assert.Equal([]f.Variable{a}, DFSOrder(fac, p.ParseUnsafe("A")))
	assert.Equal([]f.Variable{a, b}, DFSOrder(fac, p.ParseUnsafe("A => ~B")))
	assert.Equal([]f.Variable{a, b}, DFSOrder(fac, p.ParseUnsafe("A <=> B")))
	assert.Equal([]f.Variable{a, b}, DFSOrder(fac, p.ParseUnsafe("~(A <=> ~B)")))
	assert.Equal([]f.Variable{a, c, b, d}, DFSOrder(fac, p.ParseUnsafe("A | ~C | B | D")))
	assert.Equal([]f.Variable{a, c, b, d}, DFSOrder(fac, p.ParseUnsafe("A & ~C & B & D")))
	assert.Equal([]f.Variable{a, c, b, d}, DFSOrder(fac, p.ParseUnsafe("A + C + B + D < 2")))

	formula := p.ParseUnsafe("(A => ~B) & ((A & C) | (D & ~C)) & (A | Y | X) & (Y <=> (X | (W + A + F < 1)))")
	expected := []f.Variable{
		fac.Var("A"),
		fac.Var("B"),
		fac.Var("C"),
		fac.Var("D"),
		fac.Var("Y"),
		fac.Var("X"),
		fac.Var("W"),
		fac.Var("F"),
	}
	assert.Equal(expected, DFSOrder(fac, formula))
}

func TestMin2MaxOrder(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	a := fac.Var("A")
	b := fac.Var("B")
	c := fac.Var("C")
	d := fac.Var("D")

	assert.Equal([]f.Variable{}, MinToMaxOrder(fac, p.ParseUnsafe("$true")))
	assert.Equal([]f.Variable{}, MinToMaxOrder(fac, p.ParseUnsafe("$false")))
	assert.Equal([]f.Variable{a}, MinToMaxOrder(fac, p.ParseUnsafe("A")))
	assert.Equal([]f.Variable{a, b}, MinToMaxOrder(fac, p.ParseUnsafe("A => ~B")))
	assert.Equal([]f.Variable{a, b}, MinToMaxOrder(fac, p.ParseUnsafe("A <=> B")))
	assert.Equal([]f.Variable{a, b}, MinToMaxOrder(fac, p.ParseUnsafe("~(A <=> ~B)")))
	assert.Equal([]f.Variable{a, c, b, d}, MinToMaxOrder(fac, p.ParseUnsafe("A | ~C | B | D")))
	assert.Equal([]f.Variable{a, c, b, d}, MinToMaxOrder(fac, p.ParseUnsafe("A & ~C & B & D")))
	assert.Equal([]f.Variable{a, c, b, d}, MinToMaxOrder(fac, p.ParseUnsafe("A + C + B + D < 2")))

	formula := p.ParseUnsafe("(A => ~B) & ((A & C) | (D & ~C)) & (A | Y | X) & (Y <=> (X | (W + A + F < 1)))")
	expected := []f.Variable{
		fac.Var("B"),
		fac.Var("D"),
		fac.Var("W"),
		fac.Var("F"),
		fac.Var("C"),
		fac.Var("Y"),
		fac.Var("X"),
		fac.Var("A"),
	}
	assert.Equal(expected, MinToMaxOrder(fac, formula))
}

func TestMax2MinOrder(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	a := fac.Var("A")
	b := fac.Var("B")
	c := fac.Var("C")
	d := fac.Var("D")

	assert.Equal([]f.Variable{}, MaxToMinOrder(fac, p.ParseUnsafe("$true")))
	assert.Equal([]f.Variable{}, MaxToMinOrder(fac, p.ParseUnsafe("$false")))
	assert.Equal([]f.Variable{a}, MaxToMinOrder(fac, p.ParseUnsafe("A")))
	assert.Equal([]f.Variable{a, b}, MaxToMinOrder(fac, p.ParseUnsafe("A => ~B")))
	assert.Equal([]f.Variable{a, b}, MaxToMinOrder(fac, p.ParseUnsafe("A <=> B")))
	assert.Equal([]f.Variable{a, b}, MaxToMinOrder(fac, p.ParseUnsafe("~(A <=> ~B)")))
	assert.Equal([]f.Variable{a, c, b, d}, MaxToMinOrder(fac, p.ParseUnsafe("A | ~C | B | D")))
	assert.Equal([]f.Variable{a, c, b, d}, MaxToMinOrder(fac, p.ParseUnsafe("A & ~C & B & D")))
	assert.Equal([]f.Variable{a, c, b, d}, MaxToMinOrder(fac, p.ParseUnsafe("A + C + B + D < 2")))

	formula := p.ParseUnsafe("(A => ~B) & ((A & C) | (D & ~C)) & (A | Y | X) & (Y <=> (X | (W + A + F < 1)))")
	expected := []f.Variable{
		fac.Var("A"),
		fac.Var("C"),
		fac.Var("Y"),
		fac.Var("X"),
		fac.Var("B"),
		fac.Var("D"),
		fac.Var("W"),
		fac.Var("F"),
	}
	assert.Equal(expected, MaxToMinOrder(fac, formula))
}

func TestBDDOrderings(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("(A => ~B) & ((A & C) | (D & ~C)) & (A | Y | X) & (Y <=> (X | (W + A + F < 1)))")

	bddNoOrder := Build(fac, formula)
	bddBfs := BuildWithVarOrder(fac, formula, BFSOrder(fac, formula))
	bddDfs := BuildWithVarOrder(fac, formula, DFSOrder(fac, formula))
	bddMin2Max := BuildWithVarOrder(fac, formula, MinToMaxOrder(fac, formula))
	bddMax2Min := BuildWithVarOrder(fac, formula, MaxToMinOrder(fac, formula))
	bddForce := BuildWithVarOrder(fac, formula, ForceOrder(fac, formula))

	assert.Equal(13, bddNoOrder.NodeCount())
	assert.Equal(14, bddBfs.NodeCount())
	assert.Equal(13, bddDfs.NodeCount())
	assert.Equal(24, bddMin2Max.NodeCount())
	assert.Equal(17, bddMax2Min.NodeCount())

	assert.True(sat.IsEquivalent(fac, bddNoOrder.CNF(), formula))
	assert.True(sat.IsEquivalent(fac, bddBfs.CNF(), formula))
	assert.True(sat.IsEquivalent(fac, bddDfs.CNF(), formula))
	assert.True(sat.IsEquivalent(fac, bddMin2Max.CNF(), formula))
	assert.True(sat.IsEquivalent(fac, bddMax2Min.CNF(), formula))
	assert.True(sat.IsEquivalent(fac, bddForce.CNF(), formula))

	assert.True(sat.IsEquivalent(fac, bddNoOrder.DNF(), formula))
	assert.True(sat.IsEquivalent(fac, bddBfs.DNF(), formula))
	assert.True(sat.IsEquivalent(fac, bddDfs.DNF(), formula))
	assert.True(sat.IsEquivalent(fac, bddMin2Max.DNF(), formula))
	assert.True(sat.IsEquivalent(fac, bddMax2Min.DNF(), formula))
	assert.True(sat.IsEquivalent(fac, bddForce.DNF(), formula))
}
