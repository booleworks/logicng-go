package maxsat

import (
	"booleworks.com/logicng/parser"
	"testing"

	f "booleworks.com/logicng/formula"
	"github.com/stretchr/testify/assert"
)

var pureMaxsatFiles = []string{
	"c5315-bug-gate-0.dimacs.seq.filtered.cnf",
	"c6288-bug-gate-0.dimacs.seq.filtered.cnf",
	"c7552-bug-gate-0.dimacs.seq.filtered.cnf",
	"mot_comb1._red-gate-0.dimacs.seq.filtered.cnf",
	"mot_comb2._red-gate-0.dimacs.seq.filtered.cnf",
	"mot_comb3._red-gate-0.dimacs.seq.filtered.cnf",
	"s15850-bug-onevec-gate-0.dimacs.seq.filtered.cnf",
}

var partialMaxsatFiles = []string{
	"c1355_F176gat-1278gat@1.wcnf",
	"c1355_F1001gat-1048gat@1.wcnf",
	"c1355_F1183gat-1262gat@1.wcnf",
	"c1355_F1229gat@1.wcnf",
	"normalized-s3-3-3-1pb.wcnf",
	"normalized-s3-3-3-2pb.wcnf",
	"normalized-s3-3-3-3pb.wcnf",
	"term1_gr_2pin_w4.shuffled.cnf",
}

var partialWeightedMaxsatFiles = []string{
	"8.wcsp.log.wcnf",
	"54.wcsp.log.wcnf",
	"404.wcsp.log.wcnf",
	"term1_gr_2pin_w4.shuffled.cnf",
}

var partialWeightedMaxsatBmoFiles = []string{
	"normalized-factor-size=9-P=11-Q=283.opb.wcnf",
	"normalized-factor-size=9-P=11-Q=53.opb.wcnf",
	"normalized-factor-size=9-P=13-Q=179.opb.wcnf",
	"normalized-factor-size=9-P=17-Q=347.opb.wcnf",
	"normalized-factor-size=9-P=17-Q=487.opb.wcnf",
	"normalized-factor-size=9-P=23-Q=293.opb.wcnf",
}

var (
	partialMaxsatResults            = []int{13, 21, 33, 33, 36, 36, 36, 0}
	partialWeightedMaxsatResults    = []int{2, 37, 114, 0}
	partialWeightedMaxsatBmoResults = []int{11, 11, 13, 17, 17, 23}
)

///////////////
// Linear SU //
///////////////

func TestDoc(t *testing.T) {
	fac := f.NewFactory()
	p := parser.New(fac)
	solver := OLL(fac)
	solver.AddHardFormula(p.ParseUnsafe("A & B & (C | D)"))
	solver.AddSoftFormula(p.ParseUnsafe("A => ~B"), 2)
	solver.AddSoftFormula(p.ParseUnsafe("~C"), 4)
	solver.AddSoftFormula(p.ParseUnsafe("~D"), 8)

	result := solver.Solve()
	assert.True(t, result.Satisfiable)
	assert.Equal(t, 6, result.Optimum)
}

func TestPureMaxsatLinearSU(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	configs := make([]*Config, 1)
	configs[0] = DefaultConfig()
	for _, file := range pureMaxsatFiles {
		t.Logf("Testing Pure MaxSAT %s", file)
		for _, config := range configs {
			solver := LinearSU(fac, config)
			ReadDimacsToSolver(fac, solver, "../test/data/maxsat/"+file)
			result := solver.Solve()
			assert.True(result.Satisfiable)
			assert.Equal(1, result.Optimum)
		}
	}
}

func TestPartialMaxsatLinearSU(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	configs := make([]*Config, 2)
	configs[0] = DefaultConfig()
	configs[0].BMO = false
	configs[1] = DefaultConfig()
	configs[1].BMO = true
	for i, file := range partialMaxsatFiles {
		t.Logf("Testing Partial MaxSAT %s", file)
		for _, config := range configs {
			solver := LinearSU(fac, config)
			ReadDimacsToSolver(fac, solver, "../test/data/partialmaxsat/"+file)
			result := solver.Solve()
			assert.True(result.Satisfiable)
			assert.Equal(partialMaxsatResults[i], result.Optimum)
		}
	}
}

