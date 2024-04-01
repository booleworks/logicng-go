package formula

// TransformationCacheSort encodes a formula transformation sort for which the
// result can be cached.
type TransformationCacheSort byte

const (
	TransNNF TransformationCacheSort = iota
	TransCNFFactorization
	TransDNFFactorization
	TransAIG
	TransDistrSimpl
	TransUnitPropagation
	TransBDDCNF
	TransBDDNF
)

//go:generate stringer -type=TransformationCacheSort

// LookupTransformationCache searches a cache entry for a given transformation
// sort and formula.  It returns the optional result and a flag whether there
// was a cache entry.
func LookupTransformationCache(fac Factory, sort TransformationCacheSort, formula Formula) (Formula, bool) {
	c, ok := (*fac.transformationCacheEntry(sort))[formula]
	return c, ok
}

// SetTransformationCache sets a cache entry for a given transformation sort,
// formula and value to cache.
func SetTransformationCache(fac Factory, sort TransformationCacheSort, formula, value Formula) {
	(*fac.transformationCacheEntry(sort))[formula] = value
}
