package count

import (
	"math/big"
	"testing"

	"booleworks.com/logicng/model/enum"

	"booleworks.com/logicng/normalform"
	"booleworks.com/logicng/parser"
	"booleworks.com/logicng/randomizer"
	"booleworks.com/logicng/sat"

	f "booleworks.com/logicng/formula"
	"github.com/stretchr/testify/assert"
)

func TestModelCounterSimple(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	a := fac.Var("a")
	b := fac.Var("b")
	d := f.NewTestData(fac)

	assert.Zero(big.NewInt(0).Cmp(cnt(t, fac, f.NewVarSet(), d.False)))
	assert.Zero(big.NewInt(0).Cmp(cnt(t, fac, f.NewVarSet(a, b), d.False)))
	assert.Zero(big.NewInt(1).Cmp(cnt(t, fac, f.NewVarSet(), d.True)))
	assert.Zero(big.NewInt(4).Cmp(cnt(t, fac, f.NewVarSet(a, b), d.True)))

	formula := p.ParseUnsafe("(~v1 => ~v0) | ~v1 | v0")
	variables := f.Variables(fac, formula)
	assert.Zero(big.NewInt(4).Cmp(cnt(t, fac, variables, formula)))

	formula = p.ParseUnsafe("(a & b) | ~b")
	variables = f.Variables(fac, formula)
	assert.Zero(big.NewInt(2).Cmp(cnt(t, fac, variables, formula, a.AsFormula())))

	formula = p.ParseUnsafe("a & b & c")
	formula2 := p.ParseUnsafe("c & d")
	variables = f.Variables(fac, formula, formula2)
	assert.Zero(big.NewInt(1).Cmp(cnt(t, fac, variables, formula, formula2)))

	formula = p.ParseUnsafe("a & b & (a + b + c + d <= 1)")
	formula2 = p.ParseUnsafe("a | b")
	variables = f.Variables(fac, formula, formula2)
	assert.Zero(big.NewInt(0).Cmp(cnt(t, fac, variables, formula, formula2)))

	formula = p.ParseUnsafe("a & (a + b + c + d <= 1)")
	formula2 = p.ParseUnsafe("a | b")
	variables = f.Variables(fac, formula, formula2)
	assert.Zero(big.NewInt(1).Cmp(cnt(t, fac, variables, formula, formula2)))

	formula = p.ParseUnsafe("a & (a + b + c + d = 1)")
	formula2 = p.ParseUnsafe("a | b")
	variables = f.Variables(fac, formula, formula2)
	assert.Zero(big.NewInt(1).Cmp(cnt(t, fac, variables, formula, formula2)))
}

func TestModelCounterNQueen(t *testing.T) {
	fac := f.NewFactory()
	testQueens(t, fac, 4, 2)
	testQueens(t, fac, 5, 10)
	testQueens(t, fac, 6, 4)
	testQueens(t, fac, 7, 40)
	testQueens(t, fac, 8, 92)
}

func TestModelCounterCornerCases(t *testing.T) {
	fac := f.NewFactory()
	for _, formula := range f.NewCornerCases(fac) {
		if formula.Sort() != f.SortCC && formula.Sort() != f.SortPBC {
			expCount := enumerationBasedModelCount(fac, formula)
			count, _ := Count(fac, f.Variables(fac, formula).Content(), formula)
			assert.Zero(t, count.Cmp(expCount))
		}
	}
}

func TestModelCounterRandom(t *testing.T) {
	fac := f.NewFactory()
	config := normalform.DefaultCNFConfig()
	config.Algorithm = normalform.CNFPlaistedGreenbaum
	fac.PutConfiguration(config)

	numTests := 1000
	if testing.Short() {
		numTests = 100
	}

	for i := 0; i < numTests; i++ {
		randomizerConfig := randomizer.DefaultConfig()
		randomizerConfig.NumVars = 6
		randomizerConfig.WeightAMO = 5
		randomizerConfig.WeightEXO = 5
		randomizerConfig.Seed = int64(i * 42)
		randomizer := randomizer.New(fac, randomizerConfig)

		formula := randomizer.Formula(4)
		expCount := enumerationBasedModelCount(fac, formula)
		count, _ := Count(fac, f.Variables(fac, formula).Content(), formula)
		assert.Zero(t, count.Cmp(expCount))
	}
}

func testQueens(t *testing.T, fac f.Factory, size int, models int64) {
	queens := sat.GenerateNQueens(fac, size)
	assert.Zero(t, big.NewInt(models).Cmp(cnt(t, fac, f.Variables(fac, queens), queens)))
}

func enumerationBasedModelCount(fac f.Factory, formulas ...f.Formula) *big.Int {
	solver := sat.NewSolver(fac)
	solver.Add(formulas...)
	variables := f.Variables(fac, formulas...)
	models := enum.OnSolver(solver, variables.Content())
	return big.NewInt(int64(len(models)))
}

func cnt(t *testing.T, fac f.Factory, variables *f.VarSet, formulas ...f.Formula) *big.Int {
	c, err := Count(fac, variables.Content(), formulas...)
	assert.Nil(t, err)
	return c
}
