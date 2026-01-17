package formula

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormulaDepthAtoms(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()
	d := NewTestData(fac)

	assert.Equal(0, Depth(fac, fac.Falsum()))
	assert.Equal(0, Depth(fac, fac.Verum()))
	assert.Equal(0, Depth(fac, d.A))
	assert.Equal(0, Depth(fac, d.NA))
	assert.Equal(0, Depth(fac, d.CC1))
	assert.Equal(0, Depth(fac, d.CC2))
	assert.Equal(0, Depth(fac, d.PBC1))
	assert.Equal(0, Depth(fac, d.PBC2))
	assert.Equal(0, Depth(fac, d.PBC3))
	assert.Equal(0, Depth(fac, d.PBC4))
	assert.Equal(0, Depth(fac, d.PBC5))
}

func TestFormulaDepthDeep(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()
	d := NewTestData(fac)

	assert.Equal(1, Depth(fac, d.AND1))
	assert.Equal(1, Depth(fac, d.AND2))
	assert.Equal(2, Depth(fac, d.AND3))
	assert.Equal(1, Depth(fac, d.OR1))
	assert.Equal(1, Depth(fac, d.OR2))
	assert.Equal(2, Depth(fac, d.OR3))
	assert.Equal(2, Depth(fac, d.NOT1))
	assert.Equal(2, Depth(fac, d.NOT2))
	assert.Equal(1, Depth(fac, d.IMP1))
	assert.Equal(1, Depth(fac, d.IMP2))
	assert.Equal(2, Depth(fac, d.IMP3))
	assert.Equal(2, Depth(fac, d.IMP4))
	assert.Equal(1, Depth(fac, d.EQ1))
	assert.Equal(1, Depth(fac, d.EQ2))
	assert.Equal(2, Depth(fac, d.EQ3))
	assert.Equal(2, Depth(fac, d.EQ4))
}

func TestFormulaTestDeeper(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()

	formula := fac.Variable("Y")
	for i := 0; i < 10; i++ {
		vari := fac.Variable(fmt.Sprintf("X%d", i))
		if i%2 == 0 {
			formula = fac.Or(formula, vari)
		} else {
			formula = fac.And(formula, vari)
		}
	}
	assert.Equal(10, Depth(fac, formula))
}
