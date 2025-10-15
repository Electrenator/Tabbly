package internal_storage

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	internal_setting "github.com/electrenator/tabbly/internal/setting"
)

// Copies file from targetFileName path to copyFileName path.
func copyFile(targetFileName string, copyFileName string) error {
	targetFile, err := os.Open(targetFileName)
	if err != nil {
		return err
	}
	defer closeFile(targetFile)

	copyFile, err := os.OpenFile(copyFileName, os.O_WRONLY|os.O_CREATE, internal_setting.DefaultFilePerms)
	if err != nil {
		return err
	}
	defer closeFile(copyFile)

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

// Close a file with error handling
func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		slog.Error("Error closing file", "error", err)
	}
}

func asciiHostnameToCamelCase(input string) string {
	const capitalDifference byte = 'a' - 'A'
	words := strings.FieldsFunc(strings.ToLower(input), isHostnameSeparator)

	for i, word := range words {
		if i == 0 {
			continue
		}
		wordParts := []byte(word)
		wordParts[0] = wordParts[0] - capitalDifference
		words[i] = string(wordParts)
	}

	return strings.Join(words, "")
}

// Returns true for ASCII hostname separators
func isHostnameSeparator(r rune) bool {
	return r == '-' || r == '_'
}
