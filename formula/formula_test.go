package formula

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFalsum(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()

	falsum := fac.Falsum()
	c := fac.Constant(false)

	assert.Equal(SortFalse, falsum.Sort())
	assert.Equal(uint32(0), falsum.ID())
	assert.Equal(falsum, c)
}

func TestVerum(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()

	verum := fac.Verum()
	c := fac.Constant(true)

	assert.Equal(SortTrue, verum.Sort())
	assert.Equal(uint32(1), verum.ID())
	assert.Equal(verum, c)
}

func TestVariable(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()

	varA := fac.Variable("a")
	name, phase, _ := fac.LiteralNamePhase(varA)
	assert.Equal(SortLiteral, varA.Sort())
	assert.Equal(uint32(3), varA.ID())
	assert.Equal("a", name)
	assert.Equal(true, phase)
	assert.Equal(1, len(fac.(*CachingFactory).posLitCache))

	varB := fac.Variable("b")
	name, phase, _ = fac.LiteralNamePhase(varB)
	assert.Equal(SortLiteral, varB.Sort())
	assert.Equal(uint32(5), varB.ID())
	assert.Equal("b", name)
	assert.Equal(true, phase)
	assert.Equal(2, len(fac.(*CachingFactory).posLitCache))

	varA = fac.Variable("a")
	name, phase, _ = fac.LiteralNamePhase(varA)
	assert.Equal(SortLiteral, varA.Sort())
	assert.Equal(uint32(3), varA.ID())
	assert.Equal("a", name)
	assert.Equal(true, phase)
	assert.Equal(2, len(fac.(*CachingFactory).posLitCache))
}

func TestLiteral(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()

	litA := fac.Literal("a", true)
	name, phase, _ := fac.LiteralNamePhase(litA)
	assert.Equal(SortLiteral, litA.Sort())
	assert.Equal(uint32(3), litA.ID())
	assert.Equal("a", name)
	assert.Equal(true, phase)
	assert.Equal(1, len(fac.(*CachingFactory).posLitCache))
	assert.Equal(0, len(fac.(*CachingFactory).negLitCache))

	litA = fac.Literal("a", false)
	name, phase, _ = fac.LiteralNamePhase(litA)
	assert.Equal(SortLiteral, litA.Sort())
	assert.Equal(uint32(2), litA.ID())
	assert.Equal("a", name)
	assert.Equal(false, phase)
	assert.Equal(1, len(fac.(*CachingFactory).posLitCache))
	assert.Equal(1, len(fac.(*CachingFactory).negLitCache))

	litA = fac.Literal("a", false)
	name, phase, _ = fac.LiteralNamePhase(litA)
	assert.Equal(SortLiteral, litA.Sort())
	assert.Equal(uint32(2), litA.ID())
	assert.Equal("a", name)
	assert.Equal(false, phase)
	assert.Equal(1, len(fac.(*CachingFactory).posLitCache))
	assert.Equal(1, len(fac.(*CachingFactory).negLitCache))
}

func TestNot(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()

	assert.Equal(fac.Falsum(), fac.Not(fac.Verum()))
	assert.Equal(fac.Verum(), fac.Not(fac.Falsum()))
	assert.Equal(fac.Literal("a", false), fac.Not(fac.Variable("a")))
	assert.Equal(fac.Literal("a", true), fac.Not(fac.Literal("a", false)))
	assert.Equal(fac.Literal("a", false), fac.Not(fac.Not(fac.Not(fac.Variable("a")))))

	impl := fac.Implication(fac.Variable("a"), fac.Literal("b", false))
	not := fac.Not(impl)
	op, _ := fac.NotOperand(not)

	assert.Equal(SortNot, not.Sort())
	assert.Equal(impl.ID()^1, not.ID())
	assert.Equal(impl, op)

	not = fac.Not(impl)
	assert.Equal(impl.ID()^1, not.ID())
	assert.Equal(1, len(fac.(*CachingFactory).notCache))
}

