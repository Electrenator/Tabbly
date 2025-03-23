package util

import (
	"os"
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
	if err == nil {
		if strings.HasPrefix(path, "~/") {
			path = filepath.Join(usr.HomeDir, path[2:])
		} else if strings.HasPrefix(path, "%APPDATA%") {
			appDataPath, err := os.UserConfigDir()
			if err != nil {
				return path, err
			}
			path = strings.Replace(path, "%APPDATA%", appDataPath, 1)
		}
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
// Note: this isn't tested on slices with slices in them and likely isn't accurate for that case.
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
