package bdd

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/booleworks/logicng-go/assignment"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/model"
	"github.com/booleworks/logicng-go/normalform"
	"github.com/booleworks/logicng-go/parser"
	"github.com/booleworks/logicng-go/randomizer"
	"github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

func TestBDDTrue(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	bdd := Compile(fac, fac.Verum())
	assert.True(bdd.IsTautology())
	assert.False(bdd.IsContradiction())
	assert.Equal(fac.Verum(), bdd.CNF())
	assert.Equal(fac.Verum(), bdd.DNF())
	assert.Equal(*big.NewInt(1), *bdd.ModelCount())
	assert.Equal(*big.NewInt(0), *bdd.NumberOfCNFClauses())
	assert.Equal("<$true>", bdd.NodeRepresentation().String())
}

func TestBDDFalse(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	bdd := Compile(fac, fac.Falsum())
	assert.False(bdd.IsTautology())
	assert.True(bdd.IsContradiction())
	assert.Equal(fac.Falsum(), bdd.CNF())
	assert.Equal(fac.Falsum(), bdd.DNF())
	assert.Equal(*big.NewInt(0), *bdd.ModelCount())
	assert.Equal(*big.NewInt(1), *bdd.NumberOfCNFClauses())
	assert.Equal("<$false>", bdd.NodeRepresentation().String())
}

func TestBDDVariable(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	va := fac.Var("A")

	bdd := Compile(fac, fac.Variable("A"))
	assert.False(bdd.IsTautology())
	assert.False(bdd.IsContradiction())
	assert.Equal(fac.Variable("A"), bdd.CNF())
	assert.Equal(fac.Variable("A"), bdd.DNF())
	assert.Equal(*big.NewInt(1), *bdd.ModelCount())
	assert.Equal(*big.NewInt(1), *bdd.NumberOfCNFClauses())
	modelsEqual(t, []*model.Model{model.New(va.AsLiteral())}, bdd.ModelEnumeration(va))
	assert.Equal("<A | low=<$false> high=<$true>>", bdd.NodeRepresentation().String())
}

func TestBDDNegativeVariable(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	va := fac.Var("A")
	vna := fac.Lit("A", false)

	bdd := Compile(fac, fac.Literal("A", false))
	assert.False(bdd.IsTautology())
	assert.False(bdd.IsContradiction())
	assert.Equal(fac.Literal("A", false), bdd.CNF())
	assert.Equal(fac.Literal("A", false), bdd.DNF())
	assert.Equal(*big.NewInt(1), *bdd.ModelCount())
	assert.Equal(*big.NewInt(1), *bdd.NumberOfCNFClauses())
	modelsEqual(t, []*model.Model{model.New(vna)}, bdd.ModelEnumeration(va))
	assert.Equal("<A | low=<$true> high=<$false>>", bdd.NodeRepresentation().String())
}

func TestBDDImplication(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	va := fac.Var("A")
	vna := fac.Lit("A", false)
	vb := fac.Var("B")
	vnb := fac.Lit("B", false)

	p := parser.New(fac)
	bdd := Compile(fac, p.ParseUnsafe("A => ~B"))
	assert.False(bdd.IsTautology())
	assert.False(bdd.IsContradiction())
	assert.Equal(p.ParseUnsafe("~A | ~B"), bdd.CNF())
	assert.Equal(p.ParseUnsafe("~A | A & ~B"), bdd.DNF())
	assert.Equal(*big.NewInt(3), *bdd.ModelCount())
	assert.Equal(*big.NewInt(1), *bdd.NumberOfCNFClauses())
	expected := []*model.Model{
		model.New(vna, vnb),
		model.New(vna, vb.AsLiteral()),
		model.New(va.AsLiteral(), vnb),
	}
	modelsEqual(t, expected, bdd.ModelEnumeration(va, vb))
	assert.Equal("<A | low=<$true> high=<B | low=<$true> high=<$false>>>", bdd.NodeRepresentation().String())
}

func TestBDDEquivalence(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	bdd := Compile(fac, p.ParseUnsafe("A <=> ~B"))
	va := fac.Var("A")
	vna := fac.Lit("A", false)
	vb := fac.Var("B")
	vnb := fac.Lit("B", false)

	assert.False(bdd.IsTautology())
	assert.False(bdd.IsContradiction())
	assert.Equal(p.ParseUnsafe("(A | B) & (~A | ~B)"), bdd.CNF())
	assert.Equal(p.ParseUnsafe("~A & B | A & ~B"), bdd.DNF())
	assert.Equal(*big.NewInt(2), *bdd.ModelCount())
	assert.Equal(*big.NewInt(2), *bdd.NumberOfCNFClauses())
	expected := []*model.Model{
		model.New(vna, vb.AsLiteral()),
		model.New(va.AsLiteral(), vnb),
	}
	modelsEqual(t, expected, bdd.ModelEnumeration(va, vb))
	assert.Equal("<A | low=<B | low=<$false> high=<$true>> high=<B | low=<$true> high=<$false>>>", bdd.NodeRepresentation().String())
}

