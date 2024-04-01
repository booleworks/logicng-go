package sat

import f "booleworks.com/logicng/formula"

// FormulasOnSolver returns the current formulas on the solver.
//
// Note that this formula is usually syntactically different to the formulas
// which were actually added to the solver, since the formulas are added as CNF
// and may be simplified or even removed depending on the state of the solver.
// Furthermore, the solver might add learnt clauses or propagate literals.
//
// If the formula on the solver is known to be unsatisfiable, this function
// will add false to the returned set of formulas. However, as long as Sat was
// not called on the current solver state, the absence of false does not imply
// that the formula is satisfiable.
//
// Also note that formulas are not added to the solver as soon as the solver is
// known be unsatisfiable.
func (s *Solver) FormulasOnSolver() []f.Formula {
	formulas := f.NewFormulaSet()
	for _, clause := range s.core.clauses {
		lits := make([]f.Literal, clause.size())
		for i := 0; i < clause.size(); i++ {
			litInt := clause.get(i)
			lits[i] = s.fac.Lit(s.core.idx2name[litInt>>1], (litInt&1) != 1)
		}
		if !clause.isAtMost {
			formulas.Add(s.fac.Clause(lits...))
		} else {
			rhs := clause.size() + 1 - clause.atMostWatchers
			vars := make([]f.Variable, len(lits))
			for i, lit := range lits {
				vars[i] = lit.Variable()
			}
			formulas.Add(s.fac.CC(f.LE, uint32(rhs), vars...))
		}
	}
	for i := 0; i < len(s.core.vars); i++ {
		variable := s.core.vars[i]
		if variable.level == 0 {
			formulas.Add(s.fac.Literal(s.core.idx2name[int32(i)], variable.assignment == f.TristateTrue))
		}
	}
	if !s.core.ok {
		formulas.Add(s.fac.Falsum())
	}
	return formulas.Content()
}