func TestImplication(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()
	a := fac.Variable("a")
	na := fac.Literal("a", false)
	equiv := fac.Equivalence(a, fac.Variable("b"))

	assert.Equal(a, fac.Implication(fac.Verum(), a))
	assert.Equal(fac.Verum(), fac.Implication(fac.Falsum(), a))
	assert.Equal(fac.Verum(), fac.Implication(a, fac.Verum()))
	assert.Equal(na, fac.Implication(a, fac.Falsum()))
	assert.Equal(fac.Verum(), fac.Implication(a, a))
	assert.Equal(na, fac.Implication(a, na))
	assert.Equal(a, fac.Implication(na, a))
	assert.Equal(fac.Not(equiv), fac.Implication(equiv, fac.Not(equiv)))
	assert.Equal(equiv, fac.Implication(fac.Not(equiv), equiv))

	impl := fac.Implication(a, fac.Literal("b", false))
	left, right, _ := fac.BinaryLeftRight(impl)
	assert.Equal(SortImpl, impl.Sort())
	assert.Equal(uint32(9), impl.ID())
	assert.Equal(a, left)
	assert.Equal(fac.Literal("b", false), right)
	assert.Equal(1, len(fac.(*CachingFactory).implCache))

	impl = fac.Implication(fac.Literal("b", false), a)
	left, right, _ = fac.BinaryLeftRight(impl)
	assert.Equal(SortImpl, impl.Sort())
	assert.Equal(uint32(11), impl.ID())
	assert.Equal(fac.Literal("b", false), left)
	assert.Equal(a, right)
	assert.Equal(2, len(fac.(*CachingFactory).implCache))

	impl = fac.Implication(a, fac.Literal("b", false))
	left, right, _ = fac.BinaryLeftRight(impl)
	assert.Equal(SortImpl, impl.Sort())
	assert.Equal(uint32(9), impl.ID())
	assert.Equal(a, left)
	assert.Equal(fac.Literal("b", false), right)
	assert.Equal(2, len(fac.(*CachingFactory).implCache))
}

func TestEquivalence(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()
	a := fac.Variable("a")
	na := fac.Literal("a", false)
	impl := fac.Implication(a, fac.Variable("b"))

	assert.Equal(a, fac.Equivalence(fac.Verum(), a))
	assert.Equal(na, fac.Equivalence(fac.Falsum(), a))
	assert.Equal(a, fac.Equivalence(a, fac.Verum()))
	assert.Equal(na, fac.Equivalence(a, fac.Falsum()))
	assert.Equal(fac.Verum(), fac.Equivalence(a, a))
	assert.Equal(fac.Falsum(), fac.Equivalence(a, na))
	assert.Equal(fac.Falsum(), fac.Equivalence(na, a))
	assert.Equal(fac.Falsum(), fac.Equivalence(impl, fac.Not(impl)))
	assert.Equal(fac.Falsum(), fac.Equivalence(fac.Not(impl), impl))

	equiv := fac.Equivalence(a, fac.Literal("b", false))
	left, right, _ := fac.BinaryLeftRight(equiv)
	assert.Equal(SortEquiv, equiv.Sort())
	assert.Equal(uint32(9), equiv.ID())
	assert.Equal(a, left)
	assert.Equal(fac.Literal("b", false), right)
	assert.Equal(1, len(fac.(*CachingFactory).equivCache))

	equiv = fac.Equivalence(fac.Literal("b", false), a)
	left, right, _ = fac.BinaryLeftRight(equiv)
	assert.Equal(SortEquiv, equiv.Sort())
	assert.Equal(uint32(11), equiv.ID())
	assert.Equal(fac.Literal("b", false), left)
	assert.Equal(a, right)
	assert.Equal(2, len(fac.(*CachingFactory).equivCache))
}

func TestBinaryOperand(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()
	a := fac.Variable("a")
	b := fac.Variable("b")

	imp, err := fac.BinaryOperator(SortImpl, a, b)
	assert.Nil(err)
	eq, err := fac.BinaryOperator(SortEquiv, a, b)
	assert.Nil(err)
	assert.Equal(fac.Implication(a, b), imp)
	assert.Equal(fac.Equivalence(a, b), eq)
}

