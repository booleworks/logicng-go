// Code generated by "stringer -type=CNFAlgorithm"; DO NOT EDIT.

package normalform

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CNFFactorization-0]
	_ = x[CNFTseitin-1]
	_ = x[CNFPlaistedGreenbaum-2]
	_ = x[CNFAdvanced-3]
}

const _CNFAlgorithm_name = "CnfFactorizationCnfTseitinCnfPlaistedGreenbaumCnfAdvanced"

var _CNFAlgorithm_index = [...]uint8{0, 16, 26, 46, 57}

func (i CNFAlgorithm) String() string {
	if i >= CNFAlgorithm(len(_CNFAlgorithm_index)-1) {
		return "CNFAlgorithm(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _CNFAlgorithm_name[_CNFAlgorithm_index[i]:_CNFAlgorithm_index[i+1]]
}
