package enum

import (
	"testing"

	"github.com/booleworks/logicng-go/model/iter"
	"github.com/booleworks/logicng-go/sat"

	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/normalform"
	"github.com/booleworks/logicng-go/parser"
	"github.com/booleworks/logicng-go/randomizer"
	"github.com/stretchr/testify/assert"
)

func TestCanonicalCNFSimple(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	assert.Equal(fac.Falsum(), CanonicalCNF(fac, fac.Falsum()))
	assert.Equal(fac.Verum(), CanonicalCNF(fac, fac.Verum()))
	assert.Equal(p.ParseUnsafe("a"), CanonicalCNF(fac, p.ParseUnsafe("a")))
	assert.Equal(p.ParseUnsafe("~a"), CanonicalCNF(fac, p.ParseUnsafe("~a")))
	assert.Equal(p.ParseUnsafe("(a | b) & (~a | b) & (~a | ~b)"), CanonicalCNF(fac, p.ParseUnsafe("~a & b")))
	assert.Equal(p.ParseUnsafe("~a | b"), CanonicalCNF(fac, p.ParseUnsafe("~a | b")))
	assert.Equal(p.ParseUnsafe("~a | b"), CanonicalCNF(fac, p.ParseUnsafe("a => b")))
	assert.Equal(p.ParseUnsafe("(a | ~b) & (~a | b)"), CanonicalCNF(fac, p.ParseUnsafe("a <=> b")))
	assert.Equal(p.ParseUnsafe("(~a | ~b) & (a | b)"), CanonicalCNF(fac, p.ParseUnsafe("a + b = 1")))
	assert.Equal(p.ParseUnsafe("$true"), CanonicalCNF(fac, p.ParseUnsafe("a | b | ~a & ~b")))
	assert.Equal(
		p.ParseUnsafe("(a | b | c) & (a | b | ~c) & (a | ~b | c) & (a | ~b | ~c) & (~a | b | ~c)"),
		CanonicalCNF(fac, p.ParseUnsafe("a & (b | ~c)")),
	)
	assert.Equal(
		p.ParseUnsafe("(a | b) & (~a | b) & (~a | ~b) & (a | ~b)"),
		CanonicalCNF(fac, p.ParseUnsafe("a & b & (~a | ~b)")),
	)
}

func TestCanonicalCNFCornerCases(t *testing.T) {
	fac := f.NewFactory()
	for _, cc := range f.NewCornerCases(fac) {
		testCNF(t, fac, cc)
	}
}

func TestCanonicalCNFRandom(t *testing.T) {
	fac := f.NewFactory()
	config := randomizer.DefaultConfig()
	config.NumVars = 5
	config.WeightPBC = 0.5
	config.Seed = 42
	randomizer := randomizer.New(fac, config)
	for range 1000 {
		testCNF(t, fac, randomizer.Formula(3))
	}
}

func TestCanonicalCNFHandler(t *testing.T) {
	fac := f.NewFactory()
	formula := randomizer.New(fac).Formula(5)
	handler := iter.HandlerWithLimit(1)
	cnf, state := CanonicalCNFWithHandler(fac, formula, handler)
	assert.False(t, state.Success)
	assert.Equal(t, fac.Falsum(), cnf)
}

