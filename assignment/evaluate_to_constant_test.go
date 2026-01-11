package assignment

import (
	"testing"

	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/parser"
	"github.com/booleworks/logicng-go/randomizer"
	"github.com/stretchr/testify/assert"
)

func TestEvaluatesToFalseConstants(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	m := maps(fac)

	assert.True(EvaluatesToFalse(fac, fac.Falsum(), m.empty))
	assert.True(EvaluatesToFalse(fac, fac.Falsum(), m.a))
	assert.True(EvaluatesToFalse(fac, fac.Falsum(), m.aNotB))

	assert.False(EvaluatesToFalse(fac, fac.Verum(), m.empty))
	assert.False(EvaluatesToFalse(fac, fac.Verum(), m.a))
	assert.False(EvaluatesToFalse(fac, fac.Verum(), m.aNotB))
}

func TestEvaluatesToTrueConstants(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	m := maps(fac)

	assert.False(EvaluatesToTrue(fac, fac.Falsum(), m.empty))
	assert.False(EvaluatesToTrue(fac, fac.Falsum(), m.a))
	assert.False(EvaluatesToTrue(fac, fac.Falsum(), m.aNotB))

	assert.True(EvaluatesToTrue(fac, fac.Verum(), m.empty))
	assert.True(EvaluatesToTrue(fac, fac.Verum(), m.a))
	assert.True(EvaluatesToTrue(fac, fac.Verum(), m.aNotB))
}

func TestEvaluatesToFalseLiterals(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	m := maps(fac)

	assert.False(EvaluatesToFalse(fac, d.A, m.empty))
	assert.False(EvaluatesToFalse(fac, d.A, m.a))
	assert.False(EvaluatesToFalse(fac, d.A, m.aNotB))
	assert.False(EvaluatesToFalse(fac, d.NA, m.empty))
	assert.True(EvaluatesToFalse(fac, d.NA, m.a))
	assert.True(EvaluatesToFalse(fac, d.NA, m.aNotB))
	assert.False(EvaluatesToFalse(fac, d.B, m.empty))
	assert.False(EvaluatesToFalse(fac, d.B, m.a))
	assert.True(EvaluatesToFalse(fac, d.B, m.aNotB))
	assert.False(EvaluatesToFalse(fac, d.NB, m.empty))
	assert.False(EvaluatesToFalse(fac, d.NB, m.a))
	assert.False(EvaluatesToFalse(fac, d.NB, m.aNotB))
}

func TestEvaluatesToTrueLiterals(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	m := maps(fac)

	assert.False(EvaluatesToTrue(fac, d.A, m.empty))
	assert.True(EvaluatesToTrue(fac, d.A, m.a))
	assert.True(EvaluatesToTrue(fac, d.A, m.aNotB))
	assert.False(EvaluatesToTrue(fac, d.NA, m.empty))
	assert.False(EvaluatesToTrue(fac, d.NA, m.a))
	assert.False(EvaluatesToTrue(fac, d.NA, m.aNotB))
	assert.False(EvaluatesToTrue(fac, d.B, m.empty))
	assert.False(EvaluatesToTrue(fac, d.B, m.a))
	assert.False(EvaluatesToTrue(fac, d.B, m.aNotB))
	assert.False(EvaluatesToTrue(fac, d.NB, m.empty))
	assert.False(EvaluatesToTrue(fac, d.NB, m.a))
	assert.True(EvaluatesToTrue(fac, d.NB, m.aNotB))
}

func TestEvaluatesToFalseNot(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	m := maps(fac)

	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~~a"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~~a"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~~a"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~~~a"), m.empty))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~~~a"), m.a))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~~~a"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~(a & b)"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~(a & b)"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~(a & b)"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~(~a & b)"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~(~a & b)"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~(~a & b)"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~(a & ~b)"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~(a & ~b)"), m.a))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~(a & ~b)"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~(~a & ~b)"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~(~a & ~b)"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~(~a & ~b)"), m.aNotB))
}

