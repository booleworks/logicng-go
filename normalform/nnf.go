package normalform

import (
	"github.com/booleworks/logicng-go/encoding"
	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
)

// NNF returns the negation normal form of the given formula.  In an NNF only
// negation, conjunction, and disjunction are allowed and negations must only
// appear before variables.
func NNF(fac f.Factory, formula f.Formula) f.Formula {
	return nnfRec(fac, formula, true)
}

// IsNNF reports whether the given formula is in negation normal form.   In an
// NNF only negation, conjunction, and disjunction are allowed and negations
// must only appear before variables.
func IsNNF(fac f.Factory, formula f.Formula) bool {
	result, ok := f.LookupPredicateCache(fac, f.PredNNF, formula)
	if ok {
		return result
	}
	switch fsort := formula.Sort(); fsort {
	case f.SortFalse, f.SortTrue, f.SortLiteral:
		result = true
	case f.SortAnd, f.SortOr:
		result = true
		nary, _ := fac.NaryOperands(formula)
		for _, op := range nary {
			if !IsNNF(fac, op) {
				result = false
				break
			}
		}
	case f.SortNot, f.SortImpl, f.SortEquiv, f.SortCC, f.SortPBC:
		result = false
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
	f.SetPredicateCache(fac, f.PredNNF, formula, result)
	return result
}

func nnfRec(fac f.Factory, formula f.Formula, polarity bool) f.Formula {
	if polarity {
		nnf, ok := f.LookupTransformationCache(fac, f.TransNNF, formula)
		if ok {
			return nnf
		}
	}
	var nnf f.Formula
	switch fsort := formula.Sort(); fsort {
	case f.SortTrue, f.SortFalse, f.SortLiteral:
		if polarity {
			nnf = formula
		} else {
			nnf = fac.Not(formula)
		}
	case f.SortNot:
		op, _ := fac.NotOperand(formula)
		nnf = nnfRec(fac, op, !polarity)
	case f.SortOr, f.SortAnd:
		ops, _ := fac.NaryOperands(formula)
		nnf = nnfRecNary(fac, fsort, polarity, ops...)
	case f.SortEquiv:
		left, right, _ := fac.BinaryLeftRight(formula)
		if polarity {
			nnf = fac.And(fac.Or(nnfRec(fac, left, false), nnfRec(fac, right, true)),
				fac.Or(nnfRec(fac, left, true), nnfRec(fac, right, false)))
		} else {
			nnf = fac.And(fac.Or(nnfRec(fac, left, false), nnfRec(fac, right, false)),
				fac.Or(nnfRec(fac, left, true), nnfRec(fac, right, true)))
		}
	case f.SortImpl:
		left, right, _ := fac.BinaryLeftRight(formula)
		if polarity {
			nnf = fac.Or(nnfRec(fac, left, false), nnfRec(fac, right, true))
		} else {
			nnf = fac.And(nnfRec(fac, left, true), nnfRec(fac, right, false))
		}
	case f.SortCC, f.SortPBC:
		if polarity {
			pbcEncoding, err := encoding.EncodePBC(fac, formula)
			if err != nil {
				panic(err) // we are sure the formula is a PBC, so the rhs must be MaxInt, nothing we can do here
			}
			nnf = nnfRecNary(fac, f.SortAnd, true, pbcEncoding...)
		} else {
			nnf = nnfRec(fac, encoding.NegatePBC(fac, formula), true)
		}
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
	if polarity {
		f.SetTransformationCache(fac, f.TransNNF, formula, nnf)
	}
	return nnf
}

func nnfRecNary(fac f.Factory, fsort f.FSort, polarity bool, operands ...f.Formula) f.Formula {
	nops := make([]f.Formula, 0, len(operands))
	for _, op := range operands {
		nops = append(nops, nnfRec(fac, op, polarity))
	}
	var sort f.FSort
	if polarity {
		sort = fsort
	} else if fsort == f.SortAnd {
		sort = f.SortOr
	} else {
		sort = f.SortAnd
	}
	naryOp, _ := fac.NaryOperator(sort, nops...)
	return naryOp
}
