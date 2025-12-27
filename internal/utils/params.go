package utils

import (
	"strconv"
)

// ParseIDParam parses a string parameter into a uint64.
func ParseIDParam(param string) (uint64, error) {
	return strconv.ParseUint(param, 10, 64)
}
