// Code generated by "stringer -type=IncrementalStrategy"; DO NOT EDIT.

package maxsat

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[IncNone-0]
	_ = x[IncIterative-1]
}

const _IncrementalStrategy_name = "IncNoneIncIterative"

var _IncrementalStrategy_index = [...]uint8{0, 7, 19}

func (i IncrementalStrategy) String() string {
	if i >= IncrementalStrategy(len(_IncrementalStrategy_index)-1) {
		return "IncrementalStrategy(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _IncrementalStrategy_name[_IncrementalStrategy_index[i]:_IncrementalStrategy_index[i+1]]
}
