// Code generated by "stringer -type=PBCEncoder"; DO NOT EDIT.

package encoding

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[PBCSWC-0]
	_ = x[PBCBinaryMerge-1]
	_ = x[PBCAdderNetworks-2]
	_ = x[PBCBest-3]
}

const _PBCEncoder_name = "PBCSWCPBCBinaryMergePBCAdderNetworksPBCBest"

var _PBCEncoder_index = [...]uint8{0, 6, 20, 36, 43}

func (i PBCEncoder) String() string {
	if i >= PBCEncoder(len(_PBCEncoder_index)-1) {
		return "PBCEncoder(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _PBCEncoder_name[_PBCEncoder_index[i]:_PBCEncoder_index[i+1]]
}
