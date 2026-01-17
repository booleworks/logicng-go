package bdd

import (
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/normalform"
	"github.com/booleworks/logicng-go/simplification"
)

// CNF transforms a given formula into a CNF by using a BDD. The
// resulting CNF does not contain any auxiliary variables, but can have
// quite a large size.
func CNF(fac f.Factory, formula f.Formula) f.Formula {
	cnf, _ := compute(fac, formula, true, handler.NopHandler)
	return cnf
}

// CNFWithHandler transforms a given formula into a CNF by using a BDD. The
// resulting CNF does not contain any auxiliary variables, but can have quite a
// large size.  The bddHandler can be used to cancel the BDD compilation.
func CNFWithHandler(fac f.Factory, formula f.Formula, hdl handler.Handler) (f.Formula, handler.State) {
	return compute(fac, formula, true, hdl)
}

// DNF transforms a given formula into a DNF by using a BDD. The
// resulting DNF does not contain any auxiliary variables, but can have
// quite a large size.
func DNF(fac f.Factory, formula f.Formula) f.Formula {
	dnf, _ := compute(fac, formula, false, handler.NopHandler)
	return dnf
}

// DNFWithHandler transforms a given formula into a DNF by using a BDD. The
// resulting DNF does not contain any auxiliary variables, but can have quite a
// large size.  The bddHandler can be used to cancel the BDD compilation.
func DNFWithHandler(fac f.Factory, formula f.Formula, hdl handler.Handler) (f.Formula, handler.State) {
	return compute(fac, formula, false, hdl)
}

func compute(fac f.Factory, formula f.Formula, cnf bool, hdl handler.Handler) (f.Formula, handler.State) {
	var cacheEntry f.TransformationCacheSort
	if cnf {
		cacheEntry = f.TransBDDCNF
	} else {
		cacheEntry = f.TransBDDNF
	}
	if formula.Sort() <= f.SortLiteral {
		return formula, succ
	}
	if hasNormalform(fac, formula, cnf) {
		return formula, succ
	}
	cached, ok := f.LookupTransformationCache(fac, cacheEntry, formula)
	if ok {
		return cached, succ
	}
	order := ForceOrder(fac, formula)
	bdd, state := CompileWithVarOrderAndHandler(fac, formula, order, hdl)
	if !state.Success {
		return 0, state
	}
	var normalForm f.Formula
	if cnf {
		normalForm = bdd.CNF()
	} else {
		normalForm = bdd.DNF()
	}
	var simplifiedNormalForm f.Formula
	if cnf {
		simplifiedNormalForm = simplification.PropagateUnits(fac, normalForm)
	} else {
		negatedDnf := normalform.NNF(fac, normalForm.Negate(fac))
		simplifiedNormalForm = normalform.NNF(fac, simplification.PropagateUnits(fac, negatedDnf).Negate(fac))
	}
	f.SetTransformationCache(fac, cacheEntry, formula, simplifiedNormalForm)
	return simplifiedNormalForm, succ
}

func hasNormalform(fac f.Factory, formula f.Formula, cnf bool) bool {
	if cnf {
		return normalform.IsCNF(fac, formula)
	}
	return normalform.IsDNF(fac, formula)
}
