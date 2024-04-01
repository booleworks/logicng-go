package cc_test

import (
	"fmt"
	"testing"

	"github.com/booleworks/logicng-go/model/enum"

	"github.com/booleworks/logicng-go/encoding"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/sat"
	"github.com/stretchr/testify/assert"
)

var amoConfigs = []encoding.Config{
	{AMOEncoder: encoding.AMOPure},
	{AMOEncoder: encoding.AMOLadder},
	{AMOEncoder: encoding.AMOBinary},
	{AMOEncoder: encoding.AMOProduct},
	{AMOEncoder: encoding.AMOProduct, ProductRecursiveBound: 10},
	{AMOEncoder: encoding.AMONested},
	{AMOEncoder: encoding.AMONested, NestingGroupSize: 5},
	{AMOEncoder: encoding.AMOCommander},
	{AMOEncoder: encoding.AMOCommander, CommanderGroupSize: 7},
	{AMOEncoder: encoding.AMOBimander, BimanderGroupSize: encoding.BimanderFixed},
	{AMOEncoder: encoding.AMOBimander, BimanderGroupSize: encoding.BimanderHalf},
	{AMOEncoder: encoding.AMOBimander, BimanderGroupSize: encoding.BimanderSqrt},
	{AMOEncoder: encoding.AMOBimander, BimanderGroupSize: encoding.BimanderFixed, BimanderFixedGroupSize: 2},
	{AMOEncoder: encoding.AMOBest},
}

func TestAmoFormulas(t *testing.T) {
	fac := f.NewFactory()
	for _, config := range amoConfigs {
		testAmoFormula(t, fac, &config, 2, false, false)
		testAmoFormula(t, fac, &config, 10, false, false)
		testAmoFormula(t, fac, &config, 100, false, false)
		testAmoFormula(t, fac, &config, 250, false, false)
		testAmoFormula(t, fac, &config, 500, false, false)
	}
}

func TestAmoSolver(t *testing.T) {
	fac := f.NewFactory()
	for _, config := range amoConfigs {
		testAmoFormula(t, fac, &config, 2, false, true)
		testAmoFormula(t, fac, &config, 10, false, true)
		testAmoFormula(t, fac, &config, 100, false, true)
		testAmoFormula(t, fac, &config, 250, false, true)
		testAmoFormula(t, fac, &config, 500, false, true)
	}
}

func TestExoFormulas(t *testing.T) {
	fac := f.NewFactory()
	for _, config := range amoConfigs {
		testAmoFormula(t, fac, &config, 2, true, false)
		testAmoFormula(t, fac, &config, 10, true, false)
		testAmoFormula(t, fac, &config, 100, true, false)
		testAmoFormula(t, fac, &config, 250, true, false)
		testAmoFormula(t, fac, &config, 500, true, false)
	}
}

func TestExoSolvers(t *testing.T) {
	fac := f.NewFactory()
	for _, config := range amoConfigs {
		testAmoFormula(t, fac, &config, 2, true, true)
		testAmoFormula(t, fac, &config, 10, true, true)
		testAmoFormula(t, fac, &config, 100, true, true)
		testAmoFormula(t, fac, &config, 250, true, true)
		testAmoFormula(t, fac, &config, 500, true, true)
	}
}

func testAmoFormula(
	t *testing.T, fac f.Factory, config *encoding.Config, numLits int, exo, useSolver bool,
) {
	problemLits := make([]f.Variable, numLits)
	for i := 0; i < numLits; i++ {
		problemLits[i] = fac.Var(fmt.Sprintf("v%d", i))
	}
	var cc f.Formula
	if exo {
		cc = fac.EXO(problemLits...)
	} else {
		cc = fac.AMO(problemLits...)
	}
	satsolver := sat.NewSolver(fac)
	if useSolver {
		satsolver.Add(cc)
	} else {
		encoding, err := encoding.EncodeCC(fac, cc, config)
		assert.Nil(t, err)
		satsolver.Add(encoding...)
	}

	assert.True(t, satsolver.Sat())
	models := enum.OnSolver(satsolver, problemLits)
	var expected int
	if exo {
		expected = numLits
	} else {
		expected = numLits + 1
	}
	assert.Equal(t, expected, len(models))
	for _, model := range models {
		if exo {
			assert.True(t, len(model.PosVars()) == 1)
		} else {
			assert.True(t, len(model.PosVars()) <= 1)
		}
	}
}
