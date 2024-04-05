package sat

import (
	"math"
	"slices"

	"github.com/booleworks/logicng-go/errorx"
	e "github.com/booleworks/logicng-go/explanation"
	f "github.com/booleworks/logicng-go/formula"
)

func (s *Solver) computeUnsatCore() *e.UnsatCore {
	clause2proposition := make(map[f.Formula]f.Proposition, 0)
	clauses := make([][]int32, len(s.core.pgOriginalClauses))

	for i, pi := range s.core.pgOriginalClauses {
		clauses[i] = pi.clause
		clause := getFormulaForVector(s, pi.clause)
		proposition := pi.proposition
		if proposition == nil {
			proposition = f.NewStandardProposition(clause)
		}
		clause2proposition[clause] = proposition
	}

	if containsEmptyClause(&clauses) {
		emptyClause := clause2proposition[s.fac.Falsum()]
		return e.NewUnsatCore([]f.Proposition{emptyClause}, true)
	}

	result := drupCompute(&clauses, &s.core.pgProof)

	if result.trivialUnsat {
		return handleTrivialCase(s)
	}

	props := make(map[f.Proposition]present, 0)
	for _, slice := range result.unsatCore {
		props[clause2proposition[getFormulaForVector(s, slice)]] = present{}
	}
	propositions := make([]f.Proposition, 0, len(props))
	for prop := range props {
		propositions = append(propositions, prop)
	}
	return e.NewUnsatCore(propositions, false)
}

func getFormulaForVector(solver *Solver, slice []int32) f.Formula {
	literals := make([]f.Formula, len(slice))
	slices.Sort(slice)
	for i := 0; i < len(slice); i++ {
		lit := slice[i]
		varName := solver.core.idx2name[int32(math.Abs(float64(lit))-1)]
		literals = append(literals, solver.fac.Literal(varName, lit > 0))
	}
	return solver.fac.Or(literals...)
}

func containsEmptyClause(clauses *[][]int32) bool {
	for _, clause := range *clauses {
		if len(clause) == 0 {
			return true
		}
	}
	return false
}

func handleTrivialCase(solver *Solver) *e.UnsatCore {
	clauses := solver.core.pgOriginalClauses
	for i := 0; i < len(clauses); i++ {
		for j := i + 1; j < len(clauses); j++ {
			if len(clauses[i].clause) == 1 && len(clauses[j].clause) == 1 &&
				clauses[i].clause[0]+clauses[j].clause[0] == 0 {
				propositions := make([]f.Proposition, 1, 2)
				pi := clauses[i].proposition
				if pi != nil {
					propositions[0] = pi
				} else {
					propositions[0] = f.NewStandardProposition(getFormulaForVector(solver, clauses[i].clause))
				}
				pj := clauses[j].proposition
				var pjp f.Proposition
				if pj != nil {
					pjp = pj
				} else {
					pjp = f.NewStandardProposition(getFormulaForVector(solver, clauses[j].clause))
				}
				if propositions[0] != pjp {
					propositions = append(propositions, pjp)
				}
				return e.NewUnsatCore(propositions, false)
			}
		}
	}
	panic(errorx.IllegalState("found no trivial unsat core"))
}

type present struct{}