func TestBDDOr(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	bdd := Compile(fac, p.ParseUnsafe("A | B | ~C"))
	va := fac.Var("A")
	vna := fac.Lit("A", false)
	vb := fac.Var("B")
	vnb := fac.Lit("B", false)
	vc := fac.Var("C")
	vnc := fac.Lit("C", false)

	assert.False(bdd.IsTautology())
	assert.False(bdd.IsContradiction())
	assert.Equal(p.ParseUnsafe("A | B | ~C"), bdd.CNF())
	assert.Equal(p.ParseUnsafe("~A & ~B & ~C | ~A & B | A"), bdd.DNF())
	assert.Equal(*big.NewInt(7), *bdd.ModelCount())
	assert.Equal(*big.NewInt(1), *bdd.NumberOfCNFClauses())

	expected := []*model.Model{
		model.New(vna, vnb, vnc),
		model.New(vna, vb.AsLiteral(), vnc),
		model.New(vna, vb.AsLiteral(), vc.AsLiteral()),
		model.New(va.AsLiteral(), vnb, vnc),
		model.New(va.AsLiteral(), vnb, vc.AsLiteral()),
		model.New(va.AsLiteral(), vb.AsLiteral(), vnc),
		model.New(va.AsLiteral(), vb.AsLiteral(), vc.AsLiteral()),
	}
	modelsEqual(t, expected, bdd.ModelEnumeration(va, vb, vc))
	assert.Equal("<A | low=<B | low=<C | low=<$true> high=<$false>> high=<$true>> high=<$true>>", bdd.NodeRepresentation().String())
}

func TestBDDAnd(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	bdd := Compile(fac, p.ParseUnsafe("A & B & ~C"))
	va := fac.Var("A")
	vb := fac.Var("B")
	vc := fac.Var("C")
	vnc := fac.Lit("C", false)

	assert.False(bdd.IsTautology())
	assert.False(bdd.IsContradiction())
	assert.Equal(p.ParseUnsafe("A & (~A | B) & (~A | ~B | ~C)"), bdd.CNF())
	assert.Equal(p.ParseUnsafe("A & B & ~C"), bdd.DNF())
	assert.Equal(*big.NewInt(1), *bdd.ModelCount())
	assert.Equal(*big.NewInt(3), *bdd.NumberOfCNFClauses())

	expected := []*model.Model{model.New(va.AsLiteral(), vb.AsLiteral(), vnc)}
	modelsEqual(t, expected, bdd.ModelEnumeration(va, vb, vc))
	assert.Equal("<A | low=<$false> high=<B | low=<$false> high=<C | low=<$true> high=<$false>>>>", bdd.NodeRepresentation().String())
}

func TestBDDFormula(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	va := fac.Var("A")
	vb := fac.Var("B")
	vc := fac.Var("C")
	p := parser.New(fac)
	formula := p.ParseUnsafe("(A => ~C) | (B & ~C)")
	bdd := Compile(fac, formula)

	assert.False(bdd.IsTautology())
	assert.False(bdd.IsContradiction())
	assert.True(normalform.IsCNF(fac, bdd.CNF()))
	assert.True(sat.IsEquivalent(fac, formula, bdd.CNF()))
	assert.True(normalform.IsDNF(fac, bdd.DNF()))
	assert.True(sat.IsEquivalent(fac, formula, bdd.DNF()))
	assert.Equal(*big.NewInt(6), *bdd.ModelCount())
	assert.Equal(6, len(bdd.ModelEnumeration(va, vb, vc)))
	assert.Equal(2, len(bdd.ModelEnumeration(va)))
}