func TestPartialWeightedMaxsatLinearSU(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	configs := make([]*Config, 1)
	configs[0] = DefaultConfig()
	configs[0].BMO = false
	for i, file := range partialWeightedMaxsatFiles {
		t.Logf("Testing Partial Weighted MaxSAT %s", file)
		for _, config := range configs {
			solver := LinearSU(fac, config)
			ReadDimacsToSolver(fac, solver, "../test/data/partialweightedmaxsat/"+file)
			result := solver.Solve()
			assert.True(result.Satisfiable)
			assert.Equal(partialWeightedMaxsatResults[i], result.Optimum)
		}
	}
	configs[0].BMO = true
	for i, file := range partialWeightedMaxsatBmoFiles {
		t.Logf("Testing Partial Weighted MaxSAT BMO %s", file)
		for _, config := range configs {
			solver := LinearSU(fac, config)
			ReadDimacsToSolver(fac, solver, "../test/data/partialweightedmaxsat/bmo/"+file)
			result := solver.Solve()
			assert.True(result.Satisfiable)
			assert.Equal(partialWeightedMaxsatBmoResults[i], result.Optimum)
		}
	}
}

///////////////
// Linear US //
///////////////

func TestPureMaxsatLinearUS(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	configs := make([]*Config, 2)
	configs[0] = DefaultConfig()
	configs[0].IncrementalStrategy = IncNone
	configs[1] = DefaultConfig()
	configs[1].IncrementalStrategy = IncIterative
	for _, file := range pureMaxsatFiles {
		t.Logf("Testing Pure MaxSAT %s", file)
		for _, config := range configs {
			solver := LinearUS(fac, config)
			ReadDimacsToSolver(fac, solver, "../test/data/maxsat/"+file)
			result := solver.Solve()
			assert.True(result.Satisfiable)
			assert.Equal(1, result.Optimum)
		}
	}
}

func TestPartialMaxsatLinearUS(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	configs := make([]*Config, 2)
	configs[0] = DefaultConfig()
	configs[0].IncrementalStrategy = IncNone
	configs[1] = DefaultConfig()
	configs[1].IncrementalStrategy = IncIterative
	for i, file := range partialMaxsatFiles {
		t.Logf("Testing Partial MaxSAT %s", file)
		for _, config := range configs {
			solver := LinearUS(fac, config)
			ReadDimacsToSolver(fac, solver, "../test/data/partialmaxsat/"+file)
			result := solver.Solve()
			assert.True(result.Satisfiable)
			assert.Equal(partialMaxsatResults[i], result.Optimum)
		}
	}
}

//////////
// MSU3 //
//////////

func TestPureMaxsatMsu3(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	configs := make([]*Config, 2)
	configs[0] = DefaultConfig()
	configs[0].IncrementalStrategy = IncNone
	configs[1] = DefaultConfig()
	configs[1].IncrementalStrategy = IncIterative
	for _, file := range pureMaxsatFiles {
		t.Logf("Testing Pure MaxSAT %s", file)
		for _, config := range configs {
			solver := MSU3(fac, config)
			ReadDimacsToSolver(fac, solver, "../test/data/maxsat/"+file)
			result := solver.Solve()
			assert.True(result.Satisfiable)
			assert.Equal(1, result.Optimum)
		}
	}
}

func TestPartialMaxsatMsu3(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	configs := make([]*Config, 2)
	configs[0] = DefaultConfig()
	configs[0].IncrementalStrategy = IncNone
	configs[1] = DefaultConfig()
	configs[1].IncrementalStrategy = IncIterative
	for i, file := range partialMaxsatFiles {
		t.Logf("Testing Partial MaxSAT %s", file)
		for _, config := range configs {
			solver := MSU3(fac, config)
			ReadDimacsToSolver(fac, solver, "../test/data/partialmaxsat/"+file)
			result := solver.Solve()
			assert.True(result.Satisfiable)
			assert.Equal(partialMaxsatResults[i], result.Optimum)
		}
	}
}

///////////
// WMSU3 //
///////////

