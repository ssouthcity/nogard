// Code generated by "stringer -type Availability"; DO NOT EDIT.

package nogard

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Permanent-0]
	_ = x[Available-1]
	_ = x[Unavailable-2]
}

const _Availability_name = "PermanentAvailableUnavailable"

var _Availability_index = [...]uint8{0, 9, 18, 29}

func (i Availability) String() string {
	if i >= Availability(len(_Availability_index)-1) {
		return "Availability(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Availability_name[_Availability_index[i]:_Availability_index[i+1]]
}
