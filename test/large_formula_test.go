package test

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/function"
	"github.com/booleworks/logicng-go/io"
	"github.com/booleworks/logicng-go/normalform"
	"github.com/stretchr/testify/assert"
)

func TestMidFormula(t *testing.T) {
	fac := f.NewFactory()
	start := time.Now()
	formula, _ := io.ReadFormula(fac, "data/formulas/mid.txt")
	elapsed := time.Since(start) / 1_000_000
	atoms := function.NumberOfAtoms(fac, formula)
	t.Logf("Read formula (%d atoms): %d ms", atoms, elapsed)
	assert.Equal(t, 151963, atoms)

	start = time.Now()
	nnf := normalform.NNF(fac, formula)
	elapsed = time.Since(start) / 1_000_000
	atoms = function.NumberOfAtoms(fac, nnf)
	t.Logf("NNF (%d atoms): %d ms", atoms, elapsed)
	assert.Equal(t, 152289, atoms)

	start = time.Now()
	cnf := normalform.PGCNFDefault(fac, nnf)
	elapsed = time.Since(start) / 1_000_000
	atoms = function.NumberOfAtoms(fac, cnf)
	t.Logf("CNF (PG) (%d atoms): %d ms", atoms, elapsed)
	assert.Equal(t, 98830, atoms)
}

func TestLargeFormula(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	start := time.Now()
	formula, _ := io.ReadFormula(fac, "data/formulas/large.txt")
	elapsed := time.Since(start) / 1_000_000
	atoms := function.NumberOfAtoms(fac, formula)
	t.Logf("Read formula (%d atoms): %d ms", atoms, elapsed)
	assert.Equal(1035698, atoms)
	PrintMemUsage()

	start = time.Now()
	vars := f.Variables(fac, formula)
	elapsed = time.Since(start) / 1_000_000
	t.Logf("%d Variables: %d ms", vars.Size(), elapsed)
	PrintMemUsage()

	start = time.Now()
	vars = f.Variables(fac, formula)
	elapsed = time.Since(start) / 1_000_000
	t.Logf("%d Variables: %d ms", vars.Size(), elapsed)
	PrintMemUsage()

	start = time.Now()
	vs := f.Variables(fac, formula).Content()
	elapsed = time.Since(start) / 1_000_000
	t.Logf("%d Variables: %d ms", len(vs), elapsed)
	PrintMemUsage()

	start = time.Now()
	nnf := normalform.NNF(fac, formula)
	elapsed = time.Since(start) / 1_000_000
	atoms = function.NumberOfAtoms(fac, nnf)
	t.Logf("NNF (%d atoms): %d ms", atoms, elapsed)
	assert.Equal(1037491, atoms)
	PrintMemUsage()

	start = time.Now()
	cnf := normalform.PGCNFDefault(fac, nnf)
	elapsed = time.Since(start) / 1_000_000
	atoms = function.NumberOfAtoms(fac, cnf)
	t.Logf("CNF (PG) (%d atoms): %d ms", atoms, elapsed)
	assert.Equal(419374, atoms)
	PrintMemUsage()
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
