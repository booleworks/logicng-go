package sat

import f "booleworks.com/logicng/formula"

// VarOccurrences counts all occurrences of variables on the solver and returns
// the mapping from variable to number of occurrences.  Note that these are
// usually not the same occurrences as in the original formula, since the
// formula might have been converted to CNF and/or variables in clauses might
// have been subsumed.  If variables are given as parameters, only the
// occurrences of these variables are counted, otherwise all variables are
// considered.
func (s *Solver) VarOccurrences(variables *f.VarSet) map[f.Variable]int {
	fac := s.fac
	solver := s.core
	var relevantVars *f.VarSet
	if variables != nil {
		relevantVars = f.NewVarSet()
		relevantVars.AddAll(variables)
	}
	counts := initResultMap(fac, solver, relevantVars)
	for _, clause := range solver.clauses {
		for i := 0; i < clause.size(); i++ {
			key := solver.idx2name[Vari(clause.get(i))]
			if cnt, ok := counts[key]; ok {
				counts[key] = cnt + 1
			}
		}
	}
	result := make(map[f.Variable]int, len(counts))
	for k, v := range counts {
		result[fac.Var(k)] = v
	}
	return result
}

func initResultMap(fac f.Factory, solver *CoreSolver, relevantVars *f.VarSet) map[string]int {
	counts := make(map[string]int)
	if relevantVars != nil {
		for _, v := range relevantVars.Content() {
			name, _ := fac.VarName(v)
			counts[name] = 0
		}
	}
	variables := solver.vars
	for i, v := range variables {
		name := solver.idx2name[int32(i)]
		variable := fac.Var(name)
		if relevantVars == nil || relevantVars.Contains(variable) {
			if v.level == 0 {
				counts[name] = 1
			} else {
				counts[name] = 0
			}
		}
	}
	return counts
}
