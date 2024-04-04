package dnnf

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/booleworks/logicng-go/bdd"
	"github.com/booleworks/logicng-go/encoding"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/io"
	"github.com/booleworks/logicng-go/parser"
	"github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

func TestDNNFTrivialFormulas(t *testing.T) {
	fac := f.NewFactory()
	parser := parser.New(fac)
	testFormula(t, fac, parser.ParseUnsafe("$true"), true)
	testFormula(t, fac, parser.ParseUnsafe("$false"), true)
	testFormula(t, fac, parser.ParseUnsafe("a"), true)
	testFormula(t, fac, parser.ParseUnsafe("~a"), true)
	testFormula(t, fac, parser.ParseUnsafe("a & b"), true)
	testFormula(t, fac, parser.ParseUnsafe("a | b"), true)
	testFormula(t, fac, parser.ParseUnsafe("a => b"), true)
	testFormula(t, fac, parser.ParseUnsafe("a <=> b"), true)
	testFormula(t, fac, parser.ParseUnsafe("a | b | c"), true)
	testFormula(t, fac, parser.ParseUnsafe("a & b & c"), true)
	testFormula(t, fac, parser.ParseUnsafe("f & ((~b | c) <=> ~a & ~c)"), true)
	testFormula(t, fac, parser.ParseUnsafe("a | ((b & ~c) | (c & (~d | ~a & b)) & e)"), true)
	testFormula(t, fac, parser.ParseUnsafe("a + b + c + d <= 1"), true)
	testFormula(t, fac, parser.ParseUnsafe("a + b + c + d <= 3"), false)
	testFormula(t, fac, parser.ParseUnsafe("2*a + 3*b + -2*c + d < 5"), false)
	testFormula(t, fac, parser.ParseUnsafe("2*a + 3*b + -2*c + d >= 5"), false)
	testFormula(t, fac, parser.ParseUnsafe("~a & (~a | b | c | d)"), true)
}

func TestDNNFLargeFormulas(t *testing.T) {
	fac := f.NewFactory()
	dimacs, _ := io.ReadDimacs(fac, "../test/data/dnnf/both_bdd_dnnf_1.cnf")
	testFormula(t, fac, fac.And(*dimacs...), true)
	dimacs, _ = io.ReadDimacs(fac, "../test/data/dnnf/both_bdd_dnnf_2.cnf")
	testFormula(t, fac, fac.And(*dimacs...), true)
	dimacs, _ = io.ReadDimacs(fac, "../test/data/dnnf/both_bdd_dnnf_3.cnf")
	testFormula(t, fac, fac.And(*dimacs...), true)
	dimacs, _ = io.ReadDimacs(fac, "../test/data/dnnf/both_bdd_dnnf_4.cnf")
	testFormula(t, fac, fac.And(*dimacs...), true)
	dimacs, _ = io.ReadDimacs(fac, "../test/data/dnnf/both_bdd_dnnf_5.cnf")
	testFormula(t, fac, fac.And(*dimacs...), true)
}

func TestDNNFLargestFormula(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	fac := f.NewFactory()
	ccConfig := encoding.DefaultConfig()
	ccConfig.AMOEncoder = encoding.AMOPure
	fac.PutConfiguration(ccConfig)
	formula, _ := io.ReadFormula(fac, "../test/data/formulas/small.txt")
	dnnf := Compile(fac, formula)
	modelCount := dnnf.ModelCount()
	assert.Equal(t, "48394912530540796831369012282361118720", fmt.Sprintf("%s", modelCount))
}

func TestDNNFSmallFormulas(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	fac := f.NewFactory()
	formulas, _ := io.ReadFormulas(fac, "../test/data/formulas/small_formulas.txt")
	for _, form := range formulas[:500] {
		testFormula(t, fac, form, false)
	}
}

func testFormula(t *testing.T, fac f.Factory, formula f.Formula, withEquivalence bool) {
	dnnf := Compile(fac, formula)
	dnnfCount := dnnf.ModelCount()
	if withEquivalence {
		assert.True(t, sat.IsEquivalent(fac, formula, dnnf.Formula))
	}
	bddCount := countWithBDD(fac, formula)
	assert.Equal(t, bddCount, dnnfCount)
}

func countWithBDD(fac f.Factory, formula f.Formula) *big.Int {
	if formula.Sort() == f.SortTrue {
		return big.NewInt(1)
	} else if formula.Sort() == f.SortFalse {
		return big.NewInt(0)
	}
	order := bdd.ForceOrder(fac, formula)
	kernel := bdd.NewKernelWithOrdering(fac, order, 100000, 1000000)
	bdd := bdd.CompileWithKernel(fac, formula, kernel)
	return bdd.ModelCount()
}
