package maxsat

import (
	"github.com/booleworks/logicng-go/encoding"
	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/normalform"
)

type pgOnSolver struct {
	fac           f.Factory
	performNNF    bool
	variableCache map[f.Formula]*varCacheEntry
	solver        algorithm
}

type varCacheEntry struct {
	pgVar             int32
	posPolarityCached bool
	negPolarityCached bool
}

func newCacheEntry(pgVar int32) *varCacheEntry {
	return &varCacheEntry{pgVar, false, false}
}

func (e *varCacheEntry) setPolarityCached(polarity bool) bool {
	var wasCached bool
	if polarity {
		wasCached = e.posPolarityCached
		e.posPolarityCached = true
	} else {
		wasCached = e.negPolarityCached
		e.negPolarityCached = true
	}
	return wasCached
}

func newPGOnSolver(fac f.Factory, performNNF bool, solver algorithm) *pgOnSolver {
	return &pgOnSolver{fac, performNNF, make(map[f.Formula]*varCacheEntry), solver}
}

func (p *pgOnSolver) addCNFToSolver(formula f.Formula, weight int) {
	workingFormula := formula
	if p.performNNF {
		workingFormula = normalform.NNF(p.fac, formula)
	}
	containsPbc := encoding.ContainsPBC(p.fac, workingFormula)
	var withoutPbc f.Formula
	if !p.performNNF && containsPbc {
		withoutPbc = normalform.NNF(p.fac, workingFormula)
	} else {
		withoutPbc = workingFormula
	}
	if normalform.IsCNF(p.fac, withoutPbc) {
		p.addCNF(withoutPbc, weight)
	} else {
		topLevelVars := p.computeTransformation(withoutPbc, weight, true, true)
		if topLevelVars != nil {
			p.solver.addClauseVec(topLevelVars, weight)
		}
	}
}

func (p *pgOnSolver) clearCache() {
	clear(p.variableCache)
}

func (p *pgOnSolver) addCNF(cnf f.Formula, weight int) {
	switch cnf.Sort() {
	case f.SortTrue:
		break
	case f.SortFalse, f.SortLiteral, f.SortOr:
		p.solver.addClause(cnf, weight)
	case f.SortAnd:
		ops, _ := p.fac.NaryOperands(cnf)
		for _, clause := range ops {
			p.solver.addClause(clause, weight)
		}
	default:
		panic(errorx.IllegalState("formula not in CNF"))
	}
}

func (p *pgOnSolver) computeTransformation(
	formula f.Formula, weight int, polarity, topLevel bool,
) []int32 {
	switch fsort := formula.Sort(); fsort {
	case f.SortLiteral:
		lit, _ := formula.AsLiteral()
		if polarity {
			return []int32{p.solver.literal(lit)}
		} else {
			return []int32{p.solver.literal(lit) ^ 1}
		}
	case f.SortNot:
		op, _ := p.fac.NotOperand(formula)
		return p.computeTransformation(op, weight, !polarity, topLevel)
	case f.SortOr, f.SortAnd:
		return p.pgNary(formula, weight, polarity, topLevel)
	case f.SortImpl:
		return p.pgImpl(formula, weight, polarity, topLevel)
	case f.SortEquiv:
		return p.pgEquiv(formula, weight, polarity, topLevel)
	default:
		panic(errorx.BadFormulaSort(&fsort))
	}
}

func (p *pgOnSolver) pgImpl(formula f.Formula, weight int, polarity, topLevel bool) []int32 {
	skipPg := polarity || topLevel
	var wasCached bool
	var pgVar int32
	if skipPg {
		wasCached, pgVar = false, -1
	} else {
		wasCached, pgVar = p.getPGVar(formula, polarity)
	}
	if wasCached {
		if polarity {
			return []int32{pgVar}
		} else {
			return []int32{pgVar ^ 1}
		}
	}

	left, right, _ := p.fac.BinaryLeftRight(formula)

	if polarity {
		// pg => (~left | right) = ~pg | ~left | right
		// Speed-Up: Skip pg var
		leftPgVarNeg := p.computeTransformation(left, weight, false, false)
		rightPgVarPos := p.computeTransformation(right, weight, true, false)
		return sliceVV(leftPgVarNeg, rightPgVarPos)
	} else {
		// (~left | right) => pg = (left & ~right) | pg = (left | pg) & (~right | pg)
		leftPgVarPos := p.computeTransformation(left, weight, true, topLevel)
		rightPgVarNeg := p.computeTransformation(right, weight, false, topLevel)
		if topLevel {
			if leftPgVarPos != nil {
				p.solver.addClauseVec(leftPgVarPos, weight)
			}
			if rightPgVarNeg != nil {
				p.solver.addClauseVec(rightPgVarNeg, weight)
			}
			return nil
		} else {
			p.solver.addClauseVec(sliceEV(pgVar, leftPgVarPos), weight)
			p.solver.addClauseVec(sliceEV(pgVar, rightPgVarNeg), weight)
			return []int32{pgVar ^ 1}
		}
	}
}

