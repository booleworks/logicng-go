package enum

import (
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/model/iter"
	"github.com/booleworks/logicng-go/sat"
)

// CanonicalCNF returns a canonical CNF of the given formula.
func CanonicalCNF(fac f.Factory, formula f.Formula) f.Formula {
	cnf, _ := canonicalEnumeration(fac, formula, true, nil)
	return cnf
}

// CanonicalCNF returns a canonical CNF of the given formula.  The given
// iterHandler can be used to abort the computation.  If the enumeration was
// aborted, the ok flag is false.
func CanonicalCNFWithHandler(fac f.Factory, formula f.Formula, iterHandler iter.Handler) (cnf f.Formula, ok bool) {
	return canonicalEnumeration(fac, formula, true, iterHandler)
}

// CanonicalDNF returns a canonical DNF of the given formula.
func CanonicalDNF(fac f.Factory, formula f.Formula) f.Formula {
	dnf, _ := canonicalEnumeration(fac, formula, false, nil)
	return dnf
}

// CanonicalDNF returns a canonical DNF of the given formula.  The given
// iterHandler can be used to abort the computation.  If the enumeration was
// aborted, the ok flag is false.
func CanonicalDNFWithHandler(fac f.Factory, formula f.Formula, iterHandler iter.Handler) (cnf f.Formula, ok bool) {
	return canonicalEnumeration(fac, formula, false, iterHandler)
}

func canonicalEnumeration(
	fac f.Factory,
	formula f.Formula,
	cnf bool,
	iterHandler iter.Handler,
) (f.Formula, bool) {
	solver := sat.NewSolver(fac)
	if cnf {
		solver.Add(formula.Negate(fac))
	} else {
		solver.Add(formula)
	}
	config := iter.DefaultConfig()
	config.Handler = iterHandler
	enumeration, ok := OnSolverWithConfig(solver, f.Variables(fac, formula).Content(), config)
	if !ok {
		return 0, false
	}
	if len(enumeration) == 0 {
		return fac.Constant(cnf), true
	}
	ops := make([]f.Formula, len(enumeration))
	for i, m := range enumeration {
		if cnf {
			neg := make([]f.Literal, len(m.Literals))
			for i, l := range m.Literals {
				neg[i] = l.Negate(fac)
			}
			ops[i] = fac.Clause(neg...)
		} else {
			ops[i] = fac.Minterm(m.Literals...)
		}
	}
	if cnf {
		return fac.And(ops...), true
	} else {
		return fac.Or(ops...), true
	}
}
