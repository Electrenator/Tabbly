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