func (p *pgOnSolver) pgEquiv(formula f.Formula, weight int, polarity, topLevel bool) []int32 {
	var wasCached bool
	var pgVar int32
	if topLevel {
		wasCached, pgVar = false, -1
	} else {
		wasCached, pgVar = p.getPGVar(formula, polarity)
	}
	if wasCached {
		if polarity {
			return []int32{pgVar}
		} else {
			return []int32{pgVar ^ 1}
		}
	}

	left, right, _ := p.fac.BinaryLeftRight(formula)
	leftPgVarPos := p.computeTransformation(left, weight, true, false)
	leftPgVarNeg := p.computeTransformation(left, weight, false, false)
	rightPgVarPos := p.computeTransformation(right, weight, true, false)
	rightPgVarNeg := p.computeTransformation(right, weight, false, false)
	if polarity {
		// pg => (left => right) & (right => left)
		// = (pg & left => right) & (pg & right => left)
		// = (~pg | ~left | right) & (~pg | ~right | left)
		if topLevel {
			p.solver.addClauseVec(sliceVV(leftPgVarNeg, rightPgVarPos), weight)
			p.solver.addClauseVec(sliceVV(leftPgVarPos, rightPgVarNeg), weight)
			return nil
		} else {
			p.solver.addClauseVec(sliceEVV(pgVar^1, leftPgVarNeg, rightPgVarPos), weight)
			p.solver.addClauseVec(sliceEVV(pgVar^1, leftPgVarPos, rightPgVarNeg), weight)
		}
	} else {
		// (left => right) & (right => left) => pg
		// = ~(left => right) | ~(right => left) | pg
		// = left & ~right | right & ~left | pg
		// = (left | right | pg) & (~right | ~left | pg)
		if topLevel {
			p.solver.addClauseVec(sliceVV(leftPgVarPos, rightPgVarPos), weight)
			p.solver.addClauseVec(sliceVV(leftPgVarNeg, rightPgVarNeg), weight)
			return nil
		} else {
			p.solver.addClauseVec(sliceEVV(pgVar, leftPgVarPos, rightPgVarPos), weight)
			p.solver.addClauseVec(sliceEVV(pgVar, leftPgVarNeg, rightPgVarNeg), weight)
		}
	}
	if polarity {
		return []int32{pgVar}
	} else {
		return []int32{pgVar ^ 1}
	}
}

func (p *pgOnSolver) pgNary(formula f.Formula, weight int, polarity, topLevel bool) []int32 {
	skipPg := topLevel || formula.Sort() == f.SortAnd && !polarity || formula.Sort() == f.SortOr && polarity

	var wasCached bool
	var pgVar int32
	if skipPg {
		wasCached, pgVar = false, -1
	} else {
		wasCached, pgVar = p.getPGVar(formula, polarity)
	}
	if wasCached {
		if polarity {
			return []int32{pgVar}
		} else {
			return []int32{pgVar ^ 1}
		}
	}

	ops, _ := p.fac.NaryOperands(formula)

	switch fsort := formula.Sort(); fsort {
	case f.SortAnd:
		if polarity {
			// pg => (v1 & ... & vk) = (~pg | v1) & ... & (~pg | vk)
			for _, op := range ops {
				opPgVars := p.computeTransformation(op, weight, true, topLevel)
				if topLevel {
					if opPgVars != nil {
						p.solver.addClauseVec(opPgVars, weight)
					}
				} else {
					p.solver.addClauseVec(sliceEV(pgVar^1, opPgVars), weight)
				}
			}
			if topLevel {
				return nil
			}
		} else {
			// (v1 & ... & vk) => pg = ~v1 | ... | ~vk | pg
			// Speed-Up: Skip pg var
			singleClause := make([]int32, 0, len(ops))
			for _, op := range ops {
				opPgVars := p.computeTransformation(op, weight, false, false)
				singleClause = append(singleClause, opPgVars...)
			}
			return singleClause
		}
	case f.SortOr:
		if polarity {
			// pg => (v1 | ... | vk) = ~pg | v1 | ... | vk
			// Speed-Up: Skip pg var
			singleClause := make([]int32, 0, len(ops))
			for _, op := range ops {
				opPgVars := p.computeTransformation(op, weight, true, false)
				singleClause = append(singleClause, opPgVars...)
			}
			return singleClause
		} else {
			// (v1 | ... | vk) => pg = (~v1 | pg) & ... & (~vk | pg)
			for _, op := range ops {
				opPgVars := p.computeTransformation(op, weight, false, topLevel)
				if topLevel {
					if opPgVars != nil {
						p.solver.addClauseVec(opPgVars, weight)
					}
				} else {
					p.solver.addClauseVec(sliceEV(pgVar, opPgVars), weight)
				}
			}
			if topLevel {
				return nil
			}
		}
	default:
		panic(errorx.BadFormulaSort(&fsort))
	}
	if polarity {
		return []int32{pgVar}
	} else {
		return []int32{pgVar ^ 1}
	}
}

func (p *pgOnSolver) getPGVar(formula f.Formula, polarity bool) (bool, int32) {
	entry, ok := p.variableCache[formula]
	if !ok {
		entry = newCacheEntry(p.newSolverVariable())
		p.variableCache[formula] = entry
	}
	wasCached := entry.setPolarityCached(polarity)
	pgVar := entry.pgVar
	return wasCached, pgVar
}

func (p *pgOnSolver) newSolverVariable() int32 {
	return p.solver.newVar() * 2
}

func sliceVV(a, b []int32) []int32 {
	result := make([]int32, len(a)+len(b))
	copy(result, a)
	offset := len(a)
	for i, elem := range b {
		result[i+offset] = elem
	}
	return result
}

func sliceEV(elt int32, a []int32) []int32 {
	result := make([]int32, 1+len(a))
	result[0] = elt
	for i, elem := range a {
		result[i+1] = elem
	}
	return result
}

func sliceEVV(elt int32, a, b []int32) []int32 {
	result := make([]int32, 1+len(a)+len(b))
	result[0] = elt
	for i, elem := range a {
		result[i+1] = elem
	}
	offset := 1 + len(a)
	for i, elem := range b {
		result[i+offset] = elem
	}
	return result
}
