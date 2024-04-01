package formula

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultPrinter(t *testing.T) {
	fac := NewFactory()
	d := NewTestData(fac)
	assert := assert.New(t)
	assert.Equal("$false", d.False.Sprint(fac))
	assert.Equal("$true", d.True.Sprint(fac))
	assert.Equal("x", d.X.Sprint(fac))
	assert.Equal("~a", d.NA.Sprint(fac))
	assert.Equal("~(a & b)", fac.Not(d.AND1).Sprint(fac))
	assert.Equal("~a => ~b", d.IMP2.Sprint(fac))
	assert.Equal("a & b => x | y", d.IMP3.Sprint(fac))
	assert.Equal("a => b <=> ~a => ~b", d.EQ4.Sprint(fac))
	assert.Equal("(x | y) & (~x | ~y)", d.AND3.Sprint(fac))
	assert.Equal("a & b & c & x", fac.And(d.A, d.B, d.C, d.X).Sprint(fac))
	assert.Equal("a | b | c | x", fac.Or(d.A, d.B, d.C, d.X).Sprint(fac))
	assert.Equal("a | b & c | x", fac.Or(d.A, fac.And(d.B, d.C), d.X).Sprint(fac))
	assert.Equal("a & (b | ~a) & x", fac.And(d.A, fac.Or(d.B, d.NA), d.X).Sprint(fac))

	assert.Equal("a < 1", d.CC1.Sprint(fac))
	assert.Equal("a + b + c >= 2", d.CC2.Sprint(fac))
	assert.Equal("a <= 1", d.AMO1.Sprint(fac))
	assert.Equal("a + b + c <= 1", d.AMO2.Sprint(fac))
	assert.Equal("a = 1", d.EXO1.Sprint(fac))
	assert.Equal("a + b + c = 1", d.EXO2.Sprint(fac))
	assert.Equal("3*a <= 2", d.PB1.Sprint(fac))
	assert.Equal("3*a + -2*~b + 7*c <= 8", d.PB2.Sprint(fac))
	assert.Equal("2*a + -4*b + 3*x = 2", d.PBC1.Sprint(fac))
	assert.Equal("2*a + -4*b + 3*x > 2", d.PBC2.Sprint(fac))
	assert.Equal("2*a + -4*b + 3*x >= 2", d.PBC3.Sprint(fac))
	assert.Equal("2*a + -4*b + 3*x < 2", d.PBC4.Sprint(fac))
	assert.Equal("2*a + -4*b + 3*x <= 2", d.PBC5.Sprint(fac))
}
