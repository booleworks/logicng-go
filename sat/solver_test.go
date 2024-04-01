package sat

import (
	"os"
	"testing"
	"time"

	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/io"
	"booleworks.com/logicng/model"
	"booleworks.com/logicng/parser"
	"github.com/stretchr/testify/assert"
)

func getSolvers(fac f.Factory) []*Solver {
	return []*Solver{
		NewSolver(fac),
		NewSolver(fac, DefaultConfig().ClauseMin(ClauseMinNone)),
		NewSolver(fac, DefaultConfig().ClauseMin(ClauseMinBasic)),
		NewSolver(fac, DefaultConfig().ClauseMin(ClauseMinDeep)),
		NewSolver(fac, DefaultConfig().CNF(CNFFactorization)),
		NewSolver(fac, DefaultConfig().CNF(CNFPG)),
		NewSolver(fac, DefaultConfig().CNF(CNFFullPG)),
	}
}

func TestSolverTrue(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	for _, s := range getSolvers(fac) {
		s.Add(p.ParseUnsafe("$true"))
		assert.True(s.Sat())
	}
}

func TestSolverFalse(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	for _, s := range getSolvers(fac) {
		s.Add(p.ParseUnsafe("$false"))
		assert.False(s.Sat())
	}
}

func TestSolverLiterals(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	a := fac.Lit("a", true)
	for _, s := range getSolvers(fac) {
		s.Add(p.ParseUnsafe("a"))
		assert.True(s.Sat())
		mdl, _ := s.Model(fac.Vars("a"))
		assert.Equal(model.New(a), mdl)

		s.Add(p.ParseUnsafe("~a"))
		assert.False(s.Sat())
		_, err := s.Model(fac.Vars("a"))
		assert.NotNil(err)
	}
}

func TestSolverAnd1(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	a := fac.Lit("a", true)
	b := fac.Lit("b", true)
	for _, s := range getSolvers(fac) {
		s.Add(p.ParseUnsafe("a & b"))
		assert.True(s.Sat())
		mdl, _ := s.Model(fac.Vars("a", "b"))
		assert.Equal(model.New(a, b), mdl)

		s.Add(p.ParseUnsafe("~a"))
		assert.False(s.Sat())
		_, err := s.Model(fac.Vars("a", "b"))
		assert.NotNil(err)
	}
}

func TestSolverAnd2(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	a := fac.Lit("a", true)
	nb := fac.Lit("b", false)
	c := fac.Lit("c", true)
	nd := fac.Lit("d", false)

	for _, s := range getSolvers(fac) {
		s.Add(p.ParseUnsafe("a & ~b & c & ~d"))
		assert.True(s.Sat())
		mdl, _ := s.Model(fac.Vars("a", "b", "c", "d"))
		assert.Equal(model.New(a, nb, c, nd), mdl)

		s.Add(p.ParseUnsafe("d"))
		assert.False(s.Sat())
		_, err := s.Model(fac.Vars("a", "b", "c", "d"))
		assert.NotNil(err)
	}
}

func TestSolverAnd3(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	a := fac.Lit("a", true)
	nb := fac.Lit("b", false)
	c := fac.Lit("c", true)
	nd := fac.Lit("d", false)

	for _, s := range getSolvers(fac) {
		s.Add(a.AsFormula())
		s.Add(nb.AsFormula())
		s.Add(c.AsFormula())
		s.Add(nd.AsFormula())
		assert.True(s.Sat())
		mdl, _ := s.Model(fac.Vars("a", "b", "c", "d"))
		assert.Equal(model.New(a, nb, c, nd), mdl)

		s.Add(p.ParseUnsafe("d"))
		assert.False(s.Sat())
		_, err := s.Model(fac.Vars("a", "b", "c", "d"))
		assert.NotNil(err)
	}
}

func TestSolverFormula1(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	nx := fac.Lit("x", false)
	y := fac.Lit("y", true)
	z := fac.Lit("z", true)

	for _, s := range getSolvers(fac) {
		s.Add(p.ParseUnsafe("(x => y) & (~x => y) & (y => z) & (z => ~x)"))
		assert.True(s.Sat())
		mdl, _ := s.Model(fac.Vars("x", "y", "z"))
		assert.Equal(model.New(nx, y, z), mdl)

		s.Add(p.ParseUnsafe("~y"))
		assert.False(s.Sat())
		_, err := s.Model(fac.Vars("x", "y", "z"))
		assert.NotNil(err)
	}
}

