package util

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Capitalization insensitive string.Contains.
func StringContains(haystack string, needle string) bool {
	return strings.Contains(
		strings.ToLower(haystack),
		strings.ToLower(needle),
	)
}

// Expands user specific properties of a directory path. Paths which will be
// expanded might start with "~" or "%APPDATA%". Returns an unaltered path
// when unsuccessful with error set.
func ExpandUserDirectory(path string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return path, err
	}
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(usr.HomeDir, path[2:])
	} else if strings.HasPrefix(path, "%APPDATA%") {
		appDataPath, err := os.UserConfigDir()
		if err != nil {
			return path, err
		}
		path = strings.Replace(path, "%APPDATA%", appDataPath, 1)
	}
	return path, nil
}

// Sums all entries of a list together. Will return -1 if no slice is given.
func SumSlice[T ~int](slice []T) T {
	if slice == nil {
		return -1
	}
	var sum T
	for i := 0; i < len(slice); i++ {
		sum += slice[i]
	}

	return sum
}

// Checks if two slices are the same both in length and there values.
//
// Note: this isn't tested on slices with slices in them and likely isn't
// accurate for that case.
// Should be replaced with std function (slices.Compare)
func SameSlice[T comparable](a, b []T) bool {
	if a == nil {
		return b == nil
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
