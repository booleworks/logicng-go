// Code generated by "stringer --type CNFMethod"; DO NOT EDIT.

package sat

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CNFFactory-0]
	_ = x[CNFPG-1]
	_ = x[CNFFullPG-2]
}

const _CNFMethod_name = "CNFFactoryCNFPGCNFFullPG"

var _CNFMethod_index = [...]uint8{0, 10, 15, 24}

func (i CNFMethod) String() string {
	if i >= CNFMethod(len(_CNFMethod_index)-1) {
		return "CNFMethod(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _CNFMethod_name[_CNFMethod_index[i]:_CNFMethod_index[i+1]]
}
