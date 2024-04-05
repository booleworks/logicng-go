package sat

import (
	"testing"
	"time"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
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

	sResult := solver.Call(Params().Handler(satHandler))

	assert.False(sResult.OK())
	assert.True(sResult.Aborted())
	assert.True(satHandler.Aborted())
	assert.False(sResult.Sat())

	ph = GeneratePigeonHole(fac, 2)
	solver = NewSolver(fac)
	solver.Add(ph)
	duration, _ = time.ParseDuration("2s")
	satHandler = HandlerWithTimeout(*handler.NewTimeoutWithDuration(duration))

	sResult = solver.Call(Params().Handler(satHandler))

	assert.True(sResult.OK())
	assert.False(sResult.Aborted())
	assert.False(satHandler.Aborted())
	assert.False(sResult.Sat())
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

	sResult := solver.Call(Params().Handler(handler))

	assert.False(sResult.OK())
	assert.True(handler.Aborted())
	assert.False(sResult.Sat())
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
