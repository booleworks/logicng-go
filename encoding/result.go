package encoding

import (
	f "github.com/booleworks/logicng-go/formula"
)

// A Result is used to abstract CNF encodings for cardinality and
// pseudo-Boolean constraints.  It provides methods for adding clauses,
// creating new auxiliary variables and accessing the result.
//
// In LogicNG there are two implementations of this abstraction: one backed by
// a formula factory where the result is a list of clauses, and one backed by
// a SAT solver, where the clauses are added directly to the solver.
type Result interface {
	AddClause(literals ...f.Literal)
	NewAuxVar(sort f.AuxVarSort) f.Variable
	Factory() f.Factory
	Formulas() []f.Formula
}

// FormulaEncoding implements a Result backed by a formula factory.
// The resulting CNF is therefore stored as a list of clauses in the Result
// field.
type FormulaEncoding struct {
	fac    f.Factory
	result []f.Formula
}

// ResultForFormula creates a new encoding result backed by the given formula
// factory.
func ResultForFormula(fac f.Factory) *FormulaEncoding {
	return &FormulaEncoding{fac, make([]f.Formula, 0)}
}

// AddClause adds a set of literals as a clause to the encoding result.
func (r *FormulaEncoding) AddClause(literals ...f.Literal) {
	r.result = append(r.result, r.fac.Clause(literals...))
}

// NewAuxVar returns a new auxiliary variable of the given sort from the
// formula factory.
func (r *FormulaEncoding) NewAuxVar(sort f.AuxVarSort) f.Variable {
	return r.fac.NewAuxVar(sort)
}

// Factory returns the backing formula factory from the encoding result.
func (r *FormulaEncoding) Factory() f.Factory {
	return r.fac
}

// Formulas returns the encoding as a list of formulas.
func (r *FormulaEncoding) Formulas() []f.Formula {
	return r.result
}
