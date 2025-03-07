package util

import "strings"

func StringContains(haystack string, needle string) bool {
	return strings.Contains(
		strings.ToLower(haystack),
		strings.ToLower(needle),
	)
}
