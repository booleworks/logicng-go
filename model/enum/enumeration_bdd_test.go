package enum

import (
	"slices"
	"testing"

	"github.com/booleworks/logicng-go/bdd"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/model/iter"
	"github.com/booleworks/logicng-go/parser"
	"github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

func TestBDDMEContradiction(t *testing.T) {
	for _, cfg := range cfgs {
		assert := assert.New(t)
		fac := f.NewFactory()
		solver := sat.NewSolver(fac)
		solver.Add(fac.Literal("A", true))
		solver.Add(fac.Literal("A", false))
		result, _ := ToBDDOnSolverWithConfig(solver, []f.Variable{}, cfg)
		kernel := result.Kernel
		exp := bdd.BuildWithKernel(fac, fac.Falsum(), kernel)
		assert.Equal(exp, result)
	}
}

func TestBDDMETautology(t *testing.T) {
	for _, cfg := range cfgs {
		assert := assert.New(t)
		fac := f.NewFactory()
		solver := sat.NewSolver(fac)
		result, _ := ToBDDOnSolverWithConfig(solver, []f.Variable{}, cfg)
		kernel := result.Kernel
		exp := bdd.BuildWithKernel(fac, fac.Verum(), kernel)
		assert.Equal(exp, result)
	}
}

func TestBDDMEEmptyVars(t *testing.T) {
	for _, cfg := range cfgs {
		assert := assert.New(t)
		fac := f.NewFactory()
		p := parser.New(fac)
		solver := sat.NewSolver(fac)
		formula := p.ParseUnsafe("A & (B | C)")
		solver.Add(formula)
		result, _ := ToBDDOnSolverWithConfig(solver, nil, cfg)
		kernel := result.Kernel
		exp := bdd.BuildWithKernel(fac, fac.Verum(), kernel)
		assert.Equal(exp, result)
	}
}

func TestBDDMESimple1OnFormula(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("A & (B | C)")
	for _, cfg := range cfgs {
		result, _ := ToBDDOnFormulaWithConfig(fac, formula, fac.Vars("A", "B", "C"), cfg)
		kernel := result.Kernel
		exp := bdd.BuildWithKernel(fac, formula, kernel)
		assert.Equal(exp, result)
	}

	result := ToBDDOnFormula(fac, formula, fac.Vars("A", "B", "C"))
	kernel := result.Kernel
	exp := bdd.BuildWithKernel(fac, formula, kernel)
	assert.Equal(exp, result)
}

func TestBDDMESimple2(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("(~A | C) & (~B | C)")
	for _, cfg := range cfgs {
		solver := sat.NewSolver(fac)
		solver.Add(formula)
		result, _ := ToBDDOnSolverWithConfig(solver, fac.Vars("A", "B", "C"), cfg)
		assert.Equal(int64(5), result.ModelCount().Int64())
		assert.Equal(3, len(result.Support()))
		assert.True(slices.Contains(result.Support(), fac.Var("A")))
		assert.True(slices.Contains(result.Support(), fac.Var("B")))
		assert.True(slices.Contains(result.Support(), fac.Var("C")))
	}
}

func TestBDDMEDuplicate(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("A & (B | C)")
	for _, cfg := range cfgs {
		solver := sat.NewSolver(fac)
		solver.Add(formula)
		result, _ := ToBDDOnSolverWithConfig(solver, fac.Vars("A", "A", "B"), cfg)
		assert.Equal(int64(2), result.ModelCount().Int64())
		assert.Equal(1, len(result.Support()))
		assert.True(slices.Contains(result.Support(), fac.Var("A")))
	}
}

func TestBDDMEMultiple(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("(~A | C) & (~B | C)")
	for _, cfg := range cfgs {
		solver := sat.NewSolver(fac)
		solver.Add(formula)
		firstRun, _ := ToBDDOnSolverWithConfig(solver, fac.Vars("A", "B", "C"), cfg)
		secondRun, _ := ToBDDOnSolverWithConfig(solver, fac.Vars("A", "B", "C"), cfg)
		assert.Equal(int64(5), firstRun.ModelCount().Int64())
		assert.Equal(int64(5), secondRun.ModelCount().Int64())
	}
}

func TestBDDMEDontCares1(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("(~A | C) & (~B | C)")
	for _, cfg := range cfgs {
		solver := sat.NewSolver(fac)
		solver.Add(formula)
		result, _ := ToBDDOnSolverWithConfig(solver, fac.Vars("A", "B", "C", "D"), cfg)
		assert.Equal(int64(10), result.ModelCount().Int64())
		assert.Equal(3, len(result.Support()))
		assert.True(slices.Contains(result.Support(), fac.Var("A")))
		assert.True(slices.Contains(result.Support(), fac.Var("B")))
		assert.True(slices.Contains(result.Support(), fac.Var("C")))
	}
}

func TestBDDMEDontCares2(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("(~A | C) & (~B | C)")
	for _, cfg := range cfgs {
		solver := sat.NewSolver(fac)
		solver.Add(formula)
		result, _ := ToBDDOnSolverWithConfig(solver, fac.Vars("A", "C", "D", "E"), cfg)
		assert.Equal(int64(12), result.ModelCount().Int64())
		assert.Equal(2, len(result.Support()))
		assert.True(slices.Contains(result.Support(), fac.Var("A")))
		assert.True(slices.Contains(result.Support(), fac.Var("C")))
	}
}

func TestBDDMEDontCares3(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("A | B | (X & ~X)")
	for _, cfg := range cfgs {
		result, _ := ToBDDOnFormulaWithConfig(fac, formula, fac.Vars("A", "B", "X"), cfg)
		assert.Equal(int64(6), result.ModelCount().Int64())
		assert.Equal(2, len(result.Support()))
		assert.True(slices.Contains(result.Support(), fac.Var("A")))
		assert.True(slices.Contains(result.Support(), fac.Var("B")))
	}
}

func TestBDDMEWithLimit(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("(~A | C) & (~B | C)")
	for _, cfg := range cfgs {
		newCfg := &iter.Config{
			Handler:  iter.HandlerWithLimit(3),
			Strategy: cfg.Strategy,
		}
		result, ok := ToBDDOnFormulaWithConfig(fac, formula, fac.Vars("A", "B", "C"), newCfg)
		assert.False(ok)
		assert.True(newCfg.Handler.Aborted())
		assert.Equal(int64(3), result.ModelCount().Int64())
	}
}
