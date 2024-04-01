package sat

import (
	"booleworks.com/logicng/errorx"
	f "booleworks.com/logicng/formula"
)

// UpZeroLits returns all unit propagated literals on level 0 of the current
// formula on the solver.  Returns an error if the solver is not yet solced or
// the formula is UpZeroLits returns all unit propagated literals on level 0 of
// the current UNSAT.
func (s *Solver) UpZeroLits() ([]f.Literal, error) {
	if s.result == f.TristateUndef {
		return nil, errorx.IllegalState("SAT solver is not yet solved")
	}
	if s.result == f.TristateFalse {
		return nil, errorx.IllegalState("SAT problem was not satisfiable")
	}
	litIdxs := s.core.upZeroLiterals()
	lits := make([]f.Literal, len(litIdxs))
	for i, lit := range litIdxs {
		name := s.core.idx2name[Vari(lit)]
		lits[i] = s.fac.Lit(name, !Sign(lit))
	}
	return lits, nil
}
