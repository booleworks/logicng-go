// Code generated by "stringer -type=FSort"; DO NOT EDIT.

package formula

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SortFalse-0]
	_ = x[SortTrue-1]
	_ = x[SortLiteral-2]
	_ = x[SortNot-3]
	_ = x[SortAnd-4]
	_ = x[SortOr-5]
	_ = x[SortImpl-6]
	_ = x[SortEquiv-7]
	_ = x[SortCC-8]
	_ = x[SortPBC-9]
}

const _FSort_name = "FalsumVerumLiteralNotAndOrImplEquivCCPBC"

var _FSort_index = [...]uint8{0, 6, 11, 18, 21, 24, 26, 30, 35, 37, 40}

func (i FSort) String() string {
	if i >= FSort(len(_FSort_index)-1) {
		return "FSort(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FSort_name[_FSort_index[i]:_FSort_index[i+1]]
}
