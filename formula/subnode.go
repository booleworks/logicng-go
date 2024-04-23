package formula

import "github.com/emirpasic/gods/sets/linkedhashset"

// SubNodes returns all sub-nodes of the given formula.  The order of the
// sub-nodes is bottom-up, i.e. a sub-node only appears in the result when all
// of its sub-nodes are already listed.
func SubNodes(fac Factory, formula Formula) []Formula {
	cached, ok := LookupFunctionCache(fac, FuncSubnodes, formula)
	if ok {
		return cached.([]Formula)
	}
	result := linkedhashset.New()
	for _, op := range fac.Operands(formula) {
		if !result.Contains(op) {
			for _, recOp := range SubNodes(fac, op) {
				result.Add(recOp)
			}
		}
	}
	result.Add(formula)
	slice := make([]Formula, result.Size())
	result.Each(func(i int, formula interface{}) { slice[i] = formula.(Formula) })
	SetFunctionCache(fac, FuncSubnodes, formula, slice)
	return slice
}
