package randomizer

import (
	"fmt"
	"testing"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/function"
	"github.com/stretchr/testify/assert"
)

func TestRandomizerDeterminism(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	expected := NewWithSeed(fac, 42).Formula(3)

	assert.Equal(expected, NewWithSeed(fac, 42).Formula(3))
	assert.NotEqual(expected, NewWithSeed(fac, 43).Formula(3))
	assert.NotEqual(expected, New(fac).Formula(3))

	expectedList := randomFormulas(fac)
	for i := 0; i < 10; i++ {
		assert.Equal(expectedList, randomFormulas(fac))
	}
}

func TestRandomizerConstant(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	random := NewWithSeed(fac, 42)
	numTrue := 0

	for i := 0; i < 100; i++ {
		constant := random.Constant()
		assert.True(constant.Sort() <= f.SortTrue)
		if constant == fac.Verum() {
			numTrue++
		}
	}
	assert.True(numTrue >= 40 && numTrue <= 60)
}

func TestRandomizerVariable(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	vars := fac.Vars("A", "B", "C")
	config := DefaultConfig()
	config.Seed = 42
	config.Variables = vars
	random := New(fac, config)
	numA := 0
	numB := 0
	numC := 0
	for i := 0; i < 100; i++ {
		variable := random.Variable()
		name, _ := fac.VarName(variable)
		assert.Contains([]string{"A", "B", "C"}, name)
		switch name {
		case "A":
			numA++
		case "B":
			numB++
		case "C":
			numC++
		}
	}
	assert.True(numA >= 20 && numA <= 40)
	assert.True(numB >= 20 && numB <= 40)
	assert.True(numC >= 20 && numC <= 40)

	vars2 := make([]f.Variable, 20)
	for i := 0; i < len(vars2); i++ {
		vars2[i] = fac.Var(fmt.Sprintf("TEST_VAR_%d", i))
	}
	config = DefaultConfig()
	config.Variables = vars2
	config.WeightPBC = 1
	config.WeightCC = 1
	config.WeightAMO = 1
	config.WeightEXO = 1
	config.Seed = 42
	random = New(fac, config)
	for i := 0; i < 100; i++ {
		formula := random.Formula(4)
		for _, v := range f.Variables(fac, formula).Content() {
			assert.Contains(vars2, v)
		}
	}
}

func TestRandomizerLiteral(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	config := DefaultConfig()
	config.Seed = 42
	config.WeightPosLit = 40
	config.WeightNegLit = 60
	random := New(fac, config)

	numPos := 0
	for i := 0; i < 100; i++ {
		literal := random.Literal()
		if literal.IsPos() {
			numPos++
		}
	}
	assert.True(numPos >= 30 && numPos <= 50)
}

func TestRandomizerAtom(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	config := DefaultConfig()
	config.Seed = 42
	config.WeightConstant = 1
	config.WeightPosLit = 2
	config.WeightNegLit = 3
	config.WeightPBC = 4
	config.WeightCC = 5
	config.WeightAMO = 6
	config.WeightEXO = 7
	random := New(fac, config)

	numConst, numPos, numNeg, numPbc, numCc, numAmo, numExo := 0, 0, 0, 0, 0, 0, 0
	for i := 0; i < 1000; i++ {
		formula := random.Atom()
		assert.True(formula.IsAtomic())
		if formula.Sort() <= f.SortTrue {
			numConst++
		} else if formula.Sort() == f.SortLiteral {
			if formula.IsPos() {
				numPos++
			} else {
				numNeg++
			}
		} else if formula.Sort() == f.SortCC {
			sort, rhs, _, _, _ := fac.PBCOps(formula)
			if rhs == 1 && sort == f.LE {
				numAmo++
			} else if rhs == 1 && sort == f.EQ {
				numExo++
			} else {
				numCc++
			}
		} else {
			numPbc++
		}
	}

	assert.True(numExo >= int(.7*7/6*float64(numAmo)) && numExo <= int(1.3*7/6*float64(numAmo)))
	assert.True(numAmo >= int(.7*6/5*float64(numCc)) && numAmo <= int(1.3*6/5*float64(numCc)))
	assert.True(numCc >= int(.7*5/4*float64(numPbc)) && numCc <= int(1.3*5/4*float64(numPbc)))
	assert.True(numPbc >= int(.7*4/3*float64(numNeg)) && numPbc <= int(1.3*4/3*float64(numNeg)))
	assert.True(numNeg >= int(.7*3/2*float64(numPos)) && numNeg <= int(1.3*3/2*float64(numPos)))
	assert.True(numPos >= int(.7*2/1*float64(numConst)) && numPos <= int(1.3*2/1*float64(numConst)))

	config = DefaultConfig()
	config.Seed = 42
	config.WeightConstant = 0
	config.WeightPosLit = 3
	config.WeightNegLit = 6
	random = New(fac, config)
	for i := 0; i < 100; i++ {
		assert.Equal(f.SortLiteral, random.Atom().Sort())
	}
}

