package encoding

import (
	"booleworks.com/logicng/errorx"
	f "booleworks.com/logicng/formula"
)

// ContainsPBC reports whether the given formula contains a pseudo-Boolean
// constraint of any sort.
func ContainsPBC(fac f.Factory, formula f.Formula) bool {
	switch fsort := formula.Sort(); fsort {
	case f.SortFalse, f.SortTrue, f.SortLiteral:
		return false
	case f.SortAnd, f.SortOr:
		ops, _ := fac.NaryOperands(formula)
		for _, op := range ops {
			if ContainsPBC(fac, op) {
				return true
			}
		}
		return false
	case f.SortNot:
		op, _ := fac.NotOperand(formula)
		return ContainsPBC(fac, op)
	case f.SortImpl, f.SortEquiv:
		left, right, _ := fac.BinaryLeftRight(formula)
		return ContainsPBC(fac, left) || ContainsPBC(fac, right)
	case f.SortCC, f.SortPBC:
		return true
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
}
