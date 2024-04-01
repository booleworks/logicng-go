package sat

import (
	"slices"
	"testing"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/parser"
	"github.com/stretchr/testify/assert"
)

func TestFormulaOnSolver(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	formulas := []f.Formula{
		p.ParseUnsafe("A | B | C"),
		p.ParseUnsafe("~A | ~B | ~C"),
		p.ParseUnsafe("A | ~B"),
		p.ParseUnsafe("A"),
	}
	solver := NewSolver(fac)
	solver.Add(formulas...)
	compareFormulas(t, fac, formulas, solver.FormulasOnSolver())

	formulas = append(formulas, p.ParseUnsafe("~A | C"))
	solver = NewSolver(fac)
	solver.Add(formulas...)
	compareFormulas(t, fac, formulas, solver.FormulasOnSolver())
}

func TestFormulaOnSolverContradiction(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	solver := NewSolver(fac)
	solver.Add(fac.Variable("A"))
	solver.Add(fac.Variable("B"))
	solver.Add(p.ParseUnsafe("C & (~A | ~B)"))
	res := solver.FormulasOnSolver()
	assert.Equal(4, len(res))
	assert.True(slices.Contains(res, fac.Variable("A")))
	assert.True(slices.Contains(res, fac.Variable("B")))
	assert.True(slices.Contains(res, fac.Variable("C")))
	assert.True(slices.Contains(res, fac.Falsum()))

	solver = NewSolver(fac)
	solver.Add(p.ParseUnsafe("A <=> B"))
	solver.Add(p.ParseUnsafe("B <=> ~A"))
	res = solver.FormulasOnSolver()
	assert.Equal(4, len(res))
	assert.True(slices.Contains(res, p.ParseUnsafe("A | ~B")))
	assert.True(slices.Contains(res, p.ParseUnsafe("~A | B")))
	assert.True(slices.Contains(res, p.ParseUnsafe("~A | ~B")))
	assert.True(slices.Contains(res, p.ParseUnsafe("A | B")))

	solver.Sat()
	res = solver.FormulasOnSolver()
	assert.Equal(7, len(res))
	assert.True(slices.Contains(res, p.ParseUnsafe("~B | A")))
	assert.True(slices.Contains(res, p.ParseUnsafe("~A | ~B")))
	assert.True(slices.Contains(res, p.ParseUnsafe("~A | ~B")))
	assert.True(slices.Contains(res, p.ParseUnsafe("A | B")))
	assert.True(slices.Contains(res, p.ParseUnsafe("A")))
	assert.True(slices.Contains(res, p.ParseUnsafe("B")))
	assert.True(slices.Contains(res, fac.Falsum()))
}

func compareFormulas(t *testing.T, fac f.Factory, original, solver []f.Formula) {
	orig := fac.And(original...)
	onSolver := fac.And(solver...)
	assert.True(t, IsEquivalent(fac, orig, onSolver))
}