func TestSolverDimacsSat(t *testing.T) {
	fac := f.NewFactory()
	folder := "../test/data/dimacs/sat/"
	items, _ := os.ReadDir(folder)
	for _, item := range items {
		if !item.IsDir() {
			testSingleDimacs(t, fac, folder+item.Name(), true)
		}
	}
}

func TestSolverDimacsUnsat(t *testing.T) {
	fac := f.NewFactory()
	folder := "../test/data/dimacs/unsat/"
	items, _ := os.ReadDir(folder)
	for _, item := range items {
		if !item.IsDir() {
			testSingleDimacs(t, fac, folder+item.Name(), false)
		}
	}
}

func testSingleDimacs(t *testing.T, fac f.Factory, filename string, expected bool) {
	t.Logf("Testing DIMACS file %s", filename)
	assert := assert.New(t)
	cnf, err := io.ReadDimacs(fac, filename)
	if err != nil {
		t.Errorf("could not read file %s, error: %s", filename, err)
	}
	for _, s := range getSolvers(fac) {
		s.Add(*cnf...)
		sat := s.Sat()

		if expected {
			assert.True(sat, "Did not get SAT for file %s", filename)
		} else {
			assert.False(sat, "Did not get UNSAT for file %s", filename)
		}
	}
}

func TestPigeonHole(t *testing.T) {
	for n := 1; n <= 7; n++ {
		fac := f.NewFactory()
		for _, s := range getSolvers(fac) {
			s.Add(GeneratePigeonHole(fac, n))

			start := time.Now()
			sat := s.Sat()
			elapsed := time.Since(start) / 1_000_000
			t.Logf("Pigeon Hole of size %d took %d ms", n, elapsed)
			assert.False(t, sat)
		}
	}
}

func TestSolverIncDecSimple(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	solver := NewSolver(fac)

	solver.Add(fac.Variable("a"))
	state1 := solver.SaveState()
	assert.Equal(int32(0), state1.id)
	assert.Equal([]int{1, 1, 0, 1, 0, 0}, state1.state)
	assert.True(solver.Sat())
	solver.Add(GeneratePigeonHole(fac, 5))
	assert.False(solver.Sat())
	err := solver.LoadState(state1)
	assert.Nil(err)
	assert.True(solver.Sat())
	solver.Add(fac.Literal("a", false))
	assert.False(solver.Sat())
	err = solver.LoadState(state1)
	assert.Nil(err)
	assert.True(solver.Sat())
	solver.Add(GeneratePigeonHole(fac, 5))
	state2 := solver.SaveState()
	assert.Equal(int32(1), state2.id)
	assert.Equal([]int{1, 31, 81, 1, 0, 0}, state2.state)
	solver.Add(GeneratePigeonHole(fac, 4))
	assert.False(solver.Sat())
	err = solver.LoadState(state2)
	assert.Nil(err)
	assert.False(solver.Sat())
	err = solver.LoadState(state1)
	assert.Nil(err)
	assert.True(solver.Sat())
}

func TestSolverIncDecDeep(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	solver := NewSolver(fac)

	solver.Add(fac.Variable("a"))
	state1 := solver.SaveState()
	solver.Add(fac.Variable("b"))
	assert.True(solver.Sat())
	solver.Add(fac.Literal("a", false))
	assert.False(solver.Sat())
	err := solver.LoadState(state1)
	assert.Nil(err)
	solver.Add(fac.Literal("b", false))
	assert.True(solver.Sat())
	state3 := solver.SaveState()
	solver.Add(fac.Literal("a", false))
	assert.False(solver.Sat())
	err = solver.LoadState(state3)
	assert.Nil(err)
	solver.Add(fac.Variable("c"))
	state4 := solver.SaveState()
	solver.SaveState()
	err = solver.LoadState(state4)
	assert.Nil(err)
	assert.True(solver.Sat())
	err = solver.LoadState(state1)
	assert.Nil(err)
	assert.True(solver.Sat())
}