func TestPartialWeightedMaxsatWmsu3(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	configs := make([]*Config, 2)
	configs[0] = DefaultConfig()
	configs[0].IncrementalStrategy = IncNone
	configs[0].BMO = false
	configs[1] = DefaultConfig()
	configs[1].IncrementalStrategy = IncIterative
	configs[1].BMO = false
	for i, file := range partialWeightedMaxsatFiles {
		t.Logf("Testing Partial Weighted MaxSAT %s", file)
		for _, config := range configs {
			solver := WMSU3(fac, config)
			ReadDimacsToSolver(fac, solver, "../test/data/partialweightedmaxsat/"+file)
			result := solver.Solve()
			assert.True(result.Satisfiable)
			assert.Equal(partialWeightedMaxsatResults[i], result.Optimum)
		}
	}

	configs = make([]*Config, 1)
	configs[0] = DefaultConfig()
	configs[0].IncrementalStrategy = IncIterative
	configs[0].BMO = true

	for i, file := range partialWeightedMaxsatBmoFiles {
		t.Logf("Testing Partial Weighted MaxSAT BMO %s", file)
		for _, config := range configs {
			solver := WMSU3(fac, config)
			ReadDimacsToSolver(fac, solver, "../test/data/partialweightedmaxsat/bmo/"+file)
			result := solver.Solve()
			assert.True(result.Satisfiable)
			assert.Equal(partialWeightedMaxsatBmoResults[i], result.Optimum)
		}
	}
}

/////////
// WBO //
/////////

func TestPureMaxsatWbo(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	configs := make([]*Config, 2)
	configs[0] = DefaultConfig()
	configs[0].WeightStrategy = WeightNone
	configs[0].Symmetry = false
	configs[1] = DefaultConfig()
	configs[1].WeightStrategy = WeightNone
	configs[1].Symmetry = true
	for _, file := range pureMaxsatFiles {
		t.Logf("Testing Pure MaxSAT %s", file)
		for _, config := range configs {
			solver := WBO(fac, config)
			ReadDimacsToSolver(fac, solver, "../test/data/maxsat/"+file)
			result := solver.Solve()
			assert.True(result.Satisfiable)
			assert.Equal(1, result.Optimum)
		}
	}
}

func TestPartialMaxsatWbo(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	configs := make([]*Config, 1)
	configs[0] = DefaultConfig()
	for i, file := range partialMaxsatFiles {
		t.Logf("Testing Partial MaxSAT %s", file)
		for _, config := range configs {
			solver := WBO(fac, config)
			ReadDimacsToSolver(fac, solver, "../test/data/partialmaxsat/"+file)
			result := solver.Solve()
			assert.True(result.Satisfiable)
			assert.Equal(partialMaxsatResults[i], result.Optimum)
		}
	}
}

func TestPartialWeightedMaxsatWbo(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	configs := make([]*Config, 3)
	configs[0] = DefaultConfig()
	configs[0].WeightStrategy = WeightNone
	configs[1] = DefaultConfig()
	configs[1].WeightStrategy = WeightNormal
	configs[2] = DefaultConfig()
	configs[2].WeightStrategy = WeightDiversify
	for i, file := range partialWeightedMaxsatFiles {
		t.Logf("Testing Partial Weighted MaxSAT %s", file)
		for _, config := range configs {
			solver := WBO(fac, config)
			ReadDimacsToSolver(fac, solver, "../test/data/partialweightedmaxsat/"+file)
			result := solver.Solve()
			assert.True(result.Satisfiable)
			assert.Equal(partialWeightedMaxsatResults[i], result.Optimum)
		}
	}
}

////////////
// IncWBO //
////////////

func TestPureMaxsatIncWbo(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	configs := make([]*Config, 2)
	configs[0] = DefaultConfig()
	configs[0].WeightStrategy = WeightNone
	configs[0].Symmetry = false
	configs[1] = DefaultConfig()
	configs[1].WeightStrategy = WeightNone
	configs[1].Symmetry = true
	for _, file := range pureMaxsatFiles {
		t.Logf("Testing Pure MaxSAT %s", file)
		for _, config := range configs {
			solver := IncWBO(fac, config)
			ReadDimacsToSolver(fac, solver, "../test/data/maxsat/"+file)
			result := solver.Solve()
			assert.True(result.Satisfiable)
			assert.Equal(1, result.Optimum)
		}
	}
}

