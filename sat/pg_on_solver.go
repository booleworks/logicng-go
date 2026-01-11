package sat

import (
	"fmt"
	"slices"

	"github.com/booleworks/logicng-go/encoding"
	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/normalform"
)

type pgOnSolver struct {
	fac           f.Factory
	performNNF    bool
	variableCache map[f.Formula]*varCacheEntry
	solver        *CoreSolver
	initialPhase  bool
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

func newPGOnSolver(fac f.Factory, performNNF bool, solver *CoreSolver, initialPhase bool) *pgOnSolver {
	return &pgOnSolver{fac, performNNF, make(map[f.Formula]*varCacheEntry), solver, initialPhase}
}

func (p *pgOnSolver) addCNFToSolver(formula f.Formula, proposition f.Proposition) {
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
		p.addCNF(withoutPbc, proposition)
	} else {
		topLevelVars := p.computeTransformation(withoutPbc, proposition, true, true)
		if topLevelVars != nil {
			p.solver.AddClause(topLevelVars, proposition)
		}
	}
}

func (p *pgOnSolver) clearCache() {
	clear(p.variableCache)
}

func (p *pgOnSolver) addCNF(cnf f.Formula, proposition f.Proposition) {
	switch cnf.Sort() {
	case f.SortTrue:
		break
	case f.SortFalse, f.SortLiteral, f.SortOr:
		p.solver.AddClause(p.generateClauseSlice(f.Literals(p.fac, cnf).Content()), proposition)
	case f.SortAnd:
		ops, _ := p.fac.NaryOperands(cnf)
		for _, clause := range ops {
			p.solver.AddClause(p.generateClauseSlice(f.Literals(p.fac, clause).Content()), proposition)
		}
	default:
		panic(errorx.IllegalState("formula not in CNF"))
	}
}

func (p *pgOnSolver) computeTransformation(
	formula f.Formula, prop f.Proposition, polarity, topLevel bool,
) []int32 {
	switch fsort := formula.Sort(); fsort {
	case f.SortLiteral:
		name, phase, _ := p.fac.LiteralNamePhase(formula)
		if polarity {
			return []int32{p.solverLiteral(name, phase)}
		} else {
			return []int32{p.solverLiteral(name, phase) ^ 1}
		}
	case f.SortNot:
		op, _ := p.fac.NotOperand(formula)
		return p.computeTransformation(op, prop, !polarity, topLevel)
	case f.SortOr, f.SortAnd:
		return p.pgNary(formula, prop, polarity, topLevel)
	case f.SortImpl:
		return p.pgImpl(formula, prop, polarity, topLevel)
	case f.SortEquiv:
		return p.pgEquiv(formula, prop, polarity, topLevel)
	default:
		panic(errorx.BadFormulaSort(&fsort))
	}
}

func (p *pgOnSolver) pgImpl(formula f.Formula, prop f.Proposition, polarity, topLevel bool) []int32 {
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
		leftPgVarNeg := p.computeTransformation(left, prop, false, false)
		rightPgVarPos := p.computeTransformation(right, prop, true, false)
		return sliceVV(leftPgVarNeg, rightPgVarPos)
	} else {
		// (~left | right) => pg = (left & ~right) | pg = (left | pg) & (~right | pg)
		leftPgVarPos := p.computeTransformation(left, prop, true, topLevel)
		rightPgVarNeg := p.computeTransformation(right, prop, false, topLevel)
		if topLevel {
			if leftPgVarPos != nil {
				p.solver.AddClause(leftPgVarPos, prop)
			}
			if rightPgVarNeg != nil {
				p.solver.AddClause(rightPgVarNeg, prop)
			}
			return nil
		} else {
			p.solver.AddClause(sliceEV(pgVar, leftPgVarPos), prop)
			p.solver.AddClause(sliceEV(pgVar, rightPgVarNeg), prop)
			return []int32{pgVar ^ 1}
		}
	}
}

