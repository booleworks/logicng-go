package pbc_test

import (
	"fmt"
	"testing"

	"github.com/booleworks/logicng-go/model/enum"
	"github.com/booleworks/logicng-go/sat"

	"github.com/booleworks/logicng-go/assignment"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/stretchr/testify/assert"
)

func TestPbOnSolverLess(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	for _, solver := range getSolvers(fac) {
		coeffs10 := []int{3, 2, 2, 2, 2, 2, 2, 2, 2, 2}
		vars10 := make([]f.Variable, 10)
		for i := 0; i < 10; i++ {
			vars10[i] = fac.Var(fmt.Sprintf("v%d", i))
		}
		literals10 := f.VariablesAsLiterals(vars10)
		solver.Add(fac.PBC(f.LE, 6, literals10, coeffs10))
		assert.True(solver.Sat())
		models := enum.OnSolver(solver, vars10)
		assert.Len(models, 140)
		for _, model := range models {
			assert.True(len(model.PosVars()) <= 3)
		}

		solver.Reset()
		solver.Add(fac.PBC(f.LT, 7, literals10, coeffs10))
		assert.True(solver.Sat())
		models = enum.OnSolver(solver, vars10)
		assert.Len(models, 140)
		for _, model := range models {
			assert.True(len(model.PosVars()) <= 3)
		}

		solver.Reset()
		solver.Add(fac.PBC(f.LE, 0, literals10, coeffs10))
		assert.True(solver.Sat())
		models = enum.OnSolver(solver, vars10)
		assert.Len(models, 1)

		solver.Reset()
		solver.Add(fac.PBC(f.LE, 1, literals10, coeffs10))
		assert.True(solver.Sat())
		models = enum.OnSolver(solver, vars10)
		assert.Len(models, 1)

		solver.Reset()
		solver.Add(fac.PBC(f.LT, 2, literals10, coeffs10))
		assert.True(solver.Sat())
		models = enum.OnSolver(solver, vars10)
		assert.Len(models, 1)

		solver.Reset()
		solver.Add(fac.PBC(f.LT, 1, literals10, coeffs10))
		assert.True(solver.Sat())
		models = enum.OnSolver(solver, vars10)
		assert.Len(models, 1)
	}
}

func TestPbOnSolverGreater(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	for _, solver := range getSolvers(fac) {
		coeffs10 := []int{3, 2, 2, 2, 2, 2, 2, 2, 2, 2}
		vars10 := make([]f.Variable, 10)
		for i := 0; i < 10; i++ {
			vars10[i] = fac.Var(fmt.Sprintf("v%d", i))
		}
		literals10 := f.VariablesAsLiterals(vars10)
		solver.Add(fac.PBC(f.GE, 17, literals10, coeffs10))
		assert.True(solver.Sat())
		models := enum.OnSolver(solver, vars10)
		assert.Len(models, 47)
		for _, model := range models {
			assert.True(len(model.PosVars()) >= 8)
		}

		solver.Reset()
		solver.Add(fac.PBC(f.GT, 16, literals10, coeffs10))
		assert.True(solver.Sat())
		models = enum.OnSolver(solver, vars10)
		assert.Len(models, 47)
		for _, model := range models {
			assert.True(len(model.PosVars()) >= 8)
		}

		solver.Reset()
		solver.Add(fac.PBC(f.GE, 21, literals10, coeffs10))
		assert.True(solver.Sat())
		models = enum.OnSolver(solver, vars10)
		assert.Len(models, 1)

		solver.Reset()
		solver.Add(fac.PBC(f.GE, 22, literals10, coeffs10))
		assert.False(solver.Sat())

		solver.Reset()
		solver.Add(fac.PBC(f.GT, 42, literals10, coeffs10))
		assert.False(solver.Sat())
	}
}

