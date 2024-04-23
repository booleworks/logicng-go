package formula

// FormulaDepth returns the depth of the given formula. The depth of an atomic
// formula is defined as 0, all other operators increase the depth by 1.
func FormulaDepth(fac Factory, formula Formula) int {
	cached, ok := LookupFunctionCache(fac, FuncDepth, formula)
	if ok {
		return cached.(int)
	}
	var result int
	if formula.IsAtomic() {
		result = 0
	} else {
		maxDepth := 0
		ops := fac.Operands(formula)
		for _, op := range ops {
			maxDepth = max(maxDepth, FormulaDepth(fac, op))
		}
		result = maxDepth + 1
	}
	SetFunctionCache(fac, FuncDepth, formula, result)
	return result
}
