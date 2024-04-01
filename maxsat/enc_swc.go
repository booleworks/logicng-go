package maxsat

import (
	"math"

	"github.com/booleworks/logicng-go/errorx"
	"github.com/booleworks/logicng-go/sat"
)

type swc struct {
	pbOutlits          []int32
	unitLits           []int32
	unitCoeffs         []int
	currentPbRhs       int
	currentLitBlocking int32
	seqAuxiliaryInc    [][]int32
	litsInc            []int32
	coeffsInc          []int
	hasEncoding        bool
}

func newSwc() *swc {
	return &swc{
		currentPbRhs:       -1,
		currentLitBlocking: sat.LitUndef,
		pbOutlits:          []int32{},
		unitLits:           []int32{},
		unitCoeffs:         []int{},
		seqAuxiliaryInc:    [][]int32{},
		litsInc:            []int32{},
		coeffsInc:          []int{},
		hasEncoding:        false,
	}
}

func (m *swc) updateAssumptions(assumptions *[]int32) {
	*assumptions = append(*assumptions, sat.Not(m.currentLitBlocking))
	for i := 0; i < len(m.unitLits); i++ {
		*assumptions = append(*assumptions, sat.Not(m.unitLits[i]))
	}
}

func (m *swc) encode(s *sat.CoreSolver, lits *[]int32, coeffs *[]int, rhs int) {
	if rhs == math.MaxInt {
		panic(errorx.IllegalState("overflow in the encoding"))
	}
	m.hasEncoding = false
	simpLits := make([]int32, len(*lits))
	copy(simpLits, *lits)
	simpCoeffs := make([]int, len(*coeffs))
	copy(simpCoeffs, *coeffs)
	*lits = []int32{}
	*coeffs = []int{}
	for i := 0; i < len(simpLits); i++ {
		if simpCoeffs[i] <= rhs {
			*lits = append(*lits, simpLits[i])
			*coeffs = append(*coeffs, simpCoeffs[i])
		} else {
			addUnitClause(s, sat.Not(simpLits[i]))
		}
	}
	if len(*lits) == 1 {
		addUnitClause(s, sat.Not((*lits)[0]))
		return
	}
	if len(*lits) == 0 {
		return
	}
	n := len(*lits)
	seqAuxiliary := make([][]int32, n+1)
	for i := 0; i < n+1; i++ {
		seqAuxiliary[i] = make([]int32, rhs+1)
		for j := 0; j < rhs+1; j++ {
			seqAuxiliary[i][j] = -1
		}
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= rhs; j++ {
			seqAuxiliary[i][j] = sat.MkLit(s.NVars(), false)
			newSatVariable(s)
		}
	}
	for i := 1; i <= rhs; i++ {
		m.pbOutlits = append(m.pbOutlits, seqAuxiliary[n][i])
	}
	for i := 1; i <= n; i++ {
		wi := (*coeffs)[i-1]
		for j := 1; j <= rhs; j++ {
			if i >= 2 && i <= n && j <= rhs {
				addBinaryClause(s, sat.Not(seqAuxiliary[i-1][j]), seqAuxiliary[i][j])
			}
			if i <= n && j <= wi {
				addBinaryClause(s, sat.Not((*lits)[i-1]), seqAuxiliary[i][j])
			}
			if i >= 2 && i <= n && j <= rhs-wi {
				addTernaryClause(s, sat.Not(seqAuxiliary[i-1][j]), sat.Not((*lits)[i-1]), seqAuxiliary[i][j+wi])
			}
		}
		if i >= 2 {
			addBinaryClause(s, sat.Not(seqAuxiliary[i-1][rhs+1-wi]), sat.Not((*lits)[i-1]))
		}
	}
	m.currentPbRhs = rhs
	m.hasEncoding = true
}

