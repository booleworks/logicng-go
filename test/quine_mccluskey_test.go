package test

import (
	"testing"

	"github.com/booleworks/logicng-go/model/enum"
	"github.com/booleworks/logicng-go/simplification"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/io"
	"github.com/booleworks/logicng-go/normalform"
	"github.com/booleworks/logicng-go/parser"
	"github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

func TestQMCTrivialCases(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()

	assert.Equal(fac.Verum(), simplification.QMC(fac, fac.Verum()))
	assert.Equal(fac.Falsum(), simplification.QMC(fac, fac.Falsum()))
	assert.Equal(fac.Literal("a", true), simplification.QMC(fac, fac.Literal("a", true)))
	assert.Equal(fac.Literal("a", false), simplification.QMC(fac, fac.Literal("a", false)))
}

func TestQMCSimple1(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("(~a & ~b & ~c) | (~a & ~b & c) | (~a & b & ~c) | (a & ~b & c) | (a & b & ~c) | (a & b & c)")
	dnf := simplification.QMC(fac, formula)

	assert.True(normalform.IsDNF(fac, dnf))
	assert.True(sat.IsEquivalent(fac, formula, dnf))
}

func TestQMCSimple2(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("(~a & ~b & ~c) | (~a & b & ~c) | (a & ~b & c) | (a & b & c)")
	dnf := simplification.QMC(fac, formula)

	assert.True(normalform.IsDNF(fac, dnf))
	assert.True(sat.IsEquivalent(fac, formula, dnf))
}

func TestQMCSimple3(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("~5 & ~4 & 3 & 2 & 1 | ~3 & ~7 & ~2 & 1 | ~6 & 1 & ~3 & 2 | ~9 & 6 & 8 & ~1 | 3 & 4 & 2 & 1 | ~2 & 7 & 1 | ~10 & ~8 & ~1")
	dnf := simplification.QMC(fac, formula)

	assert.True(normalform.IsDNF(fac, dnf))
	assert.True(sat.IsEquivalent(fac, formula, dnf))
	assert.True(f.Variables(fac, formula).ContainsAll(f.Variables(fac, dnf)))
}

func TestQMCLarge1(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("A => B & ~((D | E | I | J) & ~K) & L")
	dnf := simplification.QMC(fac, formula)

	assert.True(normalform.IsDNF(fac, dnf))
	assert.True(sat.IsEquivalent(fac, formula, dnf))
}

func TestQMCLarge2(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	formula, _ := io.ReadFormula(fac, "../test/data/formulas/mid.txt")
	solver := sat.NewSolver(fac)
	solver.Add(formula)
	models := enum.OnSolver(solver, []f.Variable{
		fac.Var("v111"),
		fac.Var("v410"),
		fac.Var("v434"),
		fac.Var("v35"),
		fac.Var("v36"),
		fac.Var("v78"),
		fac.Var("v125"),
		fac.Var("v125"),
		fac.Var("v58"),
		fac.Var("v61"),
	})
	operands := make([]f.Formula, len(models))
	for i, m := range models {
		operands[i] = m.Formula(fac)
	}
	canonicalDnf := fac.Or(operands...)

	dnf := simplification.QMC(fac, canonicalDnf)
	assert.True(normalform.IsDNF(fac, dnf))
	assert.True(sat.IsEquivalent(fac, canonicalDnf, dnf))
}

func TestQMCLarge3(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	formula, _ := io.ReadFormula(fac, "../test/data/formulas/mid.txt")
	solver := sat.NewSolver(fac)
	solver.Add(formula)
	models := enum.OnSolver(solver, []f.Variable{
		fac.Var("v111"),
		fac.Var("v410"),
		fac.Var("v434"),
		fac.Var("v35"),
		fac.Var("v36"),
		fac.Var("v78"),
		fac.Var("v125"),
		fac.Var("v125"),
		fac.Var("v58"),
		fac.Var("v27"),
		fac.Var("v462"),
		fac.Var("v463"),
		fac.Var("v280"),
		fac.Var("v61"),
	})
	operands := make([]f.Formula, len(models))
	for i, m := range models {
		operands[i] = m.Formula(fac)
	}
	canonicalDnf := fac.Or(operands...)

	dnf := simplification.QMC(fac, canonicalDnf)
	assert.True(normalform.IsDNF(fac, dnf))
	assert.True(sat.IsEquivalent(fac, canonicalDnf, dnf))
}

func TestQMCSmallFormulas(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	formulas, _ := io.ReadFormulas(fac, "../test/data/formulas/small_formulas.txt")
	for _, formula := range formulas {
		variables := f.Variables(fac, formula).Content()
		projectedVars := variables[:min(6, len(variables))]
		solver := sat.NewSolver(fac)
		solver.Add(formula)
		models := enum.OnSolver(solver, projectedVars)
		operands := make([]f.Formula, len(models))
		for i, m := range models {
			operands[i] = m.Formula(fac)
		}
		canonicalDnf := fac.Or(operands...)
		dnf := simplification.QMC(fac, canonicalDnf)
		assert.True(normalform.IsDNF(fac, dnf))
		assert.True(sat.IsEquivalent(fac, canonicalDnf, dnf))
	}
}
