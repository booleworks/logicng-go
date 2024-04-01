// Code generated by "stringer -type=Sort"; DO NOT EDIT.

package configuration

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[FormulaFactory-0]
	_ = x[CNF-1]
	_ = x[Sat-2]
	_ = x[MaxSat-3]
	_ = x[Encoder-4]
	_ = x[FormulaRandomizer-5]
	_ = x[AdvancedSimplifier-6]
	_ = x[ModelIteration-7]
}

const _Sort_name = "FormulaFactoryCNFSatMaxSatEncoderFormulaRandomizerAdvancedSimplifierModelIteration"

var _Sort_index = [...]uint8{0, 14, 17, 20, 26, 33, 50, 68, 82}

func (i Sort) String() string {
	if i >= Sort(len(_Sort_index)-1) {
		return "Sort(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Sort_name[_Sort_index[i]:_Sort_index[i+1]]
}
