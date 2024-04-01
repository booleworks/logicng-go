package parser

import (
	"errors"
	"testing"

	"booleworks.com/logicng/errorx"
	f "booleworks.com/logicng/formula"
	"github.com/stretchr/testify/assert"
)

func TestParseConstants(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := New(fac)

	assert.Equal(fac.Verum(), p.ParseUnsafe("$true"))
	assert.Equal(fac.Falsum(), p.ParseUnsafe("$false"))
}

func TestParseLiterals(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := New(fac)

	assert.Equal(p.ParseUnsafe("A"), fac.Variable("A"))
	assert.Equal(p.ParseUnsafe("a"), fac.Variable("a"))
	assert.Equal(p.ParseUnsafe("a1"), fac.Variable("a1"))
	assert.Equal(p.ParseUnsafe("aA_Bb_Cc_12_3"), fac.Variable("aA_Bb_Cc_12_3"))
	assert.Equal(p.ParseUnsafe("~A"), fac.Literal("A", false))
	assert.Equal(p.ParseUnsafe("~a"), fac.Literal("a", false))
	assert.Equal(p.ParseUnsafe("~a1"), fac.Literal("a1", false))
	assert.Equal(p.ParseUnsafe("~aA_Bb_Cc_12_3"), fac.Literal("aA_Bb_Cc_12_3", false))
	assert.Equal(p.ParseUnsafe("~@aA_Bb_Cc_12_3"), fac.Literal("@aA_Bb_Cc_12_3", false))
	assert.Equal(p.ParseUnsafe("#"), fac.Literal("#", true))
	assert.Equal(p.ParseUnsafe("~#"), fac.Literal("#", false))
	assert.Equal(p.ParseUnsafe("~A#B"), fac.Literal("A#B", false))
	assert.Equal(p.ParseUnsafe("A#B"), fac.Literal("A#B", true))
	assert.Equal(p.ParseUnsafe("~A#B"), fac.Literal("A#B", false))
	assert.Equal(p.ParseUnsafe("#A#B_"), fac.Literal("#A#B_", true))
	assert.Equal(p.ParseUnsafe("~#A#B_"), fac.Literal("#A#B_", false))
}

func TestParseOperators(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := New(fac)

	assert.Equal(p.ParseUnsafe("~a"), fac.Not(fac.Variable("a")))
	assert.Equal(p.ParseUnsafe("~Var"), fac.Not(fac.Variable("Var")))
	assert.Equal(p.ParseUnsafe("a & b"), fac.And(fac.Variable("a"), fac.Variable("b")))
	assert.Equal(p.ParseUnsafe("~a & ~b"), fac.And(fac.Literal("a", false), fac.Literal("b", false)))
	assert.Equal(p.ParseUnsafe("~a & b & ~c & d"), fac.And(fac.Literal("a", false), fac.Variable("b"), fac.Literal("c", false), fac.Variable("d")))
	assert.Equal(p.ParseUnsafe("a | b"), fac.Or(fac.Variable("a"), fac.Variable("b")))
	assert.Equal(p.ParseUnsafe("~a | ~b"), fac.Or(fac.Literal("a", false), fac.Literal("b", false)))
	assert.Equal(p.ParseUnsafe("~a | b | ~c | d"), fac.Or(fac.Literal("a", false), fac.Variable("b"), fac.Literal("c", false), fac.Variable("d")))
	assert.Equal(p.ParseUnsafe("a => b"), fac.Implication(fac.Variable("a"), fac.Variable("b")))
	assert.Equal(p.ParseUnsafe("~a => ~b"), fac.Implication(fac.Literal("a", false), fac.Literal("b", false)))
	assert.Equal(p.ParseUnsafe("a <=> b"), fac.Equivalence(fac.Variable("a"), fac.Variable("b")))
	assert.Equal(p.ParseUnsafe("~a <=> ~b"), fac.Equivalence(fac.Literal("a", false), fac.Literal("b", false)))
}