func TestPbOnSolverEq(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	for _, solver := range getSolvers(fac) {
		coeffs10 := []int{3, 2, 2, 2, 2, 2, 2, 2, 2, 2}
		vars10 := make([]f.Variable, 10)
		for i := 0; i < 10; i++ {
			vars10[i] = fac.Var(fmt.Sprintf("v%d", i))
		}
		literals10 := f.VariablesAsLiterals(vars10)
		solver.Add(fac.PBC(f.EQ, 5, literals10, coeffs10))
		assert.True(solver.Sat())
		models := enum.OnSolver(solver, vars10)
		assert.Len(models, 9)
		for _, model := range models {
			assert.True(len(model.PosVars()) == 2)
			assert.Contains(model.PosVars(), fac.Var("v0"))
		}

		solver.Reset()
		solver.Add(fac.PBC(f.EQ, 7, literals10, coeffs10))
		assert.True(solver.Sat())
		models = enum.OnSolver(solver, vars10)
		assert.Len(models, 36)
		for _, model := range models {
			assert.True(len(model.PosVars()) == 3)
			assert.Contains(model.PosVars(), fac.Var("v0"))
		}

		solver.Reset()
		solver.Add(fac.PBC(f.EQ, 0, literals10, coeffs10))
		assert.True(solver.Sat())
		models = enum.OnSolver(solver, vars10)
		assert.Len(models, 1)

		solver.Reset()
		solver.Add(fac.PBC(f.EQ, 1, literals10, coeffs10))
		assert.False(solver.Sat())

		solver.Reset()
		solver.Add(fac.PBC(f.EQ, 22, literals10, coeffs10))
		assert.False(solver.Sat())
	}
}

func TestPbOnSolverNegative(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	for _, solver := range getSolvers(fac) {
		coeffs10 := []int{2, 2, 2, 2, 2, 2, 2, 2, 2, -2}
		vars10 := make([]f.Variable, 10)
		for i := 0; i < 10; i++ {
			vars10[i] = fac.Var(fmt.Sprintf("v%d", i))
		}
		literals10 := f.VariablesAsLiterals(vars10)
		solver.Add(fac.PBC(f.EQ, 2, literals10, coeffs10))
		assert.True(solver.Sat())
		models := enum.OnSolver(solver, vars10)
		assert.Len(models, 45)

		solver.Reset()
		solver.Add(fac.PBC(f.EQ, 4, literals10, coeffs10))
		assert.True(solver.Sat())
		models = enum.OnSolver(solver, vars10)
		assert.Len(models, 120)

		solver.Reset()
		coeffs10 = []int{2, 2, -3, 2, -7, 2, 2, 2, 2, -2}
		solver.Add(fac.PBC(f.EQ, 4, literals10, coeffs10))
		assert.True(solver.Sat())
		models = enum.OnSolver(solver, vars10)
		assert.Len(models, 57)

		solver.Reset()
		coeffs10 = []int{2, 2, -3, 2, -7, 2, 2, 2, 2, -2}
		solver.Add(fac.PBC(f.EQ, -10, literals10, coeffs10))
		assert.True(solver.Sat())
		models = enum.OnSolver(solver, vars10)
		assert.Len(models, 8)

		solver.Reset()
		coeffs10 = []int{2, 2, -4, 2, -6, 2, 2, 2, 2, -2}
		solver.Add(fac.PBC(f.EQ, -12, literals10, coeffs10))
		assert.True(solver.Sat())
		models = enum.OnSolver(solver, vars10)
		assert.Len(models, 1)
	}
}

func TestPbOnSolverLarge(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	solver := getSolvers(fac)[0]
	numLits := 100
	coeffs := make([]int, numLits)
	vars := make([]f.Variable, numLits)
	for i := 0; i < numLits; i++ {
		vars[i] = fac.Var(fmt.Sprintf("v%d", i))
		coeffs[i] = i + 1
	}
	lits := f.VariablesAsLiterals(vars)
	pbc := fac.PBC(f.GE, 5000, lits, coeffs)
	solver.Add(pbc)
	sResult := solver.Call(sat.Params().ModelIfSat(vars))
	assert.True(sResult.Sat())
	ass, _ := sResult.Model().Assignment(fac)
	assert.True(assignment.Evaluate(fac, pbc, ass))
}
