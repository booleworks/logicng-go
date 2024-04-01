package bdd

import (
	"fmt"
	"math/big"
	"slices"
	"sort"
	"testing"
	"time"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/function"
	"github.com/booleworks/logicng-go/parser"
	"github.com/booleworks/logicng-go/randomizer"
	"github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

var reorderMethods = []ReorderingMethod{
	ReorderWin2,
	ReorderWin2Ite,
	ReorderWin3,
	ReorderWin3Ite,
	ReorderSift,
	ReorderSiftIte,
	ReorderRandom,
}

func TestBDDSwapping(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	a := fac.Var("a")
	b := fac.Var("b")
	c := fac.Var("c")
	order := []f.Variable{a, b, c}

	kernel := NewKernelWithOrdering(fac, order, 100, 100)
	formula := p.ParseUnsafe("a | b | c")
	bdd := BuildWithKernel(fac, formula, kernel)
	assert.Equal([]f.Variable{a, b, c}, bdd.VariableOrder())
	kernel.SwapVariables(a, b)
	assert.Equal([]f.Variable{b, a, c}, bdd.VariableOrder())
	kernel.SwapVariables(a, b)
	assert.Equal([]f.Variable{a, b, c}, bdd.VariableOrder())
	kernel.SwapVariables(a, a)
	assert.Equal([]f.Variable{a, b, c}, bdd.VariableOrder())
	kernel.SwapVariables(a, c)
	assert.Equal([]f.Variable{c, b, a}, bdd.VariableOrder())
	kernel.SwapVariables(b, c)
	assert.Equal([]f.Variable{b, c, a}, bdd.VariableOrder())
	assert.True(sat.IsEquivalent(fac, formula, bdd.CNF()))
}

func TestBDDSwappingMultipleBdds(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	a := fac.Var("a")
	b := fac.Var("b")
	c := fac.Var("c")
	order := []f.Variable{a, b, c}

	kernel := NewKernelWithOrdering(fac, order, 100, 100)
	formula1 := p.ParseUnsafe("a | b | c")
	formula2 := p.ParseUnsafe("a & b")
	bdd1 := BuildWithKernel(fac, formula1, kernel)
	bdd2 := BuildWithKernel(fac, formula2, kernel)
	assert.Equal([]f.Variable{a, b, c}, bdd1.VariableOrder())
	assert.Equal([]f.Variable{a, b, c}, bdd2.VariableOrder())
	kernel.SwapVariables(a, b)
	assert.Equal([]f.Variable{b, a, c}, bdd1.VariableOrder())
	assert.Equal([]f.Variable{b, a, c}, bdd2.VariableOrder())
}

func TestBDDRandomReorderingQuick(t *testing.T) {
	stats := &swapStats{}
	testRandomReordering(t, 25, 30, false, stats)
}

func TestBDDRandomReorderingLongRunning(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	stats := &swapStats{}
	testRandomReordering(t, 25, 45, true, stats)
}

func TestBDDReorderOnBuildQuick(t *testing.T) {
	stats := &swapStats{}
	testReorderOnBuild(t, 25, 30, false, stats)
}

func TestBDDReorderOnBuildLongRunning(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	stats := &swapStats{}
	testReorderOnBuild(t, 25, 47, true, stats)
}

func testRandomReordering(t *testing.T, minVars, maxVars int, verbose bool, stats *swapStats) {
	for vars := minVars; vars <= maxVars; vars++ {
		for depth := 4; depth <= 6; depth++ {
			fac := f.NewFactory()
			formula := randomFormula(fac, vars, depth)
			if verbose {
				fmt.Printf("vars = %2d, depth = %2d, nodes = %5d\n", vars, depth, function.NumberOfNodes(fac, formula))
			}
			for _, method := range reorderMethods {
				performReorder(t, fac, formula, method, true, verbose, stats)
			}
			for _, method := range reorderMethods {
				performReorder(t, fac, formula, method, false, verbose, stats)
			}
		}
	}
}

func performReorder(
	t *testing.T,
	fac f.Factory,
	formula f.Formula,
	reorderMethod ReorderingMethod,
	withBlocks, verbose bool,
	stats *swapStats,
) {
	assert := assert.New(t)
	order := order(fac, formula)
	kernel := NewKernelWithOrdering(fac, order, 1000, 10000)
	bdd := BuildWithKernel(fac, formula, kernel)
	count := bdd.ModelCount()
	usedBefore := bdd.NodeCount()
	start := time.Now()
	addVariableBlocks(f.Variables(fac, formula).Size(), withBlocks, kernel)
	kernel.reordering.reorder(reorderMethod)
	duration := time.Since(start) / 1_000_000
	usedAfter := bdd.NodeCount()
	assert.True(verifyBddConsistency(fac, formula, bdd, count, stats))
	verifyVariableBlocks(t, fac, formula, withBlocks, bdd)
	if reorderMethod != ReorderRandom {
		assert.LessOrEqual(usedAfter, usedBefore)
	}
	reduction := (float64(usedBefore-usedAfter) / float64(usedBefore)) * 100
	if verbose {
		var blocks string
		if withBlocks {
			blocks = "with"
		} else {
			blocks = "without"
		}
		fmt.Printf("%-20s: Reduced %7s blocks in %5dms by %.2f%% from %d to %d\n",
			reorderMethod, blocks, duration, reduction, usedBefore, usedAfter)
	}
}

func testReorderOnBuild(t *testing.T, minVars, maxVars int, verbose bool, stats *swapStats) {
	for vars := minVars; vars <= maxVars; vars++ {
		for depth := 4; depth <= 6; depth++ {
			fac := f.NewFactory()
			formula := randomFormula(fac, vars, depth)
			if verbose {
				fmt.Printf("vars = %2d, depth = %2d, nodes = %5d\n", vars, depth, function.NumberOfNodes(fac, formula))
			}
			order := order(fac, formula)
			kernel := NewKernelWithOrdering(fac, order, 1000, 10000)
			bdd := BuildWithKernel(fac, formula, kernel)
			nodeCount := bdd.NodeCount()
			modelCount := bdd.ModelCount()
			for _, method := range reorderMethods {
				reorderOnBuild(t, fac, formula, method, modelCount, nodeCount, true, verbose, stats)
			}
			for _, method := range reorderMethods {
				reorderOnBuild(t, fac, formula, method, modelCount, nodeCount, false, verbose, stats)
			}
		}
	}
}

func reorderOnBuild(
	t *testing.T,
	fac f.Factory,
	formula f.Formula,
	method ReorderingMethod,
	originalCount *big.Int,
	originalUsedNodes int,
	withBlocks, verbose bool,
	stats *swapStats,
) {
	order := order(fac, formula)
	kernel := NewKernelWithOrdering(fac, order, 1000, 10000)
	addVariableBlocks(len(order), withBlocks, kernel)
	kernel.reordering.setReorderDuringConstruction(method, 10000)
	start := time.Now()
	bdd := BuildWithKernel(fac, formula, kernel)
	duration := time.Since(start) / 1_000_000
	usedAfter := bdd.NodeCount()
	verifyVariableBlocks(t, fac, formula, withBlocks, bdd)
	verifyBddConsistency(fac, formula, bdd, originalCount, stats)
	reduction := (float64(originalUsedNodes-usedAfter) / float64(originalUsedNodes)) * 100
	if verbose {
		fmt.Printf("%-20s: Built in %5d ms, reduction by %6.2f%% from %6d to %6d\n", method,
			duration, reduction, originalUsedNodes, usedAfter)
	}
}

func order(fac f.Factory, formula f.Formula) []f.Variable {
	order := f.Variables(fac, formula).Content()
	sort.Slice(order, func(i, j int) bool {
		n1, _ := fac.VarName(order[i])
		n2, _ := fac.VarName(order[j])
		return n1 < n2
	})
	return order
}

func randomFormula(fac f.Factory, vars, depth int) f.Formula {
	config := randomizer.DefaultConfig()
	config.NumVars = vars
	config.Seed = int64(vars * depth * 42)
	config.WeightEquiv = 0
	config.WeightImpl = 0
	config.WeightNot = 0
	randomizer := randomizer.New(fac, config)
	for {
		formula := randomizer.And(depth)
		if f.Variables(fac, formula).Size() == vars && sat.IsSatisfiable(fac, formula) {
			return formula
		}
	}
}

func addVariableBlocks(numVars int, withBlocks bool, kernel *Kernel) {
	reordering := kernel.reordering
	if withBlocks {
		reordering.addVariableBlockAll()
		reordering.addVariableBlock(0, 20, true)
		reordering.addVariableBlock(0, 10, false)
		reordering.addVariableBlock(11, 20, false)
		reordering.addVariableBlock(15, 19, false)
		reordering.addVariableBlock(15, 17, true)
		reordering.addVariableBlock(18, 19, false)
		reordering.addVariableBlock(21, int32(numVars-1), false)
		if numVars > 33 {
			reordering.addVariableBlock(30, 33, false)
		}
	} else {
		reordering.addVariableBlockAll()
	}
}

func verifyBddConsistency(fac f.Factory, f1 f.Formula, bdd *BDD, modelCount *big.Int, stats *swapStats) bool {
	if !verify(bdd.Kernel, bdd.Index) {
		return false
	}
	nodes := verifyTree(bdd.Kernel, bdd.Index)
	if nodes < 0 {
		return false
	}
	stats.newBddSize(nodes)
	if modelCount != nil && modelCount.Cmp(bdd.ModelCount()) != 0 {
		fmt.Println("Model count changed!")
		return false
	}
	if modelCount == nil && !sat.IsEquivalent(fac, f1, bdd.CNF()) {
		fmt.Println("Not equal")
		return false
	}
	return true
}

func verifyVariableBlocks(t *testing.T, fac f.Factory, formula f.Formula, withBlocks bool, bdd *BDD) {
	assert := assert.New(t)
	if withBlocks {
		assert.True(findSequence(fac, bdd, getVars(0, 21)))
		assert.True(findSequence(fac, bdd, getVars(0, 11)))
		assert.True(findSequence(fac, bdd, getVars(11, 21)))
		assert.True(findSequence(fac, bdd, getVars(15, 20)))
		assert.True(findSequence(fac, bdd, getVars(15, 18)))
		assert.True(findSequence(fac, bdd, getVars(18, 20)))
		assert.True(findSequence(fac, bdd, getVars(21, f.Variables(fac, formula).Size())))
		if f.Variables(fac, formula).Size() > 33 {
			assert.True(findSequence(fac, bdd, getVars(30, 34)))
		}
		order := bdd.VariableOrder()
		assert.Less(slices.Index(order, fac.Var("v00")), slices.Index(order, fac.Var("v11")))
		assert.Equal(slices.Index(order, fac.Var("v16")), slices.Index(order, fac.Var("v15"))+1)
		assert.Equal(slices.Index(order, fac.Var("v17")), slices.Index(order, fac.Var("v16"))+1)
	}
}

func getVars(first, last int) []string {
	vars := make([]string, last-first)
	for i := first; i < last; i++ {
		vars[i-first] = fmt.Sprintf("v%02d", i)
	}
	return vars
}

func findSequence(fac f.Factory, bdd *BDD, vars []string) bool {
	order := bdd.VariableOrder()
	for i, it := range order {
		name, _ := fac.VarName(it)
		if slices.Contains(vars, name) {
			numFound := 1
			for numFound < len(vars) && i < len(order) {
				i++
				name, _ := fac.VarName(order[i])
				if !slices.Contains(vars, name) {
					return false
				} else {
					numFound++
				}
			}
			return true
		}
	}
	return false
}

type swapStats struct {
	testedFormulas int
	numSwaps       int
	maxFormulaSize int
	maxBddNodes    int
	maxBddSize     int32
}

func (s *swapStats) newFormula(fac f.Factory, formula f.Formula) {
	s.maxFormulaSize = max(s.maxFormulaSize, function.NumberOfNodes(fac, formula))
}

func (s *swapStats) newBdd(bdd *BDD) {
	s.maxBddNodes = max(s.maxBddNodes, bdd.NodeCount())
}

func (s *swapStats) newBddSize(size int32) {
	s.maxBddSize = max(s.maxBddSize, size)
}

/////////////////// Verification /////////////////////////////////

func verify(k *Kernel, root int32) bool {
	varnum := int32(len(k.level2var) - 1)
	for i := int32(0); i < varnum*2+2; i++ {
		if k.refcou(i) != maxref {
			fmt.Printf("Constant or Variable without MAXREF count: %d\n", i)
			return false
		}
		if i == 0 && (k.low(i) != 0 || k.high(i) != 0 || k.level(i) != varnum) {
			fmt.Println("Illegal FALSE node")
			return false
		}
		if i == 1 && (k.low(i) != 1 || k.high(i) != 1 || k.level(i) != varnum) {
			fmt.Println("Illegal TRUE node")
			return false
		}
		if i > 1 && i%2 == 0 {
			if k.low(i) != 0 {
				fmt.Println("VAR Low wrong")
				return false
			} else if k.high(i) != 1 {
				fmt.Println("VAR High wrong")
				return false
			}
		}
		if i > 1 && i%2 == 1 {
			if k.low(i) != 1 {
				fmt.Println("VAR Low wrong")
				return false
			} else if k.high(i) != 0 {
				fmt.Println("VAR High wrong")
				return false
			}
		}
		if i > 1 && k.level(i) >= varnum {
			fmt.Println("VAR Level wrong")
			return false
		}
	}
	if root >= 0 {
		for i := varnum*2 + 2; i < k.nodesize; i++ {
			if k.refcou(i) > 1 {
				fmt.Println("Refcou > 1")
				return false
			} else if k.refcou(i) == 1 && i != root {
				fmt.Println("Wrong refcou")
				return false
			} else if k.refcou(i) == 0 && i == root {
				fmt.Println("Entry point not marked")
				return false
			}
		}
	}
	return true
}

func verifyTree(k *Kernel, root int32) int32 {
	return verifyTreeRec(k, root, make([]int32, len(k.nodes)))
}

func verifyTreeRec(k *Kernel, root int32, cache []int32) int32 {
	if cache[root] > 0 {
		return cache[root]
	}
	low := k.low(root)
	high := k.high(root)

	nodeLevel := k.level(root)
	lowLevel := k.level(low)
	highLevel := k.level(high)

	if root == 0 || root == 1 {
		cache[root] = 1
		return 1
	}
	if nodeLevel > lowLevel && nodeLevel > highLevel {
		fmt.Printf("%d inconsistent!\n", root)
		return -1
	}
	lowRec := verifyTreeRec(k, low, cache)
	highRec := verifyTreeRec(k, high, cache)
	var result int32
	if lowRec < 0 || highRec < 0 {
		result = -1
	} else {
		result = lowRec + highRec
	}
	if result >= 0 {
		cache[root] = result
	}
	return result
}
