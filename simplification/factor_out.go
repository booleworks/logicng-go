package simplification

import (
	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
)

// FactorOut simplifies a formula by applying factor out operations.  For
// example, given the formula A & B & C | A & D, both conjunction terms have
// the common factor A. Thus, the method returns A & (B & C | D).  The optional
// ratingFunction can be used to rate functions between simplification steps
// and choose the one with the lower rating.
func FactorOut(fac f.Factory, formula f.Formula, ratingFunction ...RatingFunction) f.Formula {
	rf := DefaultRatingFunction
	if len(ratingFunction) > 0 {
		rf = ratingFunction[0]
	}
	var last f.Formula
	simplified := formula
	for condition := true; condition; condition = simplified != last {
		last = simplified
		simplified = factorOutRec(fac, last, rf)
	}
	return simplified
}

func factorOutRec(fac f.Factory, formula f.Formula, rf RatingFunction) f.Formula {
	switch formula.Sort() {
	case f.SortOr, f.SortAnd:
		ops, _ := fac.NaryOperands(formula)
		newOps := make([]f.Formula, len(ops))
		for i, op := range ops {
			newOps[i] = FactorOut(fac, op, rf)
		}
		newFormula, _ := fac.NaryOperator(formula.Sort(), newOps...)
		if formula.Sort() == f.SortAnd || formula.Sort() == f.SortOr {
			return factorOutSimplify(fac, newFormula, rf)
		} else {
			return newFormula
		}
	case f.SortNot:
		op, _ := fac.NotOperand(formula)
		return FactorOut(fac, op, rf).Negate(fac)
	case f.SortTrue, f.SortFalse, f.SortLiteral, f.SortImpl, f.SortEquiv, f.SortCC, f.SortPBC:
		return formula
	default:
		panic(errorx.UnknownEnumValue(formula.Sort()))
	}
}

func factorOutSimplify(fac f.Factory, formula f.Formula, rf RatingFunction) f.Formula {
	simplified, ok := factorOut(fac, formula)
	if !ok {
		return formula
	}
	if rf(fac, formula) < rf(fac, simplified) {
		return formula
	} else {
		return simplified
	}
}

func factorOut(fac f.Factory, formula f.Formula) (f.Formula, bool) {
	factorOutFormula, ok := computeMaxOccurringSubformula(fac, formula)
	if !ok {
		return 0, false
	}
	sort := formula.Sort()
	var formulasWithRemoved, unchangedFormulas []f.Formula
	for _, operand := range fac.Operands(formula) {
		if operand.Sort() == f.SortLiteral {
			if operand == factorOutFormula {
				formulasWithRemoved = append(formulasWithRemoved, fac.Constant(sort == f.SortOr))
			} else {
				unchangedFormulas = append(unchangedFormulas, operand)
			}
		} else if operand.Sort() == f.SortAnd || operand.Sort() == f.SortOr {
			removed := false
			ops, _ := fac.NaryOperands(operand)
			newOps := make([]f.Formula, 0, len(ops))
			for _, op := range ops {
				if op != factorOutFormula {
					newOps = append(newOps, op)
				} else {
					removed = true
				}
			}
			app, _ := fac.NaryOperator(operand.Sort(), newOps...)
			if removed {
				formulasWithRemoved = append(formulasWithRemoved, app)
			} else {
				unchangedFormulas = append(unchangedFormulas, app)
			}
		} else {
			unchangedFormulas = append(unchangedFormulas, operand)
		}
	}
	unchanged, _ := fac.NaryOperator(sort, unchangedFormulas...)
	removed, _ := fac.NaryOperator(sort, formulasWithRemoved...)
	dualSort, _ := f.DualSort(sort)
	factorOut, _ := fac.NaryOperator(dualSort, factorOutFormula, removed)
	result, _ := fac.NaryOperator(sort, unchanged, factorOut)
	return result, true
}

func computeMaxOccurringSubformula(fac f.Factory, formula f.Formula) (f.Formula, bool) {
	formulaCounts := make(map[f.Formula]int)
	maxCount := 0
	maxFormula := fac.Falsum()
	for _, operand := range fac.Operands(formula) {
		if operand.Sort() == f.SortLiteral {
			formulaCounts[operand]++
			if newCount := formulaCounts[operand]; newCount > maxCount {
				maxCount = newCount
				maxFormula = operand
			}
		} else if operand.Sort() == f.SortAnd || operand.Sort() == f.SortOr {
			for _, subOperand := range fac.Operands(operand) {
				formulaCounts[subOperand]++
				if newCount := formulaCounts[subOperand]; newCount > maxCount {
					maxCount = newCount
					maxFormula = subOperand
				}
			}
		}
	}
	if maxCount < 2 {
		return 0, false
	} else {
		return maxFormula, true
	}
}
