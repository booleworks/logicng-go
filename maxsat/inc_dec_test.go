package maxsat

import (
	"slices"
	"testing"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/parser"
	"github.com/stretchr/testify/assert"
)

func TestMaxsatIncrementalityPartial(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	solvers := []*Solver{
		WBO(fac),
		IncWBO(fac),
		OLL(fac),
		LinearSU(fac),
		LinearUS(fac),
		MSU3(fac),
	}
	for _, solver := range solvers {
		solver.AddHardFormula(p.ParseUnsafe("(~a | ~b) & (~b | ~c) & ~d"))
		solver.AddSoftFormula(p.ParseUnsafe("a"), 1)
		solver.AddSoftFormula(p.ParseUnsafe("b"), 1)
		res := solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(1, res.Optimum)
		assert.True(
			slices.Contains(res.Model.Literals, fac.Lit("a", true)) && !slices.Contains(res.Model.Literals, fac.Lit("b", true)) ||
				slices.Contains(res.Model.Literals, fac.Lit("b", true)) && !slices.Contains(res.Model.Literals, fac.Lit("a", true)))

		solver.AddSoftFormula(p.ParseUnsafe("c"), 1)
		solver.AddSoftFormula(p.ParseUnsafe("d"), 1)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(2, res.Optimum)
		assert.True(
			slices.Contains(res.Model.Literals, fac.Lit("a", true)) && !slices.Contains(res.Model.Literals, fac.Lit("b", true)) ||
				slices.Contains(res.Model.Literals, fac.Lit("b", true)) && !slices.Contains(res.Model.Literals, fac.Lit("a", true)))
		assert.True(
			slices.Contains(res.Model.Literals, fac.Lit("c", true)) && !slices.Contains(res.Model.Literals, fac.Lit("d", true)) ||
				slices.Contains(res.Model.Literals, fac.Lit("d", true)) && !slices.Contains(res.Model.Literals, fac.Lit("c", true)))

		solver.AddSoftFormula(p.ParseUnsafe("~a"), 1)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(3, res.Optimum)
	}
}

func TestMaxsatIncrementalityWeighted(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	solvers := []*Solver{
		WBO(fac),
		IncWBO(fac),
		OLL(fac),
		LinearSU(fac),
		WMSU3(fac),
	}

	for _, solver := range solvers {
		solver.AddHardFormula(p.ParseUnsafe("(~a | ~b) & (~b | ~c) & ~d"))
		solver.AddSoftFormula(p.ParseUnsafe("a"), 1)
		solver.AddSoftFormula(p.ParseUnsafe("b"), 2)
		res := solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(1, res.Optimum)
		assert.True(
			slices.Contains(res.Model.Literals, fac.Lit("b", true)) && !slices.Contains(res.Model.Literals, fac.Lit("a", true)))

		solver.AddSoftFormula(p.ParseUnsafe("c"), 1)
		solver.AddSoftFormula(p.ParseUnsafe("d"), 2)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(4, res.Optimum)

		solver.AddSoftFormula(p.ParseUnsafe("~a"), 4)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(4, res.Optimum)
	}
}

