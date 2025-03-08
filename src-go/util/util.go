package util

import (
	"log/slog"
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
	if err == nil && strings.HasPrefix(path, "~/") {
		path = filepath.Join(usr.HomeDir, path[2:])
	}
	return path, err
}

func FileSize(filePath string) int64 {
	stats, err := os.Stat(filePath)
	if err != nil {
		slog.Error("Can't read file size!", "error", err)
		return -1
	}
	return stats.Size()
}
