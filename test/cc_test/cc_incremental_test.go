package cc

import (
	"fmt"
	"testing"

	e "booleworks.com/logicng/encoding"
	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/sat"
	"github.com/stretchr/testify/assert"
)

func cfgs() []*e.Config {
	cfgs := make([]*e.Config, 3)
	cfgs[0] = e.DefaultConfig()
	cfgs[0].AMKEncoder = e.AMKTotalizer
	cfgs[0].ALKEncoder = e.ALKTotalizer
	cfgs[1] = e.DefaultConfig()
	cfgs[1].AMKEncoder = e.AMKCardinalityNetwork
	cfgs[1].ALKEncoder = e.ALKCardinalityNetwork
	cfgs[2] = e.DefaultConfig()
	cfgs[2].AMKEncoder = e.AMKModularTotalizer
	cfgs[2].ALKEncoder = e.ALKModularTotalizer
	return cfgs
}

func TestIncrementalEncodingSimpleAMK(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	for _, config := range cfgs() {
		numLits := 10
		vars := make([]f.Variable, numLits)
		for i := 0; i < numLits; i++ {
			vars[i] = fac.Var(fmt.Sprintf("v%d", i))
		}
		result := e.ResultForFormula(fac)
		incData, err := e.EncodeIncremental(fac, fac.CC(f.LE, 9, vars...), result, config)
		assert.Nil(err)

		solver := sat.NewSolver(fac)
		solver.Add(fac.CC(f.GE, 4, vars...))
		solver.Add(fac.CC(f.LE, 7, vars...))

		solver.Add(incData.Result.Formulas()...)
		assert.True(solver.Sat())
		solver.Add(incData.NewUpperBound(8)...)
		assert.True(solver.Sat())
		solver.Add(incData.NewUpperBound(7)...)
		assert.True(solver.Sat())
		solver.Add(incData.NewUpperBound(6)...)
		assert.True(solver.Sat())
		solver.Add(incData.NewUpperBound(5)...)
		assert.True(solver.Sat())
		solver.Add(incData.NewUpperBound(4)...)
		assert.True(solver.Sat())

		state := solver.SaveState()
		solver.Add(incData.NewUpperBound(3)...)
		assert.False(solver.Sat())
		solver.LoadState(state)
		assert.True(solver.Sat())

		solver.Add(incData.NewUpperBound(2)...)
		assert.False(solver.Sat())
	}
}

func TestIncrementalEncodingSimpleALK(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	for _, config := range cfgs() {
		numLits := 10
		vars := make([]f.Variable, numLits)
		for i := 0; i < numLits; i++ {
			vars[i] = fac.Var(fmt.Sprintf("v%d", i))
		}
		result := e.ResultForFormula(fac)
		incData, err := e.EncodeIncremental(fac, fac.CC(f.GE, 2, vars...), result, config)
		assert.Nil(err)

		solver := sat.NewSolver(fac)
		solver.Add(fac.CC(f.GE, 4, vars...))
		solver.Add(fac.CC(f.LE, 7, vars...))

		solver.Add(incData.Result.Formulas()...)
		assert.True(solver.Sat())

		solver.Add(incData.NewLowerBound(3)...)
		assert.True(solver.Sat())
		solver.Add(incData.NewLowerBound(4)...)
		assert.True(solver.Sat())
		solver.Add(incData.NewLowerBound(5)...)
		assert.True(solver.Sat())
		solver.Add(incData.NewLowerBound(6)...)
		assert.True(solver.Sat())
		solver.Add(incData.NewLowerBound(7)...)
		assert.True(solver.Sat())

		state := solver.SaveState()
		solver.Add(incData.NewLowerBound(8)...)
		assert.False(solver.Sat())
		solver.LoadState(state)
		assert.True(solver.Sat())
		solver.Add(incData.NewLowerBound(9)...)
		assert.False(solver.Sat())
	}
}

func TestIncrementalEncodingLargeTotalizerAMK(t *testing.T) {
	assert := assert.New(t)
	config := cfgs()[0]
	fac := f.NewFactory()
	numLits := 100
	currentBound := numLits - 1
	vars := make([]f.Variable, numLits)
	for i := 0; i < numLits; i++ {
		vars[i] = fac.Var(fmt.Sprintf("v%d", i))
	}
	result := e.ResultForFormula(fac)
	incData, err := e.EncodeIncremental(fac, fac.CC(f.LE, uint32(currentBound), vars...), result, config)
	assert.Nil(err)

	solver := sat.NewSolver(fac)
	solver.Add(fac.CC(f.GE, 42, vars...))
	solver.Add(incData.Result.Formulas()...)

	// search the lower bound
	for solver.Sat() {
		currentBound--
		solver.Add(incData.NewUpperBound(currentBound)...)
	}
	assert.Equal(41, currentBound)
}

func TestIncrementalEncodingLargeTotalizerALK(t *testing.T) {
	assert := assert.New(t)
	config := cfgs()[0]
	fac := f.NewFactory()
	numLits := 100
	currentBound := 2
	vars := make([]f.Variable, numLits)
	for i := 0; i < numLits; i++ {
		vars[i] = fac.Var(fmt.Sprintf("v%d", i))
	}

	result := e.ResultForFormula(fac)
	incData, err := e.EncodeIncremental(fac, fac.CC(f.GE, uint32(currentBound), vars...), result, config)
	assert.Nil(err)

	solver := sat.NewSolver(fac)
	solver.Add(fac.CC(f.LE, 87, vars...))
	solver.Add(incData.Result.Formulas()...)

	// search the lower bound
	for solver.Sat() {
		currentBound++
		solver.Add(incData.NewLowerBound(currentBound)...)
	}
	assert.Equal(88, currentBound)
}

