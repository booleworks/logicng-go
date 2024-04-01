package maxsat

import (
	"booleworks.com/logicng/sat"
)

type encoder struct {
	totalizer           *totalizer
	swc                 *swc
	incrementalStrategy IncrementalStrategy
}

func newEncoder() *encoder {
	return &encoder{
		incrementalStrategy: IncNone,
		totalizer:           newTotalizer(IncNone),
		swc:                 newSwc(),
	}
}

func (e *encoder) setIncremental(incremental IncrementalStrategy) {
	e.incrementalStrategy = incremental
	e.totalizer.incrementalStrategy = incremental
}

func (e *encoder) encodeAMO(s *sat.CoreSolver, lits []int32) {
	encodeLadder(s, lits)
}

func (e *encoder) encodeCardinality(s *sat.CoreSolver, lits []int32, rhs int) {
	e.totalizer.build(s, lits, rhs)
	if e.totalizer.hasEncoding {
		e.totalizer.update(s, rhs)
	}
}

func (e *encoder) updateCardinality(s *sat.CoreSolver, rhs int) {
	e.totalizer.update(s, rhs)
}

func (e *encoder) buildCardinality(s *sat.CoreSolver, lits []int32, rhs int) {
	e.totalizer.build(s, lits, rhs)
}

func (e *encoder) incUpdateCardinality(s *sat.CoreSolver, join []int32, rhs int, assumptions *[]int32) {
	if len(join) > 0 {
		e.totalizer.join(s, join, rhs)
	}
	e.totalizer.updateWithAssumptions(s, rhs, assumptions)
}

func (e *encoder) encodePB(s *sat.CoreSolver, lits *[]int32, coeffs *[]int, rhs int) {
	e.swc.encode(s, lits, coeffs, rhs)
}

func (e *encoder) updatePB(s *sat.CoreSolver, rhs int) {
	e.swc.update(s, rhs)
}

func (e *encoder) incEncodePB(s *sat.CoreSolver, lits *[]int32, coeffs *[]int, rhs int, assumptions *[]int32, size int) {
	e.swc.encodeWithAssumptions(s, lits, coeffs, rhs, assumptions, size)
}

func (e *encoder) incUpdatePB(s *sat.CoreSolver, lits []int32, coeffs []int, rhs int) {
	e.swc.updateInc(s, rhs)
	e.swc.join(s, lits, coeffs)
}

func (e *encoder) incUpdatePBAssumptions(assumptions *[]int32) {
	e.swc.updateAssumptions(assumptions)
}

func (e *encoder) hasCardEncoding() bool {
	return e.totalizer.hasEncoding
}

func (e *encoder) hasPBEncoding() bool {
	return e.swc.hasEncoding
}

func (e *encoder) lits() []int32 {
	return e.totalizer.lits()
}

func (e *encoder) outputs() []int32 {
	return e.totalizer.outputs()
}

func addUnitClause(s *sat.CoreSolver, a int32) {
	addUnitClauseWithBlocking(s, a, sat.LitUndef)
}

func addUnitClauseWithBlocking(s *sat.CoreSolver, a, blocking int32) {
	clause := []int32{a}
	if blocking != sat.LitUndef {
		clause = append(clause, blocking)
	}
	s.AddClause(clause, nil)
}

func addBinaryClause(s *sat.CoreSolver, a, b int32) {
	addBinaryClauseWithBlocking(s, a, b, sat.LitUndef)
}

func addBinaryClauseWithBlocking(s *sat.CoreSolver, a, b, blocking int32) {
	clause := []int32{a, b}
	if blocking != sat.LitUndef {
		clause = append(clause, blocking)
	}
	s.AddClause(clause, nil)
}

func addTernaryClause(s *sat.CoreSolver, a, b, c int32) {
	addTernaryClauseWithBlocking(s, a, b, c, sat.LitUndef)
}

func addTernaryClauseWithBlocking(s *sat.CoreSolver, a, b, c, blocking int32) {
	clause := []int32{a, b, c}
	if blocking != sat.LitUndef {
		clause = append(clause, blocking)
	}
	s.AddClause(clause, nil)
}

func pop[T any](slice *[]T) {
	*slice = (*slice)[:len(*slice)-1]
}
