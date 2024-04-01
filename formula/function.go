package formula

// FunctionCacheSort encodes a formula function sort for which the result can
// be cached.
type FunctionCacheSort byte

const (
	FuncVariables FunctionCacheSort = iota
	FuncLiterals
	FuncDepth
	FuncSubnodes
	FuncNumberOfAtoms
	FuncNumberOfNodes
	FuncDNNFModelCount
)

//go:generate stringer -type=FunctionCacheSort

// LookupFunctionCache searches a cache entry for a given function sort and
// formula.  It returns the optional result and a flag whether there was a
// cache entry.
func LookupFunctionCache(fac Factory, sort FunctionCacheSort, formula Formula) (any, bool) {
	c, ok := (*fac.functionCacheEntry(sort))[formula]
	return c, ok
}

// SetFunctionCache sets a cache entry for a given function sort, formula and
// value to cache.
func SetFunctionCache(fac Factory, sort FunctionCacheSort, formula Formula, value any) {
	(*fac.functionCacheEntry(sort))[formula] = value
}