func TestParsePrecedences(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := New(fac)

	assert.Equal(p.ParseUnsafe("x | y & z"), fac.Or(fac.Variable("x"), fac.And(fac.Variable("y"), fac.Variable("z"))))
	assert.Equal(p.ParseUnsafe("x & y | z"), fac.Or(fac.And(fac.Variable("x"), fac.Variable("y")), fac.Variable("z")))
	assert.Equal(p.ParseUnsafe("x => y & z"), fac.Implication(fac.Variable("x"), fac.And(fac.Variable("y"), fac.Variable("z"))))
	assert.Equal(p.ParseUnsafe("x & y => z"), fac.Implication(fac.And(fac.Variable("x"), fac.Variable("y")), fac.Variable("z")))
	assert.Equal(p.ParseUnsafe("x <=> y & z"), fac.Equivalence(fac.Variable("x"), fac.And(fac.Variable("y"), fac.Variable("z"))))
	assert.Equal(p.ParseUnsafe("x & y <=> z"), fac.Equivalence(fac.And(fac.Variable("x"), fac.Variable("y")), fac.Variable("z")))
	assert.Equal(p.ParseUnsafe("x => y | z"), fac.Implication(fac.Variable("x"), fac.Or(fac.Variable("y"), fac.Variable("z"))))
	assert.Equal(p.ParseUnsafe("x | y => z"), fac.Implication(fac.Or(fac.Variable("x"), fac.Variable("y")), fac.Variable("z")))
	assert.Equal(p.ParseUnsafe("x <=> y | z"), fac.Equivalence(fac.Variable("x"), fac.Or(fac.Variable("y"), fac.Variable("z"))))
	assert.Equal(p.ParseUnsafe("x | y <=> z"), fac.Equivalence(fac.Or(fac.Variable("x"), fac.Variable("y")), fac.Variable("z")))
	assert.Equal(p.ParseUnsafe("x => y => z"), fac.Implication(fac.Variable("x"), fac.Implication(fac.Variable("y"), fac.Variable("z"))))
	assert.Equal(p.ParseUnsafe("x <=> y <=> z"), fac.Equivalence(fac.Variable("x"), fac.Equivalence(fac.Variable("y"), fac.Variable("z"))))
	assert.Equal(p.ParseUnsafe("(x | y) & z"), fac.And(fac.Or(fac.Variable("x"), fac.Variable("y")), fac.Variable("z")))
	assert.Equal(p.ParseUnsafe("x & (y | z)"), fac.And(fac.Variable("x"), fac.Or(fac.Variable("y"), fac.Variable("z"))))
	assert.Equal(p.ParseUnsafe("(x => y) & z"), fac.And(fac.Implication(fac.Variable("x"), fac.Variable("y")), fac.Variable("z")))
	assert.Equal(p.ParseUnsafe("x & (y => z)"), fac.And(fac.Variable("x"), fac.Implication(fac.Variable("y"), fac.Variable("z"))))
	assert.Equal(p.ParseUnsafe("(x => y) | z"), fac.Or(fac.Implication(fac.Variable("x"), fac.Variable("y")), fac.Variable("z")))
	assert.Equal(p.ParseUnsafe("x | (y => z)"), fac.Or(fac.Variable("x"), fac.Implication(fac.Variable("y"), fac.Variable("z"))))
	assert.Equal(p.ParseUnsafe("(x <=> y) & z"), fac.And(fac.Equivalence(fac.Variable("x"), fac.Variable("y")), fac.Variable("z")))
	assert.Equal(p.ParseUnsafe("x & (y <=> z)"), fac.And(fac.Variable("x"), fac.Equivalence(fac.Variable("y"), fac.Variable("z"))))
	assert.Equal(p.ParseUnsafe("(x <=> y) | z"), fac.Or(fac.Equivalence(fac.Variable("x"), fac.Variable("y")), fac.Variable("z")))
	assert.Equal(p.ParseUnsafe("x | (y <=> z)"), fac.Or(fac.Variable("x"), fac.Equivalence(fac.Variable("y"), fac.Variable("z"))))
	assert.Equal(p.ParseUnsafe("x => y <=> z"), fac.Equivalence(fac.Implication(fac.Variable("x"), fac.Variable("y")), fac.Variable("z")))
	assert.Equal(p.ParseUnsafe("x => (y <=> z)"), fac.Implication(fac.Variable("x"), fac.Equivalence(fac.Variable("y"), fac.Variable("z"))))
}

func TestParseEmptyString(t *testing.T) {
	fac := f.NewFactory()
	p := New(fac)

	assert.Equal(t, p.ParseUnsafe(""), fac.Verum())
}

func TestParseMul(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := New(fac)
	abc := fac.Lit("abc", true)
	nabc := fac.Lit("abc", false)

	assert.Equal(fac.PBC(f.EQ, 4, []f.Literal{abc}, []int{13}), p.ParseUnsafe("13*abc = 4"))
	assert.Equal(fac.PBC(f.EQ, 4, []f.Literal{nabc}, []int{-13}), p.ParseUnsafe("-13*~abc = 4"))
	assert.Equal(fac.PBC(f.EQ, -442, []f.Literal{nabc}, []int{13}), p.ParseUnsafe("13 * ~abc = -442"))
	assert.Equal(fac.PBC(f.EQ, -442, []f.Literal{nabc}, []int{-13}), p.ParseUnsafe("-13 * ~abc = -442"))

	assert.Equal(fac.PBC(f.GT, 4, []f.Literal{abc}, []int{13}), p.ParseUnsafe("13 * abc > 4"))
	assert.Equal(fac.PBC(f.GE, 4, []f.Literal{abc}, []int{13}), p.ParseUnsafe("13 * abc >= 4"))
	assert.Equal(fac.PBC(f.LT, 4, []f.Literal{abc}, []int{13}), p.ParseUnsafe("13 * abc < 4"))
	assert.Equal(fac.PBC(f.LE, 4, []f.Literal{abc}, []int{13}), p.ParseUnsafe("13 * abc <= 4"))
}

