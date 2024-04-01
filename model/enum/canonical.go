package enum

import (
	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/sat"
)

// CanonicalCNF returns a canonical CNF of the given formula.
func CanonicalCNF(fac f.Factory, formula f.Formula) f.Formula {
	return canonicalEnumeration(fac, formula, true)
}

// CanonicalDNF returns a canonical DNF of the given formula.
func CanonicalDNF(fac f.Factory, formula f.Formula) f.Formula {
	return canonicalEnumeration(fac, formula, false)
}

func canonicalEnumeration(fac f.Factory, formula f.Formula, cnf bool) f.Formula {
	solver := sat.NewSolver(fac)
	if cnf {
		solver.Add(formula.Negate(fac))
	} else {
		solver.Add(formula)
	}
	enumeration := OnSolver(solver, f.Variables(fac, formula).Content())
	if len(enumeration) == 0 {
		return fac.Constant(cnf)
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
		return fac.And(ops...)
	} else {
		return fac.Or(ops...)
	}
}
