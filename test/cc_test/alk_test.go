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

var alkConfigs = []encoding.Config{
	{ALKEncoder: encoding.ALKTotalizer},
	{ALKEncoder: encoding.ALKModularTotalizer},
	{ALKEncoder: encoding.ALKCardinalityNetwork},
}

func TestAlkFormulas(t *testing.T) {
	fac := f.NewFactory()
	for _, config := range alkConfigs {
		testAlkFormula(t, fac, &config, 10, 0, 1024, false)
		testAlkFormula(t, fac, &config, 10, 1, 1023, false)
		testAlkFormula(t, fac, &config, 10, 2, 1013, false)
		testAlkFormula(t, fac, &config, 10, 3, 968, false)
		testAlkFormula(t, fac, &config, 10, 4, 848, false)
		testAlkFormula(t, fac, &config, 10, 5, 638, false)
		testAlkFormula(t, fac, &config, 10, 6, 386, false)
		testAlkFormula(t, fac, &config, 10, 7, 176, false)
		testAlkFormula(t, fac, &config, 10, 8, 56, false)
		testAlkFormula(t, fac, &config, 10, 9, 11, false)
		testAlkFormula(t, fac, &config, 10, 10, 1, false)
		testAlkFormula(t, fac, &config, 10, 12, 0, false)
	}
}

func TestAlkSolver(t *testing.T) {
	fac := f.NewFactory()
	for _, config := range alkConfigs {
		testAlkFormula(t, fac, &config, 10, 0, 1024, true)
		testAlkFormula(t, fac, &config, 10, 1, 1023, true)
		testAlkFormula(t, fac, &config, 10, 2, 1013, true)
		testAlkFormula(t, fac, &config, 10, 3, 968, true)
		testAlkFormula(t, fac, &config, 10, 4, 848, true)
		testAlkFormula(t, fac, &config, 10, 5, 638, true)
		testAlkFormula(t, fac, &config, 10, 6, 386, true)
		testAlkFormula(t, fac, &config, 10, 7, 176, true)
		testAlkFormula(t, fac, &config, 10, 8, 56, true)
		testAlkFormula(t, fac, &config, 10, 9, 11, true)
		testAlkFormula(t, fac, &config, 10, 10, 1, true)
		testAlkFormula(t, fac, &config, 10, 12, 0, true)
	}
}

func testAlkFormula(
	t *testing.T, fac f.Factory, config *encoding.Config, numLits, rhs, expected int, useSolver bool,
) {
	problemLits := make([]f.Variable, numLits)
	for i := 0; i < numLits; i++ {
		problemLits[i] = fac.Var(fmt.Sprintf("v%d", i))
	}

	cc := fac.CC(f.GE, uint32(rhs), problemLits...)
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
		assert.True(t, len(model.PosVars()) >= rhs)
	}
}
