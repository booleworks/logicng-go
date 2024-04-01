package function

import (
	f "github.com/booleworks/logicng-go/formula"
	"github.com/emirpasic/gods/sets/linkedhashset"
)

// SubNodes returns all sub-nodes of the given formula.  The order of the
// sub-nodes is bottom-up, i.e. a sub-node only appears in the result when all
// of its sub-nodes are already listed.
func SubNodes(fac f.Factory, formula f.Formula) []f.Formula {
	cached, ok := f.LookupFunctionCache(fac, f.FuncSubnodes, formula)
	if ok {
		return cached.([]f.Formula)
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
	slice := make([]f.Formula, result.Size())
	result.Each(func(i int, formula interface{}) { slice[i] = formula.(f.Formula) })
	f.SetFunctionCache(fac, f.FuncSubnodes, formula, slice)
	return slice
}
