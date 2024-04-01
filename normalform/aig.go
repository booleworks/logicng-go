package normalform

import (
	"booleworks.com/logicng/errorx"
	f "booleworks.com/logicng/formula"
)

// IsAIG reports whether the given formula is in AIG (and-inverter-graph)
// normal-form, therefore only containing conjunctions and negations.
func IsAIG(fac f.Factory, formula f.Formula) bool {
	cached, ok := f.LookupPredicateCache(fac, f.PredAIG, formula)
	if ok {
		return cached
	}
	var result bool
	switch formula.Sort() {
	case f.SortFalse, f.SortTrue, f.SortLiteral:
		result = true
	case f.SortImpl, f.SortEquiv, f.SortOr, f.SortCC, f.SortPBC:
		result = false
	case f.SortNot:
		op, _ := fac.NotOperand(formula)
		result = IsAIG(fac, op)
	case f.SortAnd:
		result = true
		for _, op := range fac.Operands(formula) {
			if !IsAIG(fac, op) {
				result = false
				break
			}
		}
	default:
		panic(errorx.UnknownEnumValue(formula.Sort()))
	}
	f.SetPredicateCache(fac, f.PredAIG, formula, result)
	return result
}

// AIG returns the given formula as AIG (and-inverter-graph), therefore only
// containing conjunctions and negations.
func AIG(fac f.Factory, formula f.Formula) f.Formula {
	cached, ok := f.LookupTransformationCache(fac, f.TransAIG, formula)
	if ok {
		return cached
	}
	var result f.Formula
	switch formula.Sort() {
	case f.SortFalse, f.SortTrue, f.SortLiteral:
		result = formula
	case f.SortNot:
		op, _ := fac.NotOperand(formula)
		result = fac.Not(AIG(fac, op))
	case f.SortImpl:
		left, right, _ := fac.BinaryLeftRight(formula)
		result = fac.Not(fac.And(AIG(fac, left), fac.Not(AIG(fac, right))))
	case f.SortEquiv:
		left, right, _ := fac.BinaryLeftRight(formula)
		result = fac.And(
			fac.Not(fac.And(AIG(fac, left), fac.Not(AIG(fac, right)))),
			fac.Not(fac.And(fac.Not(left), right)),
		)
	case f.SortAnd:
		ops, _ := fac.NaryOperands(formula)
		nops := make([]f.Formula, len(ops))
		for i, op := range ops {
			nops[i] = AIG(fac, op)
		}
		result = fac.And(nops...)
	case f.SortOr:
		ops, _ := fac.NaryOperands(formula)
		nops := make([]f.Formula, len(ops))
		for i, op := range ops {
			nops[i] = fac.Not(AIG(fac, op))
		}
		result = fac.Not(fac.And(nops...))
	case f.SortCC, f.SortPBC:
		result = AIG(fac, CNF(fac, formula))
	default:
		panic(errorx.UnknownEnumValue(formula.Sort()))
	}
	f.SetTransformationCache(fac, f.TransAIG, formula, result)
	return result
}
