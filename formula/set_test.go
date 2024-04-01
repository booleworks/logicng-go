package formula

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariables(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()
	d := NewTestData(fac)

	assert.Equal(NewVarSet(), Variables(fac, d.True))
	assert.Equal(NewVarSet(), Variables(fac, d.False))
	assert.Equal(NewVarSet(d.VA), Variables(fac, d.A))
	assert.Equal(NewVarSet(d.VA), Variables(fac, d.NA))
	assert.Equal(NewVarSet(d.VA, d.VB), Variables(fac, d.AND1))
	assert.Equal(NewVarSet(d.VA, d.VB), Variables(fac, d.AND2))
	assert.Equal(NewVarSet(d.VX, d.VY), Variables(fac, d.AND3))
	assert.Equal(NewVarSet(d.VX, d.VY), Variables(fac, d.OR1))
	assert.Equal(NewVarSet(d.VX, d.VY), Variables(fac, d.OR2))
	assert.Equal(NewVarSet(d.VA, d.VB), Variables(fac, d.OR3))
	assert.Equal(NewVarSet(d.VA, d.VB), Variables(fac, d.NOT1))
	assert.Equal(NewVarSet(d.VX, d.VY), Variables(fac, d.NOT2))
	assert.Equal(NewVarSet(d.VA, d.VB), Variables(fac, d.IMP1))
	assert.Equal(NewVarSet(d.VA, d.VB), Variables(fac, d.IMP2))
	assert.Equal(NewVarSet(d.VA, d.VB, d.VX, d.VY), Variables(fac, d.IMP3))
	assert.Equal(NewVarSet(d.VA, d.VB, d.VX, d.VY), Variables(fac, d.IMP4))
	assert.Equal(NewVarSet(d.VA, d.VB), Variables(fac, d.EQ1))
	assert.Equal(NewVarSet(d.VA, d.VB), Variables(fac, d.EQ2))
	assert.Equal(NewVarSet(d.VA, d.VB, d.VX, d.VY), Variables(fac, d.EQ3))
	assert.Equal(NewVarSet(d.VA, d.VB), Variables(fac, d.EQ4))
}

func TestLiterals(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()
	d := NewTestData(fac)

	assert.Equal(NewLitSet(), Literals(fac, d.True))
	assert.Equal(NewLitSet(), Literals(fac, d.False))
	assert.Equal(NewLitSet(d.LA), Literals(fac, d.A))
	assert.Equal(NewLitSet(d.LNA), Literals(fac, d.NA))
	assert.Equal(NewLitSet(d.LA, d.LB), Literals(fac, d.AND1))
	assert.Equal(NewLitSet(d.LNA, d.LNB), Literals(fac, d.AND2))
	assert.Equal(NewLitSet(d.LX, d.LY, d.LNX, d.LNY), Literals(fac, d.AND3))
	assert.Equal(NewLitSet(d.LX, d.LY), Literals(fac, d.OR1))
	assert.Equal(NewLitSet(d.LNX, d.LNY), Literals(fac, d.OR2))
	assert.Equal(NewLitSet(d.LA, d.LB, d.LNA, d.LNB), Literals(fac, d.OR3))
	assert.Equal(NewLitSet(d.LA, d.LB), Literals(fac, d.NOT1))
	assert.Equal(NewLitSet(d.LX, d.LY), Literals(fac, d.NOT2))
	assert.Equal(NewLitSet(d.LA, d.LB), Literals(fac, d.IMP1))
	assert.Equal(NewLitSet(d.LNA, d.LNB), Literals(fac, d.IMP2))
	assert.Equal(NewLitSet(d.LA, d.LB, d.LX, d.LY), Literals(fac, d.IMP3))
	assert.Equal(NewLitSet(d.LA, d.LB, d.LNX, d.LNY), Literals(fac, d.IMP4))
	assert.Equal(NewLitSet(d.LA, d.LB), Literals(fac, d.EQ1))
	assert.Equal(NewLitSet(d.LNA, d.LNB), Literals(fac, d.EQ2))
	assert.Equal(NewLitSet(d.LA, d.LB, d.LX, d.LY), Literals(fac, d.EQ3))
	assert.Equal(NewLitSet(d.LA, d.LB, d.LNA, d.LNB), Literals(fac, d.EQ4))
}
