// Code generated by "stringer -type=CSort"; DO NOT EDIT.

package formula

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EQ-0]
	_ = x[LE-1]
	_ = x[LT-2]
	_ = x[GE-3]
	_ = x[GT-4]
}

const _CSort_name = "EQLELTGEGT"

var _CSort_index = [...]uint8{0, 2, 4, 6, 8, 10}

func (i CSort) String() string {
	if i >= CSort(len(_CSort_index)-1) {
		return "CSort(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _CSort_name[_CSort_index[i]:_CSort_index[i+1]]
}
