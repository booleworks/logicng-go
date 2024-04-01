package simplification

import (
	"testing"

	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/parser"
	"github.com/stretchr/testify/assert"
)

func TestUnitPropagationSimple(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)

	assert.Equal(fac.Falsum(), PropagateUnits(fac, fac.Falsum()))
	assert.Equal(fac.Verum(), PropagateUnits(fac, fac.Verum()))
	assert.Equal(d.A, PropagateUnits(fac, d.A))
	assert.Equal(d.NA, PropagateUnits(fac, d.NA))
	assert.Equal(d.AND1, PropagateUnits(fac, d.AND1))
	assert.Equal(d.AND2, PropagateUnits(fac, d.AND2))
	assert.Equal(d.OR1, PropagateUnits(fac, d.OR1))
	assert.Equal(d.OR2, PropagateUnits(fac, d.OR2))
}

func TestUnitPropagationPropagations(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	p := parser.New(fac)

	assert.Equal(d.AND1, PropagateUnits(fac, fac.And(d.AND1, d.A)))
	assert.Equal(d.False, PropagateUnits(fac, fac.And(d.AND2, d.A)))
	assert.Equal(d.X, PropagateUnits(fac, fac.And(d.OR1, d.X)))
	assert.Equal(d.A, PropagateUnits(fac, fac.Or(d.AND1, d.A)))
	assert.Equal(d.OR1, PropagateUnits(fac, fac.Or(d.OR1, d.X)))
	assert.Equal(
		p.ParseUnsafe("(e | g) & (e | ~g | h) & f & c & d & ~a & b"),
		PropagateUnits(fac, p.ParseUnsafe("(a | b | ~c) & (~a | ~d) & (~c | d) & (~b | e | ~f | g) & (e | f | g | h) & (e | ~f | ~g | h) & f & c")),
	)
}
