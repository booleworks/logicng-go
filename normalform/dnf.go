package normalform

import (
	"github.com/booleworks/logicng-go/errorx"
	"github.com/booleworks/logicng-go/event"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
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
// handler in order to be able to cancel it.
func FactorizedDNF(fac f.Factory, formula f.Formula) f.Formula {
	cnf, _ := factorizedDNFRec(fac, formula, handler.NopHandler)
	return cnf
}

// FactorizedDNFWithHandler returns the given formula in disjunctive normal
// form.  A DNF is a disjunction of conjunctions of literals.  The given
// handler can be used to cancel the factorization.  Returns the DNF and
// the handler state.
func FactorizedDNFWithHandler(fac f.Factory, formula f.Formula, hdl handler.Handler) (f.Formula, handler.State) {
	if e := event.FactorizationStarted; !hdl.ShouldResume(e) {
		return 0, handler.Cancelation(e)
	}
	return factorizedDNFRec(fac, formula, hdl)
}

func factorizedDNFRec(fac f.Factory, formula f.Formula, hdl handler.Handler) (f.Formula, handler.State) {
	if formula.Sort() <= f.SortLiteral {
		return formula, succ
	}
	cached, ok := f.LookupTransformationCache(fac, f.TransDNFFactorization, formula)
	if ok {
		return cached, succ
	}
	state := handler.Success()
	switch fsort := formula.Sort(); fsort {
	case f.SortNot, f.SortImpl, f.SortEquiv, f.SortCC, f.SortPBC:
		cached, state = factorizedDNFRec(fac, NNF(fac, formula), hdl)
	case f.SortOr:
		nary, _ := fac.NaryOperands(formula)
		nops := make([]f.Formula, 0, len(nary))
		for _, op := range nary {
			var apply f.Formula
			apply, state = factorizedDNFRec(fac, op, hdl)
			if !state.Success {
				return 0, state
			}
			nops = append(nops, apply)
		}
		cached = fac.Or(nops...)
	case f.SortAnd:
		nary, _ := fac.NaryOperands(formula)
		nops := make([]f.Formula, 0, len(nary))
		for _, op := range nary {
			if !state.Success {
				return 0, state
			}
			var nop f.Formula
			nop, state = factorizedDNFRec(fac, op, hdl)
			nops = append(nops, nop)
		}
		cached = nops[0]
		for i := 1; i < len(nops); i++ {
			if !state.Success {
				return 0, state
			}
			cached, state = distributeDNF(fac, cached, nops[i], hdl)
		}
	default:
		panic(errorx.UnknownEnumValue(fsort))
	}

	if state.Success {
		f.SetTransformationCache(fac, f.TransDNFFactorization, formula, cached)
		return cached, succ
	}
	return 0, state
}

func distributeDNF(fac f.Factory, f1, f2 f.Formula, hdl handler.Handler) (f.Formula, handler.State) {
	if e := event.DistributionPerformed; !hdl.ShouldResume(e) {
		return 0, handler.Cancelation(e)
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
			distribute, state := distributeDNF(fac, op, form, hdl)
			if !state.Success {
				return 0, state
			}
			nops = append(nops, distribute)
		}
		return fac.Or(nops...), succ
	}
	clause := fac.And(f1, f2)
	if e := event.FactorizationCreatedClause; !hdl.ShouldResume(e) {
		return clause, handler.Cancelation(e)
	}
	return clause, succ
}
