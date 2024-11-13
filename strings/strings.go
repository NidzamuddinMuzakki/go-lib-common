package strings

import (
	"strconv"
	"strings"
)

// Numeric get all numeric string from the string provided
func Numeric(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if '0' <= b && b <= '9' {
			result.WriteByte(b)
		}
	}
	return result.String()
}

// ToUint64 convert string to uint64 with default value 0
func ToUint64(s string) uint64 {
	number := Numeric(s)
	if len(number) > 0 {
		if s, err := strconv.ParseUint(number, 10, 64); err == nil {
			return s
		}
	}
	return 0
}

// ToFloat64 convert string to float64 with default value 0
func ToFloat64(s string) float64 {
	number := Numeric(s)
	if len(number) > 0 {
		if s, err := strconv.ParseFloat(number, 64); err == nil {
			return s
		}
	}
	return 0
}
