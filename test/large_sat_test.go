package test

import (
	"math/big"
	"testing"
	"time"

	"github.com/booleworks/logicng-go/model/count"
	"github.com/booleworks/logicng-go/model/enum"

	"github.com/booleworks/logicng-go/assignment"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/io"
	"github.com/booleworks/logicng-go/model"
	"github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

func TestSolverMidFormula(t *testing.T) {
	fac := f.NewFactory()
	start := time.Now()
	formula, _ := io.ReadFormula(fac, "data/formulas/mid.txt")
	elapsed := time.Since(start) / 1_000_000
	atoms := f.NumberOfAtoms(fac, formula)
	t.Logf("Read formula (%d atoms): %d ms", atoms, elapsed)

	start = time.Now()
	solver := sat.NewSolver(fac)
	solver.Add(formula)
	elapsed = time.Since(start) / 1_000_000
	t.Logf("Added to solver: %d ms", elapsed)

	start = time.Now()
	sat := solver.Call(sat.WithModel(f.Variables(fac, formula).Content()))
	elapsed = time.Since(start) / 1_000_000
	t.Logf("Solved: %d ms", elapsed)

	assert.True(t, sat.Sat())
	validateModel(t, fac, formula, sat.Model())
}

func TestSolverLargeFormula(t *testing.T) {
	fac := f.NewFactory()
	start := time.Now()
	formula, _ := io.ReadFormula(fac, "data/formulas/large.txt")
	elapsed := time.Since(start) / 1_000_000
	atoms := f.NumberOfAtoms(fac, formula)
	t.Logf("Read formula (%d atoms): %d ms", atoms, elapsed)

	start = time.Now()
	solver := sat.NewSolver(fac)
	ops, _ := fac.NaryOperands(formula)
	for _, op := range ops {
		solver.Add(op)
	}
	elapsed = time.Since(start) / 1_000_000
	t.Logf("Added to solver: %d ms", elapsed)

	start = time.Now()
	satResult := solver.Call(sat.Params())
	elapsed = time.Since(start) / 1_000_000
	t.Logf("Solved: %d ms", elapsed)

	assert.True(t, satResult.Sat())
	satResult = solver.Call(sat.WithModel(f.Variables(fac, formula).Content()))
	validateModel(t, fac, formula, satResult.Model())

	if !testing.Short() {
		start = time.Now()
		vars := fac.Vars("v1697", "v3043", "v30", "v183", "v2222", "v1", "v2902", "v1111", "v77", "v690",
			"v711", "v999", "v3111", "v21")
		models := enum.OnSolver(solver, vars)
		elapsed = time.Since(start) / 1_000_000
		t.Logf("Enumeration (%d models): %d ms", len(models), elapsed)
		assert.Equal(t, 4608, len(models))

		start = time.Now()
		count := count.OnSolver(solver, vars)
		elapsed = time.Since(start) / 1_000_000
		t.Logf("Counting (%s models): %d ms", count, elapsed)
		assert.Equal(t, big.NewInt(4608), count)
	}

	start = time.Now()
	bb := solver.ComputeBackbone(fac, f.Variables(fac, formula).Content())
	elapsed = time.Since(start) / 1_000_000
	t.Logf("Backbone (%d pos, %d neg): %d ms", len(bb.Positive), len(bb.Negative), elapsed)
	assert.Equal(t, 11, len(bb.Positive))
	assert.Equal(t, 30, len(bb.Negative))
}

func validateModel(t *testing.T, fac f.Factory, formula f.Formula, model *model.Model) {
	ass, _ := model.Assignment(fac)
	eval := assignment.Evaluate(fac, formula, ass)
	assert.True(t, eval)
}
