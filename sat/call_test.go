package sat

import (
	"testing"

	"github.com/booleworks/logicng-go/assignment"
	"github.com/booleworks/logicng-go/event"
	"github.com/booleworks/logicng-go/model"
	"github.com/booleworks/logicng-go/parser"

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

func TestSolverCallSequence(t *testing.T) {
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
		assert.False(result.Canceled())
		assert.Nil(result.Model())
		assert.Nil(result.UnsatCore())
		assert.Nil(result.UpZeroLits())

		result = solver.Call(WithModel(vars))
		assert.True(result.OK())
		assert.True(result.Sat())
		assert.False(result.Canceled())
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

func TestSolverCall(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	for _, config := range allConfigs {
		solver := NewSolver(fac, config)
		solver.Add(p.ParseUnsafe("a | b"))
		solver.Add(p.ParseUnsafe("c & (~c | ~a)"))
		assert.True(solver.Call().satisfiable)
		assert.False(solver.Call(Params().Literal(fac.Lit("b", false))).satisfiable)
		assert.True(solver.Call().satisfiable)
	}
}

func TestSolverCallModel(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	for _, config := range allConfigs {
		solver := NewSolver(fac, config)
		solver.Add(p.ParseUnsafe("a | b"))
		solver.Add(p.ParseUnsafe("c & (~c | ~a)"))
		abc := fac.Vars("a", "b", "c")
		abcd := fac.Vars("a", "b", "c", "d")

		result := solver.Call(WithModel(abc))
		assert.True(result.Sat())
		assert.Equal(model.New(fac.Lit("a", false), fac.Lit("b", true), fac.Lit("c", true)), result.Model())
		result = solver.Call(WithModel(abcd).Formula(p.ParseUnsafe("c | d")))
		assert.True(result.Sat())
		assert.Equal(4, result.Model().Size())

		result = solver.Call(WithModel(abc))
		assert.True(result.Sat())
		assert.Equal(model.New(fac.Lit("a", false), fac.Lit("b", true), fac.Lit("c", true)), result.Model())
	}
}

func TestSolverCallCore(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	for _, config := range allConfigs {
		if config.ProofGeneration {
			solver := NewSolver(fac, config)
			solver.Add(p.ParseUnsafe("a | b"))
			solver.Add(p.ParseUnsafe("c & (~c | ~a)"))

			result := solver.Call(WithCore().Literal(fac.Lit("b", false)))
			assert.False(result.Sat())
			verifyUnsatCore(t, fac, result.UnsatCore())
			assert.Equal(4, len(result.UnsatCore().Propositions))

			result = solver.Call(WithCore().Formula(p.ParseUnsafe("~b | a")))
			assert.False(result.Sat())
			verifyUnsatCore(t, fac, result.UnsatCore())
			assert.Equal(4, len(result.UnsatCore().Propositions))

			assert.True(solver.Call().Sat())
			assert.True(solver.Sat())

			solver.Add(p.ParseUnsafe("~b"))
			result = solver.Call(WithCore().Formula(p.ParseUnsafe("~b | a")))
			assert.False(result.Sat())
			verifyUnsatCore(t, fac, result.UnsatCore())
			assert.Equal(4, len(result.UnsatCore().Propositions))
		}
	}
}

func TestSolverCallHandler(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	formula, _ := io.ReadFormula(fac, "../test/data/formulas/small_formulas.txt")
	for _, config := range allConfigs {
		solver := NewSolver(fac, config)
		solver.Add(formula)

		result := solver.Call(WithHandler(newMaxConflictHandler(0)))
		assert.False(result.OK())
		assert.True(result.Canceled())

		result = solver.Call(WithHandler(newMaxConflictHandler(0)))
		assert.False(result.OK())
		assert.True(result.Canceled())

		result = solver.Call(WithHandler(newMaxConflictHandler(100)))
		assert.True(result.OK())
		assert.True(result.Sat())
		assert.False(result.Canceled())
	}
}

func TestSolverCallAdditionalFormulas(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	for _, config := range allConfigs {
		if config.ProofGeneration {
			config.CNFMethod = CNFFactorization
			solver := NewSolver(fac, config)
			solver.Add(p.ParseUnsafe("a | b | c | d"))
			solver.AddProposition(f.NewStandardProposition(p.ParseUnsafe("a => b")))
			solver.AddProposition(f.NewStandardProposition(p.ParseUnsafe("c => d")))
			solver.AddProposition(f.NewStandardProposition(p.ParseUnsafe("e => ~a & ~b")))
			solver.AddProposition(f.NewStandardProposition(p.ParseUnsafe("~f => ~c & ~d")))

			result := solver.Call(WithCore().Formula(p.ParseUnsafe("e <=> ~f"), p.ParseUnsafe("~f")))
			assert.False(result.Sat())
			assert.Equal(5, len(result.UnsatCore().Propositions))
			verifyUnsatCore(t, fac, result.UnsatCore())
			assert.True(propsContain(result.UnsatCore().Propositions, p.ParseUnsafe("~f")))

			result = solver.Call(WithCore().Formula(p.ParseUnsafe("e <=> ~f"), p.ParseUnsafe("e")))
			assert.False(result.Sat())
			assert.Equal(5, len(result.UnsatCore().Propositions))
			verifyUnsatCore(t, fac, result.UnsatCore())
			assert.False(propsContain(result.UnsatCore().Propositions, p.ParseUnsafe("~f")))
			assert.True(propsContain(result.UnsatCore().Propositions, p.ParseUnsafe("e")))

			result = solver.Call(WithCore().Formula(p.ParseUnsafe("e <=> ~f"), p.ParseUnsafe("~f")))
			assert.False(result.Sat())
			assert.Equal(5, len(result.UnsatCore().Propositions))
			verifyUnsatCore(t, fac, result.UnsatCore())
			assert.True(propsContain(result.UnsatCore().Propositions, p.ParseUnsafe("~f")))

			result = solver.Call(WithCore().Proposition(
				f.NewStandardProposition(p.ParseUnsafe("e <=> ~f")),
				f.NewStandardProposition(p.ParseUnsafe("e"))))
			assert.False(result.Sat())
			assert.Equal(5, len(result.UnsatCore().Propositions))
			verifyUnsatCore(t, fac, result.UnsatCore())
			assert.False(propsContain(result.UnsatCore().Propositions, p.ParseUnsafe("~f")))
			assert.True(propsContain(result.UnsatCore().Propositions, p.ParseUnsafe("e")))

			result = solver.Call(WithCore().
				WithModel(fac.Vars("a", "b", "c", "d", "e", "f")).
				Formula(p.ParseUnsafe("e <=> ~f"), p.ParseUnsafe("~e")),
			)
			assert.True(result.Sat())
			assert.Nil(result.UnsatCore())
			assert.NotNil(result.Model())
			assert.Equal(6, result.Model().Size())

			result = solver.Call(WithCore().
				WithModel(fac.Vars("a", "b", "c")).
				Formula(p.ParseUnsafe("e <=> ~f"), p.ParseUnsafe("~e")),
			)
			assert.True(result.Sat())
			assert.Nil(result.UnsatCore())
			assert.NotNil(result.Model())
			assert.Equal(3, result.Model().Size())

			result = solver.Call(WithCore().Proposition(
				f.NewStandardProposition(p.ParseUnsafe("e <=> ~f")),
				f.NewStandardProposition(p.ParseUnsafe("e"))))
			assert.False(result.Sat())
			assert.Equal(5, len(result.UnsatCore().Propositions))
			verifyUnsatCore(t, fac, result.UnsatCore())
			assert.False(propsContain(result.UnsatCore().Propositions, p.ParseUnsafe("~f")))
			assert.True(propsContain(result.UnsatCore().Propositions, p.ParseUnsafe("e")))
		}
	}
}

type maxConflictHander struct {
	maxConflicts int
	numConflicts int
	canceled     bool
}

func newMaxConflictHandler(maxConflicts int) *maxConflictHander {
	return &maxConflictHander{maxConflicts, 0, false}
}

func (h *maxConflictHander) ShouldResume(e event.Event) bool {
	h.canceled = h.numConflicts > h.maxConflicts
	h.numConflicts++
	return !h.canceled
}

func verifyUnsatCore(t *testing.T, fac f.Factory, unsatCore *explanation.UnsatCore) {
	solver := NewSolver(fac)
	for _, prop := range unsatCore.Propositions {
		solver.AddProposition(prop)
	}
	assert.False(t, solver.Sat())
}

func propsContain(props []f.Proposition, formula f.Formula) bool {
	for _, prop := range props {
		if prop.Formula() == formula {
			return true
		}
	}
	return false
}
