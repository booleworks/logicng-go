package maxsat

import "booleworks.com/logicng/sat"

type totalizer struct {
	totalizerIterativeLeft   [][]int32
	totalizerIterativeRight  [][]int32
	totalizerIterativeOutput [][]int32
	totalizerIterativeRhs    []int
	blocking                 int32
	cardinalityOutlits       []int32
	cardinalityInlits        []int32
	incrementalStrategy      IncrementalStrategy
	currentCardinalityRhs    int
	joinMode                 bool
	ilits                    []int32
	hasEncoding              bool
}

func newTotalizer(strategy IncrementalStrategy) *totalizer {
	return &totalizer{
		blocking:                 sat.LitUndef,
		joinMode:                 false,
		currentCardinalityRhs:    -1,
		incrementalStrategy:      strategy,
		totalizerIterativeLeft:   [][]int32{},
		totalizerIterativeRight:  [][]int32{},
		totalizerIterativeOutput: [][]int32{},
		totalizerIterativeRhs:    []int{},
		cardinalityInlits:        []int32{},
		cardinalityOutlits:       []int32{},
		ilits:                    []int32{},
		hasEncoding:              false,
	}
}

func (t *totalizer) update(s *sat.CoreSolver, rhs int) {
	t.updateWithAssumptions(s, rhs, &[]int32{})
}

func (t *totalizer) join(s *sat.CoreSolver, lits []int32, rhs int) {
	leftCardinalityOutlits := make([]int32, len(t.cardinalityOutlits))
	copy(leftCardinalityOutlits, t.cardinalityOutlits)
	oldCardinality := t.currentCardinalityRhs
	if len(lits) > 1 {
		t.build(s, lits, min(rhs, len(lits)))
	} else {
		t.cardinalityOutlits = []int32{lits[0]}
	}
	rightCardinalityOutlits := make([]int32, len(t.cardinalityOutlits))
	copy(rightCardinalityOutlits, t.cardinalityOutlits)
	t.cardinalityOutlits = []int32{}
	for i := 0; i < len(leftCardinalityOutlits)+len(rightCardinalityOutlits); i++ {
		p := sat.MkLit(s.NVars(), false)
		newSatVariable(s)
		t.cardinalityOutlits = append(t.cardinalityOutlits, p)
	}
	t.currentCardinalityRhs = rhs
	t.adder(s, leftCardinalityOutlits, rightCardinalityOutlits, t.cardinalityOutlits)
	t.currentCardinalityRhs = oldCardinality
	for i := 0; i < len(lits); i++ {
		t.ilits = append(t.ilits, lits[i])
	}
}

func (t *totalizer) updateWithAssumptions(s *sat.CoreSolver, rhs int, assumptions *[]int32) {
	switch t.incrementalStrategy {
	case IncNone:
		for i := rhs; i < len(t.cardinalityOutlits); i++ {
			addUnitClause(s, sat.Not(t.cardinalityOutlits[i]))
		}
	case IncIterative:
		t.incremental(s, rhs)
		*assumptions = make([]int32, len(t.cardinalityOutlits)-rhs)
		for i := rhs; i < len(t.cardinalityOutlits); i++ {
			(*assumptions)[i-rhs] = sat.Not(t.cardinalityOutlits[i])
		}
	}
}

func (t *totalizer) build(s *sat.CoreSolver, lits []int32, rhs int) {
	t.cardinalityOutlits = []int32{}
	t.hasEncoding = false
	if rhs == 0 {
		for i := 0; i < len(lits); i++ {
			addUnitClause(s, sat.Not(lits[i]))
		}
		return
	}
	if t.incrementalStrategy == IncNone && rhs == len(lits) {
		return
	}
	if rhs == len(lits) && !t.joinMode {
		return
	}
	for i := 0; i < len(lits); i++ {
		p := sat.MkLit(s.NVars(), false)
		newSatVariable(s)
		t.cardinalityOutlits = append(t.cardinalityOutlits, p)
	}
	t.cardinalityInlits = make([]int32, len(lits))
	copy(t.cardinalityInlits, lits)
	t.currentCardinalityRhs = rhs
	t.toCnf(s, t.cardinalityOutlits)
	if !t.joinMode {
		t.joinMode = true
	}
	t.hasEncoding = true
	t.ilits = make([]int32, len(lits))
	copy(t.ilits, lits)
}