func TestBDDCC(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	va := fac.Var("A")
	vna := fac.Lit("A", false)
	vb := fac.Var("B")
	vnb := fac.Lit("B", false)
	vc := fac.Var("C")
	vnc := fac.Lit("C", false)
	p := parser.New(fac)
	formula := p.ParseUnsafe("A + B + C = 1")
	bdd := Compile(fac, formula)

	assert.False(bdd.IsTautology())
	assert.False(bdd.IsContradiction())
	assert.True(normalform.IsCNF(fac, bdd.CNF()))
	assert.True(sat.IsEquivalent(fac, formula, bdd.CNF()))
	assert.True(normalform.IsDNF(fac, bdd.DNF()))
	assert.True(sat.IsEquivalent(fac, formula, bdd.DNF()))
	assert.Equal(*big.NewInt(3), *bdd.ModelCount())
	assert.Equal(*big.NewInt(4), *bdd.NumberOfCNFClauses())

	expected := []*model.Model{
		model.New(vna, vnb, vc.AsLiteral()),
		model.New(vna, vb.AsLiteral(), vnc),
		model.New(va.AsLiteral(), vnb, vnc),
	}
	modelsEqual(t, expected, bdd.ModelEnumeration(va, vb, vc))
}

func modelsEqual(t *testing.T, m1, m2 []*model.Model) {
	assert.Equal(t, len(m1), len(m2))
	for i := 0; i < len(m1); i++ {
		assert.Equal(t, *m1[i], *m2[i])
	}
}

func TestBddPigeonHole(t *testing.T) {
	fac := f.NewFactory()
	testPigeonHole(t, fac, 2)
	testPigeonHole(t, fac, 3)
	testPigeonHole(t, fac, 4)
	testPigeonHole(t, fac, 5)
	testPigeonHole(t, fac, 6)
	testPigeonHole(t, fac, 7)
	testPigeonHole(t, fac, 8)
	testPigeonHole(t, fac, 9)
}

func testPigeonHole(t *testing.T, fac f.Factory, size int) {
	pigeon := sat.GeneratePigeonHole(fac, size)
	numVars := f.Variables(fac, pigeon).Size()
	kernel := NewKernel(fac, int32(numVars), 10000, 10000)
	bdd := CompileWithKernel(fac, pigeon, kernel)
	assert.True(t, bdd.IsContradiction())
}

func TestBDDNQueens(t *testing.T) {
	fac := f.NewFactory()
	testQueens(t, fac, 4, 2)
	testQueens(t, fac, 5, 10)
	testQueens(t, fac, 6, 4)
	testQueens(t, fac, 7, 40)
	testQueens(t, fac, 8, 92)
}

func testQueens(t *testing.T, fac f.Factory, size, models int) {
	queens := sat.GenerateNQueens(fac, size)
	numVars := f.Variables(fac, queens).Size()
	kernel := NewKernel(fac, int32(numVars), 10000, 10000)
	bdd := CompileWithKernel(fac, queens, kernel)
	cnf := bdd.CNF()
	assert.True(t, normalform.IsCNF(fac, cnf))
	cnfBdd := CompileWithKernel(fac, cnf, kernel)
	assert.Equal(t, bdd.Index, cnfBdd.Index)
	assert.Equal(t, *big.NewInt(int64(models)), *bdd.ModelCount())
	assert.Equal(t, numVars, len(bdd.Support()))
}

type testBDDs struct {
	kernel    *Kernel
	bddVerum  *BDD
	bddFalsum *BDD
	bddPosLit *BDD
	bddNegLit *BDD
	bddImpl   *BDD
	bddEquiv  *BDD
	bddOr     *BDD
	bddAnd    *BDD
}

func testData(fac f.Factory) testBDDs {
	parser := parser.New(fac)
	kernel := NewKernel(fac, 3, 100, 100)
	return testBDDs{
		kernel:    kernel,
		bddVerum:  CompileWithKernel(fac, fac.Verum(), kernel),
		bddFalsum: CompileWithKernel(fac, fac.Falsum(), kernel),
		bddPosLit: CompileWithKernel(fac, fac.Literal("A", true), kernel),
		bddNegLit: CompileWithKernel(fac, fac.Literal("A", false), kernel),
		bddImpl:   CompileWithKernel(fac, parser.ParseUnsafe("A => ~B"), kernel),
		bddEquiv:  CompileWithKernel(fac, parser.ParseUnsafe("A <=> ~B"), kernel),
		bddOr:     CompileWithKernel(fac, parser.ParseUnsafe("A | B | ~C"), kernel),
		bddAnd:    CompileWithKernel(fac, parser.ParseUnsafe("A & B & ~C"), kernel),
	}
}

