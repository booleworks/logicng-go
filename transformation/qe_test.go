package transformation

import (
	"testing"

	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/parser"
	"github.com/stretchr/testify/assert"
)

func TestQEConstants(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	x := fac.Var("x")
	y := fac.Var("y")

	assert.Equal(fac.Falsum(), ExistentialQE(fac, fac.Falsum()))
	assert.Equal(fac.Falsum(), ExistentialQE(fac, fac.Falsum(), x))
	assert.Equal(fac.Falsum(), ExistentialQE(fac, fac.Falsum(), x, y))
	assert.Equal(fac.Falsum(), UniversalQE(fac, fac.Falsum()))
	assert.Equal(fac.Falsum(), UniversalQE(fac, fac.Falsum(), x))
	assert.Equal(fac.Falsum(), UniversalQE(fac, fac.Falsum(), x, y))

	assert.Equal(fac.Verum(), ExistentialQE(fac, fac.Verum()))
	assert.Equal(fac.Verum(), ExistentialQE(fac, fac.Verum(), x))
	assert.Equal(fac.Verum(), ExistentialQE(fac, fac.Verum(), x, y))
	assert.Equal(fac.Verum(), UniversalQE(fac, fac.Verum()))
	assert.Equal(fac.Verum(), UniversalQE(fac, fac.Verum(), x))
	assert.Equal(fac.Verum(), UniversalQE(fac, fac.Verum(), x, y))
}

func TestQELiterals(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	x := fac.Var("x")
	y := fac.Var("y")

	fx := fac.Variable("x")
	fy := fac.Literal("y", false)
	fz := fac.Variable("z")

	assert.Equal(fx, ExistentialQE(fac, fx))
	assert.Equal(fac.Verum(), ExistentialQE(fac, fx, x))
	assert.Equal(fac.Verum(), ExistentialQE(fac, fx, x, y))
	assert.Equal(fx, UniversalQE(fac, fx))
	assert.Equal(fac.Falsum(), UniversalQE(fac, fx, x))
	assert.Equal(fac.Falsum(), UniversalQE(fac, fx, x, y))

	assert.Equal(fy, ExistentialQE(fac, fy))
	assert.Equal(fy, ExistentialQE(fac, fy, x))
	assert.Equal(fac.Verum(), ExistentialQE(fac, fy, x, y))
	assert.Equal(fy, UniversalQE(fac, fy))
	assert.Equal(fy, UniversalQE(fac, fy, x))
	assert.Equal(fac.Falsum(), UniversalQE(fac, fy, x, y))

	assert.Equal(fz, ExistentialQE(fac, fz))
	assert.Equal(fz, ExistentialQE(fac, fz, x))
	assert.Equal(fz, ExistentialQE(fac, fz, x, y))
	assert.Equal(fz, UniversalQE(fac, fz))
	assert.Equal(fz, UniversalQE(fac, fz, x))
	assert.Equal(fz, UniversalQE(fac, fz, x, y))
}

func TestQEFormulas(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	x := fac.Var("x")
	y := fac.Var("y")

	f1 := p.ParseUnsafe("a & (b | ~c)")
	f2 := p.ParseUnsafe("x & (b | ~c)")
	f3 := p.ParseUnsafe("x & (y | ~c)")

	assert.Equal(f1, ExistentialQE(fac, f1))
	assert.Equal(f1, ExistentialQE(fac, f1, x))
	assert.Equal(f1, ExistentialQE(fac, f1, x, y))
	assert.Equal(f1, UniversalQE(fac, f1))
	assert.Equal(f1, UniversalQE(fac, f1, x))
	assert.Equal(f1, UniversalQE(fac, f1, x, y))

	assert.Equal(f2, ExistentialQE(fac, f2))
	assert.Equal(p.ParseUnsafe("b | ~c"), ExistentialQE(fac, f2, x))
	assert.Equal(p.ParseUnsafe("b | ~c"), ExistentialQE(fac, f2, x, y))
	assert.Equal(f2, UniversalQE(fac, f2))
	assert.Equal(fac.Falsum(), UniversalQE(fac, f2, x))
	assert.Equal(fac.Falsum(), UniversalQE(fac, f2, x, y))

	assert.Equal(f3, ExistentialQE(fac, f3))
	assert.Equal(p.ParseUnsafe("y | ~c"), ExistentialQE(fac, f3, x))
	assert.Equal(fac.Verum(), ExistentialQE(fac, f3, x, y))
	assert.Equal(f3, UniversalQE(fac, f3))
	assert.Equal(fac.Falsum(), UniversalQE(fac, f3, x))
	assert.Equal(fac.Falsum(), UniversalQE(fac, f3, x, y))
}
