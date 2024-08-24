package enum

import (
	"testing"
	"time"

	"github.com/booleworks/logicng-go/event"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/model/iter"
	"github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

func TestModelIterationTimeoutHandler(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	nq := sat.GenerateNQueens(fac, 25)
	vars := f.Variables(fac, nq).Content()
	solver := sat.NewSolver(fac)
	solver.Add(nq)
	duration, _ := time.ParseDuration("500ms")
	end := time.Now().Add(duration)

	meConfig := iter.DefaultConfig()
	meConfig.Handler = handler.NewTimeoutWithEnd(end)
	_, state := OnSolverWithConfig(solver, vars, meConfig)
	assert.False(state.Success)
	assert.NotEqual(event.Nothing, state.CancelCause)

	nq = sat.GenerateNQueens(fac, 5)
	vars = f.Variables(fac, nq).Content()
	solver = sat.NewSolver(fac)
	solver.Add(nq)
	duration, _ = time.ParseDuration("1h")
	meConfig.Handler = handler.NewTimeoutWithDuration(duration)

	result, state := OnSolverWithConfig(solver, vars, meConfig)
	assert.True(state.Success)
	assert.Equal(event.Nothing, state.CancelCause)
	assert.Equal(10, len(result))
}