func TestBDDToFormula(t *testing.T) {
	fac := f.NewFactory()
	assert := assert.New(t)
	p := parser.New(fac)
	d := testData(fac)
	assert.Equal(fac.Verum(), d.bddVerum.ToFormula(fac))
	assert.Equal(fac.Falsum(), d.bddFalsum.ToFormula(fac))
	assert.Equal(fac.Literal("A", true), d.bddPosLit.ToFormula(fac))
	assert.Equal(fac.Literal("A", false), d.bddNegLit.ToFormula(fac))
	compareFormula(t, fac, d.bddImpl, p.ParseUnsafe("A => ~B"))
	compareFormula(t, fac, d.bddEquiv, p.ParseUnsafe("A <=> ~B"))
	compareFormula(t, fac, d.bddOr, p.ParseUnsafe("A | B | ~C"))
	compareFormula(t, fac, d.bddAnd, p.ParseUnsafe("A & B & ~C"))
}

func TestBDDToFormulaStyles(t *testing.T) {
	fac := f.NewFactory()
	assert := assert.New(t)
	p := parser.New(fac)
	bdd := Compile(fac, p.ParseUnsafe("~A | ~B | ~C"))
	expFollowPathsToTrue := p.ParseUnsafe("~A | A & (~B | B & ~C)")
	assert.True(sat.IsEquivalent(fac, bdd.ToFormula(fac), expFollowPathsToTrue))
	assert.True(sat.IsEquivalent(fac, bdd.ToFormula(fac, true), expFollowPathsToTrue))
	assert.True(sat.IsEquivalent(fac, bdd.ToFormula(fac, false), expFollowPathsToTrue))
}

func TestBDDToFormulaRandom(t *testing.T) {
	fac := f.NewFactory()
	numTests := 100
	if testing.Short() {
		numTests = 10
	}
	for i := 0; i < numTests; i++ {
		rand := randomizer.NewWithSeed(fac, int64(i))
		formula := rand.Formula(5)
		bdd := Compile(fac, formula)
		compareFormula(t, fac, bdd, formula)
	}
}

func TestRestriction(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	d := testData(fac)
	a := fac.Lit("A", true)
	na := fac.Lit("A", false)
	b := fac.Lit("B", true)
	nb := fac.Lit("B", false)

	equalRestrict(t, d.bddVerum, d.bddVerum, a)
	equalRestrict(t, d.bddVerum, d.bddVerum, na)
	equalRestrict(t, d.bddVerum, d.bddVerum, a, b)
	equalRestrict(t, d.bddFalsum, d.bddFalsum, a)
	equalRestrict(t, d.bddFalsum, d.bddFalsum, na)
	equalRestrict(t, d.bddFalsum, d.bddFalsum, a, b)
	equalRestrict(t, d.bddVerum, d.bddPosLit, a)
	equalRestrict(t, d.bddFalsum, d.bddPosLit, na)
	equalRestrict(t, d.bddVerum, d.bddPosLit, a, b)
	equalRestrict(t, d.bddFalsum, d.bddNegLit, a)
	equalRestrict(t, d.bddVerum, d.bddNegLit, na)
	equalRestrict(t, d.bddFalsum, d.bddNegLit, a, b)
	equalRestrict(t, CompileWithKernel(fac, nb.AsFormula(), d.kernel), d.bddImpl, a)
	equalRestrict(t, d.bddVerum, d.bddImpl, na)
	equalRestrict(t, d.bddFalsum, d.bddImpl, a, b)
	equalRestrict(t, CompileWithKernel(fac, nb.AsFormula(), d.kernel), d.bddEquiv, a)
	equalRestrict(t, CompileWithKernel(fac, b.AsFormula(), d.kernel), d.bddEquiv, na)
	equalRestrict(t, d.bddFalsum, d.bddEquiv, a, b)
	equalRestrict(t, d.bddVerum, d.bddOr, a)
	equalRestrict(t, CompileWithKernel(fac, p.ParseUnsafe("B | ~C"), d.kernel), d.bddOr, na)
	equalRestrict(t, d.bddVerum, d.bddOr, a, b)
	equalRestrict(t, CompileWithKernel(fac, p.ParseUnsafe("B & ~C"), d.kernel), d.bddAnd, a)
	equalRestrict(t, d.bddFalsum, d.bddAnd, na)
	equalRestrict(t, CompileWithKernel(fac, fac.Literal("C", false), d.kernel), d.bddAnd, a, b)
}

func equalRestrict(t *testing.T, expected, bdd *BDD, vars ...f.Literal) {
	restriction := bdd.Restrict(vars...)
	assert.Equal(t, expected, restriction)
}

