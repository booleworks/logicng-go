package sat

import (
	"github.com/booleworks/logicng-go/assignment"
	"github.com/booleworks/logicng-go/parser"
	"testing"

	"github.com/booleworks/logicng-go/explanation"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/io"
	"github.com/stretchr/testify/assert"
)

var allConfigs = []*Config{

	DefaultConfig().ClauseMin(ClauseMinNone).InitPhase(true).UseAtMost(true).Proofs(true),
	DefaultConfig().ClauseMin(ClauseMinNone).InitPhase(true).UseAtMost(true).Proofs(false),
	DefaultConfig().ClauseMin(ClauseMinNone).InitPhase(true).UseAtMost(false).Proofs(true),
	DefaultConfig().ClauseMin(ClauseMinNone).InitPhase(true).UseAtMost(false).Proofs(false),
	DefaultConfig().ClauseMin(ClauseMinNone).InitPhase(false).UseAtMost(true).Proofs(true),
	DefaultConfig().ClauseMin(ClauseMinNone).InitPhase(false).UseAtMost(true).Proofs(false),
	DefaultConfig().ClauseMin(ClauseMinNone).InitPhase(false).UseAtMost(false).Proofs(true),
	DefaultConfig().ClauseMin(ClauseMinNone).InitPhase(false).UseAtMost(false).Proofs(false),

	DefaultConfig().ClauseMin(ClauseMinBasic).InitPhase(true).UseAtMost(true).Proofs(true),
	DefaultConfig().ClauseMin(ClauseMinBasic).InitPhase(true).UseAtMost(true).Proofs(false),
	DefaultConfig().ClauseMin(ClauseMinBasic).InitPhase(true).UseAtMost(false).Proofs(true),
	DefaultConfig().ClauseMin(ClauseMinBasic).InitPhase(true).UseAtMost(false).Proofs(false),
	DefaultConfig().ClauseMin(ClauseMinBasic).InitPhase(false).UseAtMost(true).Proofs(true),
	DefaultConfig().ClauseMin(ClauseMinBasic).InitPhase(false).UseAtMost(true).Proofs(false),
	DefaultConfig().ClauseMin(ClauseMinBasic).InitPhase(false).UseAtMost(false).Proofs(true),
	DefaultConfig().ClauseMin(ClauseMinBasic).InitPhase(false).UseAtMost(false).Proofs(false),

	DefaultConfig().ClauseMin(ClauseMinDeep).InitPhase(true).UseAtMost(true).Proofs(true),
	DefaultConfig().ClauseMin(ClauseMinDeep).InitPhase(true).UseAtMost(true).Proofs(false),
	DefaultConfig().ClauseMin(ClauseMinDeep).InitPhase(true).UseAtMost(false).Proofs(true),
	DefaultConfig().ClauseMin(ClauseMinDeep).InitPhase(true).UseAtMost(false).Proofs(false),
	DefaultConfig().ClauseMin(ClauseMinDeep).InitPhase(false).UseAtMost(true).Proofs(true),
	DefaultConfig().ClauseMin(ClauseMinDeep).InitPhase(false).UseAtMost(true).Proofs(false),
	DefaultConfig().ClauseMin(ClauseMinDeep).InitPhase(false).UseAtMost(false).Proofs(true),
	DefaultConfig().ClauseMin(ClauseMinDeep).InitPhase(false).UseAtMost(false).Proofs(false),
}

func TestSolverCall(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	formula, _ := io.ReadFormula(fac, "../test/data/formulas/small.txt")
	vars := f.Variables(fac, formula).Content()
	p := parser.New(fac)
	for _, config := range allConfigs {
		solver := NewSolver(fac, config)
		solver.Add(formula)
		state := solver.SaveState()

		result := solver.Call()
		assert.True(result.OK())
		assert.True(result.Sat())
		assert.False(result.Aborted())
		assert.Nil(result.Model())
		assert.Nil(result.UnsatCore())
		assert.Nil(result.UpZeroLits())

		result = solver.Call(WithModel(vars))
		assert.True(result.OK())
		assert.True(result.Sat())
		assert.False(result.Aborted())
		assert.NotNil(result.Model())
		assert.Nil(result.UnsatCore())
		assert.Nil(result.UpZeroLits())
		ass, _ := result.Model().Assignment(fac)
		assert.True(assignment.Evaluate(fac, formula, ass))

		solver.Add(p.ParseUnsafe("v328 & v443 & v447"))
		if config.ProofGeneration {
			result = solver.Call(WithCore())
		} else {
			result = solver.Call()
		}
		assert.True(result.OK())
		assert.False(result.Sat())
		assert.Nil(result.Model())
		if config.ProofGeneration {
			assert.NotNil(result.UnsatCore())
			verifyUnsatCore(t, fac, result.UnsatCore())
		}
		err := solver.LoadState(state)
		assert.Nil(err)

		if config.ProofGeneration {
			result = solver.Call(WithCore().Variable(fac.Vars("v328", "v443", "v447")...))
		} else {
			result = solver.Call(Params().Variable(fac.Vars("v328", "v443", "v447")...))
		}
		assert.True(result.OK())
		assert.False(result.Sat())
		assert.Nil(result.Model())
		if config.ProofGeneration {
			assert.NotNil(result.UnsatCore())
			verifyUnsatCore(t, fac, result.UnsatCore())
		}

		if config.ProofGeneration {
			result = solver.Call(WithCore().WithModel(vars).Variable(fac.Vars("v328")...))
		} else {
			result = solver.Call(Params().WithModel(vars).Variable(fac.Vars("v328")...))
		}
		assert.True(result.OK())
		assert.True(result.Sat())
		assert.NotNil(result.Model())
		ass, _ = result.Model().Assignment(fac)
		assert.True(assignment.Evaluate(fac, formula, ass))
		assert.Nil(result.UnsatCore())
	}
}

func verifyUnsatCore(t *testing.T, fac f.Factory, unsatCore *explanation.UnsatCore) {
	solver := NewSolver(fac)
	for _, prop := range unsatCore.Propositions {
		solver.AddProposition(prop)
	}
	assert.False(t, solver.Sat())
}
