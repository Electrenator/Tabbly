package util

import (
	"io"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
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

func CopyFile(targetFileName string, copyFileName string) error {
	targetFile, err := os.Open(targetFileName)
	if err != nil {
		return err
	}
	defer CloseFile(targetFile)

	copyFile, err := os.OpenFile(copyFileName, os.O_WRONLY|os.O_CREATE, DefaultFilePerms)
	if err != nil {
		return err
	}
	defer CloseFile(copyFile)

	copiedBytes, err := io.Copy(copyFile, targetFile)
	if err == nil {
		err = copyFile.Sync()

		if err == nil {
			slog.Info("Copied over files",
				"Source file", targetFileName,
				"Destination file", copyFileName,
				"Size", humanize.IBytes(uint64(copiedBytes)),
			)
		}
	}
	return err
}

func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		slog.Error("Error closing file", "error", err)
	}
}
