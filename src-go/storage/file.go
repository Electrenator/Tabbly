package storage

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/Electrenator/Tabbly/src-go/util"
	"github.com/dustin/go-humanize"
)

func copyFile(targetFileName string, copyFileName string) error {
	targetFile, err := os.Open(targetFileName)
	if err != nil {
		return err
	}
	defer closeFile(targetFile)

	copyFile, err := os.OpenFile(copyFileName, os.O_WRONLY|os.O_CREATE, util.DefaultFilePerms)
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

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		slog.Error("Error closing file", "error", err)
	}
}

func asciiHostnameToPascalCase(input string) string {
	const capitalDifference byte = 'a' - 'A'
	parts := strings.FieldsFunc(strings.ToLower(input), isHostnameSperator)

	for i, word := range parts {
		wordParts := []byte(word)
		wordParts[0] = wordParts[0] - capitalDifference
		parts[i] = string(wordParts)
	}

	return strings.Join(parts, "")
}

// Returns true for ASCII hostname separators
func isHostnameSperator(r rune) bool {
	return r == '-' || r == '_'
}
