package enum

import (
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/handler"
	"github.com/booleworks/logicng-go/model/iter"
	"github.com/booleworks/logicng-go/sat"
)

// CanonicalCNF returns a canonical CNF of the given formula.
func CanonicalCNF(fac f.Factory, formula f.Formula) f.Formula {
	cnf, _ := canonicalEnumeration(fac, formula, true, handler.NopHandler)
	return cnf
}

// CanonicalCNFWithHandler returns a canonical CNF of the given formula.
// The given  iterHandler can be used to cancel the computation.
func CanonicalCNFWithHandler(fac f.Factory, formula f.Formula, hdl handler.Handler) (f.Formula, handler.State) {
	return canonicalEnumeration(fac, formula, true, hdl)
}

// CanonicalDNF returns a canonical DNF of the given formula.
func CanonicalDNF(fac f.Factory, formula f.Formula) f.Formula {
	dnf, _ := canonicalEnumeration(fac, formula, false, handler.NopHandler)
	return dnf
}

// CanonicalDNFWithHandler returns a canonical DNF of the given formula.
// The given  iterHandler can be used to cancel the computation.
func CanonicalDNFWithHandler(fac f.Factory, formula f.Formula, hdl handler.Handler) (f.Formula, handler.State) {
	return canonicalEnumeration(fac, formula, false, hdl)
}

func canonicalEnumeration(
	fac f.Factory,
	formula f.Formula,
	cnf bool,
	hdl handler.Handler,
) (f.Formula, handler.State) {
	solver := sat.NewSolver(fac)
	if cnf {
		solver.Add(formula.Negate(fac))
	} else {
		solver.Add(formula)
	}
	config := iter.DefaultConfig()
	config.Handler = hdl
	enumeration, state := OnSolverWithConfig(solver, f.Variables(fac, formula).Content(), config)
	if !state.Success {
		return 0, state
	}
	if len(enumeration) == 0 {
		return fac.Constant(cnf), succ
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
		return fac.And(ops...), succ
	} else {
		return fac.Or(ops...), succ
	}
}
