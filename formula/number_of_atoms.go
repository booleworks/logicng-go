package formula

// NumberOfAtoms returns the number of atomic formulas of the given formula. An
// atomic formula is a constant or a literal.
func NumberOfAtoms(fac Factory, formula Formula) int {
	cached, ok := LookupFunctionCache(fac, FuncNumberOfAtoms, formula)
	if ok {
		return cached.(int)
	}
	result := 1
	switch formula.Sort() {
	case SortNot, SortImpl, SortEquiv, SortOr, SortAnd:
		result = 0
		for _, op := range fac.Operands(formula) {
			result += NumberOfAtoms(fac, op)
		}
	}
	SetFunctionCache(fac, FuncNumberOfAtoms, formula, result)
	return result
}
