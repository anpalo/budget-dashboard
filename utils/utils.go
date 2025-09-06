package utils

import (
	"strconv"
	"strings"
)

func IsDate(s string) bool {
	return len(s) > 0 && s[0] >= '0' && s[0] <= '9'
}

func ParseToFloat(s string) float64 {
	s = strings.ReplaceAll(s, ",", "")
	if s == "" {
		return 0
	}
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return num
}