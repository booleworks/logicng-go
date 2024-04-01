package maxsat

import "github.com/booleworks/logicng-go/sat"

func encodeLadder(s *sat.CoreSolver, lits []int32) {
	if len(lits) == 1 {
		addUnitClause(s, lits[0])
	} else {
		seqAuxiliary := make([]int32, len(lits)-1)
		for i := 0; i < len(lits)-1; i++ {
			seqAuxiliary[i] = sat.MkLit(s.NVars(), false)
			newSatVariable(s)
		}
		for i := 0; i < len(lits); i++ {
			if i == 0 {
				addBinaryClause(s, lits[i], sat.Not(seqAuxiliary[i]))
				addBinaryClause(s, sat.Not(lits[i]), seqAuxiliary[i])
			} else if i == len(lits)-1 {
				addBinaryClause(s, lits[i], seqAuxiliary[i-1])
				addBinaryClause(s, sat.Not(lits[i]), sat.Not(seqAuxiliary[i-1]))
			} else {
				addBinaryClause(s, sat.Not(seqAuxiliary[i-1]), seqAuxiliary[i])
				addTernaryClause(s, lits[i], sat.Not(seqAuxiliary[i]), seqAuxiliary[i-1])
				addBinaryClause(s, sat.Not(lits[i]), seqAuxiliary[i])
				addBinaryClause(s, sat.Not(lits[i]), sat.Not(seqAuxiliary[i-1]))
			}
		}
	}
}
