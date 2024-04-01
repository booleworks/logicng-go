package simplification

import (
	f "github.com/booleworks/logicng-go/formula"
	"github.com/booleworks/logicng-go/sat"
)

// PropagateUnits performs unit propagation on the given formula. Unit
// propagation works the following way: If a formula is such that a literal is
// forced for the formula to be satisfied, then this literal is propagated
// through the formula and thus simplifies the formula. For example, consider
// the formula (A | C) & ~C & (B | C) & (A | ~C) then the literal ~C is forced
// in the formula. Thus, the simplified formula (created by unit propagation)
// yields ~C & A & B.
func PropagateUnits(fac f.Factory, formula f.Formula) f.Formula {
	cached, ok := f.LookupTransformationCache(fac, f.TransUnitPropagation, formula)
	if ok {
		return cached
	}
	miniSatPropagator := sat.NewUnitPropagator(fac)
	miniSatPropagator.Add(formula)
	result := miniSatPropagator.PropagateFormula()
	f.SetTransformationCache(fac, f.TransUnitPropagation, formula, result)
	return result
}