func TestPartialMaxsatIncWbo(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	configs := make([]*Config, 1)
	configs[0] = DefaultConfig()
	for i, file := range partialMaxsatFiles {
		for _, config := range configs {
			solver := IncWBO(fac, config)
			ReadDimacsToSolver(fac, solver, "../test/data/partialmaxsat/"+file)
			result := solver.Solve()
			assert.True(result.Satisfiable)
			assert.Equal(partialMaxsatResults[i], result.Optimum)
		}
	}
}

func TestPartialWeightedMaxsatIncWbo(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	configs := make([]*Config, 3)
	configs[0] = DefaultConfig()
	configs[0].WeightStrategy = WeightNone
	configs[1] = DefaultConfig()
	configs[1].WeightStrategy = WeightNormal
	configs[2] = DefaultConfig()
	configs[2].WeightStrategy = WeightDiversify
	for i, file := range partialWeightedMaxsatFiles {
		t.Logf("Testing Partial Weighted MaxSAT %s", file)
		for _, config := range configs {
			solver := IncWBO(fac, config)
			ReadDimacsToSolver(fac, solver, "../test/data/partialweightedmaxsat/"+file)
			result := solver.Solve()
			assert.True(result.Satisfiable)
			assert.Equal(partialWeightedMaxsatResults[i], result.Optimum)
		}
	}
}

/////////
// OLL //
/////////

func TestPureMaxsatOll(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	for _, file := range pureMaxsatFiles {
		t.Logf("Testing Pure MaxSAT %s", file)
		solver := OLL(fac)
		ReadDimacsToSolver(fac, solver, "../test/data/maxsat/"+file)
		result := solver.Solve()
		assert.True(result.Satisfiable)
		assert.Equal(1, result.Optimum)
	}
}

func TestPartialMaxsatOll(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	for i, file := range partialMaxsatFiles {
		t.Logf("Testing Partial MaxSAT %s", file)
		solver := OLL(fac)
		ReadDimacsToSolver(fac, solver, "../test/data/partialmaxsat/"+file)
		result := solver.Solve()
		assert.True(result.Satisfiable)
		assert.Equal(partialMaxsatResults[i], result.Optimum)
	}
}

func TestPartialWeightedMaxsatOll(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	for i, file := range partialWeightedMaxsatFiles {
		t.Logf("Testing Partial Weighted MaxSAT %s", file)
		solver := OLL(fac)
		ReadDimacsToSolver(fac, solver, "../test/data/partialweightedmaxsat/"+file)
		result := solver.Solve()
		assert.True(result.Satisfiable)
		assert.Equal(partialWeightedMaxsatResults[i], result.Optimum)
	}
	for i, file := range partialWeightedMaxsatBmoFiles {
		t.Logf("Testing Partial Weighted MaxSAT %s", file)
		solver := OLL(fac)
		ReadDimacsToSolver(fac, solver, "../test/data/partialweightedmaxsat/bmo/"+file)
		result := solver.Solve()
		assert.True(result.Satisfiable)
		assert.Equal(partialWeightedMaxsatBmoResults[i], result.Optimum)
	}
}

func TestMaxsatOllLarge1(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	file := "large_industrial.wcnf"
	t.Logf("Testing Partial MaxSAT %s", file)
	solver := OLL(fac)
	ReadDimacsToSolver(fac, solver, "../test/data/partialweightedmaxsat/large/"+file)
	result := solver.Solve()
	assert.True(result.Satisfiable)
	assert.Equal(68974, result.Optimum)
}

func TestMaxsatOllLarge2(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	file := "t3g3-5555.spn.wcnf"
	t.Logf("Testing Partial MaxSAT %s", file)
	solver := OLL(fac)
	ReadDimacsToSolver(fac, solver, "../test/data/partialweightedmaxsat/large/"+file)
	result := solver.Solve()
	assert.True(result.Satisfiable)
	assert.Equal(1100610, result.Optimum)
}

func TestMaxsatOllLargeWeights(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	file := "large_weights.wcnf"
	t.Logf("Testing Partial MaxSAT %s", file)
	solver := OLL(fac)
	ReadDimacsToSolver(fac, solver, "../test/data/partialweightedmaxsat/large/"+file)
	result := solver.Solve()
	assert.True(result.Satisfiable)
	assert.Equal(90912, result.Optimum)
}
