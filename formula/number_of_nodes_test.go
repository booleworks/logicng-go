package formula

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumberOfNodes(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()
	d := NewTestData(fac)

	assert.Equal(1, NumberOfNodes(fac, d.False))
	assert.Equal(1, NumberOfNodes(fac, d.True))
	assert.Equal(1, NumberOfNodes(fac, d.A))
	assert.Equal(1, NumberOfNodes(fac, d.NA))
	assert.Equal(3, NumberOfNodes(fac, d.OR1))
	assert.Equal(3, NumberOfNodes(fac, d.OR2))
	assert.Equal(7, NumberOfNodes(fac, d.OR3))
	assert.Equal(3, NumberOfNodes(fac, d.AND1))
	assert.Equal(3, NumberOfNodes(fac, d.AND2))
	assert.Equal(7, NumberOfNodes(fac, d.AND3))
	assert.Equal(3, NumberOfNodes(fac, d.IMP1))
	assert.Equal(3, NumberOfNodes(fac, d.IMP2))
	assert.Equal(7, NumberOfNodes(fac, d.IMP3))
	assert.Equal(7, NumberOfNodes(fac, d.IMP4))
	assert.Equal(3, NumberOfNodes(fac, d.EQ1))
	assert.Equal(3, NumberOfNodes(fac, d.EQ2))
	assert.Equal(7, NumberOfNodes(fac, d.EQ3))
	assert.Equal(7, NumberOfNodes(fac, d.EQ4))
}