func (p *pgOnSolver) pgEquiv(formula f.Formula, prop f.Proposition, polarity, topLevel bool) []int32 {
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
	leftPgVarPos := p.computeTransformation(left, prop, true, false)
	leftPgVarNeg := p.computeTransformation(left, prop, false, false)
	rightPgVarPos := p.computeTransformation(right, prop, true, false)
	rightPgVarNeg := p.computeTransformation(right, prop, false, false)
	if polarity {
		// pg => (left => right) & (right => left)
		// = (pg & left => right) & (pg & right => left)
		// = (~pg | ~left | right) & (~pg | ~right | left)
		if topLevel {
			p.solver.AddClause(sliceVV(leftPgVarNeg, rightPgVarPos), prop)
			p.solver.AddClause(sliceVV(leftPgVarPos, rightPgVarNeg), prop)
			return nil
		} else {
			p.solver.AddClause(sliceEVV(pgVar^1, leftPgVarNeg, rightPgVarPos), prop)
			p.solver.AddClause(sliceEVV(pgVar^1, leftPgVarPos, rightPgVarNeg), prop)
		}
	} else {
		// (left => right) & (right => left) => pg
		// = ~(left => right) | ~(right => left) | pg
		// = left & ~right | right & ~left | pg
		// = (left | right | pg) & (~right | ~left | pg)
		if topLevel {
			p.solver.AddClause(sliceVV(leftPgVarPos, rightPgVarPos), prop)
			p.solver.AddClause(sliceVV(leftPgVarNeg, rightPgVarNeg), prop)
			return nil
		} else {
			p.solver.AddClause(sliceEVV(pgVar, leftPgVarPos, rightPgVarPos), prop)
			p.solver.AddClause(sliceEVV(pgVar, leftPgVarNeg, rightPgVarNeg), prop)
		}
	}
	if polarity {
		return []int32{pgVar}
	} else {
		return []int32{pgVar ^ 1}
	}
}

func (p *pgOnSolver) pgNary(formula f.Formula, prop f.Proposition, polarity, topLevel bool) []int32 {
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
				opPgVars := p.computeTransformation(op, prop, true, topLevel)
				if topLevel {
					if opPgVars != nil {
						p.solver.AddClause(opPgVars, prop)
					}
				} else {
					p.solver.AddClause(sliceEV(pgVar^1, opPgVars), prop)
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
				opPgVars := p.computeTransformation(op, prop, false, false)
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
				opPgVars := p.computeTransformation(op, prop, true, false)
				singleClause = append(singleClause, opPgVars...)
			}
			return singleClause
		} else {
			// (v1 | ... | vk) => pg = (~v1 | pg) & ... & (~vk | pg)
			for _, op := range ops {
				opPgVars := p.computeTransformation(op, prop, false, topLevel)
				if topLevel {
					if opPgVars != nil {
						p.solver.AddClause(opPgVars, prop)
					}
				} else {
					p.solver.AddClause(sliceEV(pgVar, opPgVars), prop)
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

func (p *pgOnSolver) generateClauseSlice(literals []f.Literal) []int32 {
	clause := make([]int32, len(literals))
	slices.Sort(literals)
	for i, lit := range literals {
		name, phase, _ := p.fac.LitNamePhase(lit)
		clause[i] = p.solverLiteral(name, phase)
	}
	return clause
}

func (p *pgOnSolver) solverLiteral(name string, phase bool) int32 {
	index := p.solver.IdxForName(name)
	if index == -1 {
		index = p.solver.NewVar(!p.initialPhase, true)
		p.solver.addName(name, index)
	}
	if phase {
		return index * 2
	} else {
		return (index * 2) ^ 1
	}
}

func (p *pgOnSolver) newSolverVariable() int32 {
	index := p.solver.NewVar(!p.initialPhase, true)
	name := fmt.Sprintf("%sSOLVER_%d", f.AuxCNF, index)
	p.solver.addName(name, index)
	return index * 2
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
