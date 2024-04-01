package encoding

import (
	"testing"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/stretchr/testify/assert"
)

func TestPbcNormalizationTrivial(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()

	lits := []f.Literal{fac.Lit("a", true), fac.Lit("b", false), fac.Lit("c", true), fac.Lit("d", true)}
	coeffs := []int{2, -2, 3, 0}

	pb1 := fac.PBC(f.LE, 4, lits, coeffs)
	pb2 := fac.PBC(f.LE, 5, lits, coeffs)
	pb3 := fac.PBC(f.LE, 7, lits, coeffs)
	pb4 := fac.PBC(f.LE, 10, lits, coeffs)
	pb5 := fac.PBC(f.LE, -3, lits, coeffs)

	assert.Equal("2*a + 2*b + 3*c <= 6", Normalize(fac, pb1).Sprint(fac))
	assert.Equal(fac.Verum(), Normalize(fac, pb2))
	assert.Equal(fac.Verum(), Normalize(fac, pb3))
	assert.Equal(fac.Verum(), Normalize(fac, pb4))
	assert.Equal(fac.Falsum(), Normalize(fac, pb5))
}

func TestPbcNormalization(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()

	lits := []f.Literal{fac.Lit("a", true), fac.Lit("b", false), fac.Lit("c", true), fac.Lit("d", true), fac.Lit("b", false)}
	coeffs := []int{2, -3, 3, 0, 1}
	pb1 := fac.PBC(f.EQ, 2, lits, coeffs)
	pb2 := fac.PBC(f.GE, 1, lits, coeffs)
	pb3 := fac.PBC(f.GT, 0, lits, coeffs)
	pb4 := fac.PBC(f.LE, 1, lits, coeffs)
	pb5 := fac.PBC(f.LT, 2, lits, coeffs)

	assert.Equal("(2*a + 2*b + 3*c <= 4) & (2*~a + 2*~b + 3*~c <= 3)", Normalize(fac, pb1).Sprint(fac))
	assert.Equal("2*~a + 2*~b + 3*~c <= 4", Normalize(fac, pb2).Sprint(fac))
	assert.Equal("2*~a + 2*~b + 3*~c <= 4", Normalize(fac, pb3).Sprint(fac))
	assert.Equal("2*a + 2*b + 3*c <= 3", Normalize(fac, pb4).Sprint(fac))
	assert.Equal("2*a + 2*b + 3*c <= 3", Normalize(fac, pb5).Sprint(fac))
}

func TestPbcNormalizationSimplification(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()

	lits := []f.Literal{fac.Lit("a", true), fac.Lit("a", true), fac.Lit("c", true), fac.Lit("d", true)}
	coeffs := []int{2, -2, 4, 4}
	pb1 := fac.PBC(f.LE, 4, lits, coeffs)
	assert.Equal("c + d <= 1", Normalize(fac, pb1).Sprint(fac))

	lits = []f.Literal{fac.Lit("a", true), fac.Lit("a", false), fac.Lit("c", true), fac.Lit("d", true)}
	coeffs = []int{2, 2, 4, 2}
	pb2 := fac.PBC(f.LE, 4, lits, coeffs)
	assert.Equal("2*c + d <= 1", Normalize(fac, pb2).Sprint(fac))
}