func TestEvaluatesToTrueNot(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	m := maps(fac)

	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~~a"), m.empty))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("~~a"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("~~a"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~~~a"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~~~a"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~~~a"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~(a & b)"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~(a & b)"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("~(a & b)"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~(~a & b)"), m.empty))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("~(~a & b)"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("~(~a & b)"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~(a & ~b)"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~(a & ~b)"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~(a & ~b)"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~(~a & ~b)"), m.empty))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("~(~a & ~b)"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("~(~a & ~b)"), m.aNotB))
}

func TestEvaluatesToFalseOr(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	m := maps(fac)

	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a | b"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a | b"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a | b"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a | b"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a | b"), m.a))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~a | b"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a | ~b"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a | ~b"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a | ~b"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a | ~b"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a | ~b"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a | ~b"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a | ~b | c | ~d"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a | ~b | c | ~d"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a | ~b | c | ~d"), m.aNotB))
}

func TestEvaluatesToTrueOr(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	m := maps(fac)

	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a | b"), m.empty))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("a | b"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("a | b"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a | b"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a | b"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a | b"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a | ~b"), m.empty))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("a | ~b"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("a | ~b"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a | ~b"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a | ~b"), m.empty))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("~a | ~b"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a | ~b | c | ~d"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a | ~b | c | ~d"), m.empty))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("~a | ~b | c | ~d"), m.aNotB))
}

func TestEvaluatesToFalseAnd(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	m := maps(fac)

	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("a & ~a"), m.empty))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("a & ~a"), m.a))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("a & ~a"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & b"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & b"), m.a))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("a & b"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a & b"), m.empty))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~a & b"), m.a))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~a & b"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & ~b"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & ~b"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & ~b"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a & ~b"), m.empty))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~a & ~b"), m.a))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~a & ~b"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a & ~b & c & ~d"), m.empty))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~a & ~b & c & ~d"), m.a))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~a & ~b & c & ~d"), m.aNotB))
}

func TestEvaluatesToTrueAnd(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	m := maps(fac)

	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & ~a"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & ~a"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & ~a"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & b"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & b"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & b"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a & b"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a & b"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a & b"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & ~b"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & ~b"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("a & ~b"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a & ~b"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a & ~b"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a & ~b"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a & ~b & c & ~d"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a & ~b & c & ~d"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a & ~b & c & ~d"), m.aNotB))
}

func TestEvaluatesToFalseImplication(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	m := maps(fac)

	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a => a"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a => a"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a => a"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("b => b"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("b => b"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("b => b"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a => b"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a => b"), m.a))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("a => b"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a => b"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a => b"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a => b"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a => ~b"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a => ~b"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a => ~b"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a => ~b"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a => ~b"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a => ~b"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("b => a"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("b => a"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("b => a"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~b => a"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~b => a"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~b => a"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("b => ~a"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("b => ~a"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("b => ~a"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~b => ~a"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~b => ~a"), m.a))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~b => ~a"), m.aNotB))
}

