package utils

import (
	"strconv"
)

// StringToInt converts a string to an integer and returns an error if the conversion fails.
func StringToInt(s string, def int) int {
	// Attempt to convert the string to an integer
	num, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return num
}

// StringToInt64 converts a string to an integer64 and returns an error if the conversion fails.
func StringToInt64(s string, def int64) int64 {
	// Attempt to convert the string to an integer
	num, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return def
	}
	return num
}
