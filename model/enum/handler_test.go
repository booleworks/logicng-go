package enum

import (
	"testing"
	"time"

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
	meConfig.Handler = iter.HandlerWithTimeout(*handler.NewTimeoutWithEnd(end))
	_, ok := OnSolverWithConfig(solver, vars, meConfig)
	assert.False(ok)
	assert.True(meConfig.Handler.Aborted())

	nq = sat.GenerateNQueens(fac, 5)
	vars = f.Variables(fac, nq).Content()
	solver = sat.NewSolver(fac)
	solver.Add(nq)
	duration, _ = time.ParseDuration("1h")
	meConfig.Handler = iter.HandlerWithTimeout(*handler.NewTimeoutWithDuration(duration))

	result, ok := OnSolverWithConfig(solver, vars, meConfig)
	assert.True(ok)
	assert.False(meConfig.Handler.Aborted())
	assert.Equal(10, len(result))
}
