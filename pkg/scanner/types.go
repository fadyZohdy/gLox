package scanner

import "strconv"

type LoxFloat64 float64

func (f LoxFloat64) String() string {
	// print the fewest digits necessary to accurately represent the float.
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}