func (t *totalizer) toCnf(s *sat.CoreSolver, lits []int32) {
	var left []int32
	var right []int32
	split := len(lits) / 2
	for i := 0; i < len(lits); i++ {
		if i < split {
			if split == 1 {
				left = append(left, t.cardinalityInlits[len(t.cardinalityInlits)-1])
				pop(&t.cardinalityInlits)
			} else {
				p := sat.MkLit(s.NVars(), false)
				newSatVariable(s)
				left = append(left, p)
			}
		} else {
			if len(lits)-split == 1 {
				right = append(right, t.cardinalityInlits[len(t.cardinalityInlits)-1])
				pop(&t.cardinalityInlits)
			} else {
				p := sat.MkLit(s.NVars(), false)
				newSatVariable(s)
				right = append(right, p)
			}
		}
	}
	t.adder(s, left, right, lits)
	if len(left) > 1 {
		t.toCnf(s, left)
	}
	if len(right) > 1 {
		t.toCnf(s, right)
	}
}

func (t *totalizer) adder(s *sat.CoreSolver, left, right, output []int32) {
	if t.incrementalStrategy == IncIterative {
		t.totalizerIterativeLeft = append(t.totalizerIterativeLeft, left)
		t.totalizerIterativeRight = append(t.totalizerIterativeRight, right)
		t.totalizerIterativeOutput = append(t.totalizerIterativeOutput, output)
		t.totalizerIterativeRhs = append(t.totalizerIterativeRhs, t.currentCardinalityRhs)
	}
	for i := 0; i <= len(left); i++ {
		for j := 0; j <= len(right); j++ {
			if i == 0 && j == 0 {
				continue
			}
			if i+j > t.currentCardinalityRhs+1 {
				continue
			}
			if i == 0 {
				addBinaryClauseWithBlocking(s, sat.Not(right[j-1]), output[j-1], t.blocking)
			} else if j == 0 {
				addBinaryClauseWithBlocking(s, sat.Not(left[i-1]), output[i-1], t.blocking)
			} else {
				addTernaryClauseWithBlocking(s, sat.Not(left[i-1]), sat.Not(right[j-1]), output[i+j-1], t.blocking)
			}
		}
	}
}

func (t *totalizer) incremental(s *sat.CoreSolver, rhs int) {
	for z := 0; z < len(t.totalizerIterativeRhs); z++ {
		for i := 0; i <= len(t.totalizerIterativeLeft[z]); i++ {
			for j := 0; j <= len(t.totalizerIterativeRight[z]); j++ {
				if i == 0 && j == 0 {
					continue
				}
				if i+j > rhs+1 || i+j <= t.totalizerIterativeRhs[z]+1 {
					continue
				}
				if i == 0 {
					addBinaryClause(s, sat.Not(t.totalizerIterativeRight[z][j-1]), t.totalizerIterativeOutput[z][j-1])
				} else if j == 0 {
					addBinaryClause(s, sat.Not(t.totalizerIterativeLeft[z][i-1]), t.totalizerIterativeOutput[z][i-1])
				} else {
					addTernaryClause(
						s, sat.Not(t.totalizerIterativeLeft[z][i-1]),
						sat.Not(t.totalizerIterativeRight[z][j-1]),
						t.totalizerIterativeOutput[z][i+j-1],
					)
				}
			}
		}
		t.totalizerIterativeRhs[z] = rhs
	}
}

func (t *totalizer) lits() []int32 {
	return t.ilits
}

func (t *totalizer) outputs() []int32 {
	return t.cardinalityOutlits
}
