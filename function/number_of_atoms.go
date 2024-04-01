package function

import f "booleworks.com/logicng/formula"

// NumberOfAtoms returns the number of atomic formulas of the given formula. An
// atomic formula is a constant or a literal.
func NumberOfAtoms(fac f.Factory, formula f.Formula) int {
	cached, ok := f.LookupFunctionCache(fac, f.FuncNumberOfAtoms, formula)
	if ok {
		return cached.(int)
	}
	result := 1
	switch formula.Sort() {
	case f.SortNot, f.SortImpl, f.SortEquiv, f.SortOr, f.SortAnd:
		result = 0
		for _, op := range fac.Operands(formula) {
			result += NumberOfAtoms(fac, op)
		}
	}
	f.SetFunctionCache(fac, f.FuncNumberOfAtoms, formula, result)
	return result
}
