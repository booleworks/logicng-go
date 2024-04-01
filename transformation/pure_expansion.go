package transformation

import (
	"booleworks.com/logicng/encoding"
	"booleworks.com/logicng/errorx"
	f "booleworks.com/logicng/formula"
)

// ExpandAMOAndEXO expands all at-most-one and exactly-one cardinality
// constraints in the given formula by their pure encoding, meaning without the
// introduction of auxiliary variables.  Since there are no such encodings for
// pseudo-Boolean constraints or arbitrary cardinality constraints, in this
// case the function returns an error.
func ExpandAMOAndEXO(fac f.Factory, formula f.Formula) (f.Formula, error) {
	switch formula.Sort() {
	case f.SortTrue, f.SortFalse, f.SortLiteral:
		return formula, nil
	case f.SortNot:
		op, _ := fac.NotOperand(formula)
		exp, err := ExpandAMOAndEXO(fac, op)
		if err != nil {
			return 0, err
		}
		return fac.Not(exp), nil
	case f.SortImpl, f.SortEquiv:
		left, right, _ := fac.BinaryLeftRight(formula)
		newLeft, err := ExpandAMOAndEXO(fac, left)
		if err != nil {
			return 0, err
		}
		newRight, err := ExpandAMOAndEXO(fac, right)
		if err != nil {
			return 0, err
		}
		binOp, _ := fac.BinaryOperator(formula.Sort(), newLeft, newRight)
		return binOp, nil
	case f.SortOr, f.SortAnd:
		ops, _ := fac.NaryOperands(formula)
		newOps := make([]f.Formula, len(ops))
		var err error
		for i, op := range ops {
			newOps[i], err = ExpandAMOAndEXO(fac, op)
			if err != nil {
				return 0, err
			}
		}
		naryOp, _ := fac.NaryOperator(formula.Sort(), newOps...)
		return naryOp, nil
	case f.SortCC:
		op, rhs, _, _, _ := fac.PBCOps(formula)
		if isValidCC(op, rhs) {
			config := encoding.DefaultConfig()
			config.AMOEncoder = encoding.AMOPure
			cc, _ := encoding.EncodeCC(fac, formula, config)
			return fac.And(cc...), nil
		} else {
			return 0, errorx.BadInput("CC other than AMO or EXO cannot be expanded")
		}
	case f.SortPBC:
		return 0, errorx.BadInput("PBC cannot be expanded")
	default:
		panic(errorx.UnknownEnumValue(formula.Sort()))
	}
}

func isValidCC(op f.CSort, rhs int) bool {
	return op == f.LE && rhs == 1 || op == f.LT && rhs == 2 || op == f.EQ && rhs == 1
}