func TestMaxsatDecrementalityPartial(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	solvers := []*Solver{
		WBO(fac),
		IncWBO(fac),
		OLL(fac),
		LinearSU(fac),
		LinearUS(fac),
		MSU3(fac),
	}

	for _, solver := range solvers {
		res := solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(0, res.Optimum)
		state0 := solver.SaveState()
		assert.Equal(SolverState{int32(1), 0, 0, 0, 0, 1, []int{}}, *state0)

		solver.AddHardFormula(p.ParseUnsafe("(~a | ~b) & (~b | ~c) & ~d"))
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(0, res.Optimum)
		state1 := solver.SaveState()
		assert.Equal(SolverState{int32(3), 4, 3, 0, 0, 1, []int{}}, *state1)

		solver.AddSoftFormula(p.ParseUnsafe("a"), 1)
		solver.AddSoftFormula(p.ParseUnsafe("b"), 1)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(1, res.Optimum)
		state2 := solver.SaveState()
		assert.True(
			slices.Contains(res.Model.Literals, fac.Lit("a", true)) && !slices.Contains(res.Model.Literals, fac.Lit("b", true)) ||
				slices.Contains(res.Model.Literals, fac.Lit("b", true)) && !slices.Contains(res.Model.Literals, fac.Lit("a", true)))
		assert.Equal(SolverState{int32(5), 6, 7, 2, 2, 1, []int{1, 1}}, *state2)

		solver.LoadState(state1)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(0, res.Optimum)

		solver.AddSoftFormula(p.ParseUnsafe("a"), 1)
		solver.AddSoftFormula(p.ParseUnsafe("b"), 1)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(1, res.Optimum)
		state3 := solver.SaveState()
		assert.True(
			slices.Contains(res.Model.Literals, fac.Lit("a", true)) && !slices.Contains(res.Model.Literals, fac.Lit("b", true)) ||
				slices.Contains(res.Model.Literals, fac.Lit("b", true)) && !slices.Contains(res.Model.Literals, fac.Lit("a", true)))

		solver.AddSoftFormula(p.ParseUnsafe("c"), 1)
		solver.AddSoftFormula(p.ParseUnsafe("d"), 1)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(2, res.Optimum)
		state4 := solver.SaveState()

		solver.AddSoftFormula(p.ParseUnsafe("~a"), 1)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(3, res.Optimum)

		solver.LoadState(state4)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(2, res.Optimum)

		solver.LoadState(state3)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(1, res.Optimum)

		solver.LoadState(state0)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(0, res.Optimum)

		solver.AddSoftFormula(p.ParseUnsafe("a"), 1)
		solver.AddSoftFormula(p.ParseUnsafe("b"), 1)
		solver.AddSoftFormula(p.ParseUnsafe("~b"), 1)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(1, res.Optimum)
	}
}

func TestMaxsatDecrementalityWeighted(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	solvers := []*Solver{
		WBO(fac),
		IncWBO(fac),
		OLL(fac),
		LinearSU(fac),
		WMSU3(fac),
	}

	for _, solver := range solvers {
		solver.AddSoftFormula(p.ParseUnsafe("x"), 2)
		res := solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(0, res.Optimum)
		state0 := solver.SaveState()

		assert.Equal(SolverState{int32(1), 2, 2, 1, 2, 2, []int{2}}, *state0)
		solver.AddHardFormula(p.ParseUnsafe("(~a | ~b) & (~b | ~c) & ~d"))
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(0, res.Optimum)
		state1 := solver.SaveState()
		assert.Equal(SolverState{int32(3), 6, 5, 1, 2, 2, []int{2}}, *state1)

		solver.AddSoftFormula(p.ParseUnsafe("a"), 1)
		solver.AddSoftFormula(p.ParseUnsafe("b"), 2)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(1, res.Optimum)
		state2 := solver.SaveState()
		assert.Equal(SolverState{int32(5), 8, 9, 3, 5, 2, []int{2, 1, 2}}, *state2)

		solver.LoadState(state1)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(0, res.Optimum)

		solver.AddSoftFormula(p.ParseUnsafe("a"), 1)
		solver.AddSoftFormula(p.ParseUnsafe("b"), 2)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(1, res.Optimum)
		state3 := solver.SaveState()

		solver.AddSoftFormula(p.ParseUnsafe("c"), 1)
		solver.AddSoftFormula(p.ParseUnsafe("d"), 2)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(4, res.Optimum)
		state4 := solver.SaveState()

		solver.AddSoftFormula(p.ParseUnsafe("~a"), 1)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(4, res.Optimum)

		solver.LoadState(state4)
		assert.True(res.Satisfiable)
		assert.Equal(4, res.Optimum)

		solver.LoadState(state3)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(1, res.Optimum)

		solver.LoadState(state0)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(0, res.Optimum)

		solver.AddSoftFormula(p.ParseUnsafe("a"), 1)
		solver.AddSoftFormula(p.ParseUnsafe("b"), 1)
		solver.AddSoftFormula(p.ParseUnsafe("~b"), 1)
		res = solver.Solve()
		assert.True(res.Satisfiable)
		assert.Equal(1, res.Optimum)
	}
}