func TestBDDExistentialQuantification(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	d := testData(fac)
	a := fac.Var("A")
	b := fac.Var("B")

	equalExistential(t, d.bddVerum, d.bddVerum, a)
	equalExistential(t, d.bddVerum, d.bddVerum, a, b)
	equalExistential(t, d.bddFalsum, d.bddFalsum, a)
	equalExistential(t, d.bddFalsum, d.bddFalsum, a, b)
	equalExistential(t, d.bddVerum, d.bddPosLit, a)
	equalExistential(t, d.bddVerum, d.bddPosLit, a, b)
	equalExistential(t, d.bddVerum, d.bddNegLit, a)
	equalExistential(t, d.bddVerum, d.bddNegLit, a, b)
	equalExistential(t, d.bddVerum, d.bddImpl, a)
	equalExistential(t, d.bddVerum, d.bddImpl, a, b)
	equalExistential(t, d.bddVerum, d.bddEquiv, a)
	equalExistential(t, d.bddVerum, d.bddEquiv, a, b)
	equalExistential(t, d.bddVerum, d.bddOr, a)
	equalExistential(t, d.bddVerum, d.bddOr, a, b)
	equalExistential(t, CompileWithKernel(fac, p.ParseUnsafe("B & ~C"), d.kernel), d.bddAnd, a)
	equalExistential(t, CompileWithKernel(fac, p.ParseUnsafe("~C"), d.kernel), d.bddAnd, a, b)
}

func equalExistential(t *testing.T, expected, bdd *BDD, vars ...f.Variable) {
	restriction := bdd.Exists(vars...)
	assert.Equal(t, expected, restriction)
}

func TestBDDUniversalQuantification(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	d := testData(fac)
	a := fac.Var("A")
	b := fac.Var("B")

	equalForAll(t, d.bddVerum, d.bddVerum, a)
	equalForAll(t, d.bddVerum, d.bddVerum, a, b)
	equalForAll(t, d.bddFalsum, d.bddFalsum, a)
	equalForAll(t, d.bddFalsum, d.bddFalsum, a, b)
	equalForAll(t, d.bddFalsum, d.bddPosLit, a)
	equalForAll(t, d.bddFalsum, d.bddPosLit, a, b)
	equalForAll(t, d.bddFalsum, d.bddNegLit, a)
	equalForAll(t, d.bddFalsum, d.bddNegLit, a, b)
	equalForAll(t, CompileWithKernel(fac, p.ParseUnsafe("~B"), d.kernel), d.bddImpl, a)
	equalForAll(t, d.bddFalsum, d.bddImpl, a, b)
	equalForAll(t, d.bddFalsum, d.bddEquiv, a)
	equalForAll(t, d.bddFalsum, d.bddEquiv, a, b)
	equalForAll(t, CompileWithKernel(fac, p.ParseUnsafe("B | ~C"), d.kernel), d.bddOr, a)
	equalForAll(t, CompileWithKernel(fac, p.ParseUnsafe("~C"), d.kernel), d.bddOr, a, b)
	equalForAll(t, d.bddFalsum, d.bddAnd, a)
	equalForAll(t, d.bddFalsum, d.bddAnd, a, b)
}

func equalForAll(t *testing.T, expected, bdd *BDD, vars ...f.Variable) {
	restriction := bdd.ForAll(vars...)
	assert.Equal(t, expected, restriction)
}

func TestBDDModel(t *testing.T) {
	fac := f.NewFactory()
	d := testData(fac)
	va := fac.Lit("A", true)
	vna := fac.Lit("A", false)
	vb := fac.Lit("B", true)
	vnb := fac.Lit("B", false)
	vnc := fac.Lit("C", false)

	m, err := d.bddFalsum.Model()
	assert.Nil(t, m)
	assert.NotNil(t, err)

	equalModel(t, model.New(), d.bddVerum)
	equalModel(t, model.New(va), d.bddPosLit)
	equalModel(t, model.New(vna), d.bddNegLit)
	equalModel(t, model.New(vna), d.bddImpl)
	equalModel(t, model.New(vb, vna), d.bddEquiv)
	equalModel(t, model.New(vnc, vnb, vna), d.bddOr)
	equalModel(t, model.New(vnc, vb, va), d.bddAnd)
}

func equalModel(t *testing.T, expected *model.Model, bdd *BDD) {
	model, err := bdd.Model()
	assert.Nil(t, err)
	assert.Equal(t, expected, model)
}

