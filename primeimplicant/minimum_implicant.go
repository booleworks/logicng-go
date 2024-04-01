package primeimplicant

import (
	"github.com/booleworks/logicng-go/errorx"
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/normalform"
	"github.com/booleworks/logicng-go/sat"
	"github.com/booleworks/logicng-go/transformation"
)

// Minimum computes a minimum-size prime implicant for the given formula.  If
// the formula is unsatisfiable it returns an error.
func Minimum(fac f.Factory, formula f.Formula) ([]f.Literal, error) {
	nnf := normalform.NNF(fac, formula)
	newVar2oldLit := make(map[f.Variable]f.Literal)
	substitution := make(map[f.Literal]f.Literal)
	literals := f.Literals(fac, nnf)
	newLiterals := make([]f.Literal, literals.Size())

	for i, literal := range literals.Content() {
		name, phase, _ := fac.LitNamePhase(literal)
		var polarity string
		if phase {
			polarity = pos
		} else {
			polarity = neg
		}
		newVar := fac.Var(name + polarity)
		newLiterals[i] = newVar.AsLiteral()
		newVar2oldLit[newVar] = literal
		substitution[literal] = newVar.AsLiteral()
	}

	substituted := transformation.SubstituteLiterals(fac, nnf, &substitution)
	solver := sat.NewSolver(fac)
	solver.Add(substituted)

	for _, literal := range newVar2oldLit {
		if literal.IsPos() && literals.Contains(literal.Negate(fac)) {
			name, _, _ := fac.LitNamePhase(literal)
			solver.Add(fac.AMO(fac.Var(name+pos), fac.Var(name+neg)))
		}
	}
	if !solver.Sat() {
		return nil, errorx.BadInput("formula was unsatisfiable")
	}
	minimumModel := solver.Minimize(newLiterals)
	primeImplicant := make([]f.Literal, 0)

	for _, variable := range minimumModel.PosVars() {
		literal, ok := newVar2oldLit[variable]
		if ok {
			primeImplicant = append(primeImplicant, literal)
		}
	}
	return primeImplicant, nil
}
