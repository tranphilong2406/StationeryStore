package utils

import (
	"strconv"
	"time"
)

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

func ParseTime(s string) time.Time {
	if s == "" {
		return time.Now()
	}
	// Assuming the time is in a specific format, you can add parsing logic here if needed.
	timeParse, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Now()
	}
	return timeParse
}

func ParseBool(s string) bool {
	if s == "" {
		return true
	}
	
	b, err := strconv.ParseBool(s)
	if err != nil {
		return true
	}
	return b
}
