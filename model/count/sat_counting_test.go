package count

import (
	"math/big"
	"testing"

	"github.com/booleworks/logicng-go/normalform"
	"github.com/booleworks/logicng-go/randomizer"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/model/iter"
	"github.com/booleworks/logicng-go/parser"
	"github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

var cfgs = []*iter.Config{
	//iter.DefaultConfig(),
	//{handler: nil, Strategy: iter.NewNoSplitMEStrategy()},
	//{handler: nil, Strategy: iter.NewBasicStrategy(iter.DefaultLeastCommonVarProvider(), 2)},
	{Handler: nil, Strategy: iter.NewBasicStrategy(iter.DefaultMostCommonVarProvider(), 2)},
}

func TestMCContradiction(t *testing.T) {
	for _, cfg := range cfgs {
		assert := assert.New(t)
		fac := f.NewFactory()
		solver := sat.NewSolver(fac)
		solver.Add(fac.Literal("A", true))
		solver.Add(fac.Literal("A", false))
		count, _ := OnSolverWithConfig(solver, []f.Variable{}, cfg)
		assert.Equal(big.NewInt(0), count)
	}
}

func TestMCTautology(t *testing.T) {
	for _, cfg := range cfgs {
		assert := assert.New(t)
		fac := f.NewFactory()
		solver := sat.NewSolver(fac)
		count, _ := OnSolverWithConfig(solver, []f.Variable{}, cfg)
		assert.Equal(big.NewInt(1), count)
	}
}

func TestMCEmptyVars(t *testing.T) {
	for _, cfg := range cfgs {
		assert := assert.New(t)
		fac := f.NewFactory()
		p := parser.New(fac)
		solver := sat.NewSolver(fac)
		formula := p.ParseUnsafe("A & (B | C)")
		solver.Add(formula)
		count, _ := OnSolverWithConfig(solver, nil, cfg)
		assert.Equal(big.NewInt(1), count)
	}
}

func TestMCSimple1OnFormula(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("A & (B | C)")
	for _, cfg := range cfgs {
		count, _ := OnFormulaWithConfig(fac, formula, fac.Vars("A", "B", "C"), cfg)
		assert.Equal(big.NewInt(3), count)
	}

	count := OnFormula(fac, formula, fac.Vars("A", "B", "C"))
	assert.Equal(big.NewInt(3), count)
}

func TestMCSimple2(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("(~A | C) & (~B | C)")
	for _, cfg := range cfgs {
		solver := sat.NewSolver(fac)
		solver.Add(formula)
		count, _ := OnSolverWithConfig(solver, fac.Vars("A", "B", "C"), cfg)
		assert.Equal(big.NewInt(5), count)
	}
}

func TestMCDontCares1(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("(~A | C) & (~B | C)")
	for _, cfg := range cfgs {
		solver := sat.NewSolver(fac)
		solver.Add(formula)
		count, _ := OnSolverWithConfig(solver, fac.Vars("A", "B", "C", "D"), cfg)
		assert.Equal(big.NewInt(10), count)
	}
}

func TestMCDontCares2(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("(~A | C) & (~B | C)")
	for _, cfg := range cfgs {
		solver := sat.NewSolver(fac)
		solver.Add(formula)
		count, _ := OnSolverWithConfig(solver, fac.Vars("A", "C", "D", "E"), cfg)
		assert.Equal(big.NewInt(12), count)
	}
}

func TestMCDontCares3(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("A | B | (X & ~X)")
	for _, cfg := range cfgs {
		count, _ := OnFormulaWithConfig(fac, formula, fac.Vars("A", "B", "X"), cfg)
		assert.Equal(big.NewInt(6), count)
	}
}

func TestMCWithLimit(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("(~A | C) & (~B | C)")
	for _, cfg := range cfgs {
		newCfg := &iter.Config{
			Handler:  iter.HandlerWithLimit(3),
			Strategy: cfg.Strategy,
		}
		count, ok := OnFormulaWithConfig(fac, formula, fac.Vars("A", "B", "C"), newCfg)
		assert.False(ok)
		assert.True(newCfg.Handler.Aborted())
		assert.Equal(big.NewInt(3), count)
	}
}

func TestMERandom(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	for _, cfg := range cfgs {
		for i := 0; i < 100; i++ {
			config := randomizer.DefaultConfig()
			config.Seed = int64(i)
			config.NumVars = 20
			randomizer := randomizer.New(fac, config)
			formula := normalform.CNF(fac, randomizer.Formula(2))
			vars := f.Variables(fac, formula).Content()
			solver := sat.NewSolver(fac)
			solver.Add(formula)

			count, _ := OnSolverWithConfig(solver, vars, cfg)
			exp, _ := Count(fac, vars, formula)
			assert.Equal(exp, count)
		}
	}
}