func TestBDDModelWithGivenVars(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := testData(fac)
	la := fac.Lit("A", true)
	va := fac.Var("A")
	vna := fac.Lit("A", false)
	lb := fac.Lit("B", true)
	vb := fac.Var("B")
	vnb := fac.Lit("B", false)
	vnc := fac.Lit("C", false)

	m, err := d.bddFalsum.ModelWithVariables(true, va)
	assert.Nil(m)
	assert.NotNil(err)
	m, err = d.bddFalsum.ModelWithVariables(true, va, vb)
	assert.Nil(m)
	assert.NotNil(err)
	m, err = d.bddFalsum.ModelWithVariables(false, va)
	assert.Nil(m)
	assert.NotNil(err)
	m, err = d.bddFalsum.ModelWithVariables(false, va, vb)
	assert.Nil(m)
	assert.NotNil(err)

	equalModelVars(t, model.New(la), d.bddVerum, true, va)
	equalModelVars(t, model.New(lb, la), d.bddVerum, true, va, vb)
	equalModelVars(t, model.New(vna), d.bddVerum, false, va)
	equalModelVars(t, model.New(vnb, vna), d.bddVerum, false, va, vb)
	equalModelVars(t, model.New(la), d.bddPosLit, true, va)
	equalModelVars(t, model.New(lb, la), d.bddPosLit, true, va, vb)
	equalModelVars(t, model.New(la), d.bddPosLit, false, va)
	equalModelVars(t, model.New(vnb, la), d.bddPosLit, false, va, vb)
	equalModelVars(t, model.New(vna), d.bddNegLit, true, va)
	equalModelVars(t, model.New(lb, vna), d.bddNegLit, true, va, vb)
	equalModelVars(t, model.New(vna), d.bddNegLit, false, va)
	equalModelVars(t, model.New(vnb, vna), d.bddNegLit, false, va, vb)
	equalModelVars(t, model.New(vna), d.bddImpl, true, va)
	equalModelVars(t, model.New(lb, vna), d.bddImpl, true, va, vb)
	equalModelVars(t, model.New(vna), d.bddImpl, false, va)
	equalModelVars(t, model.New(vnb, vna), d.bddImpl, false, va, vb)
	equalModelVars(t, model.New(lb, vna), d.bddEquiv, true, va)
	equalModelVars(t, model.New(lb, vna), d.bddEquiv, true, va, vb)
	equalModelVars(t, model.New(lb, vna), d.bddEquiv, false, va)
	equalModelVars(t, model.New(lb, vna), d.bddEquiv, false, va, vb)
	equalModelVars(t, model.New(vnc, vnb, vna), d.bddOr, true, va)
	equalModelVars(t, model.New(vnc, vnb, vna), d.bddOr, true, va, vb)
	equalModelVars(t, model.New(vnc, vnb, vna), d.bddOr, false, va)
	equalModelVars(t, model.New(vnc, vnb, vna), d.bddOr, false, va, vb)
	equalModelVars(t, model.New(vnc, lb, la), d.bddAnd, true, va)
	equalModelVars(t, model.New(vnc, lb, la), d.bddAnd, true, va, vb)
	equalModelVars(t, model.New(vnc, lb, la), d.bddAnd, false, va)
	equalModelVars(t, model.New(vnc, lb, la), d.bddAnd, false, va, vb)
}

func equalModelVars(t *testing.T, expected *model.Model, bdd *BDD, def bool, vars ...f.Variable) {
	model, err := bdd.ModelWithVariables(def, vars...)
	assert.Nil(t, err)
	assert.Equal(t, expected, model)
}

func TestBDDFullModel(t *testing.T) {
	fac := f.NewFactory()
	d := testData(fac)
	va := fac.Lit("A", true)
	vna := fac.Lit("A", false)
	vb := fac.Lit("B", true)
	vnb := fac.Lit("B", false)
	vnc := fac.Lit("C", false)

	m, err := d.bddFalsum.FullModel()
	assert.Nil(t, m)
	assert.NotNil(t, err)

	equalFullModel(t, model.New(vnc, vnb, vna), d.bddVerum)
	equalFullModel(t, model.New(vnc, vnb, va), d.bddPosLit)
	equalFullModel(t, model.New(vnc, vnb, vna), d.bddNegLit)
	equalFullModel(t, model.New(vnc, vnb, vna), d.bddImpl)
	equalFullModel(t, model.New(vnc, vb, vna), d.bddEquiv)
	equalFullModel(t, model.New(vnc, vnb, vna), d.bddOr)
	equalFullModel(t, model.New(vnc, vb, va), d.bddAnd)
}

func equalFullModel(t *testing.T, expected *model.Model, bdd *BDD) {
	model, err := bdd.FullModel()
	assert.Nil(t, err)
	assert.Equal(t, expected, model)
}

