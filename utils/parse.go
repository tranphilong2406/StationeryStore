package utils

import "strconv"

func ParseInt(s string, def int) int {
	if s == "" {
		return def
	}
	num, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return num
}
