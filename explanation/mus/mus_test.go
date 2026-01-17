package mus

import (
	"testing"

	e "github.com/booleworks/logicng-go/explanation"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/io"
	s "github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

func TestInsertionBasedMUS(t *testing.T) {
	fac := f.NewFactory()
	for _, props := range testFiles(fac, false) {
		mus, err := ComputeInsertionBased(fac, props)
		assert.Nil(t, err)
		testMUS(t, fac, props, mus)
	}
}

func TestDeletionBasedMUS(t *testing.T) {
	fac := f.NewFactory()
	for _, props := range testFiles(fac, !testing.Short()) {
		mus, err := ComputeDeletionBased(fac, props)
		assert.Nil(t, err)
		testMUS(t, fac, props, mus)
	}
}

func testFiles(fac f.Factory, all bool) []*[]f.Proposition {
	pg2 := generatePGPropositions(fac, 2)
	pg3 := generatePGPropositions(fac, 3)
	pg4 := generatePGPropositions(fac, 4)
	pg5 := generatePGPropositions(fac, 5)
	pg6 := generatePGPropositions(fac, 6)
	pg7 := generatePGPropositions(fac, 7)
	file1 := generateDimacsPropositions(fac, "../../test/data/dimacs/unsat/3col40_5_10.shuffled.cnf")
	file2 := generateDimacsPropositions(fac, "../../test/data/dimacs/unsat/x1_16.shuffled.cnf")
	file3 := generateDimacsPropositions(fac, "../../test/data/dimacs/unsat/grid_10_20.shuffled.cnf")
	file4 := generateDimacsPropositions(fac, "../../test/data/dimacs/unsat/ca032.shuffled.cnf")
	if all {
		return []*[]f.Proposition{&pg2, &pg3, &pg4, &pg5, &pg6, &pg7, &file1, &file2, &file3, &file4}
	}
	return []*[]f.Proposition{&pg2, &pg3, &pg4, &pg5, &pg6, &file1, &file2}
}

func testMUS(t *testing.T, fac f.Factory, original *[]f.Proposition, mus *e.UnsatCore) {
	assert.True(t, mus.IsGuaranteedMUS)
	assert.True(t, len(mus.Propositions) <= len(*original))
	solver := s.NewSolver(fac)
	for _, p := range mus.Propositions {
		assert.True(t, containsProps(original, p))
		assert.True(t, solver.Sat())
		solver.AddProposition(p)
	}
	assert.False(t, solver.Sat())
}

func containsProps(props *[]f.Proposition, prop f.Proposition) bool {
	for _, p := range *props {
		if p == prop {
			return true
		}
	}
	return false
}

func generatePGPropositions(fac f.Factory, n int) []f.Proposition {
	pgf := s.GeneratePigeonHole(fac, n)
	ops, _ := fac.NaryOperands(pgf)
	result := make([]f.Proposition, len(ops))
	for i, op := range ops {
		result[i] = f.NewStandardProposition(op)
	}
	return result
}

func generateDimacsPropositions(fac f.Factory, filename string) []f.Proposition {
	ops, _ := io.ReadDimacs(fac, filename)
	result := make([]f.Proposition, len(*ops))
	for i, op := range *ops {
		result[i] = f.NewStandardProposition(op)
	}
	return result
}
