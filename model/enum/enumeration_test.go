package enum

import (
	"math"
	"testing"

	"booleworks.com/logicng/model/iter"
	"booleworks.com/logicng/randomizer"
	"booleworks.com/logicng/sat"

	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/parser"
	"github.com/stretchr/testify/assert"
)

var cfgs = []*iter.Config{
	iter.DefaultConfig(),
	{Handler: nil, Strategy: iter.NewNoSplitMEStrategy()},
	{Handler: nil, Strategy: iter.NewBasicStrategy(iter.DefaultLeastCommonVarProvider(), 2)},
	{Handler: nil, Strategy: iter.NewBasicStrategy(iter.DefaultMostCommonVarProvider(), 2)},
}

func TestMEContradiction(t *testing.T) {
	for _, cfg := range cfgs {
		assert := assert.New(t)
		fac := f.NewFactory()
		solver := sat.NewSolver(fac)
		solver.Add(fac.Literal("A", true))
		solver.Add(fac.Literal("A", false))
		models, _ := OnSolverWithConfig(solver, []f.Variable{}, cfg, fac.Vars("A", "B")...)
		assert.Equal(0, len(models))
	}
}

func TestMETautology(t *testing.T) {
	for _, cfg := range cfgs {
		assert := assert.New(t)
		fac := f.NewFactory()
		solver := sat.NewSolver(fac)
		models, _ := OnSolverWithConfig(solver, []f.Variable{}, cfg, fac.Vars("A", "B")...)
		assert.Equal(1, len(models))
		for _, m := range models {
			vars := f.Variables(fac, m.Formula(fac))
			assert.Equal(2, vars.Size())
			assert.True(vars.ContainsAll(f.NewVarSet(fac.Vars("A", "B")...)))
		}
	}
}

func TestMEEmptyVars(t *testing.T) {
	for _, cfg := range cfgs {
		assert := assert.New(t)
		fac := f.NewFactory()
		p := parser.New(fac)
		solver := sat.NewSolver(fac)
		formula := p.ParseUnsafe("A & (B | C)")
		solver.Add(formula)
		models, _ := OnSolverWithConfig(solver, nil, cfg)
		assert.Equal(1, len(models))
		for _, m := range models {
			vars := f.Variables(fac, m.Formula(fac))
			assert.Equal(0, vars.Size())
		}
		models, _ = OnSolverWithConfig(solver, []f.Variable{}, cfg, fac.Vars("A", "B", "C")...)
		assert.Equal(1, len(models))
		for _, m := range models {
			vars := f.Variables(fac, m.Formula(fac))
			assert.Equal(3, vars.Size())
		}
	}
}

func TestMESimple1OnFormula(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("A & (B | C)")
	for _, cfg := range cfgs {
		models, _ := OnFormulaWithConfig(fac, formula, fac.Vars("A", "B", "C"), cfg)
		assert.Equal(3, len(models))
		for _, m := range models {
			vars := f.Variables(fac, m.Formula(fac))
			assert.Equal(3, vars.Size())
			assert.True(vars.ContainsAll(f.NewVarSet(fac.Vars("A", "B", "C")...)))
		}
	}

	models := OnFormula(fac, formula, fac.Vars("A", "B", "C"))
	assert.Equal(3, len(models))
	for _, m := range models {
		vars := f.Variables(fac, m.Formula(fac))
		assert.Equal(3, vars.Size())
		assert.True(vars.ContainsAll(f.NewVarSet(fac.Vars("A", "B", "C")...)))
	}
}

func TestMESimple1(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("A & (B | C)")
	for _, cfg := range cfgs {
		solver := sat.NewSolver(fac)
		solver.Add(formula)
		models, _ := OnSolverWithConfig(solver, fac.Vars("A", "B", "C"), cfg)
		assert.Equal(3, len(models))
		for _, m := range models {
			vars := f.Variables(fac, m.Formula(fac))
			assert.Equal(3, vars.Size())
			assert.True(vars.ContainsAll(f.NewVarSet(fac.Vars("A", "B", "C")...)))
		}
	}
}

func TestMESimple2(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("(~A | C) & (~B | C)")
	for _, cfg := range cfgs {
		solver := sat.NewSolver(fac)
		solver.Add(formula)
		models, _ := OnSolverWithConfig(solver, fac.Vars("A", "B", "C"), cfg)
		assert.Equal(5, len(models))
		for _, m := range models {
			vars := f.Variables(fac, m.Formula(fac))
			assert.Equal(3, vars.Size())
			assert.True(vars.ContainsAll(f.NewVarSet(fac.Vars("A", "B", "C")...)))
		}
	}
}

func TestMEDuplicate(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("A & (B | C)")
	for _, cfg := range cfgs {
		solver := sat.NewSolver(fac)
		solver.Add(formula)
		models, _ := OnSolverWithConfig(solver, fac.Vars("A", "A", "B"), cfg)
		assert.Equal(2, len(models))
		for _, m := range models {
			vars := f.Variables(fac, m.Formula(fac))
			assert.Equal(2, vars.Size())
			assert.True(vars.ContainsAll(f.NewVarSet(fac.Vars("A", "B")...)))
		}
	}
}

func TestMEMultiple(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("(~A | C) & (~B | C)")
	for _, cfg := range cfgs {
		solver := sat.NewSolver(fac)
		solver.Add(formula)
		firstRun, _ := OnSolverWithConfig(solver, fac.Vars("A", "B", "C"), cfg)
		secondRun, _ := OnSolverWithConfig(solver, fac.Vars("A", "B", "C"), cfg)
		assert.Equal(5, len(firstRun))
		assert.Equal(5, len(secondRun))
	}
}

