package extensions

import (
	"fmt"
	"strconv"
)

// TryParseInt attempts a conversion ignoring any error.
func TryParseInt(val interface{}) int {
	i, _ := strconv.Atoi(val.(string))
	return i
}

// MustParseInt panics if conversion fails.
func MustParseInt(val interface{}) int {
	i, err := strconv.Atoi(val.(string))
	if err != nil {
		panic(fmt.Sprintf("error converting %v", val))
	}
	return i
}

// TryParseFloat32 attempts a conversion ignoring any error.
func TryParseFloat32(val interface{}) float32 {
	f, _ := strconv.ParseFloat(val.(string), 32)
	return float32(f)
}

// MustParseFloat32 panics if conversion fails.
func MustParseFloat32(val interface{}) float32 {
	f, err := strconv.ParseFloat(val.(string), 32)
	if err != nil {
		panic(fmt.Sprintf("error converting %v", val))
	}
	return float32(f)
}

// FormatFloat32 formats float32 with no exponent.
func FormatFloat32(val float32) string {
	return strconv.FormatFloat(float64(val), 'f', -1, 32)
}
