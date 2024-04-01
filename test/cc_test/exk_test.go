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

var exkConfigs = []encoding.Config{
	{EXKEncoder: encoding.EXKTotalizer},
	{EXKEncoder: encoding.EXKCardinalityNetwork},
}

func TestExkFormulas(t *testing.T) {
	fac := f.NewFactory()
	for _, config := range exkConfigs {
		testExkFormula(t, fac, &config, 10, 1, 10, false)
		testExkFormula(t, fac, &config, 10, 2, 45, false)
		testExkFormula(t, fac, &config, 10, 3, 120, false)
		testExkFormula(t, fac, &config, 10, 4, 210, false)
		testExkFormula(t, fac, &config, 10, 5, 252, false)
		testExkFormula(t, fac, &config, 10, 6, 210, false)
		testExkFormula(t, fac, &config, 10, 7, 120, false)
		testExkFormula(t, fac, &config, 10, 8, 45, false)
		testExkFormula(t, fac, &config, 10, 9, 10, false)
		testExkFormula(t, fac, &config, 10, 10, 1, false)
		testExkFormula(t, fac, &config, 10, 12, 0, false)
	}
}

func TestExkSolver(t *testing.T) {
	fac := f.NewFactory()
	for _, config := range exkConfigs {
		testExkFormula(t, fac, &config, 10, 1, 10, true)
		testExkFormula(t, fac, &config, 10, 2, 45, true)
		testExkFormula(t, fac, &config, 10, 3, 120, true)
		testExkFormula(t, fac, &config, 10, 4, 210, true)
		testExkFormula(t, fac, &config, 10, 5, 252, true)
		testExkFormula(t, fac, &config, 10, 6, 210, true)
		testExkFormula(t, fac, &config, 10, 7, 120, true)
		testExkFormula(t, fac, &config, 10, 8, 45, true)
		testExkFormula(t, fac, &config, 10, 9, 10, true)
		testExkFormula(t, fac, &config, 10, 10, 1, true)
		testExkFormula(t, fac, &config, 10, 12, 0, true)
	}
}

func testExkFormula(
	t *testing.T, fac f.Factory, config *encoding.Config, numLits, rhs, expected int, useSolver bool,
) {
	problemLits := make([]f.Variable, numLits)
	for i := 0; i < numLits; i++ {
		problemLits[i] = fac.Var(fmt.Sprintf("v%d", i))
	}

	cc := fac.CC(f.EQ, uint32(rhs), problemLits...)
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
		assert.True(t, len(model.PosVars()) == rhs)
	}
}
