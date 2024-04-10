package bdd

import (
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/normalform"
	"github.com/booleworks/logicng-go/simplification"
)

// CNF transforms a given formula into a CNF by using a BDD. The
// resulting CNF does not contain any auxiliary variables, but can have
// quite a large size.
func CNF(fac f.Factory, formula f.Formula) f.Formula {
	cnf, _ := compute(fac, formula, true, nil)
	return cnf
}

// CNFWithHandler transforms a given formula into a CNF by using a BDD. The
// resulting CNF does not contain any auxiliary variables, but can have quite a
// large size.  The bddHandler can be used to abort the BDD compilation.  If
// the BDD compilation was aborted, the ok flag is false.
func CNFWithHandler(fac f.Factory, formula f.Formula, bddHandler Handler) (cnf f.Formula, ok bool) {
	return compute(fac, formula, true, bddHandler)
}

// DNF transforms a given formula into a DNF by using a BDD. The
// resulting DNF does not contain any auxiliary variables, but can have
// quite a large size.
func DNF(fac f.Factory, formula f.Formula) f.Formula {
	dnf, _ := compute(fac, formula, false, nil)
	return dnf
}

// DNFWithHandler transforms a given formula into a DNF by using a BDD. The
// resulting DNF does not contain any auxiliary variables, but can have quite a
// large size.  The bddHandler can be used to abort the BDD compilation.  If
// the BDD compilation was aborted, the ok flag is false.
func DNFWithHandler(fac f.Factory, formula f.Formula, bddHandler Handler) (cnf f.Formula, ok bool) {
	return compute(fac, formula, false, bddHandler)
}

func compute(fac f.Factory, formula f.Formula, cnf bool, bddHandler Handler) (f.Formula, bool) {
	var cacheEntry f.TransformationCacheSort
	if cnf {
		cacheEntry = f.TransBDDCNF
	} else {
		cacheEntry = f.TransBDDNF
	}
	if formula.Sort() <= f.SortLiteral {
		return formula, true
	}
	if hasNormalform(fac, formula, cnf) {
		return formula, true
	}
	cached, ok := f.LookupTransformationCache(fac, cacheEntry, formula)
	if ok {
		return cached, true
	}
	order := ForceOrder(fac, formula)
	bdd, ok := CompileWithVarOrderAndHandler(fac, formula, order, bddHandler)
	if !ok {
		return 0, false
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
	return simplifiedNormalForm, true
}

func hasNormalform(fac f.Factory, formula f.Formula, cnf bool) bool {
	if cnf {
		return normalform.IsCNF(fac, formula)
	} else {
		return normalform.IsDNF(fac, formula)
	}
}
