package sat

import (
	"slices"
	"testing"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/io"
	"github.com/stretchr/testify/assert"
)

func TestSolverBackboneSimple(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	solver := NewSolver(fac)
	vx := fac.Var("x")
	vy := fac.Var("y")
	vz := fac.Var("z")
	vu := fac.Var("u")
	vv := fac.Var("v")

	x := fac.Literal("x", true)
	y := fac.Literal("y", true)
	z := fac.Literal("z", true)
	u := fac.Literal("u", true)
	v := fac.Literal("v", true)

	variables := []f.Variable{vx, vy, vz, vu, vv}
	formula := fac.Verum()

	before := solver.SaveState()
	solver.Add(formula)
	bb := solver.ComputeBackbone(fac, []f.Variable{})
	assert.Equal(true, bb.Sat)
	assert.Equal([]f.Variable{}, bb.Positive)
	assert.Equal([]f.Variable{}, bb.Negative)
	assert.Equal([]f.Variable{}, bb.Optional)
	solver.LoadState(before)

	formula = x
	before = solver.SaveState()
	solver.Add(formula)
	bb = solver.ComputeBackbone(fac, variables)
	assert.Equal(true, bb.Sat)
	assert.Equal([]f.Variable{vx}, bb.Positive)
	assert.Equal([]f.Variable{}, bb.Negative)
	assert.Equal([]f.Variable{vy, vz, vu, vv}, bb.Optional)
	solver.LoadState(before)

	formula = fac.And(x, y)
	before = solver.SaveState()
	solver.Add(formula)
	bb = solver.ComputeBackbone(fac, variables)
	assert.Equal(true, bb.Sat)
	solver.LoadState(before)

	formula = fac.Or(x, y)
	before = solver.SaveState()
	solver.Add(formula)
	bb = solver.ComputeBackbone(fac, variables)
	assert.Equal([]f.Variable{}, bb.Positive)
	assert.Equal([]f.Variable{}, bb.Negative)
	assert.Equal([]f.Variable{vx, vy, vz, vu, vv}, bb.Optional)
	solver.LoadState(before)

	formula = x.Negate(fac)
	before = solver.SaveState()
	solver.Add(formula)
	bb = solver.ComputeBackbone(fac, variables)
	assert.Equal([]f.Variable{}, bb.Positive)
	assert.Equal([]f.Variable{vx}, bb.Negative)
	assert.Equal([]f.Variable{vy, vz, vu, vv}, bb.Optional)
	solver.LoadState(before)

	formula = fac.Or(fac.And(x, y, z), fac.And(x, y, u), fac.And(x, u, z))
	before = solver.SaveState()
	solver.Add(formula)
	bb = solver.ComputeBackbone(fac, variables)
	assert.Equal([]f.Variable{vx}, bb.Positive)
	assert.Equal([]f.Variable{}, bb.Negative)
	assert.Equal([]f.Variable{vy, vz, vu, vv}, bb.Optional)
	solver.LoadState(before)

	formula = fac.And(fac.Or(x, y, z), fac.Or(x, y, u), fac.Or(x, u, z))
	before = solver.SaveState()
	solver.Add(formula)
	bb = solver.ComputeBackbone(fac, variables)
	assert.Equal([]f.Variable{}, bb.Positive)
	assert.Equal([]f.Variable{}, bb.Negative)
	assert.Equal([]f.Variable{vx, vy, vz, vu, vv}, bb.Optional)
	solver.LoadState(before)

	formula = fac.And(fac.Or(x.Negate(fac), y), x)
	before = solver.SaveState()
	solver.Add(formula)
	bb = solver.ComputeBackbone(fac, variables)
	assert.Equal([]f.Variable{vx, vy}, bb.Positive)
	assert.Equal([]f.Variable{}, bb.Negative)
	assert.Equal([]f.Variable{vz, vu, vv}, bb.Optional)
	solver.LoadState(before)

	formula = fac.And(fac.Or(x, y), fac.Or(x.Negate(fac), y))
	before = solver.SaveState()
	solver.Add(formula)
	bb = solver.ComputeBackbone(fac, variables)
	assert.Equal([]f.Variable{vy}, bb.Positive)
	assert.Equal([]f.Variable{}, bb.Negative)
	assert.Equal([]f.Variable{vx, vz, vu, vv}, bb.Optional)
	solver.LoadState(before)

	formula = fac.And(fac.And(fac.Or(x.Negate(fac), y), x.Negate(fac)), fac.And(z, fac.Or(x, y)))
	before = solver.SaveState()
	solver.Add(formula)
	bb = solver.ComputeBackbone(fac, variables)
	assert.Equal([]f.Variable{vy, vz}, bb.Positive)
	assert.Equal([]f.Variable{vx}, bb.Negative)
	assert.Equal([]f.Variable{vu, vv}, bb.Optional)
	solver.LoadState(before)

	formula = fac.And(fac.Or(x, y), fac.Or(u, v), z)
	solver.Add(formula)
	bb = solver.ComputeBackbone(fac, variables)
	assert.Equal([]f.Variable{vz}, bb.Positive)
	assert.Equal([]f.Variable{}, bb.Negative)
	assert.Equal([]f.Variable{vx, vy, vu, vv}, bb.Optional)
}

func TestSolverBackboneSmallFormulas(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	solver := NewSolver(fac)
	formula, _ := io.ReadFormula(fac, "../test/data/formulas/small_formulas.txt")
	solver.Add(formula)
	variables := f.Variables(fac, formula).Content()
	bb := solver.ComputeBackbone(fac, variables)
	assert.True(verifyBackbone(fac, bb, formula, variables))
}

func TestSolverBackboneMidFormula(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	solver := NewSolver(fac)
	formula, _ := io.ReadFormula(fac, "../test/data/formulas/mid.txt")
	solver.Add(formula)
	variables := f.Variables(fac, formula).Content()
	bb := solver.ComputeBackbone(fac, variables)
	assert.True(verifyBackbone(fac, bb, formula, variables))
}

func verifyBackbone(fac f.Factory, bb *Backbone, formula f.Formula, variables []f.Variable) bool {
	solver := NewSolver(fac)
	solver.Add(formula)
	for _, bbVar := range bb.Positive {
		if solver.Call(Params().Literal(bbVar.Negate(fac))).Sat() {
			return false
		}
	}
	for _, bbVar := range bb.Negative {
		if solver.Call(Params().Variable(bbVar)).Sat() {
			return false
		}
	}
	for _, variable := range variables {
		if !sliceContains(bb.Positive, variable) && !sliceContains(bb.Negative, variable) {
			if !solver.Call(Params().Variable(variable)).Sat() {
				return false
			}
			if !solver.Call(Params().Literal(variable.Negate(fac))).Sat() {
				return false
			}
		}
	}
	return true
}

func sliceContains(slice []f.Variable, v f.Variable) bool {
	return slices.Contains(slice, v)
}
