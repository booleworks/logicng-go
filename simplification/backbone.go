package simplification

import (
	"booleworks.com/logicng/assignment"
	f "booleworks.com/logicng/formula"
	"booleworks.com/logicng/sat"
)

// SimplifyWithBackbone simplifies the given formula by computing its backbone
// and propagating it through the formula. E.g. in the formula A & B & (A | B |
// C) & (~B | D) the backbone A, B is computed and propagated, yielding the
// simplified formula A & B & D.
func SimplifyWithBackbone(fac f.Factory, formula f.Formula) f.Formula {
	solver := sat.NewSolver(fac)
	solver.Add(formula)
	backbone := solver.ComputeBackbone(fac, f.Variables(fac, formula).Content())
	if !backbone.Sat {
		return fac.Falsum()
	}
	if len(backbone.Negative) > 0 || len(backbone.Positive) > 0 {
		backboneFormula := backbone.ToFormula(fac)
		ass := assignment.Empty()
		for _, lit := range backbone.Positive {
			_ = ass.AddLit(fac, lit.AsLiteral())
		}
		for _, lit := range backbone.Negative {
			_ = ass.AddLit(fac, lit.Negate(fac))
		}
		restrictedFormula := assignment.Restrict(fac, formula, ass)
		return fac.And(backboneFormula, restrictedFormula)
	} else {
		return formula
	}
}