func (m *swc) encodeWithAssumptions(
	s *sat.CoreSolver, lits *[]int32, coeffs *[]int, rhs int, assumptions *[]int32, size int,
) {
	if rhs == math.MaxInt {
		panic(errorx.IllegalState("overflow in the encoding"))
	}
	m.hasEncoding = false
	simpLits := make([]int32, len(*lits))
	copy(simpLits, *lits)
	simpCoeffs := make([]int, len(*coeffs))
	copy(simpCoeffs, *coeffs)
	*lits = []int32{}
	*coeffs = []int{}
	simpUnitLits := make([]int32, len(m.unitLits))
	copy(simpUnitLits, m.unitLits)
	simpUnitCoeffs := make([]int, len(m.unitCoeffs))
	copy(simpUnitCoeffs, m.unitCoeffs)
	m.unitLits = []int32{}
	m.unitCoeffs = []int{}

	for i := 0; i < len(simpUnitLits); i++ {
		if simpUnitCoeffs[i] <= rhs {
			*lits = append(*lits, simpUnitLits[i])
			*coeffs = append(*coeffs, simpUnitCoeffs[i])
		} else {
			m.unitLits = append(m.unitLits, simpUnitLits[i])
			m.unitCoeffs = append(m.unitCoeffs, simpUnitCoeffs[i])
		}
	}
	for i := 0; i < len(simpLits); i++ {
		if simpCoeffs[i] <= rhs {
			*lits = append(*lits, simpLits[i])
			*coeffs = append(*coeffs, simpCoeffs[i])
		} else {
			m.unitLits = append(m.unitLits, simpLits[i])
			m.unitCoeffs = append(m.unitCoeffs, simpCoeffs[i])
		}
	}
	if len(*lits) == 1 {
		for i := 0; i < len(m.unitLits); i++ {
			*assumptions = append(*assumptions, sat.Not(m.unitLits[i]))
		}
		m.unitLits = append(m.unitLits, (*lits)[0])
		m.unitCoeffs = append(m.unitCoeffs, (*coeffs)[0])
		return
	}
	if len(*lits) == 0 {
		for i := 0; i < len(m.unitLits); i++ {
			*assumptions = append(*assumptions, sat.Not(m.unitLits[i]))
		}
		return
	}
	n := len(*lits)
	m.seqAuxiliaryInc = make([][]int32, size+1)
	for i := 0; i <= n; i++ {
		m.seqAuxiliaryInc[i] = make([]int32, rhs+1)
		for j := 0; j < rhs+1; j++ {
			m.seqAuxiliaryInc[i][j] = -1
		}
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= rhs; j++ {
			m.seqAuxiliaryInc[i][j] = sat.MkLit(s.NVars(), false)
			newSatVariable(s)
		}
	}
	blocking := sat.MkLit(s.NVars(), false)
	newSatVariable(s)
	m.currentLitBlocking = blocking
	*assumptions = append(*assumptions, sat.Not(blocking))
	for i := 1; i <= n; i++ {
		wi := (*coeffs)[i-1]
		for j := 1; j <= rhs; j++ {
			if i >= 2 && i <= n {
				addBinaryClause(s, sat.Not(m.seqAuxiliaryInc[i-1][j]), m.seqAuxiliaryInc[i][j])
			}
			if i <= n && j <= wi {
				addBinaryClause(s, sat.Not((*lits)[i-1]), m.seqAuxiliaryInc[i][j])
			}
			if i >= 2 && i <= n && j <= rhs-wi {
				addTernaryClause(s,
					sat.Not(m.seqAuxiliaryInc[i-1][j]),
					sat.Not((*lits)[i-1]),
					m.seqAuxiliaryInc[i][j+wi],
				)
			}
		}
		if i >= 2 {
			addBinaryClauseWithBlocking(s, sat.Not(m.seqAuxiliaryInc[i-1][rhs+1-wi]), sat.Not((*lits)[i-1]), blocking)
		}
	}
	for i := 0; i < len(m.unitLits); i++ {
		*assumptions = append(*assumptions, sat.Not(m.unitLits[i]))
	}
	m.currentPbRhs = rhs
	m.hasEncoding = true
	m.litsInc = make([]int32, len(*lits))
	copy(m.litsInc, *lits)
	m.coeffsInc = make([]int, len(*coeffs))
	copy(m.coeffsInc, *coeffs)
}

func (m *swc) update(s *sat.CoreSolver, rhs int) {
	for i := rhs; i < m.currentPbRhs; i++ {
		addUnitClause(s, sat.Not(m.pbOutlits[i]))
	}
	m.currentPbRhs = rhs
}

