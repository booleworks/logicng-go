package sat

import (
	"fmt"
	"os"
	"testing"

	"booleworks.com/logicng/explanation"
	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/io"
	"booleworks.com/logicng/parser"
	"github.com/stretchr/testify/assert"
)

func TestUnsatCoreSimple(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	p1 := f.NewStandardProposition(p.ParseUnsafe("((a & b) => c) &  ((a & b) => d)"), "P1")
	p2 := f.NewStandardProposition(p.ParseUnsafe("(c & d) <=> ~e"), "P2")
	p3 := f.NewStandardProposition(p.ParseUnsafe("~e => f | g"), "P3")
	p4 := f.NewStandardProposition(p.ParseUnsafe("(f => ~a) & (g => ~b) & p & q"), "P4")
	p5 := f.NewStandardProposition(p.ParseUnsafe("a => b"), "P5")
	p6 := f.NewStandardProposition(p.ParseUnsafe("a"), "P6")
	p7 := f.NewStandardProposition(p.ParseUnsafe("g | h"), "P7")
	p8 := f.NewStandardProposition(p.ParseUnsafe("(x => ~y | z) & (z | w)"), "P8")

	config := DefaultConfig()
	config.ProofGeneration = true
	solver := NewSolver(fac, config)
	solver.AddProposition(p1, p2, p3, p4, p5, p6, p7, p8)
	assert.False(solver.Sat())

	unsatCore, err := solver.ComputeUnsatCore()
	assert.Nil(err)
	assert.False(unsatCore.IsGuaranteedMUS)
	assert.True(containsAll(&unsatCore.Propositions, &[]f.Proposition{p1, p2, p3, p4, p5, p6}))
}

func TestUnsatCoresFromDimacs(t *testing.T) {
	fac := f.NewFactory()
	folder := "../test/data/dimacs/unsat/"
	items, _ := os.ReadDir(folder)
	for _, item := range items {
		if !item.IsDir() {
			cnf, _ := io.ReadDimacs(fac, folder+item.Name())
			t.Log("Testing Core for formula ", item.Name())
			config := DefaultConfig()
			config.ProofGeneration = true
			solver := NewSolver(fac, config)
			props := make([]f.Proposition, len(*cnf))
			for i, clause := range *cnf {
				prop := f.NewStandardProposition(clause, fmt.Sprintf("P%d", i))
				solver.AddProposition(prop)
				props[i] = prop
			}
			sat := solver.Sat()
			assert.False(t, sat)
			unsatCore, err := solver.ComputeUnsatCore()
			assert.Nil(t, err)
			t.Logf("Core with size %d of %d", len(unsatCore.Propositions), len(*cnf))
			verifyCore(fac, t, unsatCore, &props)
		}
	}
}

func TestUnsatCoreIncDec(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	p1 := f.NewStandardProposition(p.ParseUnsafe("((a & b) => c) &  ((a & b) => d)"), "P1")
	p2 := f.NewStandardProposition(p.ParseUnsafe("(c & d) <=> ~e"), "P2")
	p3 := f.NewStandardProposition(p.ParseUnsafe("~e => f | g"), "P3")
	p4 := f.NewStandardProposition(p.ParseUnsafe("(f => ~a) & (g => ~b) & p & q"), "P4")
	p5 := f.NewStandardProposition(p.ParseUnsafe("a => b"), "P5")
	p6 := f.NewStandardProposition(p.ParseUnsafe("a"), "P6")
	p7 := f.NewStandardProposition(p.ParseUnsafe("g | h"), "P7")
	p8 := f.NewStandardProposition(p.ParseUnsafe("(x => ~y | z) & (z | w)"), "P8")
	p9 := f.NewStandardProposition(p.ParseUnsafe("a & b"), "P9")
	p10 := f.NewStandardProposition(p.ParseUnsafe("(p => q) & p"), "P10")
	p11 := f.NewStandardProposition(p.ParseUnsafe("a & ~q"), "P11")

	config := DefaultConfig()
	config.ProofGeneration = true
	solver := NewSolver(fac, config)

	solver.AddProposition(p1, p2, p3, p4)
	state1 := solver.SaveState()
	solver.AddProposition(p5, p6)
	state2 := solver.SaveState()
	solver.AddProposition(p7, p8)

	assert.False(solver.Sat())
	unsatCore, err := solver.ComputeUnsatCore()
	assert.Nil(err)
	assert.True(containsAll(&unsatCore.Propositions, &[]f.Proposition{p1, p2, p3, p4, p5, p6}))

	solver.LoadState(state2)
	assert.False(solver.Sat())
	unsatCore, err = solver.ComputeUnsatCore()
	assert.Nil(err)
	assert.True(containsAll(&unsatCore.Propositions, &[]f.Proposition{p1, p2, p3, p4, p5, p6}))

	solver.LoadState(state1)
	solver.AddProposition(p9)
	assert.False(solver.Sat())
	unsatCore, err = solver.ComputeUnsatCore()
	assert.Nil(err)
	assert.True(containsAll(&unsatCore.Propositions, &[]f.Proposition{p1, p2, p3, p4, p9}))

	solver.LoadState(state1)
	solver.AddProposition(p5)
	solver.AddProposition(p6)
	assert.False(solver.Sat())
	unsatCore, err = solver.ComputeUnsatCore()
	assert.Nil(err)
	assert.True(containsAll(&unsatCore.Propositions, &[]f.Proposition{p1, p2, p3, p4, p5, p6}))

	solver.LoadState(state1)
	solver.AddProposition(p10)
	solver.AddProposition(p11)
	assert.False(solver.Sat())
	unsatCore, err = solver.ComputeUnsatCore()
	assert.Nil(err)
	assert.True(containsAll(&unsatCore.Propositions, &[]f.Proposition{p4, p11}))
}

func verifyCore(fac f.Factory, t *testing.T, originalCore *explanation.UnsatCore, props *[]f.Proposition) {
	assert.True(t, containsAll(props, &originalCore.Propositions))
	solver := NewSolver(fac)
	solver.AddProposition(originalCore.Propositions...)
	assert.False(t, solver.Sat())
}

func containsAll[T comparable](super, sub *[]T) bool {
	m := make(map[T]present, len(*super))
	for _, s := range *super {
		m[s] = present{}
	}
	for _, s := range *sub {
		if _, ok := m[s]; !ok {
			return false
		}
	}
	return true
}
