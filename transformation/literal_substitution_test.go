package transformation

import (
	"testing"

	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/parser"
	"github.com/stretchr/testify/assert"
)

func TestLiteralSubstition(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	ls := make(map[f.Literal]f.Literal)
	ls[fac.Lit("a", true)] = fac.Lit("a_t", true)
	ls[fac.Lit("a", false)] = fac.Lit("a_f", true)
	ls[fac.Lit("b", false)] = fac.Lit("x", true)
	ls[fac.Lit("c", true)] = fac.Lit("y", true)

	assert.Equal(p.ParseUnsafe("$true"), SubstituteLiterals(fac, p.ParseUnsafe("$true"), &ls))
	assert.Equal(p.ParseUnsafe("$false"), SubstituteLiterals(fac, p.ParseUnsafe("$false"), &ls))
	assert.Equal(p.ParseUnsafe("m"), SubstituteLiterals(fac, p.ParseUnsafe("m"), &ls))
	assert.Equal(p.ParseUnsafe("~m"), SubstituteLiterals(fac, p.ParseUnsafe("~m"), &ls))
	assert.Equal(p.ParseUnsafe("a_t"), SubstituteLiterals(fac, p.ParseUnsafe("a"), &ls))
	assert.Equal(p.ParseUnsafe("a_f"), SubstituteLiterals(fac, p.ParseUnsafe("~a"), &ls))
	assert.Equal(p.ParseUnsafe("b"), SubstituteLiterals(fac, p.ParseUnsafe("b"), &ls))
	assert.Equal(p.ParseUnsafe("x"), SubstituteLiterals(fac, p.ParseUnsafe("~b"), &ls))
	assert.Equal(p.ParseUnsafe("y"), SubstituteLiterals(fac, p.ParseUnsafe("c"), &ls))
	assert.Equal(p.ParseUnsafe("~y"), SubstituteLiterals(fac, p.ParseUnsafe("~c"), &ls))

	assert.Equal(p.ParseUnsafe("~(a_t & b & ~y & x)"), SubstituteLiterals(fac, p.ParseUnsafe("~(a & b & ~c & x)"), &ls))
	assert.Equal(p.ParseUnsafe("a_t & b & ~y & x"), SubstituteLiterals(fac, p.ParseUnsafe("a & b & ~c & x"), &ls))
	assert.Equal(p.ParseUnsafe("(a_t | b) <=> (~y | x)"), SubstituteLiterals(fac, p.ParseUnsafe("(a | b) <=> (~c | x)"), &ls))
	assert.Equal(p.ParseUnsafe("2*a_t + 3*x + -4*~y + x <= 5"), SubstituteLiterals(fac, p.ParseUnsafe("2*a + 3*~b + -4*~c + x <= 5"), &ls))

	clear(ls)
	assert.Equal(p.ParseUnsafe("2*a + 3*~b + -4*~c + x <= 5"), SubstituteLiterals(fac, p.ParseUnsafe("2*a + 3*~b + -4*~c + x <= 5"), &ls))
}
