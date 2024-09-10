package sat

import (
	"testing"
	"time"

	"github.com/booleworks/logicng-go/event"
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
	satHandler := handler.NewTimeoutWithDuration(duration)

	sResult := solver.Call(Params().Handler(satHandler))

	assert.False(sResult.OK())
	assert.True(sResult.Canceled())
	assert.NotEqual(event.Nothing, sResult.state.CancelCause)
	assert.False(sResult.Sat())

	ph = GeneratePigeonHole(fac, 2)
	solver = NewSolver(fac)
	solver.Add(ph)
	duration, _ = time.ParseDuration("2s")
	satHandler = handler.NewTimeoutWithDuration(duration)

	sResult = solver.Call(Params().Handler(satHandler))

	assert.True(sResult.OK())
	assert.False(sResult.Canceled())
	assert.Equal(event.Nothing, sResult.state.CancelCause)
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
	handler := *handler.NewTimeoutWithEnd(end)

	sResult := solver.Call(Params().Handler(handler))

	assert.False(sResult.OK())
	assert.True(sResult.Canceled())
	assert.NotEqual(event.Nothing, sResult.state.CancelCause)
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
	optHandler := handler.NewTimeoutWithEnd(end)

	result, state := solver.MaximizeWithHandler(vars, optHandler)

	assert.False(state.Success)
	assert.NotEqual(event.Nothing, state.CancelCause)
	assert.NotNil(result)

	nq = GenerateNQueens(fac, 4)
	vars = f.VariablesAsLiterals(f.Variables(fac, nq).Content())
	solver = NewSolver(fac)
	solver.Add(nq)
	duration, _ = time.ParseDuration("2h")
	optHandler = handler.NewTimeoutWithDuration(duration)

	result, state = solver.MaximizeWithHandler(vars, optHandler)

	assert.True(state.Success)
	assert.Equal(event.Nothing, state.CancelCause)
	assert.NotNil(result)
}
