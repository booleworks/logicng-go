package formula

// PredicateCacheSort encodes a formula predicate sort for which the result can
// be cached.
type PredicateCacheSort byte

const (
	PredNNF PredicateCacheSort = iota
	PredCNF
	PredDNF
	PredAIG
)

//go:generate stringer -type=PredicateCacheSort

// LookupPredicateCache searches a cache entry for a given predicate sort and
// formula.  It returns the optional result and a flag whether there was a
// cache entry.
func LookupPredicateCache(fac Factory, sort PredicateCacheSort, formula Formula) (bool, bool) {
	c, ok := (*fac.predicateCacheEntry(sort))[formula]
	return c, ok
}

// SetPredicateCache sets a cache entry for a given predicate sort, formula and
// value to cache.
func SetPredicateCache(fac Factory, sort PredicateCacheSort, formula Formula, value bool) {
	(*fac.predicateCacheEntry(sort))[formula] = value
}
