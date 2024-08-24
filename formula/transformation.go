package formula

import "github.com/booleworks/logicng-go/handler"

// A Transformation is a function which maps a formula to another formula
type Transformation func(Factory, Formula) Formula

// An CancellableTransformation is a function which maps a formula to another
// formula.  It takes a handler which can be used to cancel the computation.
// If the computation is cancelled by the handler, the ok flag in the response is
// false.
type CancellableTransformation func(Factory, Formula, handler.Handler) (Formula, bool)

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