func TestParseAdd(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := New(fac)
	a := fac.Lit("a", true)
	c := fac.Lit("c", true)
	d := fac.Lit("d", true)
	va := fac.Var("a")
	vc := fac.Var("c")
	vd := fac.Var("d")
	nb := fac.Lit("b", false)
	nc := fac.Lit("c", false)
	nd := fac.Lit("d", false)

	assert.Equal(fac.PBC(f.LT, -4, []f.Literal{c, nd}, []int{4, -4}), p.ParseUnsafe("4 * c + -4 * ~d < -4"))
	assert.Equal(fac.PBC(f.GE, -5, []f.Literal{c, nc}, []int{5, -5}), p.ParseUnsafe("5 * c + -5 * ~c >= -5"))
	assert.Equal(fac.PBC(f.GT, -6, []f.Literal{a, nb, nc}, []int{6, -6, 12}), p.ParseUnsafe("6 * a + -6 * ~b + 12 * ~c > -6"))
	assert.Equal(fac.PBC(f.LT, -4, []f.Literal{c, nd}, []int{1, -4}), p.ParseUnsafe("c + -4 * ~d < -4"))
	assert.Equal(fac.PBC(f.GE, -5, []f.Literal{c, nc}, []int{5, 1}), p.ParseUnsafe("5 * c + ~c >= -5"))
	assert.Equal(fac.PBC(f.GE, -5, []f.Literal{c, d}, []int{1, 1}), p.ParseUnsafe("c + d >= -5"))

	assert.Equal(fac.AMO(vc, vd), p.ParseUnsafe("c + d <= 1"))
	assert.Equal(fac.EXO(va, vc, vd), p.ParseUnsafe("a + c + d = 1"))
	assert.Equal(fac.CC(f.GT, 2, va, vc, vd), p.ParseUnsafe("a + c + d > 2"))
	assert.Equal(fac.PBC(f.GE, -5, []f.Literal{nc, nd}, []int{1, 1}), p.ParseUnsafe("~c + ~d >= -5"))
	assert.Equal(fac.PBC(f.EQ, -5, []f.Literal{nc}, []int{1}), p.ParseUnsafe("~c = -5"))
	assert.Equal(fac.Not(fac.PBC(f.EQ, -5, []f.Literal{c}, []int{1})), p.ParseUnsafe("~(c = -5)"))
}

func TestParseNumericalLiteral(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := New(fac)

	assert.Equal(fac.Variable("12"), p.ParseUnsafe("12"))
	assert.Equal(fac.And(fac.Literal("12", false), fac.Variable("A")), p.ParseUnsafe("~12 & A"))
	assert.Equal(
		fac.PBC(f.LE, 25, []f.Literal{fac.Lit("12", true), fac.Lit("A", true), fac.Lit("B", true)}, []int{12, 13, 10}),
		p.ParseUnsafe("12 * 12 + 13 * A + 10 * B <= 25"),
	)
	assert.Equal(
		fac.PBC(f.LE, 25, []f.Literal{fac.Lit("12", false), fac.Lit("A", true), fac.Lit("B", true)}, []int{-12, 13, 10}),
		p.ParseUnsafe("-12 * ~12 + 13 * A + 10 * B <= 25"),
	)
}

func TestParseCombinations(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := New(fac)
	pbc := fac.PBC(f.GT, -6, []f.Literal{fac.Lit("a", true), fac.Lit("b", false), fac.Lit("c", false)}, []int{6, 7, -12})

	assert.Equal(fac.Not(pbc), p.ParseUnsafe("~(6 * a + 7 * ~b + -12 * ~c > -6)"))
	assert.Equal(
		fac.And(fac.Implication(fac.Variable("x"), fac.And(fac.Variable("y"), fac.Variable("z"))), pbc),
		p.ParseUnsafe("(x => y & z) & (6 * a + 7 * ~b + -12 * ~c > -6)"),
	)
}

func TestParseSafe(t *testing.T) {
	assert := assert.New(t)
	fac := f.NewFactory()
	p := New(fac)

	parsed, err := p.Parse("x^")
	assert.NotNil(err)
	assert.Equal(fac.Falsum(), parsed)
	assert.True(errors.Is(err, errorx.ErrBadInput))
	assert.Equal("bad input: Syntax error at line 1, column 1: token recognition error at: '^'\n", err.Error())

	parsed, err = p.Parse("A &")
	assert.NotNil(err)
	assert.Equal(fac.Falsum(), parsed)
	assert.True(errors.Is(err, errorx.ErrBadInput))
	assert.Equal("bad input: Syntax error at line 1, column 3: mismatched input '<EOF>' expecting {NUMBER, LITERAL, '$true', '$false', '(', '~'}\n", err.Error())

	parsed, err = p.Parse("(A & B")
	assert.NotNil(err)
	assert.Equal(fac.Falsum(), parsed)
	assert.True(errors.Is(err, errorx.ErrBadInput))
	assert.Equal("bad input: Syntax error at line 1, column 6: missing ')' at '<EOF>'\n", err.Error())
}
