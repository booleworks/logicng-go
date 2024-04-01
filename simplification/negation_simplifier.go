package simplification

import (
	"booleworks.com/logicng/errorx"
	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/normalform"
)

// SimplifyNegations minimizes the number of negations of the given formula by
// applying De Morgan's Law heuristically for a smaller formula. The resulting
// formula is minimized for the length of its string representation (using the
// string representation which is defined in the formula's formula factory).
// For example, the formula ~A & ~B & ~C stays this way (since ~(A | B | C) is
// of same length as the initial formula), but the formula ~A & ~B & ~C & ~D is
// being transformed to ~(A | B | C | D) since its length is 16 vs. 17 in the
// un-simplified version.
func SimplifyNegations(fac f.Factory, formula f.Formula) f.Formula {
	nnf := normalform.NNF(fac, formula)
	if nnf.IsAtomic() {
		return getSmallestFormula(fac, true, formula, nnf)
	}
	result := negationMinimize(fac, nnf, true)
	return getSmallestFormula(fac, true, formula, nnf, result.positiveResult)
}

func negationMinimize(fac f.Factory, formula f.Formula, topLevel bool) minimizationResult {
	switch fsort := formula.Sort(); fsort {
	case f.SortLiteral:
		return minimizationResult{formula, formula.Negate(fac)}
	case f.SortOr, f.SortAnd:
		ops, _ := fac.NaryOperands(formula)
		opResults := make([]minimizationResult, len(ops))
		for i, op := range ops {
			opResults[i] = negationMinimize(fac, op, false)
		}
		positiveOpResults := make([]f.Formula, len(opResults))
		negativeOpResults := make([]f.Formula, len(opResults))
		for i, result := range opResults {
			positiveOpResults[i] = result.positiveResult
			negativeOpResults[i] = result.negativeResult
		}
		smallestPositive := findSmallestPositive(fac, formula.Sort(), positiveOpResults, negativeOpResults, topLevel)
		smallestNegative := findSmallestNegative(fac, formula.Sort(), negativeOpResults, smallestPositive, topLevel)
		return minimizationResult{smallestPositive, smallestNegative}
	case f.SortFalse, f.SortTrue, f.SortNot, f.SortEquiv, f.SortImpl, f.SortCC, f.SortPBC:
		panic(errorx.IllegalState("unexpected formula in NNF: %s", fsort))
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
}

func findSmallestPositive(
	fac f.Factory, sort f.FSort, positiveOpResults, negativeOpResults []f.Formula, topLevel bool,
) f.Formula {
	allPositive, _ := fac.NaryOperator(sort, positiveOpResults...)
	smallerPositiveOps := make([]f.Formula, 0, len(positiveOpResults))
	smallerNegativeOps := make([]f.Formula, 0, len(positiveOpResults))
	for i := 0; i < len(positiveOpResults); i++ {
		positiveOp := positiveOpResults[i]
		negativeOp := negativeOpResults[i]
		if formattedLength(fac, positiveOp, false) < formattedLength(fac, negativeOp, false) {
			smallerPositiveOps = append(smallerPositiveOps, positiveOp)
		} else {
			smallerNegativeOps = append(smallerNegativeOps, negativeOp)
		}
	}
	smallerPosOp, _ := fac.NaryOperator(sort, smallerPositiveOps...)
	dualSort, _ := f.DualSort(sort)
	smallerNegOp, _ := fac.NaryOperator(dualSort, smallerNegativeOps...)
	partialNegative, _ := fac.NaryOperator(sort, smallerPosOp, fac.Not(smallerNegOp))
	return getSmallestFormula(fac, topLevel, allPositive, partialNegative)
}

func findSmallestNegative(
	fac f.Factory, sort f.FSort, negativeOpResults []f.Formula, smallestPositive f.Formula, topLevel bool,
) f.Formula {
	negation := fac.Not(smallestPositive)
	dualSort, _ := f.DualSort(sort)
	flipped, _ := fac.NaryOperator(dualSort, negativeOpResults...)
	return getSmallestFormula(fac, topLevel, negation, flipped)
}

func getSmallestFormula(fac f.Factory, topLevel bool, formulas ...f.Formula) f.Formula {
	currentLen := formattedLength(fac, formulas[0], topLevel)
	currentFormula := formulas[0]
	for i := 1; i < len(formulas); i++ {
		if length := formattedLength(fac, formulas[i], topLevel); length < currentLen {
			currentLen = length
			currentFormula = formulas[i]
		}
	}
	return currentFormula
}

func formattedLength(fac f.Factory, formula f.Formula, topLevel bool) int {
	length := len(formula.Sprint(fac))
	if !topLevel && formula.Sort() == f.SortOr {
		return length + 2
	} else {
		return length
	}
}

type minimizationResult struct {
	positiveResult f.Formula
	negativeResult f.Formula
}
