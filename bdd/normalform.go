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
	return compute(fac, formula, true)
}

// DNF transforms a given formula into a DNF by using a BDD. The
// resulting DNF does not contain any auxiliary variables, but can have
// quite a large size.
func DNF(fac f.Factory, formula f.Formula) f.Formula {
	return compute(fac, formula, false)
}

func compute(fac f.Factory, formula f.Formula, cnf bool) f.Formula {
	var cacheEntry f.TransformationCacheSort
	if cnf {
		cacheEntry = f.TransBDDCNF
	} else {
		cacheEntry = f.TransBDDNF
	}
	if formula.Sort() <= f.SortLiteral {
		return formula
	}
	if hasNormalform(fac, formula, cnf) {
		return formula
	}
	cached, ok := f.LookupTransformationCache(fac, cacheEntry, formula)
	if ok {
		return cached
	}
	order := ForceOrder(fac, formula)
	bdd := CompileWithVarOrder(fac, formula, order)
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
	return simplifiedNormalForm
}

func hasNormalform(fac f.Factory, formula f.Formula, cnf bool) bool {
	if cnf {
		return normalform.IsCNF(fac, formula)
	} else {
		return normalform.IsDNF(fac, formula)
	}
}
