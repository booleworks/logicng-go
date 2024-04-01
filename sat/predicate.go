package sat

import f "github.com/booleworks/logicng-go/formula"

// IsSatisfiable reports whether the formula is satisfiable.
func IsSatisfiable(fac f.Factory, formula f.Formula) bool {
	solver := NewSolver(fac)
	solver.Add(formula)
	return solver.Sat()
}

// IsTautology reports whether the formula is a tautology (always true).
func IsTautology(fac f.Factory, formula f.Formula) bool {
	solver := NewSolver(fac)
	solver.Add(formula.Negate(fac))
	return !solver.Sat()
}

// IsContradiction reports whether the formula is a contradiction (always false).
func IsContradiction(fac f.Factory, formula f.Formula) bool {
	solver := NewSolver(fac)
	solver.Add(formula)
	return !solver.Sat()
}

// Implies reports whether f1 implies f2 (f1 => f2 is always true).
func Implies(fac f.Factory, f1, f2 f.Formula) bool {
	solver := NewSolver(fac)
	solver.Add(fac.And(f1, f2.Negate(fac)))
	return !solver.Sat()
}

// IsEquivalent reports whether f1 is equivalent to f2 (f1 <=> f2 is always true).
func IsEquivalent(fac f.Factory, f1, f2 f.Formula) bool {
	solver := NewSolver(fac)
	solver.Add(fac.Or(fac.And(f1, f2.Negate(fac)), fac.And(f2, f1.Negate(fac))))
	return !solver.Sat()
}