func TestCanonicalDNFSimple(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)

	assert.Equal(fac.Falsum(), CanonicalDNF(fac, fac.Falsum()))
	assert.Equal(fac.Verum(), CanonicalDNF(fac, fac.Verum()))
	assert.Equal(p.ParseUnsafe("a"), CanonicalDNF(fac, p.ParseUnsafe("a")))
	assert.Equal(p.ParseUnsafe("~a"), CanonicalDNF(fac, p.ParseUnsafe("~a")))
	assert.Equal(p.ParseUnsafe("~a & b"), CanonicalDNF(fac, p.ParseUnsafe("~a & b")))
	assert.Equal(p.ParseUnsafe("~a & ~b | ~a & b | a & b"), CanonicalDNF(fac, p.ParseUnsafe("~a | b")))
	assert.Equal(p.ParseUnsafe("~a & ~b | ~a & b | a & b"), CanonicalDNF(fac, p.ParseUnsafe("a => b")))
	assert.Equal(p.ParseUnsafe("~a & ~b | a & b"), CanonicalDNF(fac, p.ParseUnsafe("a <=> b")))
	assert.Equal(p.ParseUnsafe("~a & b | a & ~b"), CanonicalDNF(fac, p.ParseUnsafe("a + b = 1")))
	assert.Equal(p.ParseUnsafe("$false"), CanonicalDNF(fac, p.ParseUnsafe("a & b & (~a | ~b)")))
	assert.Equal(
		p.ParseUnsafe("a & ~b & ~c | a & b & ~c | a & b & c"),
		CanonicalDNF(fac, p.ParseUnsafe("a & (b | ~c)")),
	)
	assert.Equal(
		p.ParseUnsafe("~a & b | a & b | a & ~b | ~a & ~b"),
		CanonicalDNF(fac, p.ParseUnsafe("a | b | ~a & ~b")),
	)
}

func TestCanonicalDNFCornerCases(t *testing.T) {
	fac := f.NewFactory()
	for _, cc := range f.NewCornerCases(fac) {
		testDNF(t, fac, cc)
	}
}

func TestCanonicalDNFRandom(t *testing.T) {
	fac := f.NewFactory()
	config := randomizer.DefaultConfig()
	config.NumVars = 5
	config.WeightPBC = 0.5
	config.Seed = 42
	randomizer := randomizer.New(fac, config)
	for range 1000 {
		testDNF(t, fac, randomizer.Formula(3))
	}
}

func testCNF(t *testing.T, fac f.Factory, formula f.Formula) {
	cnf := CanonicalCNF(fac, formula)
	assert.True(t, normalform.IsCNF(fac, cnf))
	assert.True(t, sat.IsEquivalent(fac, formula, cnf))
	if sat.IsTautology(fac, formula) {
		assert.Equal(t, cnf, fac.Verum())
	} else {
		assert.True(t, hasConstantTermSizeCNF(fac, cnf))
	}
}

func testDNF(t *testing.T, fac f.Factory, formula f.Formula) {
	dnf := CanonicalDNF(fac, formula)
	assert.True(t, normalform.IsDNF(fac, dnf))
	assert.True(t, sat.IsEquivalent(fac, formula, dnf))
	if sat.IsContradiction(fac, formula) {
		assert.Equal(t, dnf, fac.Falsum())
	} else {
		assert.True(t, hasConstantTermSizeDNF(fac, dnf))
	}
}

func hasConstantTermSizeCNF(fac f.Factory, cnf f.Formula) bool {
	switch cnf.Sort() {
	case f.SortLiteral, f.SortTrue, f.SortFalse, f.SortOr:
		return true
	case f.SortAnd:
		ops, _ := fac.NaryOperands(cnf)
		count := f.NumberOfAtoms(fac, ops[0])
		for i := 1; i < len(ops); i++ {
			if f.NumberOfAtoms(fac, ops[i]) != count {
				return false
			}
		}
		return true
	default:
		panic(errorx.BadFormulaSort(cnf.Sort()))
	}
}

func hasConstantTermSizeDNF(fac f.Factory, cnf f.Formula) bool {
	switch cnf.Sort() {
	case f.SortLiteral, f.SortTrue, f.SortFalse, f.SortAnd:
		return true
	case f.SortOr:
		ops, _ := fac.NaryOperands(cnf)
		count := f.NumberOfAtoms(fac, ops[0])
		for i := 1; i < len(ops); i++ {
			if f.NumberOfAtoms(fac, ops[i]) != count {
				return false
			}
		}
		return true
	default:
		panic(errorx.BadFormulaSort(cnf.Sort()))
	}
}

func TestCanonicalDNFHandler(t *testing.T) {
	fac := f.NewFactory()
	formula := randomizer.New(fac).Formula(5)
	handler := iter.HandlerWithLimit(1)
	dnf, state := CanonicalDNFWithHandler(fac, formula, handler)
	assert.False(t, state.Success)
	assert.Equal(t, fac.Falsum(), dnf)
}