func TestRandomizerAnd(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	random := NewWithSeed(fac, 42)

	for i := 0; i < 10; i++ {
		assert.True(random.And(0).IsAtomic())
	}
	for depth := 1; depth <= 7; depth++ {
		for i := 0; i < 10; i++ {
			formula := random.And(depth)
			assert.Equal(f.SortAnd, formula.Sort())
			assert.True(function.FormulaDepth(fac, formula) <= depth)
		}
	}
}

func TestRandomizerOr(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	random := NewWithSeed(fac, 42)

	for i := 0; i < 10; i++ {
		assert.True(random.Or(0).IsAtomic())
	}
	for depth := 1; depth <= 7; depth++ {
		for i := 0; i < 10; i++ {
			formula := random.Or(depth)
			assert.Equal(f.SortOr, formula.Sort())
			assert.True(function.FormulaDepth(fac, formula) <= depth)
		}
	}
}

func TestRandomizerNot(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	random := NewWithSeed(fac, 42)

	for i := 0; i < 10; i++ {
		assert.True(random.Not(0).IsAtomic())
		assert.True(random.Not(1).IsAtomic())
	}
	for depth := 2; depth <= 7; depth++ {
		for i := 0; i < 10; i++ {
			formula := random.Not(depth)
			assert.Equal(f.SortNot, formula.Sort())
			assert.True(function.FormulaDepth(fac, formula) <= depth)
		}
	}
}

func TestRandomizerImpl(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	random := NewWithSeed(fac, 42)

	for i := 0; i < 10; i++ {
		assert.True(random.Impl(0).IsAtomic())
	}
	for depth := 1; depth <= 7; depth++ {
		for i := 0; i < 10; i++ {
			formula := random.Impl(depth)
			assert.Equal(f.SortImpl, formula.Sort())
			assert.True(function.FormulaDepth(fac, formula) <= depth)
		}
	}
}

func TestRandomizerEquiv(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	random := NewWithSeed(fac, 42)

	for i := 0; i < 10; i++ {
		assert.True(random.Equiv(0).IsAtomic())
	}
	for depth := 1; depth <= 7; depth++ {
		for i := 0; i < 10; i++ {
			formula := random.Equiv(depth)
			assert.Equal(f.SortEquiv, formula.Sort())
			assert.True(function.FormulaDepth(fac, formula) <= depth)
		}
	}
}

func randomFormulas(fac f.Factory) []f.Formula {
	random := New(fac)
	formulas := make([]f.Formula, 10, 15)
	formulas[0] = random.Constant()
	formulas[1] = random.Variable().AsFormula()
	formulas[2] = random.Literal().AsFormula()
	formulas[3] = random.Atom()
	formulas[4] = random.And(3)
	formulas[5] = random.Or(3)
	formulas[6] = random.Not(3)
	formulas[7] = random.Impl(3)
	formulas[8] = random.Equiv(3)
	formulas[9] = random.Formula(3)
	formulas = append(formulas, random.ConstraintSet(5, 3)...)
	return formulas
}
