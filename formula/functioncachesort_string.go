// Code generated by "stringer -type=FunctionCacheSort"; DO NOT EDIT.

package formula

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[FuncVariables-0]
	_ = x[FuncLiterals-1]
	_ = x[FuncDepth-2]
	_ = x[FuncSubnodes-3]
	_ = x[FuncNumberOfAtoms-4]
	_ = x[FuncNumberOfNodes-5]
	_ = x[FuncDNNFModelCount-6]
}

const _FunctionCacheSort_name = "FuncVariablesFuncLiteralsFuncDepthFuncSubnodesFuncNumberOfAtomsFuncNumberOfNodesFuncDNNFModelCount"

var _FunctionCacheSort_index = [...]uint8{0, 13, 25, 34, 46, 63, 80, 98}

func (i FunctionCacheSort) String() string {
	if i >= FunctionCacheSort(len(_FunctionCacheSort_index)-1) {
		return "FunctionCacheSort(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FunctionCacheSort_name[_FunctionCacheSort_index[i]:_FunctionCacheSort_index[i+1]]
}
