package assignment

import (
	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
)

// Restrict restricts a formula with the given assignment.  In contrast to
// Evaluate Restrict yield no truth value but a new formula where the
// literals of the assignment (and only these) are substituted by their
// respective mapped truth value.
func Restrict(fac f.Factory, formula f.Formula, assignment *Assignment) f.Formula {
	switch fsort := formula.Sort(); fsort {
	case f.SortTrue, f.SortFalse:
		return formula
	case f.SortLiteral:
		if formula.IsPos() {
			return assignment.restrictVariable(fac, f.Variable(formula))
		} else {
			return assignment.restrictNegativeLiteral(fac, f.Literal(formula))
		}
	case f.SortNot:
		op, _ := fac.NotOperand(formula)
		return fac.Not(Restrict(fac, op, assignment))
	case f.SortImpl, f.SortEquiv:
		left, right, _ := fac.BinaryLeftRight(formula)
		binOp, _ := fac.BinaryOperator(fsort, Restrict(fac, left, assignment), Restrict(fac, right, assignment))
		return binOp
	case f.SortAnd, f.SortOr:
		ops, _ := fac.NaryOperands(formula)
		nops := make([]f.Formula, 0, len(ops))
		for _, op := range ops {
			nops = append(nops, Restrict(fac, op, assignment))
		}
		naryOp, _ := fac.NaryOperator(fsort, nops...)
		return naryOp
	case f.SortCC, f.SortPBC:
		comparator, rhs, literals, coefficients, _ := fac.PBCOps(formula)
		return restrict(fac, comparator, rhs, literals, coefficients, assignment)
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
}

func restrict(
	fac f.Factory, comparator f.CSort, rhs int, literals []f.Literal, coefficients []int, assignment *Assignment,
) f.Formula {
	newLits := make([]f.Literal, 0)
	newCoeffs := make([]int, 0)
	lhsFixed := 0
	minValue := 0
	maxValue := 0
	var restriction f.Formula
	for i := 0; i < len(literals); i++ {
		if literals[i].IsPos() {
			restriction = assignment.restrictVariable(fac, f.Variable(literals[i]))
		} else {
			restriction = assignment.restrictNegativeLiteral(fac, literals[i])
		}
		if restriction.Sort() == f.SortLiteral {
			newLits = append(newLits, literals[i])
			coeff := coefficients[i]
			newCoeffs = append(newCoeffs, coeff)
			if coeff > 0 {
				maxValue += coeff
			} else {
				minValue += coeff
			}
		} else if restriction.Sort() == f.SortTrue {
			lhsFixed += coefficients[i]
		}
	}

	if len(newLits) == 0 {
		return fac.Constant(evaluateComparator(lhsFixed, rhs, comparator))
	}

	newRHS := rhs - lhsFixed
	if comparator != f.EQ {
		fixed := evaluateCoeffs(minValue, maxValue, newRHS, comparator)
		switch fixed {
		case f.TristateTrue:
			return fac.Verum()
		case f.TristateFalse:
			return fac.Falsum()
		default:
			// do nothing
		}
	}
	return fac.PBC(comparator, newRHS, newLits, newCoeffs)
}
