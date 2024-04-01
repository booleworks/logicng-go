package test

import (
	"testing"
	"time"

	"booleworks.com/logicng/bdd"
	"booleworks.com/logicng/encoding"
	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/io"
	"booleworks.com/logicng/normalform"
	"github.com/stretchr/testify/assert"
)

func TestBddSmallFormulas(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	fac := f.NewFactory()
	start := time.Now()
	formulas, _ := io.ReadFormulas(fac, "data/formulas/small_formulas.txt")
	elapsed := time.Since(start) / 1_000_000
	t.Logf("Read %d formula : %d ms", len(formulas), elapsed)

	minNoOrder := 0
	minBfs := 0
	minDfs := 0
	minMax2Min := 0
	minForce := 0

	fastestNoOrder := 0
	fastestBfs := 0
	fastestDfs := 0
	fastestMax2Min := 0
	fastestForce := 0

	for _, formula := range formulas {
		start = time.Now()
		bddNoOrder := bdd.Build(fac, formula)
		timeNoOrder := time.Since(start).Nanoseconds()
		start = time.Now()
		bddBfs := bdd.BuildWithVarOrder(fac, formula, bdd.BFSOrder(fac, formula))
		timeBfs := time.Since(start).Nanoseconds()
		start = time.Now()
		bddDfs := bdd.BuildWithVarOrder(fac, formula, bdd.DFSOrder(fac, formula))
		timeDfs := time.Since(start).Nanoseconds()
		start = time.Now()
		bddMax2Min := bdd.BuildWithVarOrder(fac, formula, bdd.MaxToMinOrder(fac, formula))
		timeMax2Min := time.Since(start).Nanoseconds()
		start = time.Now()
		bddForce := bdd.BuildWithVarOrder(fac, formula, bdd.ForceOrder(fac, formula))
		timeForce := time.Since(start).Nanoseconds()

		nodesNoOrder := bddNoOrder.NodeCount()
		nodesBfs := bddBfs.NodeCount()
		nodesDfs := bddDfs.NodeCount()
		nodesMax2Min := bddMax2Min.NodeCount()
		nodesForce := bddForce.NodeCount()
		minNodes := min(nodesNoOrder, nodesBfs, nodesDfs, nodesMax2Min, nodesForce)
		minTime := min(timeNoOrder, timeBfs, timeDfs, timeMax2Min, timeForce)
		if nodesNoOrder == minNodes {
			minNoOrder++
		}
		if nodesBfs == minNodes {
			minBfs++
		}
		if nodesDfs == minNodes {
			minDfs++
		}
		if nodesMax2Min == minNodes {
			minMax2Min++
		}
		if nodesForce == minNodes {
			minForce++
		}
		if timeNoOrder == minTime {
			fastestNoOrder++
		}
		if timeBfs == minTime {
			fastestBfs++
		}
		if timeDfs == minTime {
			fastestDfs++
		}
		if timeMax2Min == minTime {
			fastestMax2Min++
		}
		if timeForce == minTime {
			fastestForce++
		}
	}

	t.Logf("Winner Size\n")
	t.Logf("===========\n")
	t.Logf("No Order: %d\n", minNoOrder)
	t.Logf("DFS:      %d\n", minDfs)
	t.Logf("BFS:      %d\n", minBfs)
	t.Logf("Max2Min:  %d\n", minMax2Min)
	t.Logf("Force:    %d\n", minForce)

	t.Logf("Winner Time\n")
	t.Logf("===========\n")
	t.Logf("No Order: %d\n", fastestNoOrder)
	t.Logf("DFS:      %d\n", fastestDfs)
	t.Logf("BFS:      %d\n", fastestBfs)
	t.Logf("Max2Min:  %d\n", fastestMax2Min)
	t.Logf("Force:    %d\n", fastestForce)
}

func TestBddMidFormula(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	fac := f.NewFactory()
	cnfConfig := normalform.DefaultCNFConfig()
	cnfConfig.Algorithm = normalform.CNFFactorization
	fac.PutConfiguration(cnfConfig)
	encodingConfig := encoding.DefaultConfig()
	encodingConfig.AMOEncoder = encoding.AMOPure
	fac.PutConfiguration(encodingConfig)
	formula, _ := io.ReadFormula(fac, "data/formulas/bdd_small.txt")
	cnf := normalform.FactorizedCNF(fac, formula)
	order := bdd.DFSOrder(fac, cnf)
	kernel := bdd.NewKernelWithOrdering(fac, order, 100000, 200000)
	start := time.Now()
	bdd := bdd.BuildWithKernel(fac, cnf, kernel)
	elapsed := time.Since(start) / 1_000_000
	t.Logf("#Nodes: %d\n", bdd.NodeCount())
	t.Logf("Time:   %d ms\n", elapsed)
	assert.Equal(t, 2267795, bdd.NodeCount())
}
