package formula

// NumberOfNodes returns the number of nodes (in the DAG) of the given formula.
func NumberOfNodes(fac Factory, formula Formula) int {
	cached, ok := LookupFunctionCache(fac, FuncNumberOfNodes, formula)
	if ok {
		return cached.(int)
	}
	result := 1
	switch fsort := formula.Sort(); fsort {
	case SortNot, SortImpl, SortEquiv, SortOr, SortAnd:
		for _, op := range fac.Operands(formula) {
			result += NumberOfNodes(fac, op)
		}
	case SortCC, SortPBC:
		_, _, lits, _, _ := fac.PBCOps(formula)
		result = 1 + len(lits)
	}
	SetFunctionCache(fac, FuncNumberOfNodes, formula, result)
	return result
}