func TestAnd(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()
	a := fac.Variable("a")
	b := fac.Variable("b")
	c := fac.Variable("c")
	na := fac.Literal("a", false)
	and1 := fac.And(a, b, c)
	impl := fac.Implication(a, b)

	assert.Equal(fac.Verum(), fac.And())
	assert.Equal(fac.Falsum(), fac.And(a, b, fac.Falsum()))
	assert.Equal(a, fac.And(a))
	assert.Equal(fac.Falsum(), fac.And(a, na))
	assert.Equal(fac.Falsum(), fac.And(na, a))
	assert.Equal(fac.Falsum(), fac.And(impl, fac.Not(impl)))
	assert.Equal(fac.Falsum(), fac.And(fac.Not(impl), impl))
	assert.Equal(fac.Falsum(), fac.And(a, fac.And(a, b, fac.Falsum()), fac.And(c)))
	assert.NotEqual(and1, fac.And(b, a, c))
	assert.NotEqual(and1, fac.And(c, b, a))
	assert.Equal(and1, fac.And(a, a, b, b, c))
	assert.Equal(and1, fac.And(a, a, b, b, c, c, b, a))
	assert.Equal(3, len(fac.(*CachingFactory).andCache))

	assert.Equal(and1, fac.And(a, fac.And(a, b, b), fac.And(c)))
	assert.Equal(and1, fac.And(a, fac.And(b, b), fac.And(b, fac.Verum()), c, c, fac.And(b, a, fac.Verum(), fac.Verum())))
	assert.Equal(5, len(fac.(*CachingFactory).andCache))

	ops, _ := fac.NaryOperands(and1)
	assert.Equal(a, ops[0])
	assert.Equal(b, ops[1])
	assert.Equal(c, ops[2])
}

func TestOr(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()
	a := fac.Variable("a")
	b := fac.Variable("b")
	c := fac.Variable("c")
	na := fac.Literal("a", false)
	or1 := fac.Or(a, b, c)
	impl := fac.Implication(a, b)

	assert.Equal(fac.Falsum(), fac.Or())
	assert.Equal(fac.Verum(), fac.Or(a, b, fac.Verum()))
	assert.Equal(a, fac.Or(a, fac.Falsum()))
	assert.Equal(a, fac.Or(a))
	assert.Equal(fac.Verum(), fac.Or(a, na))
	assert.Equal(fac.Verum(), fac.Or(na, a))
	assert.Equal(fac.Verum(), fac.Or(impl, fac.Not(impl)))
	assert.Equal(fac.Verum(), fac.Or(fac.Not(impl), impl))
	assert.Equal(fac.Falsum(), fac.And(impl, fac.Not(impl)))
	assert.Equal(fac.Or(a, fac.Or(c), fac.Or(a, b, fac.Verum())), fac.Verum())
	assert.NotEqual(or1, fac.Or(b, a, c))
	assert.NotEqual(or1, fac.Or(c, b, a))
	assert.Equal(or1, fac.Or(a, a, b, b, c))
	assert.Equal(or1, fac.Or(a, a, b, c, c, b, a))
	assert.Equal(3, len(fac.(*CachingFactory).orCache))

	assert.Equal(or1, fac.Or(a, fac.Or(a, b, b), fac.Or(c)))
	assert.Equal(or1, fac.Or(a, fac.Or(b, b), fac.Or(b, fac.Falsum()), c, c, fac.Or(b, a, fac.Falsum(), fac.Falsum())))
	assert.Equal(5, len(fac.(*CachingFactory).orCache))

	ops, _ := fac.NaryOperands(or1)
	assert.Equal(a, ops[0])
	assert.Equal(b, ops[1])
	assert.Equal(c, ops[2])
}

func TestNaryOperand(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()
	a := fac.Variable("a")
	b := fac.Variable("b")
	c := fac.Variable("c")

	naryOp, err := fac.NaryOperator(SortAnd, a, b, c)
	assert.Nil(err)
	assert.Equal(fac.And(a, b, c), naryOp)
	naryOp, err = fac.NaryOperator(SortOr, a, b, c)
	assert.Nil(err)
	assert.Equal(fac.Or(a, b, c), naryOp)
}

func TestCcs(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()
	a := fac.Var("a")
	b := fac.Var("b")
	c := fac.Var("c")

	assert.Equal(fac.CC(LE, 1, a, b, c), fac.AMO(a, b, c))
	assert.Equal(fac.CC(EQ, 1, a, b, c), fac.EXO(a, b, c))
	assert.Equal(fac.CC(GE, 2, a, b, c), fac.PBC(GE, 2, []Literal{a.AsLiteral(), b.AsLiteral(), c.AsLiteral()}, []int{1, 1, 1}))
	assert.Equal(fac.Falsum(), fac.CC(GE, 1))
	assert.Equal(fac.Verum(), fac.CC(LE, 4))

	assert.Equal(3, len(fac.(*CachingFactory).ccCache))

	cc1 := fac.CC(LE, 1, a, b, c)
	cc2 := fac.CC(GE, 2, a, b, c)
	pbc1 := fac.(*CachingFactory).getPBCUnsafe(cc1)
	pbc2 := fac.(*CachingFactory).getPBCUnsafe(cc2)
	assert.Equal(3, len(fac.(*CachingFactory).ccCache))

	assert.Equal(LE, pbc1.comparator)
	assert.Equal(1, pbc1.rhs)
	assert.Equal([]Literal{a.AsLiteral(), b.AsLiteral(), c.AsLiteral()}, pbc1.literals)
	assert.Equal([]int{1, 1, 1}, pbc1.coefficients)

	assert.Equal(GE, pbc2.comparator)
	assert.Equal(2, pbc2.rhs)
	assert.Equal([]Literal{a.AsLiteral(), b.AsLiteral(), c.AsLiteral()}, pbc2.literals)
	assert.Equal([]int{1, 1, 1}, pbc2.coefficients)
}

