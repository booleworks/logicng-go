package function

import (
	"fmt"
	"testing"

	f "booleworks.com/logicng/formula"
	"github.com/stretchr/testify/assert"
)

func TestFormulaDepthAtoms(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)

	assert.Equal(0, FormulaDepth(fac, fac.Falsum()))
	assert.Equal(0, FormulaDepth(fac, fac.Verum()))
	assert.Equal(0, FormulaDepth(fac, d.A))
	assert.Equal(0, FormulaDepth(fac, d.NA))
	assert.Equal(0, FormulaDepth(fac, d.CC1))
	assert.Equal(0, FormulaDepth(fac, d.CC2))
	assert.Equal(0, FormulaDepth(fac, d.PBC1))
	assert.Equal(0, FormulaDepth(fac, d.PBC2))
	assert.Equal(0, FormulaDepth(fac, d.PBC3))
	assert.Equal(0, FormulaDepth(fac, d.PBC4))
	assert.Equal(0, FormulaDepth(fac, d.PBC5))
}

func TestFormulaDepthDeep(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)

	assert.Equal(1, FormulaDepth(fac, d.AND1))
	assert.Equal(1, FormulaDepth(fac, d.AND2))
	assert.Equal(2, FormulaDepth(fac, d.AND3))
	assert.Equal(1, FormulaDepth(fac, d.OR1))
	assert.Equal(1, FormulaDepth(fac, d.OR2))
	assert.Equal(2, FormulaDepth(fac, d.OR3))
	assert.Equal(2, FormulaDepth(fac, d.NOT1))
	assert.Equal(2, FormulaDepth(fac, d.NOT2))
	assert.Equal(1, FormulaDepth(fac, d.IMP1))
	assert.Equal(1, FormulaDepth(fac, d.IMP2))
	assert.Equal(2, FormulaDepth(fac, d.IMP3))
	assert.Equal(2, FormulaDepth(fac, d.IMP4))
	assert.Equal(1, FormulaDepth(fac, d.EQ1))
	assert.Equal(1, FormulaDepth(fac, d.EQ2))
	assert.Equal(2, FormulaDepth(fac, d.EQ3))
	assert.Equal(2, FormulaDepth(fac, d.EQ4))
}

func TestFormulaTestDeeper(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()

	formula := fac.Variable("Y")
	for i := 0; i < 10; i++ {
		vari := fac.Variable(fmt.Sprintf("X%d", i))
		if i%2 == 0 {
			formula = fac.Or(formula, vari)
		} else {
			formula = fac.And(formula, vari)
		}
	}
	assert.Equal(10, FormulaDepth(fac, formula))
}
