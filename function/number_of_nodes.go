package function

import (
	f "github.com/booleworks/logicng-go/formula"
)

// NumberOfNodes returns the number of nodes (in the DAG) of the given formula.
func NumberOfNodes(fac f.Factory, formula f.Formula) int {
	cached, ok := f.LookupFunctionCache(fac, f.FuncNumberOfNodes, formula)
	if ok {
		return cached.(int)
	}
	result := 1
	switch fsort := formula.Sort(); fsort {
	case f.SortNot, f.SortImpl, f.SortEquiv, f.SortOr, f.SortAnd:
		for _, op := range fac.Operands(formula) {
			result += NumberOfNodes(fac, op)
		}
	case f.SortCC, f.SortPBC:
		_, _, lits, _, _ := fac.PBCOps(formula)
		result = 1 + len(lits)
	}
	f.SetFunctionCache(fac, f.FuncNumberOfNodes, formula, result)
	return result
}
