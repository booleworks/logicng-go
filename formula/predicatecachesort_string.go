// Code generated by "stringer -type=PredicateCacheSort"; DO NOT EDIT.

package formula

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[PredNNF-0]
	_ = x[PredCNF-1]
	_ = x[PredDNF-2]
	_ = x[PredAIG-3]
}

const _PredicateCacheSort_name = "PredNNFPredCNFPredDNFPredAIG"

var _PredicateCacheSort_index = [...]uint8{0, 7, 14, 21, 28}

func (i PredicateCacheSort) String() string {
	if i >= PredicateCacheSort(len(_PredicateCacheSort_index)-1) {
		return "PredicateCacheSort(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _PredicateCacheSort_name[_PredicateCacheSort_index[i]:_PredicateCacheSort_index[i+1]]
}