func TestMEAddVarsSimple(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("A & C | B & ~C")
	for _, cfg := range cfgs {
		solver := sat.NewSolver(fac)
		solver.Add(formula)
		models, _ := OnSolverWithConfig(solver, fac.Vars("A", "B"), cfg, fac.Var("C"))
		assert.Equal(3, len(models))
		for _, m := range models {
			vars := f.Variables(fac, m.Formula(fac))
			assert.Equal(3, vars.Size())
			assert.True(vars.ContainsAll(f.NewVarSet(fac.Vars("A", "B", "C")...)))
		}
	}
}

func TestMEAddVarsDuplicates(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("A & (B | C)")
	for _, cfg := range cfgs {
		solver := sat.NewSolver(fac)
		solver.Add(formula)
		models, _ := OnSolverWithConfig(solver, fac.Vars("A"), cfg, fac.Var("B"), fac.Var("B"))
		assert.Equal(1, len(models))
		for _, m := range models {
			vars := f.Variables(fac, m.Formula(fac))
			assert.Equal(2, vars.Size())
			assert.True(vars.ContainsAll(f.NewVarSet(fac.Vars("A", "B")...)))
		}
	}
}

func TestMEDontCares1(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("(~A | C) & (~B | C)")
	for _, cfg := range cfgs {
		solver := sat.NewSolver(fac)
		solver.Add(formula)
		models, _ := OnSolverWithConfig(solver, fac.Vars("A", "B", "C", "D"), cfg)
		assert.Equal(10, len(models))
		for _, m := range models {
			vars := f.Variables(fac, m.Formula(fac))
			assert.Equal(4, vars.Size())
			assert.True(vars.ContainsAll(f.NewVarSet(fac.Vars("A", "B", "C", "D")...)))
		}
	}
}

func TestMEDontCares2(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("(~A | C) & (~B | C)")
	for _, cfg := range cfgs {
		solver := sat.NewSolver(fac)
		solver.Add(formula)
		models, _ := OnSolverWithConfig(solver, fac.Vars("A", "C", "D", "E"), cfg)
		assert.Equal(12, len(models))
		for _, m := range models {
			vars := f.Variables(fac, m.Formula(fac))
			assert.Equal(4, vars.Size())
			assert.True(vars.ContainsAll(f.NewVarSet(fac.Vars("A", "C", "D", "E")...)))
		}
	}
}

func TestMEDontCares3(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("A | B | (X & ~X)")
	for _, cfg := range cfgs {
		models, _ := OnFormulaWithConfig(fac, formula, fac.Vars("A", "B", "X"), cfg)
		assert.Equal(6, len(models))
		for _, m := range models {
			vars := f.Variables(fac, m.Formula(fac))
			assert.Equal(3, vars.Size())
			assert.True(vars.ContainsAll(f.NewVarSet(fac.Vars("A", "B", "X")...)))
		}
	}
}

func TestMEWithLimit(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	formula := p.ParseUnsafe("(~A | C) & (~B | C)")
	for _, cfg := range cfgs {
		newCfg := &iter.Config{
			Handler:  iter.HandlerWithLimit(3),
			Strategy: cfg.Strategy,
		}
		models, ok := OnFormulaWithConfig(fac, formula, fac.Vars("A", "B", "C"), newCfg)
		assert.False(ok)
		assert.True(newCfg.Handler.Aborted())
		assert.Equal(3, len(models))
		for _, m := range models {
			vars := f.Variables(fac, m.Formula(fac))
			assert.Equal(3, vars.Size())
			assert.True(vars.ContainsAll(f.NewVarSet(fac.Vars("A", "B", "C")...)))
		}
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
			formula := randomizer.Formula(3)
			solver := sat.NewSolver(fac)
			solver.Add(formula)

			varsFormula := f.Variables(fac, formula).Content()
			numberOfVars := len(varsFormula)
			minNumberOfVars := int(math.Ceil(float64(numberOfVars) / 5))
			pmeVars := varsFormula[:minNumberOfVars]

			models, _ := OnSolverWithConfig(solver, pmeVars, cfg)
			for _, m := range models {
				assert.True(solver.Sat(m.Literals...))
			}
		}
	}
}

func TestMERandomAdditionalVars(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	for _, cfg := range cfgs {
		for i := 0; i < 500; i++ {
			config := randomizer.DefaultConfig()
			config.Seed = int64(i)
			config.NumVars = 20
			randomizer := randomizer.New(fac, config)
			formula := randomizer.Formula(4)
			solver := sat.NewSolver(fac)
			solver.Add(formula)

			varsFormula := f.Variables(fac, formula).Content()
			numberOfVars := len(varsFormula)
			minNumberOfVars := int(math.Ceil(float64(numberOfVars) / 5))
			pmeVars := varsFormula[:minNumberOfVars]
			additionalVarsStart := min(4*minNumberOfVars, numberOfVars)
			additionalVars := varsFormula[additionalVarsStart:]

			models, _ := OnSolverWithConfig(solver, pmeVars, cfg, additionalVars...)
			for _, m := range models {
				vars := f.Variables(fac, m.Formula(fac))
				assert.True(vars.ContainsAll(f.NewVarSet(additionalVars...)))
				assert.True(solver.Sat(m.Literals...))
			}
		}
	}
}
