package sat

import (
	"testing"
	"time"

	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/handler"
	"github.com/stretchr/testify/assert"
)

func TestTimeoutHandlerWithDuration(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	ph := GeneratePigeonHole(fac, 10)
	solver := NewSolver(fac)
	solver.Add(ph)
	duration, _ := time.ParseDuration("500ms")
	satHandler := HandlerWithTimeout(*handler.NewTimeoutWithDuration(duration))

	sat, ok := solver.SatWithHandler(satHandler)

	assert.False(ok)
	assert.True(satHandler.Aborted())
	assert.False(sat)

	ph = GeneratePigeonHole(fac, 2)
	solver = NewSolver(fac)
	solver.Add(ph)
	duration, _ = time.ParseDuration("2s")
	satHandler = HandlerWithTimeout(*handler.NewTimeoutWithDuration(duration))

	sat, ok = solver.SatWithHandler(satHandler)

	assert.True(ok)
	assert.False(satHandler.Aborted())
	assert.False(sat)
}

func TestTimeoutHandlerWithEnd(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	ph := GeneratePigeonHole(fac, 10)
	solver := NewSolver(fac)
	solver.Add(ph)
	duration, _ := time.ParseDuration("500ms")
	end := time.Now().Add(duration)
	handler := HandlerWithTimeout(*handler.NewTimeoutWithEnd(end))

	sat, ok := solver.SatWithHandler(handler)

	assert.False(ok)
	assert.True(handler.Aborted())
	assert.False(sat)
}

func TestOptimizationTimeoutHandler(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	nq := GenerateNQueens(fac, 25)
	vars := f.VariablesAsLiterals(f.Variables(fac, nq).Content())
	solver := NewSolver(fac)
	solver.Add(nq)
	duration, _ := time.ParseDuration("500ms")
	end := time.Now().Add(duration)
	optHandler := OptimizationHandlerWithTimeout(*handler.NewTimeoutWithEnd(end))

	result, ok := solver.MaximizeWithHandler(vars, optHandler)

	assert.False(ok)
	assert.True(optHandler.Aborted())
	assert.Nil(result)
	assert.NotNil(optHandler.IntermediateResult())

	nq = GenerateNQueens(fac, 4)
	vars = f.VariablesAsLiterals(f.Variables(fac, nq).Content())
	solver = NewSolver(fac)
	solver.Add(nq)
	duration, _ = time.ParseDuration("2h")
	optHandler = OptimizationHandlerWithTimeout(*handler.NewTimeoutWithDuration(duration))

	result, ok = solver.MaximizeWithHandler(vars, optHandler)

	assert.True(ok)
	assert.False(optHandler.Aborted())
	assert.NotNil(result)
}
