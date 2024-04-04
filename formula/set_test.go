package formula

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariables(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()
	d := NewTestData(fac)

	assert.Equal(NewVarSet().elements, Variables(fac, d.True).elements)
	assert.Equal(NewVarSet().elements, Variables(fac, d.False).elements)
	assert.Equal(NewVarSet(d.VA).elements, Variables(fac, d.A).elements)
	assert.Equal(NewVarSet(d.VA).elements, Variables(fac, d.NA).elements)
	assert.Equal(NewVarSet(d.VA, d.VB).elements, Variables(fac, d.AND1).elements)
	assert.Equal(NewVarSet(d.VA, d.VB).elements, Variables(fac, d.AND2).elements)
	assert.Equal(NewVarSet(d.VX, d.VY).elements, Variables(fac, d.AND3).elements)
	assert.Equal(NewVarSet(d.VX, d.VY).elements, Variables(fac, d.OR1).elements)
	assert.Equal(NewVarSet(d.VX, d.VY).elements, Variables(fac, d.OR2).elements)
	assert.Equal(NewVarSet(d.VA, d.VB).elements, Variables(fac, d.OR3).elements)
	assert.Equal(NewVarSet(d.VA, d.VB).elements, Variables(fac, d.NOT1).elements)
	assert.Equal(NewVarSet(d.VX, d.VY).elements, Variables(fac, d.NOT2).elements)
	assert.Equal(NewVarSet(d.VA, d.VB).elements, Variables(fac, d.IMP1).elements)
	assert.Equal(NewVarSet(d.VA, d.VB).elements, Variables(fac, d.IMP2).elements)
	assert.Equal(NewVarSet(d.VA, d.VB, d.VX, d.VY).elements, Variables(fac, d.IMP3).elements)
	assert.Equal(NewVarSet(d.VA, d.VB, d.VX, d.VY).elements, Variables(fac, d.IMP4).elements)
	assert.Equal(NewVarSet(d.VA, d.VB).elements, Variables(fac, d.EQ1).elements)
	assert.Equal(NewVarSet(d.VA, d.VB).elements, Variables(fac, d.EQ2).elements)
	assert.Equal(NewVarSet(d.VA, d.VB, d.VX, d.VY).elements, Variables(fac, d.EQ3).elements)
	assert.Equal(NewVarSet(d.VA, d.VB).elements, Variables(fac, d.EQ4).elements)
}

func TestLiterals(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()
	d := NewTestData(fac)

	assert.Equal(NewLitSet().elements, Literals(fac, d.True).elements)
	assert.Equal(NewLitSet().elements, Literals(fac, d.False).elements)
	assert.Equal(NewLitSet(d.LA).elements, Literals(fac, d.A).elements)
	assert.Equal(NewLitSet(d.LNA).elements, Literals(fac, d.NA).elements)
	assert.Equal(NewLitSet(d.LA, d.LB).elements, Literals(fac, d.AND1).elements)
	assert.Equal(NewLitSet(d.LNA, d.LNB).elements, Literals(fac, d.AND2).elements)
	assert.Equal(NewLitSet(d.LX, d.LY, d.LNX, d.LNY).elements, Literals(fac, d.AND3).elements)
	assert.Equal(NewLitSet(d.LX, d.LY).elements, Literals(fac, d.OR1).elements)
	assert.Equal(NewLitSet(d.LNX, d.LNY).elements, Literals(fac, d.OR2).elements)
	assert.Equal(NewLitSet(d.LA, d.LB, d.LNA, d.LNB).elements, Literals(fac, d.OR3).elements)
	assert.Equal(NewLitSet(d.LA, d.LB).elements, Literals(fac, d.NOT1).elements)
	assert.Equal(NewLitSet(d.LX, d.LY).elements, Literals(fac, d.NOT2).elements)
	assert.Equal(NewLitSet(d.LA, d.LB).elements, Literals(fac, d.IMP1).elements)
	assert.Equal(NewLitSet(d.LNA, d.LNB).elements, Literals(fac, d.IMP2).elements)
	assert.Equal(NewLitSet(d.LA, d.LB, d.LX, d.LY).elements, Literals(fac, d.IMP3).elements)
	assert.Equal(NewLitSet(d.LA, d.LB, d.LNX, d.LNY).elements, Literals(fac, d.IMP4).elements)
	assert.Equal(NewLitSet(d.LA, d.LB).elements, Literals(fac, d.EQ1).elements)
	assert.Equal(NewLitSet(d.LNA, d.LNB).elements, Literals(fac, d.EQ2).elements)
	assert.Equal(NewLitSet(d.LA, d.LB, d.LX, d.LY).elements, Literals(fac, d.EQ3).elements)
	assert.Equal(NewLitSet(d.LA, d.LB, d.LNA, d.LNB).elements, Literals(fac, d.EQ4).elements)
}