func TestIncrementalEncodingVeryLargeModularTotalizerAMK(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	assert := assert.New(t)
	config := cfgs()[2]
	fac := f.NewFactory()
	numLits := 300
	currentBound := numLits - 1
	vars := make([]f.Variable, numLits)
	for i := 0; i < numLits; i++ {
		vars[i] = fac.Var(fmt.Sprintf("v%d", i))
	}
	result := e.ResultForFormula(fac)
	incData, err := e.EncodeIncremental(fac, fac.CC(f.LE, uint32(currentBound), vars...), result, config)
	assert.Nil(err)

	solver := sat.NewSolver(fac)
	solver.Add(fac.CC(f.GE, 234, vars...))
	solver.Add(incData.Result.Formulas()...)

	// search the lower bound
	for solver.Sat() {
		currentBound--
		solver.Add(incData.NewUpperBound(currentBound)...)
	}
	assert.Equal(233, currentBound)
}

func TestIncrementalEncodingSimpleAMKSolver(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	for _, config := range cfgs() {
		fac.PutConfiguration(config)
		numLits := 10
		vars := make([]f.Variable, numLits)
		for i := 0; i < numLits; i++ {
			vars[i] = fac.Var(fmt.Sprintf("v%d", i))
		}

		solver := sat.NewSolver(fac)
		solver.Add(fac.CC(f.GE, 4, vars...))
		solver.Add(fac.CC(f.LE, 7, vars...))
		incData, err := solver.AddIncrementalCC(fac.CC(f.LE, 9, vars...))
		assert.Nil(err)
		assert.True(solver.Sat())

		incData.NewUpperBoundForSolver(8)
		assert.True(solver.Sat())
		incData.NewUpperBoundForSolver(7)
		assert.True(solver.Sat())
		incData.NewUpperBoundForSolver(6)
		assert.True(solver.Sat())
		incData.NewUpperBoundForSolver(5)
		assert.True(solver.Sat())
		incData.NewUpperBoundForSolver(4)
		assert.True(solver.Sat())

		state := solver.SaveState()
		incData.NewUpperBoundForSolver(3)
		assert.False(solver.Sat())
		solver.LoadState(state)
		assert.True(solver.Sat())

		incData.NewUpperBoundForSolver(2)
		assert.False(solver.Sat())
	}
}

func TestIncrementalEncodingSimpleALKSolver(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	for _, config := range cfgs() {
		fac.PutConfiguration(config)
		numLits := 10
		vars := make([]f.Variable, numLits)
		for i := 0; i < numLits; i++ {
			vars[i] = fac.Var(fmt.Sprintf("v%d", i))
		}

		solver := sat.NewSolver(fac)
		solver.Add(fac.CC(f.GE, 4, vars...))
		solver.Add(fac.CC(f.LE, 7, vars...))

		incData, err := solver.AddIncrementalCC(fac.CC(f.GE, 2, vars...))
		assert.Nil(err)
		assert.True(solver.Sat())

		incData.NewLowerBoundForSolver(3)
		assert.True(solver.Sat())
		incData.NewLowerBoundForSolver(4)
		assert.True(solver.Sat())
		incData.NewLowerBoundForSolver(5)
		assert.True(solver.Sat())
		incData.NewLowerBoundForSolver(6)
		assert.True(solver.Sat())
		incData.NewLowerBoundForSolver(7)
		assert.True(solver.Sat())

		state := solver.SaveState()
		incData.NewLowerBoundForSolver(8)
		assert.False(solver.Sat())
		solver.LoadState(state)
		assert.True(solver.Sat())
		incData.NewLowerBoundForSolver(9)
		assert.False(solver.Sat())
	}
}

func TestIncrementalEncodingLargeTotalizerAMKSolver(t *testing.T) {
	assert := assert.New(t)
	config := cfgs()[2]
	fac := f.NewFactory()
	numLits := 100
	vars := make([]f.Variable, numLits)
	for i := 0; i < numLits; i++ {
		vars[i] = fac.Var(fmt.Sprintf("v%d", i))
	}
	fac.PutConfiguration(config)

	currentBound := numLits - 1
	solver := sat.NewSolver(fac)
	solver.Add(fac.CC(f.GE, 42, vars...))
	incData, err := solver.AddIncrementalCC(fac.CC(f.LE, uint32(currentBound), vars...))
	assert.Nil(err)

	// search the lower bound
	for solver.Sat() {
		currentBound--
		incData.NewUpperBoundForSolver(currentBound)
	}
	assert.Equal(41, currentBound)
}

func TestIncrementalEncodingLargeTotalizerALKSolver(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	numLits := 100
	vars := make([]f.Variable, numLits)
	for i := 0; i < numLits; i++ {
		vars[i] = fac.Var(fmt.Sprintf("v%d", i))
	}

	for _, config := range cfgs() {
		currentBound := 2
		fac.PutConfiguration(config)

		solver := sat.NewSolver(fac)
		solver.Add(fac.CC(f.LE, 87, vars...))
		incData, err := solver.AddIncrementalCC(fac.CC(f.GE, uint32(currentBound), vars...))
		assert.Nil(err)

		// search the lower bound
		for solver.Sat() {
			currentBound++
			incData.NewLowerBoundForSolver(currentBound)
		}
		assert.Equal(88, currentBound)
	}
}
