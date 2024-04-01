package function

import (
	"testing"

	f "booleworks.com/logicng/formula"
	"github.com/stretchr/testify/assert"
)

func TestNumberOfAtoms(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)

	assert.Equal(1, NumberOfAtoms(fac, d.False))
	assert.Equal(1, NumberOfAtoms(fac, d.True))
	assert.Equal(1, NumberOfAtoms(fac, d.A))
	assert.Equal(1, NumberOfAtoms(fac, d.NA))
	assert.Equal(2, NumberOfAtoms(fac, d.OR1))
	assert.Equal(2, NumberOfAtoms(fac, d.OR2))
	assert.Equal(4, NumberOfAtoms(fac, d.OR3))
	assert.Equal(2, NumberOfAtoms(fac, d.AND1))
	assert.Equal(2, NumberOfAtoms(fac, d.AND2))
	assert.Equal(4, NumberOfAtoms(fac, d.AND3))
	assert.Equal(2, NumberOfAtoms(fac, d.IMP1))
	assert.Equal(2, NumberOfAtoms(fac, d.IMP2))
	assert.Equal(4, NumberOfAtoms(fac, d.IMP3))
	assert.Equal(4, NumberOfAtoms(fac, d.IMP4))
	assert.Equal(2, NumberOfAtoms(fac, d.EQ1))
	assert.Equal(2, NumberOfAtoms(fac, d.EQ2))
	assert.Equal(4, NumberOfAtoms(fac, d.EQ3))
	assert.Equal(4, NumberOfAtoms(fac, d.EQ4))
}
