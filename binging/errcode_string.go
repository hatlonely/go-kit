// Code generated by "stringer -type ErrCode -trimprefix Err"; DO NOT EDIT.

package binging

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ErrInvalidDstType-1]
	_ = x[ErrMissingRequiredField-2]
	_ = x[ErrInvalidFormat-3]
}

const _ErrCode_name = "InvalidDstTypeMissingRequiredFieldInvalidFormat"

var _ErrCode_index = [...]uint8{0, 14, 34, 47}

func (i ErrCode) String() string {
	i -= 1
	if i < 0 || i >= ErrCode(len(_ErrCode_index)-1) {
		return "ErrCode(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _ErrCode_name[_ErrCode_index[i]:_ErrCode_index[i+1]]
}