func (m *swc) updateInc(s *sat.CoreSolver, rhs int) {
	if m.currentLitBlocking != sat.LitUndef {
		addUnitClause(s, m.currentLitBlocking)
	}
	n := len(m.litsInc)
	offset := m.currentPbRhs + 1
	for i := 1; i <= n; i++ {
		for j := offset; j <= rhs; j++ {
			m.seqAuxiliaryInc[i] = append(m.seqAuxiliaryInc[i], sat.LitUndef)
		}
	}
	for i := 1; i <= n; i++ {
		for j := offset; j <= rhs; j++ {
			m.seqAuxiliaryInc[i][j] = sat.MkLit(s.NVars(), false)
			newSatVariable(s)
		}
	}
	m.currentLitBlocking = sat.MkLit(s.NVars(), false)
	newSatVariable(s)
	for i := 1; i <= n; i++ {
		wi := m.coeffsInc[i-1]
		for j := 1; j <= rhs; j++ {
			if i >= 2 && i <= n && j <= rhs && j >= offset {
				addBinaryClause(s, sat.Not(m.seqAuxiliaryInc[i-1][j]), m.seqAuxiliaryInc[i][j])
			}
			if i >= 2 && i <= n && j <= rhs-wi && j >= offset-wi {
				addTernaryClause(s,
					sat.Not(m.seqAuxiliaryInc[i-1][j]),
					sat.Not(m.litsInc[i-1]),
					m.seqAuxiliaryInc[i][j+wi],
				)
			}
		}
		if i >= 2 {
			addBinaryClauseWithBlocking(s,
				sat.Not(m.seqAuxiliaryInc[i-1][rhs+1-wi]),
				sat.Not(m.litsInc[i-1]),
				m.currentLitBlocking,
			)
		}
	}
	m.currentPbRhs = rhs
}

func (m *swc) join(s *sat.CoreSolver, lits []int32, coeffs []int) {
	rhs := m.currentPbRhs
	if rhs == math.MaxInt {
		panic(errorx.IllegalState("overflow in the encoding"))
	}
	simpUnitLits := make([]int32, len(m.unitLits))
	copy(simpUnitLits, m.unitLits)
	simpUnitCoeffs := make([]int, len(m.unitCoeffs))
	copy(simpUnitCoeffs, m.unitCoeffs)
	m.unitLits = []int32{}
	m.unitCoeffs = []int{}
	lhsJoin := len(m.litsInc)
	for i := 0; i < len(simpUnitLits); i++ {
		if simpUnitCoeffs[i] <= rhs {
			m.litsInc = append(m.litsInc, simpUnitLits[i])
			m.coeffsInc = append(m.coeffsInc, simpUnitCoeffs[i])
		} else {
			m.unitLits = append(m.unitLits, simpUnitLits[i])
			m.unitCoeffs = append(m.unitCoeffs, simpUnitCoeffs[i])
		}
	}
	for i := 0; i < len(lits); i++ {
		if coeffs[i] <= rhs {
			m.litsInc = append(m.litsInc, lits[i])
			m.coeffsInc = append(m.coeffsInc, coeffs[i])
		} else {
			m.unitLits = append(m.unitLits, lits[i])
			m.unitCoeffs = append(m.unitCoeffs, coeffs[i])
		}
	}
	if len(m.litsInc) == lhsJoin {
		return
	}
	n := len(m.litsInc)

	for i := lhsJoin + 1; i <= n; i++ {
		m.seqAuxiliaryInc[i] = make([]int32, rhs+1)
		for j := 0; j < rhs+1; j++ {
			m.seqAuxiliaryInc[i][j] = -1
		}
	}
	for i := lhsJoin + 1; i <= n; i++ {
		for j := 1; j <= rhs; j++ {
			m.seqAuxiliaryInc[i][j] = sat.MkLit(s.NVars(), false)
			newSatVariable(s)
		}
	}
	for i := lhsJoin; i <= n; i++ {
		wi := m.coeffsInc[i-1]
		for j := 1; j <= rhs; j++ {
			addBinaryClause(s, sat.Not(m.seqAuxiliaryInc[i-1][j]), m.seqAuxiliaryInc[i][j])
			if j <= wi {
				addBinaryClause(s, sat.Not(m.litsInc[i-1]), m.seqAuxiliaryInc[i][j])
			}
			if j <= rhs-wi {
				addTernaryClause(s,
					sat.Not(m.seqAuxiliaryInc[i-1][j]),
					sat.Not(m.litsInc[i-1]),
					m.seqAuxiliaryInc[i][j+wi],
				)
			}
		}
		if i > lhsJoin {
			addBinaryClauseWithBlocking(s,
				sat.Not(m.seqAuxiliaryInc[i-1][rhs+1-wi]),
				sat.Not(m.litsInc[i-1]),
				m.currentLitBlocking,
			)
		}
	}
}
