package maxsat

import (
	"testing"
	"time"

	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/handler"
	"booleworks.com/logicng/sat"
	"github.com/stretchr/testify/assert"
)

func TestMaxsatTimeoutHandler(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	phSmall := sat.GeneratePigeonHole(fac, 3)
	phLarge := sat.GeneratePigeonHole(fac, 10)

	for _, solver := range maxsatSolver(fac) {
		var weight int
		if solver.SupportsWeighted() {
			weight = 2
		} else {
			weight = 1
		}

		for _, clause := range fac.Operands(phLarge) {
			solver.AddSoftFormula(clause, weight)
		}
		duration, _ := time.ParseDuration("100ms")
		maxsatHandler := HandlerWithTimeout(*handler.NewTimeoutWithDuration(duration))
		result, ok := solver.SolveWithHandler(maxsatHandler)

		assert.False(ok)
		assert.True(maxsatHandler.Aborted())
		assert.Equal(Result{}, result)

		solver.Reset()

		for _, clause := range fac.Operands(phSmall) {
			solver.AddSoftFormula(clause, weight)
		}
		duration, _ = time.ParseDuration("100s")
		maxsatHandler = HandlerWithTimeout(*handler.NewTimeoutWithDuration(duration))

		result, ok = solver.SolveWithHandler(maxsatHandler)

		assert.True(ok)
		assert.False(maxsatHandler.Aborted())
		assert.True(result.Satisfiable)
		if solver.SupportsWeighted() {
			assert.Equal(2, result.Optimum)
		} else {
			assert.Equal(1, result.Optimum)
		}
	}
}

func maxsatSolver(fac f.Factory) []*Solver {
	return []*Solver{
		IncWBO(fac),
		WBO(fac),
		LinearSU(fac),
		LinearUS(fac),
		MSU3(fac),
		WMSU3(fac),
		OLL(fac),
	}
}
