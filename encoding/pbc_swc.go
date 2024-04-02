package encoding

import (
	f "github.com/booleworks/logicng-go/formula"
)

func encodePBCSWC(result Result, lits []f.Literal, coeffs []int, rhs int) {
	fac := result.Factory()
	n := len(lits)
	seqAuxiliary := make([][]f.Literal, n+1)
	for i := 0; i < n+1; i++ {
		seqAuxiliary[i] = make([]f.Literal, rhs+1)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= rhs; j++ {
			seqAuxiliary[i][j] = result.NewAuxVar(f.AuxPBC).AsLiteral()
		}
	}
	for i := 1; i <= n; i++ {
		wi := coeffs[i-1]
		for j := 1; j <= rhs; j++ {
			if i >= 2 && i <= n && j <= rhs {
				result.AddClause(seqAuxiliary[i-1][j].Negate(fac), seqAuxiliary[i][j])
			}
			if i <= n && j <= wi {
				result.AddClause(lits[i-1].Negate(fac), seqAuxiliary[i][j])
			}
			if i >= 2 && i <= n && j <= rhs-wi {
				result.AddClause(seqAuxiliary[i-1][j].Negate(fac), lits[i-1].Negate(fac), seqAuxiliary[i][j+wi])
			}
		}
		if i >= 2 {
			result.AddClause(seqAuxiliary[i-1][rhs+1-wi].Negate(fac), lits[i-1].Negate(fac))
		}
	}
}