func TestPbcs(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()
	a := fac.Var("a").AsLiteral()
	nb := fac.Lit("b", false)
	c := fac.Var("c").AsLiteral()

	assert.Equal(fac.Falsum(), fac.PBC(GT, 1, []Literal{}, []int{}))
	assert.Equal(fac.Verum(), fac.PBC(LT, 4, []Literal{}, []int{}))

	assert.Equal(0, len(fac.(*CachingFactory).pbcCache))

	lits := []Literal{a, nb, c}
	coeffs1 := []int{3, -2, 1}
	coeffs2 := []int{10, 11, -12}

	pb1 := fac.PBC(LE, 7, lits, coeffs1)
	pb2 := fac.PBC(GE, 12, lits, coeffs2)
	fac.PBC(LE, 7, lits, coeffs1)
	pbc1 := fac.(*CachingFactory).getPBCUnsafe(pb1)
	pbc2 := fac.(*CachingFactory).getPBCUnsafe(pb2)
	assert.Equal(2, len(fac.(*CachingFactory).pbcCache))

	assert.Equal(LE, pbc1.comparator)
	assert.Equal(7, pbc1.rhs)
	assert.Equal([]Literal{a, nb, c}, pbc1.literals)
	assert.Equal([]int{3, -2, 1}, pbc1.coefficients)

	assert.Equal(GE, pbc2.comparator)
	assert.Equal(12, pbc2.rhs)
	assert.Equal([]Literal{a, nb, c}, pbc2.literals)
	assert.Equal([]int{10, 11, -12}, pbc2.coefficients)
}

func TestVariableFromLit(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory()
	a := fac.Var("a")
	na := fac.Lit("a", false)
	exp := na.Variable()
	assert.Equal(a, exp)
}

func TestFactoryWithConserveVars(t *testing.T) {
	assert := assert.New(t)
	fac := NewFactory(true)
	facSimp := NewFactory(false)
	a := fac.Variable("a")
	na := fac.Literal("a", false)

	assert.Equal("a => a", fac.Implication(a, a).Sprint(fac))
	assert.Equal("~a", fac.Implication(a, na).Sprint(fac))
	assert.Equal("a", fac.Implication(na, a).Sprint(fac))
	assert.Equal("a <=> a", fac.Equivalence(a, a).Sprint(fac))
	assert.Equal("a <=> ~a", fac.Equivalence(a, na).Sprint(fac))
	assert.Equal("~a <=> a", fac.Equivalence(na, a).Sprint(fac))
	assert.Equal("~a & a", fac.And(na, a).Sprint(fac))
	assert.Equal("~a | a", fac.Or(na, a).Sprint(fac))
	assert.Equal("a & ~a", fac.And(a, na).Sprint(fac))
	assert.Equal("a | ~a", fac.Or(a, na).Sprint(fac))

	a = facSimp.Variable("a")
	na = facSimp.Literal("a", false)
	assert.Equal("$true", facSimp.Implication(a, a).Sprint(facSimp))
	assert.Equal("~a", facSimp.Implication(a, na).Sprint(facSimp))
	assert.Equal("a", facSimp.Implication(na, a).Sprint(facSimp))
	assert.Equal("$true", facSimp.Equivalence(a, a).Sprint(facSimp))
	assert.Equal("$false", facSimp.Equivalence(a, na).Sprint(facSimp))
	assert.Equal("$false", facSimp.Equivalence(na, a).Sprint(facSimp))
	assert.Equal("$false", facSimp.And(na, a).Sprint(facSimp))
	assert.Equal("$true", facSimp.Or(na, a).Sprint(facSimp))
	assert.Equal("$false", facSimp.And(a, na).Sprint(facSimp))
	assert.Equal("$true", facSimp.Or(a, na).Sprint(facSimp))
}
