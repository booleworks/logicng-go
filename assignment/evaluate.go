package assignment

import (
	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
)

// Evaluate evaluates the formula with the given assignment.  A literal not
// covered by the assignment evaluates to false if it is positive and true if
// it is negative.
func Evaluate(fac f.Factory, formula f.Formula, assignment *Assignment) bool {
	switch fsort := formula.Sort(); fsort {
	case f.SortTrue:
		return true
	case f.SortFalse:
		return false
	case f.SortLiteral:
		if formula.IsPos() {
			return assignment.evaluateVariable(f.Variable(formula))
		} else {
			return assignment.evaluateNegativeLiteral(f.Literal(formula))
		}
	case f.SortNot:
		op, _ := fac.NotOperand(formula)
		return !Evaluate(fac, op, assignment)
	case f.SortImpl:
		left, right, _ := fac.BinaryLeftRight(formula)
		return !Evaluate(fac, left, assignment) || Evaluate(fac, right, assignment)
	case f.SortEquiv:
		left, right, _ := fac.BinaryLeftRight(formula)
		return Evaluate(fac, left, assignment) == Evaluate(fac, right, assignment)
	case f.SortAnd:
		ops, _ := fac.NaryOperands(formula)
		for _, op := range ops {
			if !Evaluate(fac, op, assignment) {
				return false
			}
		}
		return true
	case f.SortOr:
		ops, _ := fac.NaryOperands(formula)
		for _, op := range ops {
			if Evaluate(fac, op, assignment) {
				return true
			}
		}
		return false
	case f.SortCC, f.SortPBC:
		comparator, rhs, literals, coefficients, _ := fac.PBCOps(formula)
		lhs := evaluateLhs(fac, literals, coefficients, assignment)
		return evaluateComparator(lhs, rhs, comparator)
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
}

func evaluateLhs(fac f.Factory, literals []f.Literal, coefficients []int, assignment *Assignment) int {
	lhs := 0
	for i := range literals {
		if Evaluate(fac, literals[i].AsFormula(), assignment) {
			lhs += coefficients[i]
		}
	}
	return lhs
}

func evaluateComparator(lhs, rhs int, comparator f.CSort) bool {
	switch comparator {
	case f.EQ:
		return lhs == rhs
	case f.LE:
		return lhs <= rhs
	case f.LT:
		return lhs < rhs
	case f.GE:
		return lhs >= rhs
	case f.GT:
		return lhs > rhs
	default:
		panic(errorx.UnknownEnumValue(comparator))
	}
}

func evaluateCoeffs(minValue, maxValue, rhs int, comparator f.CSort) f.Tristate {
	status := 0
	if rhs >= minValue {
		status++
	}
	if rhs > minValue {
		status++
	}
	if rhs >= maxValue {
		status++
	}
	if rhs > maxValue {
		status++
	}

	switch comparator {
	case f.EQ:
		if status == 0 || status == 4 {
			return f.TristateFalse
		}
	case f.LE:
		if status >= 3 {
			return f.TristateTrue
		} else if status < 1 {
			return f.TristateFalse
		}
	case f.LT:
		if status > 3 {
			return f.TristateTrue
		} else if status <= 1 {
			return f.TristateFalse
		}
	case f.GE:
		if status <= 1 {
			return f.TristateTrue
		} else if status > 3 {
			return f.TristateFalse
		}
	case f.GT:
		if status < 1 {
			return f.TristateTrue
		} else if status >= 3 {
			return f.TristateFalse
		}
	default:
		panic(errorx.UnknownEnumValue(comparator))
	}
	return f.TristateUndef
}
