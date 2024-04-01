package normalform

import (
	"booleworks.com/logicng/errorx"
	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/handler"
)

// IsDNF reports whether the given formula is in disjunctive normal form.  A
// DNF is a disjunction of conjunctions of literals.
func IsDNF(fac f.Factory, formula f.Formula) bool {
	cached, ok := f.LookupPredicateCache(fac, f.PredDNF, formula)
	if ok {
		return cached
	}
	var result bool
	switch fsort := formula.Sort(); fsort {
	case f.SortFalse, f.SortTrue, f.SortLiteral:
		return true
	case f.SortNot, f.SortImpl, f.SortEquiv, f.SortCC, f.SortPBC:
		return false
	case f.SortOr:
		result = true
		ops, _ := fac.NaryOperands(formula)
		for _, op := range ops {
			if !IsMinterm(fac, op) {
				result = false
				break
			}
		}
	case f.SortAnd:
		result = IsMinterm(fac, formula)
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}
	f.SetPredicateCache(fac, f.PredDNF, formula, result)
	return result
}

// FactorizedDNF returns the given formula in disjunctive normal form.  A DNF
// is a disjunction of conjunctions of literals.  The algorithm used is
// factorization.  The resulting DNF can grow exponentially, therefore unless
// you are sure that the input is sensible, prefer the DNF factorization with a
// handler in order to be able to abort it.
func FactorizedDNF(fac f.Factory, formula f.Formula) f.Formula {
	cnf, _ := factorizedDNFRec(fac, formula, nil)
	return cnf
}

// FactorizedDNFWithHandler returns the given formula in disjunctive normal
// form.  A DNF is a disjunction of conjunctions of literals.  The given
// handler can be used to abort the factorization.  Returns the DNF and an ok
// flag which is false when the handler aborted the computation.
func FactorizedDNFWithHandler(
	fac f.Factory, formula f.Formula, factorizationHandler FactorizationHandler,
) (dnf f.Formula, ok bool) {
	handler.Start(factorizationHandler)
	return factorizedDNFRec(fac, formula, factorizationHandler)
}

func factorizedDNFRec(fac f.Factory, formula f.Formula, handler FactorizationHandler) (f.Formula, bool) {
	if formula.Sort() <= f.SortLiteral {
		return formula, true
	}
	cached, ok := f.LookupTransformationCache(fac, f.TransDNFFactorization, formula)
	if ok {
		return cached, true
	}
	ok = true
	switch fsort := formula.Sort(); fsort {
	case f.SortNot, f.SortImpl, f.SortEquiv, f.SortCC, f.SortPBC:
		cached, ok = factorizedDNFRec(fac, NNF(fac, formula), handler)
	case f.SortOr:

		nary, _ := fac.NaryOperands(formula)
		nops := make([]f.Formula, 0, len(nary))
		for _, op := range nary {
			var apply f.Formula
			apply, ok = factorizedDNFRec(fac, op, handler)
			if !ok {
				return 0, false
			}
			nops = append(nops, apply)
		}
		cached = fac.Or(nops...)
	case f.SortAnd:
		nary, _ := fac.NaryOperands(formula)
		nops := make([]f.Formula, 0, len(nary))
		for _, op := range nary {
			if !ok {
				return 0, false
			}
			var nop f.Formula
			nop, ok = factorizedDNFRec(fac, op, handler)
			nops = append(nops, nop)
		}
		cached = nops[0]
		for i := 1; i < len(nops); i++ {
			if !ok {
				return 0, false
			}
			cached, ok = distributeDNF(fac, cached, nops[i], handler)
		}
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}

	if ok {
		f.SetTransformationCache(fac, f.TransDNFFactorization, formula, cached)
		return cached, true
	}
	return 0, false
}

func distributeDNF(fac f.Factory, f1, f2 f.Formula, handler FactorizationHandler) (f.Formula, bool) {
	proceed := true
	if handler != nil {
		proceed = handler.PerformedDistribution()
	}
	if !proceed {
		return 0, false
	}
	if f1.Sort() == f.SortOr || f2.Sort() == f.SortOr {
		nops := make([]f.Formula, 0)
		var operands []f.Formula
		var form f.Formula
		if f1.Sort() == f.SortOr {
			form = f2
			operands, _ = fac.NaryOperands(f1)
		} else {
			form = f1
			operands, _ = fac.NaryOperands(f2)
		}
		for _, op := range operands {
			distribute, ok := distributeDNF(fac, op, form, handler)
			if !ok {
				return 0, false
			}
			nops = append(nops, distribute)
		}
		return fac.Or(nops...), true
	}
	clause := fac.And(f1, f2)
	if handler != nil {
		proceed = handler.CreatedClause(clause)
	}
	return clause, proceed
}
