package formula

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperandsHashFunction(t *testing.T) {
	assert := assert.New(t)

	fac := NewFactory()

	a := fac.Variable("a")
	b := fac.Variable("b")
	c := fac.Variable("c")

	h1 := hashOperands([]Formula{a, b, c})
	h2 := hashOperands([]Formula{a, b, c})
	h3 := hashOperands([]Formula{c, b, a})
	h4 := hashOperands([]Formula{a, c, b})
	h5 := hashOperands([]Formula{b, a, c})
	h6 := hashOperands([]Formula{a, b})
	h7 := hashOperands([]Formula{a, c})
	h8 := hashOperands([]Formula{b, c})

	assert.Equal(h1, h2)
	assert.NotEqual(h1, h3)
	assert.NotEqual(h1, h4)
	assert.NotEqual(h1, h5)

	assert.NotEqual(h6, h1)
	assert.NotEqual(h7, h1)
	assert.NotEqual(h8, h1)
	assert.NotEqual(h6, h7)
	assert.NotEqual(h6, h8)
	assert.NotEqual(h7, h8)
}
