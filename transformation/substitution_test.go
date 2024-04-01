package transformation

import (
	"testing"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/parser"
	"github.com/stretchr/testify/assert"
)

func TestSubstitutionConstant(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	subst := getSubstitution(fac)

	assert.Equal(fac.Falsum(), sub(t, fac, fac.Falsum(), subst))
	assert.Equal(fac.Verum(), sub(t, fac, fac.Verum(), subst))
}

func TestSubstitutionLiteral(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	subst := getSubstitution(fac)

	assert.Equal(d.C, sub(t, fac, d.C, subst))
	assert.Equal(d.NA, sub(t, fac, d.A, subst))
	assert.Equal(d.OR1, sub(t, fac, d.B, subst))
	assert.Equal(d.AND1, sub(t, fac, d.X, subst))
	assert.Equal(d.A, sub(t, fac, d.NA, subst))
	assert.Equal(d.NOT2, sub(t, fac, d.NB, subst))
	assert.Equal(d.NOT1, sub(t, fac, d.NX, subst))
}

func TestSubstitutionNot(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	p := parser.New(fac)
	subst := getSubstitution(fac)

	assert.Equal(p.ParseUnsafe("~(~a & (x | y))"), sub(t, fac, d.NOT1, subst))
	assert.Equal(p.ParseUnsafe("~(a & b | y)"), sub(t, fac, d.NOT2, subst))
}

func TestSubstitutionBinary(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	p := parser.New(fac)
	subst := getSubstitution(fac)

	assert.Equal(p.ParseUnsafe("~a => (x | y)"), sub(t, fac, d.IMP1, subst))
	assert.Equal(p.ParseUnsafe("(~a <=> (x | y)) => (~(a & b) <=> ~y)"), sub(t, fac, d.IMP4, subst))
	assert.Equal(p.ParseUnsafe("a <=> ~(x | y)"), sub(t, fac, d.EQ2, subst))
	assert.Equal(p.ParseUnsafe("(~a & (x | y)) <=> (a & b | y)"), sub(t, fac, d.EQ3, subst))
}

func TestSubstitutionNary(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	p := parser.New(fac)
	subst := getSubstitution(fac)

	assert.Equal(p.ParseUnsafe("(a & b | y) & (~(a & b) | ~y)"), sub(t, fac, d.AND3, subst))
	assert.Equal(p.ParseUnsafe("(~a & (x | y)) | (a & ~(x | y))"), sub(t, fac, d.OR3, subst))

	assert.Equal(
		p.ParseUnsafe("~(x | y) & c & a & b & ~y"),
		sub(t, fac, fac.And(d.NB, d.C, d.X, d.NY), subst),
	)
	assert.Equal(
		p.ParseUnsafe("~a | ~(x | y) | c | a & b | ~y"),
		sub(t, fac, fac.Or(d.A, d.NB, d.C, d.X, d.NY), subst),
	)
}

func TestSubstitutionPbc(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)

	lits := []f.Literal{fac.Lit("a", true), fac.Lit("b", false), fac.Lit("c", true)}
	litsS1 := []f.Literal{fac.Lit("b", false), fac.Lit("c", true)}
	litsS2 := []f.Literal{fac.Lit("c", true)}
	litsS5 := []f.Literal{fac.Lit("a2", true), fac.Lit("b2", false), fac.Lit("c2", true)}
	litsS6 := []f.Literal{fac.Lit("a2", true), fac.Lit("c", true)}
	coeffs := []int{2, -2, 3}
	coeffS1 := []int{-2, 3}
	coeffS2 := []int{3}
	coeffS6 := []int{2, 3}

	pb := fac.PBC(f.EQ, 2, lits, coeffs)

	s1 := NewSubstitution()
	s1.AddVar(fac.Var("a"), fac.Verum())
	s2 := NewSubstitution()
	s2.AddVar(fac.Var("a"), fac.Verum())
	s2.AddVar(fac.Var("b"), fac.Falsum())
	s3 := NewSubstitution()
	s3.AddVar(fac.Var("a"), fac.Verum())
	s3.AddVar(fac.Var("b"), fac.Falsum())
	s3.AddVar(fac.Var("c"), fac.Verum())
	s4 := NewSubstitution()
	s4.AddVar(fac.Var("a"), fac.Falsum())
	s4.AddVar(fac.Var("b"), fac.Verum())
	s4.AddVar(fac.Var("c"), fac.Falsum())
	s5 := NewSubstitution()
	s5.AddVar(fac.Var("a"), fac.Variable("a2"))
	s5.AddVar(fac.Var("b"), fac.Variable("b2"))
	s5.AddVar(fac.Var("c"), fac.Variable("c2"))
	s5.AddVar(fac.Var("d"), fac.Variable("d2"))
	s6 := NewSubstitution()
	s6.AddVar(fac.Var("a"), fac.Variable("a2"))
	s6.AddVar(fac.Var("b"), fac.Falsum())

	assert.Equal(fac.PBC(f.EQ, 0, litsS1, coeffS1), sub(t, fac, pb, s1))
	assert.Equal(fac.PBC(f.EQ, 2, litsS2, coeffS2), sub(t, fac, pb, s2))
	assert.Equal(fac.Falsum(), sub(t, fac, pb, s3))
	assert.Equal(fac.Falsum(), sub(t, fac, pb, s4))
	assert.Equal(fac.PBC(f.EQ, 2, litsS5, coeffs), sub(t, fac, pb, s5))
	assert.Equal(fac.PBC(f.EQ, 4, litsS6, coeffS6), sub(t, fac, pb, s6))

	assert.Equal(fac.Verum(), sub(t, fac, d.PB2, s3))
	assert.Equal(fac.Verum(), sub(t, fac, d.PB2, s4))
}

func sub(t *testing.T, fac f.Factory, formula f.Formula, subst *Substitution) f.Formula {
	s, err := Substitute(fac, formula, subst)
	assert.Nil(t, err)
	return s
}

func getSubstitution(fac f.Factory) *Substitution {
	d := f.NewTestData(fac)
	subst := NewSubstitution()
	subst.AddVar(d.VA, d.NA)
	subst.AddVar(d.VB, d.OR1)
	subst.AddVar(d.VX, d.AND1)
	return subst
}