func TestEvaluatesToTrueImplication(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	m := maps(fac)

	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("a => a"), m.empty))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("a => a"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("a => a"), m.aNotB))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("b => b"), m.empty))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("b => b"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("b => b"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a => b"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a => b"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a => b"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a => b"), m.empty))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("~a => b"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("~a => b"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a => ~b"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a => ~b"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("a => ~b"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a => ~b"), m.empty))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("~a => ~b"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("~a => ~b"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("b => a"), m.empty))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("b => a"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("b => a"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~b => a"), m.empty))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("~b => a"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("~b => a"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("b => ~a"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("b => ~a"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("b => ~a"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~b => ~a"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~b => ~a"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~b => ~a"), m.aNotB))
}

func TestEvaluatesToFalseEquivalence(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	m := maps(fac)

	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a <=> a"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a <=> a"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a <=> a"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("b <=> b"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("b <=> b"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("b <=> b"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a <=> b"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a <=> b"), m.a))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("a <=> b"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a <=> b"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a <=> b"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a <=> b"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a <=> ~b"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a <=> ~b"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a <=> ~b"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a <=> ~b"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a <=> ~b"), m.a))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~a <=> ~b"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("b <=> a"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("b <=> a"), m.a))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("b <=> a"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~b <=> a"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~b <=> a"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~b <=> a"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("b <=> ~a"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("b <=> ~a"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("b <=> ~a"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~b <=> ~a"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~b <=> ~a"), m.a))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~b <=> ~a"), m.aNotB))
}

func TestEvaluatesToTrueEquivalence(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	m := maps(fac)

	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("a <=> a"), m.empty))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("a <=> a"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("a <=> a"), m.aNotB))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("b <=> b"), m.empty))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("b <=> b"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("b <=> b"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a <=> b"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a <=> b"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a <=> b"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a <=> b"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a <=> b"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("~a <=> b"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a <=> ~b"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a <=> ~b"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("a <=> ~b"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a <=> ~b"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a <=> ~b"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a <=> ~b"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("b <=> a"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("b <=> a"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("b <=> a"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~b <=> a"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~b <=> a"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("~b <=> a"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("b <=> ~a"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("b <=> ~a"), m.a))
	assert.True(EvaluatesToTrue(fac, p.ParseUnsafe("b <=> ~a"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~b <=> ~a"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~b <=> ~a"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~b <=> ~a"), m.aNotB))
}

func TestEvaluatesToFalsePbc(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	m := maps(fac)

	pbc01 := fac.PBC(f.EQ, 2, []f.Literal{f.Literal(d.A), f.Literal(d.B)}, []int{2, -4})
	assert.False(EvaluatesToFalse(fac, pbc01, m.empty))
	assert.False(EvaluatesToFalse(fac, pbc01, m.a))
	assert.False(EvaluatesToFalse(fac, pbc01, m.aNotB))

	pbc02 := fac.PBC(f.GT, 2, []f.Literal{f.Literal(d.B), f.Literal(d.C)}, []int{2, 1})
	assert.False(EvaluatesToFalse(fac, pbc02, m.empty))
	assert.False(EvaluatesToFalse(fac, pbc02, m.a))
	assert.True(EvaluatesToFalse(fac, pbc02, m.aNotB))

	assert.False(EvaluatesToFalse(fac, d.PBC1, m.empty))
	assert.False(EvaluatesToFalse(fac, d.PBC1, m.a))
	assert.False(EvaluatesToFalse(fac, d.PBC1, m.aNotB))

	assert.False(EvaluatesToFalse(fac, d.PBC2, m.empty))
	assert.False(EvaluatesToFalse(fac, d.PBC2, m.a))
	assert.False(EvaluatesToFalse(fac, d.PBC2, m.aNotB))
}

func TestEvaluatesToTruePbc(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	m := maps(fac)

	pbc01 := fac.PBC(f.EQ, 2, []f.Literal{f.Literal(d.A), f.Literal(d.B)}, []int{2, -4})
	assert.False(EvaluatesToTrue(fac, pbc01, m.empty))
	assert.False(EvaluatesToTrue(fac, pbc01, m.a))
	assert.True(EvaluatesToTrue(fac, pbc01, m.aNotB))

	pbc02 := fac.PBC(f.GT, 2, []f.Literal{f.Literal(d.B), f.Literal(d.C)}, []int{2, 1})
	assert.False(EvaluatesToTrue(fac, pbc02, m.empty))
	assert.False(EvaluatesToTrue(fac, pbc02, m.a))
	assert.False(EvaluatesToTrue(fac, pbc02, m.aNotB))

	assert.False(EvaluatesToTrue(fac, d.PBC1, m.empty))
	assert.False(EvaluatesToTrue(fac, d.PBC1, m.a))
	assert.False(EvaluatesToTrue(fac, d.PBC1, m.aNotB))

	assert.False(EvaluatesToTrue(fac, d.PBC2, m.empty))
	assert.False(EvaluatesToTrue(fac, d.PBC2, m.a))
	assert.False(EvaluatesToTrue(fac, d.PBC2, m.aNotB))
}

func TestEvaluatesToFalseMixed(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	m := maps(fac)

	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a & (a | ~b)"), m.empty))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~a & (a | ~b)"), m.a))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~a & (a | ~b)"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~b & (b | ~a)"), m.empty))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~b & (b | ~a)"), m.a))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~b & (b | ~a)"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a & (a | ~b) & c & (a => b | e)"), m.empty))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~a & (a | ~b) & c & (a => b | e)"), m.a))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~a & (a | ~b) & c & (a => b | e)"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("~a & ~(a | ~b) & c & (a => b | e)"), m.empty))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~a & ~(a | ~b) & c & (a => b | e)"), m.a))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("~a & ~(a | ~b) & c & (a => b | e)"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & c & (a => b | e)"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & c & (a => b | e)"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & c & (a => b | e)"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & c & (a => ~b | e)"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & c & (a => ~b | e)"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & c & (a => ~b | e)"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & (a => b | e)"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & (a => b | e)"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & (a => b | e)"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & c & (a <=> ~b | e)"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & c & (a <=> ~b | e)"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & c & (a <=> ~b | e)"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & (a <=> b | e)"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & (a <=> b | e)"), m.a))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & (a <=> b | e)"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & (a <=> b)"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & (a <=> b)"), m.a))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & (a <=> b)"), m.aNotB))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & (3 * a + 2 * b > 4)"), m.empty))
	assert.False(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & (3 * a + 2 * b > 4)"), m.a))
	assert.True(EvaluatesToFalse(fac, p.ParseUnsafe("a & (a | ~b) & (3 * a + 2 * b > 4)"), m.aNotB))
}

