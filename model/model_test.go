package model

import (
	"testing"

	f "booleworks.com/logicng/formula"
	"github.com/stretchr/testify/assert"
)

func TestModelSprint(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	assert.Equal("[]", New().Sprint(fac))
	assert.Equal("[a]", New(d.LA).Sprint(fac))
	assert.Equal("[a, b]", New(d.LA, d.LB).Sprint(fac))
	assert.Equal("[~a]", New(d.LNA).Sprint(fac))
	assert.Equal("[~a, ~b]", New(d.LNA, d.LNB).Sprint(fac))
	assert.Equal("[a, ~b]", New(d.LA, d.LNB).Sprint(fac))
	assert.Equal("[a, ~b, ~c, d]", New(d.LA, d.LNB, d.LC.Negate(fac), d.LD).Sprint(fac))
}
