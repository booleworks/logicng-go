package maxsat

import (
	"testing"
	"time"

	"github.com/booleworks/logicng-go/event"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

func TestMaxsatTimeoutHandler(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	// phSmall := sat.GeneratePigeonHole(fac, 3)
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
		maxsatHandler := handler.NewTimeoutWithDuration(duration)
		result, state := solver.SolveWithHandler(maxsatHandler)

		assert.False(state.Success)
		assert.NotEqual(event.Nothing, state.CancelCause)
		assert.Equal(Result{}, result)

		// TODO activate
		// solver.Reset()
		//
		// for _, clause := range fac.Operands(phSmall) {
		// 	solver.AddSoftFormula(clause, weight)
		// }
		// duration, _ = time.ParseDuration("100s")
		// maxsatHandler = handler.NewTimeoutWithDuration(duration)
		//
		// result, state = solver.SolveWithHandler(maxsatHandler)
		//
		// assert.True(state.Success)
		// assert.Equal(event.Nothing, state.CancelCause)
		// assert.True(result.Satisfiable)
		// if solver.SupportsWeighted() {
		// 	assert.Equal(2, result.Optimum)
		// } else {
		// 	assert.Equal(1, result.Optimum)
		// }
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
