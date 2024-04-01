package transformation

import (
	"booleworks.com/logicng/errorx"
	f "booleworks.com/logicng/formula"
)

// SubstituteLiterals performs a special substitution from literal to literal
// on the given formula.  In contrast to the standard Substitute function
// which can only map variables, this function can also map literals.
//
// Always the best fit is chosen. So if there are two mappings for e.g. a -> b
// and ~a -> c. Then ~a will be mapped to c and not to ~b. On the other hand if
// there is only the mapping a -> b, the literal ~a will be mapped to ~b.
func SubstituteLiterals(fac f.Factory, formula f.Formula, substitution *map[f.Literal]f.Literal) f.Formula {
	switch fsort := formula.Sort(); fsort {
	case f.SortTrue, f.SortFalse:
		return formula
	case f.SortLiteral:
		lit, ok := (*substitution)[f.Literal(formula)]
		if ok {
			return lit.AsFormula()
		}
		if formula.IsNeg() {
			variable := f.Literal(formula).Variable()
			lit, ok = (*substitution)[variable.AsLiteral()]
			if ok {
				return lit.Negate(fac).AsFormula()
			}
		}
		return formula
	case f.SortNot:
		op, _ := fac.NotOperand(formula)
		return fac.Not(SubstituteLiterals(fac, op, substitution))
	case f.SortEquiv, f.SortImpl:
		left, right, _ := fac.BinaryLeftRight(formula)
		binOp, _ := fac.BinaryOperator(
			fsort,
			SubstituteLiterals(fac, left, substitution),
			SubstituteLiterals(fac, right, substitution),
		)
		return binOp
	case f.SortOr, f.SortAnd:
		ops, _ := fac.NaryOperands(formula)
		operands := make([]f.Formula, len(ops))
		for i, op := range ops {
			operands[i] = SubstituteLiterals(fac, op, substitution)
		}
		naryOp, _ := fac.NaryOperator(fsort, operands...)
		return naryOp
	case f.SortCC, f.SortPBC:
		csort, rhs, lits, coeffs, _ := fac.PBCOps(formula)
		literals := make([]f.Literal, len(lits))
		for i, originalOp := range lits {
			literals[i] = f.Literal(SubstituteLiterals(fac, originalOp.AsFormula(), substitution))
		}
		return fac.PBC(csort, rhs, literals, coeffs)
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
}
