package enum

import (
	"testing"
	"time"

	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/handler"
	"booleworks.com/logicng/model/iter"
	"booleworks.com/logicng/sat"
	"github.com/stretchr/testify/assert"
)

func TestModelIterationTimeoutHandler(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	nq := sat.GenerateNQueens(fac, 25)
	vars := f.Variables(fac, nq).Content()
	solver := sat.NewSolver(fac)
	solver.Add(nq)
	duration, _ := time.ParseDuration("100ms")
	end := time.Now().Add(duration)

	meConfig := iter.DefaultConfig()
	meConfig.Handler = iter.HandlerWithTimeout(*handler.NewTimeoutWithEnd(end))
	result, ok := OnSolverWithConfig(solver, vars, meConfig)
	assert.False(ok)
	assert.True(meConfig.Handler.Aborted())
	assert.True(len(result) > 1)

	nq = sat.GenerateNQueens(fac, 5)
	vars = f.Variables(fac, nq).Content()
	solver = sat.NewSolver(fac)
	solver.Add(nq)
	duration, _ = time.ParseDuration("1h")
	meConfig.Handler = iter.HandlerWithTimeout(*handler.NewTimeoutWithDuration(duration))

	result, ok = OnSolverWithConfig(solver, vars, meConfig)
	assert.True(ok)
	assert.False(meConfig.Handler.Aborted())
	assert.Equal(10, len(result))
}