func TestEvaluatesToTrueMixed(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := parser.New(fac)
	m := maps(fac)

	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a & (a | ~b)"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a & (a | ~b)"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a & (a | ~b)"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~b & (b | ~a)"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~b & (b | ~a)"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~b & (b | ~a)"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a & (a | ~b) & c & (a => b | e)"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a & (a | ~b) & c & (a => b | e)"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a & (a | ~b) & c & (a => b | e)"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a & ~(a | ~b) & c & (a => b | e)"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a & ~(a | ~b) & c & (a => b | e)"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("~a & ~(a | ~b) & c & (a => b | e)"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & c & (a => b | e)"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & c & (a => b | e)"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & c & (a => b | e)"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & c & (a => ~b | e)"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & c & (a => ~b | e)"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & c & (a => ~b | e)"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & (a => b | e)"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & (a => b | e)"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & (a => b | e)"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & c & (a <=> ~b | e)"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & c & (a <=> ~b | e)"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & c & (a <=> ~b | e)"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & (a <=> b | e)"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & (a <=> b | e)"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & (a <=> b | e)"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & (a <=> b)"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & (a <=> b)"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & (a <=> b)"), m.aNotB))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & (3 * a + 2 * b > 4)"), m.empty))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & (3 * a + 2 * b > 4)"), m.a))
	assert.False(EvaluatesToTrue(fac, p.ParseUnsafe("a & (a | ~b) & (3 * a + 2 * b > 4)"), m.aNotB))
}

func TestEvaluatesToConstantRandom(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	assignment, _ := New(fac)
	assignment.AddLit(fac, fac.Lit("v0", false))
	assignment.AddLit(fac, fac.Lit("v1", false))
	assignment.AddLit(fac, fac.Lit("v2", true))
	assignment.AddLit(fac, fac.Lit("v3", true))

	mapping := map[f.Variable]bool{
		fac.Var("v0"): false,
		fac.Var("v1"): false,
		fac.Var("v2"): true,
		fac.Var("v3"): true,
	}

	for i := range 1000 {
		config := randomizer.DefaultConfig()
		config.NumVars = 10
		config.WeightPBC = 1
		config.Seed = int64(i * 42)
		formula := randomizer.New(fac, config).Formula(6)
		restricted := Restrict(fac, formula, assignment)
		assert.Equal(restricted.Sort() == f.SortFalse, EvaluatesToFalse(fac, formula, mapping))
		assert.Equal(restricted.Sort() == f.SortTrue, EvaluatesToTrue(fac, formula, mapping))
	}
}

type mappings struct {
	empty map[f.Variable]bool
	a     map[f.Variable]bool
	aNotB map[f.Variable]bool
}

func maps(fac f.Factory) *mappings {
	return &mappings{
		empty: make(map[f.Variable]bool),
		a:     map[f.Variable]bool{fac.Var("a"): true},
		aNotB: map[f.Variable]bool{fac.Var("a"): true, fac.Var("b"): false},
	}
}
