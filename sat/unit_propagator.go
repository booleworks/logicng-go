package sat

import (
	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/normalform"
)

// A UnitPropagator is a special solver used for just propagating unit literals
// on the solver.
type UnitPropagator struct {
	*CoreSolver
	fac f.Factory
}

// NewUnitPropagator returns a new unit propagator solver.
func NewUnitPropagator(fac f.Factory) *UnitPropagator {
	return &UnitPropagator{
		CoreSolver: NewCoreSolver(DefaultConfig(), UncheckedEnqueue),
		fac:        fac,
	}
}

// Add adds a formula to the unit propagator
func (p *UnitPropagator) Add(formula f.Formula) {
	cnf := normalform.CNF(p.fac, formula)
	switch cnf.Sort() {
	case f.SortTrue:
		break
	case f.SortFalse, f.SortLiteral, f.SortOr:
		p.AddClause(p.generateClauseVector(cnf), nil)
	case f.SortAnd:
		for _, op := range p.fac.Operands(cnf) {
			p.AddClause(p.generateClauseVector(op), nil)
		}
	default:
		panic(errorx.BadFormulaSort(cnf.Sort()))

	}
}

// PropagateFormula propagates the units on the formula and returns the result.
func (p *UnitPropagator) PropagateFormula() f.Formula {
	if !p.ok || p.propagate() != nil {
		return p.fac.Falsum()
	}
	newClauses := make([]f.Formula, len(p.clauses))
	for i, clause := range p.clauses {
		newClauses[i] = p.clauseToFormula(clause)
	}
	for i := 0; i < len(p.trail); i++ {
		newClauses = append(newClauses, p.solverLiteralToFormula(p.trail[i]))
	}
	return p.fac.And(newClauses...)
}

func (p *UnitPropagator) clauseToFormula(clause *clause) f.Formula {
	literals := make([]f.Formula, 0, clause.size())
	for i := 0; i < clause.size(); i++ {
		lit := clause.get(i)
		switch p.value(lit) {
		case f.TristateTrue:
			return p.fac.Verum()
		case f.TristateUndef:
			literals = append(literals, p.solverLiteralToFormula(lit))
		}
	}
	return p.fac.Or(literals...)
}

func (p *UnitPropagator) solverLiteralToFormula(lit int32) f.Formula {
	return p.fac.Literal(p.idx2name[Vari(lit)], !Sign(lit))
}

func (p *UnitPropagator) generateClauseVector(clause f.Formula) []int32 {
	lits := f.Literals(p.fac, clause).Content()
	clauseVec := make([]int32, 0, len(lits))
	for _, lit := range lits {
		name, phase, _ := p.fac.LitNamePhase(lit)
		index := p.IdxForName(name)
		if index == -1 {
			index = p.NewVar(false, false)
			p.addName(name, index)
		}
		var litNum int32
		if phase {
			litNum = index * 2
		} else {
			litNum = (index * 2) ^ 1
		}
		clauseVec = append(clauseVec, litNum)
	}
	return clauseVec
}
