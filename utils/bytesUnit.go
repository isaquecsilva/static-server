package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	reNumber *regexp.Regexp = regexp.MustCompile(`\d+`)

	reUnit *regexp.Regexp = regexp.MustCompile(`\D+`)

	bytesUnit = map[string]int64{
		"B":  1,
		"KB": 1024,
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
	}
)

func ParseBytesUnit(value string) (int64, error) {
	number, err := strconv.ParseInt(reNumber.FindString(value), 10, 64)
	unitStr := reUnit.FindString(strings.ToUpper(value))

	if err != nil {
		return -1, err
	}

	if unit, ok := bytesUnit[unitStr]; ok {
		return number * unit, nil
	}

	return -1, fmt.Errorf("could not find unit: %s", unitStr)
}
