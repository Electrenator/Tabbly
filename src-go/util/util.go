package util

import (
	"os/user"
	"path/filepath"
	"strings"
)

func StringContains(haystack string, needle string) bool {
	return strings.Contains(
		strings.ToLower(haystack),
		strings.ToLower(needle),
	)
}

func ExpandHomeDirectory(path string) (string, error) {
	usr, err := user.Current()
	if err == nil && strings.HasPrefix(path, "~/") {
		path = filepath.Join(usr.HomeDir, path[2:])
	}
	return path, err
}

// Converts a list of any into a list of the given type
//
// From https://stackoverflow.com/a/24454401/13042236
func ConvertSlice[E any](in []any) (out []E) {
	out = make([]E, 0, len(in))
	for _, v := range in {
		out = append(out, v.(E))
	}
	return
}
