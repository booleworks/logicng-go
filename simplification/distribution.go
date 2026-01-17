package simplification

import (
	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
)

// Distribute simplifies the given formula by applying the distributive laws.
// In contrast to the FactorOut function, the distribution step is only
// performed once.  E.g. for a formula (A | B) & (A | C & E) | B & C & D the
// Distribute function yields A | B & C & E | B & C & D whereas the FactorOut
// function also factors our B & C, yielding A | C & B & (E | D).
func Distribute(fac f.Factory, formula f.Formula) f.Formula {
	result, ok := f.LookupTransformationCache(fac, f.TransDistrSimpl, formula)
	if ok {
		return result
	}
	switch fsort := formula.Sort(); fsort {
	case f.SortFalse, f.SortTrue, f.SortLiteral, f.SortCC, f.SortPBC:
		result = formula
	case f.SortNot:
		op, _ := fac.NotOperand(formula)
		result = fac.Not(Distribute(fac, op))
	case f.SortImpl, f.SortEquiv:
		left, right, _ := fac.BinaryLeftRight(formula)
		result, _ = fac.BinaryOperator(fsort, Distribute(fac, left), Distribute(fac, right))
	case f.SortOr, f.SortAnd:
		result = distributeNary(fac, formula)
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
	f.SetTransformationCache(fac, f.TransDistrSimpl, formula, result)
	return result
}

func distributeNary(fac f.Factory, formula f.Formula) f.Formula {
	var result f.Formula
	outerSort := formula.Sort()
	innerSort, _ := f.DualSort(outerSort)
	ops, _ := fac.NaryOperands(formula)
	operands := f.NewMutableFormulaSet()
	for _, op := range ops {
		operands.Add(Distribute(fac, op))
	}
	part2Operands := make(map[f.Formula]*f.MutableFormulaSet)
	mostCommon := fac.Falsum()
	mostCommonAmount := 0
	for _, op := range operands.Content() {
		if op.Sort() == innerSort {
			for _, part := range fac.Operands(op) {
				partOperands, ok := part2Operands[part]
				if !ok {
					partOperands = f.NewMutableFormulaSet()
					part2Operands[part] = partOperands
				}
				partOperands.Add(op)
				if partOperands.Size() > mostCommonAmount {
					mostCommon = part
					mostCommonAmount = partOperands.Size()
				}
			}
		}
	}
	if mostCommon == fac.Falsum() || mostCommonAmount == 1 {
		result, _ = fac.NaryOperator(outerSort, operands.Content()...)
		return result
	}
	operands.RemoveAll(part2Operands[mostCommon].AsImmutable())
	set := part2Operands[mostCommon]
	relevantFormulas := make([]f.Formula, 0, set.Size())
	for _, preRelevantFormula := range set.Content() {
		relevantParts := make([]f.Formula, 0, len(fac.Operands(preRelevantFormula))-1)
		for _, part := range fac.Operands(preRelevantFormula) {
			if part != mostCommon {
				relevantParts = append(relevantParts, part)
			}
		}
		naryOp, _ := fac.NaryOperator(innerSort, relevantParts...)
		relevantFormulas = append(relevantFormulas, naryOp)
	}
	outerOp, _ := fac.NaryOperator(outerSort, relevantFormulas...)
	innerOp, _ := fac.NaryOperator(innerSort, mostCommon, outerOp)
	operands.Add(innerOp)
	result, _ = fac.NaryOperator(outerSort, operands.Content()...)
	return result
}