func TestBDDPathCount(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := testData(fac)
	assert.Equal(big.NewInt(1), d.bddVerum.PathCountOne())
	assert.Equal(big.NewInt(0), d.bddVerum.PathCountZero())
	assert.Equal(big.NewInt(0), d.bddFalsum.PathCountOne())
	assert.Equal(big.NewInt(1), d.bddFalsum.PathCountZero())
	assert.Equal(big.NewInt(1), d.bddPosLit.PathCountOne())
	assert.Equal(big.NewInt(1), d.bddPosLit.PathCountZero())
	assert.Equal(big.NewInt(1), d.bddNegLit.PathCountOne())
	assert.Equal(big.NewInt(1), d.bddNegLit.PathCountZero())
	assert.Equal(big.NewInt(2), d.bddImpl.PathCountOne())
	assert.Equal(big.NewInt(1), d.bddImpl.PathCountZero())
	assert.Equal(big.NewInt(2), d.bddEquiv.PathCountOne())
	assert.Equal(big.NewInt(2), d.bddEquiv.PathCountZero())
	assert.Equal(big.NewInt(3), d.bddOr.PathCountOne())
	assert.Equal(big.NewInt(1), d.bddOr.PathCountZero())
	assert.Equal(big.NewInt(1), d.bddAnd.PathCountOne())
	assert.Equal(big.NewInt(3), d.bddAnd.PathCountZero())
}

func TestBDDSupport(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := testData(fac)
	a := fac.Var("A")
	b := fac.Var("B")
	c := fac.Var("C")

	assert.Equal([]f.Variable{}, d.bddVerum.Support())
	assert.Equal([]f.Variable{}, d.bddFalsum.Support())
	assert.Equal([]f.Variable{a}, d.bddPosLit.Support())
	assert.Equal([]f.Variable{a}, d.bddNegLit.Support())
	assert.Equal([]f.Variable{b, a}, d.bddImpl.Support())
	assert.Equal([]f.Variable{b, a}, d.bddEquiv.Support())
	assert.Equal([]f.Variable{c, b, a}, d.bddOr.Support())
	assert.Equal([]f.Variable{c, b, a}, d.bddAnd.Support())
}

func TestBDDNodeCount(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := testData(fac)
	assert.Equal(0, d.bddVerum.NodeCount())
	assert.Equal(0, d.bddFalsum.NodeCount())
	assert.Equal(1, d.bddPosLit.NodeCount())
	assert.Equal(1, d.bddNegLit.NodeCount())
	assert.Equal(2, d.bddImpl.NodeCount())
	assert.Equal(3, d.bddEquiv.NodeCount())
	assert.Equal(3, d.bddOr.NodeCount())
	assert.Equal(3, d.bddAnd.NodeCount())
}

func TestBDDVariableProfile(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := testData(fac)
	a := fac.Var("A")
	b := fac.Var("B")
	c := fac.Var("C")

	profile := d.bddVerum.VariableProfile()
	assert.Equal(0, profile[a])
	assert.Equal(0, profile[b])
	assert.Equal(0, profile[c])

	profile = d.bddFalsum.VariableProfile()
	assert.Equal(0, profile[a])
	assert.Equal(0, profile[b])
	assert.Equal(0, profile[c])

	profile = d.bddPosLit.VariableProfile()
	assert.Equal(1, profile[a])
	assert.Equal(0, profile[b])
	assert.Equal(0, profile[c])

	profile = d.bddNegLit.VariableProfile()
	assert.Equal(1, profile[a])
	assert.Equal(0, profile[b])
	assert.Equal(0, profile[c])

	profile = d.bddImpl.VariableProfile()
	assert.Equal(1, profile[a])
	assert.Equal(1, profile[b])
	assert.Equal(0, profile[c])

	profile = d.bddEquiv.VariableProfile()
	assert.Equal(1, profile[a])
	assert.Equal(2, profile[b])
	assert.Equal(0, profile[c])

	profile = d.bddOr.VariableProfile()
	assert.Equal(1, profile[a])
	assert.Equal(1, profile[b])
	assert.Equal(1, profile[c])

	profile = d.bddAnd.VariableProfile()
	assert.Equal(1, profile[a])
	assert.Equal(1, profile[b])
	assert.Equal(1, profile[c])
}

func compareFormula(t *testing.T, fac f.Factory, bdd *BDD, compareFormula f.Formula) {
	bddFormulaFollowPathsToTrue := bdd.ToFormula(fac, true)
	bddFormulaFollowPathsToFalse := bdd.ToFormula(fac, false)
	assert.True(t, sat.IsEquivalent(fac, compareFormula, bddFormulaFollowPathsToTrue))
	assert.True(t, sat.IsEquivalent(fac, compareFormula, bddFormulaFollowPathsToFalse))
}

