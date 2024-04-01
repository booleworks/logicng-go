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

var amkConfigs = []encoding.Config{
	{AMKEncoder: encoding.AMKTotalizer},
	{AMKEncoder: encoding.AMKModularTotalizer},
	{AMKEncoder: encoding.AMKCardinalityNetwork},
}

func TestAmkFormulas(t *testing.T) {
	fac := f.NewFactory()
	for _, config := range amkConfigs {
		testAmkFormula(t, fac, &config, 10, 0, 1, false)
		testAmkFormula(t, fac, &config, 10, 1, 11, false)
		testAmkFormula(t, fac, &config, 10, 2, 56, false)
		testAmkFormula(t, fac, &config, 10, 3, 176, false)
		testAmkFormula(t, fac, &config, 10, 4, 386, false)
		testAmkFormula(t, fac, &config, 10, 5, 638, false)
		testAmkFormula(t, fac, &config, 10, 6, 848, false)
		testAmkFormula(t, fac, &config, 10, 7, 968, false)
		testAmkFormula(t, fac, &config, 10, 8, 1013, false)
		testAmkFormula(t, fac, &config, 10, 9, 1023, false)
	}
}

func TestLargeAmkFormula(t *testing.T) {
	fac := f.NewFactory()
	for _, config := range amkConfigs {
		testAmkFormula(t, fac, &config, 150, 2, 1+150+11175, false)
	}
}

func TestAmkSolver(t *testing.T) {
	fac := f.NewFactory()
	for _, config := range amkConfigs {
		testAmkFormula(t, fac, &config, 10, 0, 1, true)
		testAmkFormula(t, fac, &config, 10, 1, 11, true)
		testAmkFormula(t, fac, &config, 10, 2, 56, true)
		testAmkFormula(t, fac, &config, 10, 3, 176, true)
		testAmkFormula(t, fac, &config, 10, 4, 386, true)
		testAmkFormula(t, fac, &config, 10, 5, 638, true)
		testAmkFormula(t, fac, &config, 10, 6, 848, true)
		testAmkFormula(t, fac, &config, 10, 7, 968, true)
		testAmkFormula(t, fac, &config, 10, 8, 1013, true)
		testAmkFormula(t, fac, &config, 10, 9, 1023, true)
	}
}

func TestLargeAmkSolver(t *testing.T) {
	fac := f.NewFactory()
	for _, config := range amkConfigs {
		testAmkFormula(t, fac, &config, 150, 2, 1+150+11175, true)
	}
}

func testAmkFormula(
	t *testing.T, fac f.Factory, config *encoding.Config, numLits, rhs, expected int, useSolver bool,
) {
	problemLits := make([]f.Variable, numLits)
	for i := 0; i < numLits; i++ {
		problemLits[i] = fac.Var(fmt.Sprintf("v%d", i))
	}

	cc := fac.CC(f.LE, uint32(rhs), problemLits...)
	satsolver := sat.NewSolver(fac)
	if useSolver {
		satsolver.Add(cc)
	} else {
		encoding, err := encoding.EncodeCC(fac, cc, config)
		assert.Nil(t, err)
		satsolver.Add(encoding...)
	}

	if expected != 0 {
		assert.True(t, satsolver.Sat())
	} else {
		assert.False(t, satsolver.Sat())
	}

	models := enum.OnSolver(satsolver, problemLits)
	assert.Equal(t, expected, len(models))
	for _, model := range models {
		assert.True(t, len(model.PosVars()) <= rhs)
	}
}
