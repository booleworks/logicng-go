package transformation

import (
	"testing"

	"booleworks.com/logicng/errorx"
	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/parser"
	"booleworks.com/logicng/sat"
	"github.com/stretchr/testify/assert"
)

func TestPureExpansion(t *testing.T) {
	fac := f.NewFactory()
	d := f.NewTestData(fac)
	p := parser.New(fac)
	computeAndVerify(t, fac, d.False)
	computeAndVerify(t, fac, d.True)
	computeAndVerify(t, fac, d.A)
	computeAndVerify(t, fac, d.NA)
	computeAndVerify(t, fac, d.NOT1)
	computeAndVerify(t, fac, d.NOT2)
	computeAndVerify(t, fac, d.IMP1)
	computeAndVerify(t, fac, d.IMP2)
	computeAndVerify(t, fac, d.IMP3)
	computeAndVerify(t, fac, d.IMP4)
	computeAndVerify(t, fac, d.EQ1)
	computeAndVerify(t, fac, d.EQ2)
	computeAndVerify(t, fac, d.EQ3)
	computeAndVerify(t, fac, d.EQ4)
	computeAndVerify(t, fac, d.AND1)
	computeAndVerify(t, fac, d.AND2)
	computeAndVerify(t, fac, d.AND3)
	computeAndVerify(t, fac, d.OR1)
	computeAndVerify(t, fac, d.OR2)
	computeAndVerify(t, fac, d.OR3)

	exp, err := ExpandAMOAndEXO(fac, p.ParseUnsafe("a + b <= 1"))
	assert.Nil(t, err)
	assert.Equal(t, p.ParseUnsafe("~a | ~b"), exp)

	exp, err = ExpandAMOAndEXO(fac, p.ParseUnsafe("a + b < 2"))
	assert.Nil(t, err)
	assert.Equal(t, p.ParseUnsafe("~a | ~b"), exp)

	exp, err = ExpandAMOAndEXO(fac, p.ParseUnsafe("a + b = 1"))
	assert.Nil(t, err)
	assert.Equal(t, p.ParseUnsafe("(~a | ~b) & (a | b)"), exp)
}

func computeAndVerify(t *testing.T, fac f.Factory, formula f.Formula) {
	expandedFormula, err := ExpandAMOAndEXO(fac, formula)
	assert.Nil(t, err)
	verify(t, fac, formula, expandedFormula)
}

func verify(t *testing.T, fac f.Factory, formula, expandedFormula f.Formula) {
	assert.True(t, sat.IsEquivalent(fac, formula, expandedFormula))
	assert.True(t, isFreeOfPbcs(fac, expandedFormula))
}

func isFreeOfPbcs(fac f.Factory, formula f.Formula) bool {
	switch formula.Sort() {
	case f.SortTrue, f.SortFalse, f.SortLiteral:
		return true
	case f.SortNot:
		op, _ := fac.NotOperand(formula)
		return isFreeOfPbcs(fac, op)
	case f.SortImpl, f.SortEquiv:
		left, right, _ := fac.BinaryLeftRight(formula)
		return isFreeOfPbcs(fac, left) && isFreeOfPbcs(fac, right)
	case f.SortOr, f.SortAnd:
		ops, _ := fac.NaryOperands(formula)
		for _, op := range ops {
			if !isFreeOfPbcs(fac, op) {
				return false
			}
		}
		return true
	case f.SortCC, f.SortPBC:
		return false
	default:
		panic(errorx.UnknownEnumValue(formula.Sort()))
	}
}