func TestBDDConstruction(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	variables := []f.Variable{
		fac.Var("a"), fac.Var("b"), fac.Var("c"), fac.Var("d"), fac.Var("e"), fac.Var("f"), fac.Var("g"),
	}
	kernel := NewKernelWithOrdering(fac, variables, 1000, 10000)
	initFormula := p.ParseUnsafe("(a & b) => (c | d & ~e)")
	secondFormula := p.ParseUnsafe("(g & f) <=> (c | ~a | ~d)")
	initBdd := CompileWithKernel(fac, initFormula, kernel)
	secondBdd := CompileWithKernel(fac, secondFormula, kernel)

	negation := initBdd.Negate()
	expected := CompileWithKernel(fac, initFormula.Negate(fac), kernel)
	assert.Equal(expected, negation)

	implication := initBdd.Implies(secondBdd)
	expected = CompileWithKernel(fac, fac.Implication(initFormula, secondFormula), kernel)
	assert.Equal(expected, implication)

	implication = initBdd.ImpliedBy(secondBdd)
	expected = CompileWithKernel(fac, fac.Implication(secondFormula, initFormula), kernel)
	assert.Equal(expected, implication)

	equivalence := initBdd.Equivalence(secondBdd)
	expected = CompileWithKernel(fac, fac.Equivalence(secondFormula, initFormula), kernel)
	assert.Equal(expected, equivalence)

	and := initBdd.And(secondBdd)
	expected = CompileWithKernel(fac, fac.And(secondFormula, initFormula), kernel)
	assert.Equal(expected, and)

	or := initBdd.Or(secondBdd)
	expected = CompileWithKernel(fac, fac.Or(secondFormula, initFormula), kernel)
	assert.Equal(expected, or)
}

func TestBDDModelEnumerationQueens(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	expected := []*big.Int{
		big.NewInt(0),
		big.NewInt(2),
		big.NewInt(10),
		big.NewInt(4),
		big.NewInt(40),
		big.NewInt(92),
		big.NewInt(352),
	}
	formulas := make([]f.Formula, 7)
	variables := make([]*f.VarSet, 7)
	for i, problem := range []int{3, 4, 5, 6, 7, 8, 9} {
		formulas[i] = sat.GenerateNQueens(fac, problem)
		variables[i] = f.Variables(fac, formulas[i])
	}

	for i := 0; i < len(formulas); i++ {
		kernel := NewKernel(fac, int32(variables[i].Size()), 10000, 10000)
		bdd := CompileWithKernel(fac, formulas[i], kernel)
		models := bdd.ModelEnumeration(variables[i].Content()...)
		assert.Equal(expected[i].Int64(), int64(len(models)))
		for _, model := range models {
			ass, _ := model.Assignment(fac)
			assert.True(assignment.Evaluate(fac, formulas[i], ass))
		}
	}
}

func TestBDDModelCountExo(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	vars := generateVariables(fac, 100)
	constraint := normalform.NNF(fac, fac.EXO(vars...))
	numVars := f.Variables(fac, constraint).Size()
	kernel := NewKernel(fac, int32(numVars), 100000, 1000000)
	bdd := CompileWithKernel(fac, constraint, kernel)
	assert.Equal(big.NewInt(100), bdd.ModelCount())
	assert.Equal(100, len(bdd.ModelEnumeration(vars...)))
}

func TestBDDModelCountExk(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	vars := generateVariables(fac, 15)
	constraint := normalform.NNF(fac, fac.CC(f.EQ, 8, vars...))
	numVars := f.Variables(fac, constraint).Size()
	kernel := NewKernel(fac, int32(numVars), 100000, 1000000)
	bdd := CompileWithKernel(fac, constraint, kernel)
	assert.Equal(big.NewInt(6435), bdd.ModelCount())
	assert.Equal(6435, len(bdd.ModelEnumeration(vars...)))
}

func TestBDDModelCountAmo(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	vars := generateVariables(fac, 100)
	constraint := normalform.NNF(fac, fac.AMO(vars...))
	numVars := f.Variables(fac, constraint).Size()
	kernel := NewKernel(fac, int32(numVars), 100000, 1000000)
	bdd := CompileWithKernel(fac, constraint, kernel)
	assert.Equal(big.NewInt(221), bdd.ModelCount())
	assert.Equal(101, len(bdd.ModelEnumeration(vars...)))
}

func generateVariables(fac f.Factory, n int) []f.Variable {
	result := make([]f.Variable, n)
	for i := range n {
		result[i] = fac.Var(fmt.Sprintf("v%d", i))
	}
	return result
}
