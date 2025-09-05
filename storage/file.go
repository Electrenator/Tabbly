package storage

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/electrenator/tabbly/util"
)

const developmentFileAddition = "-dev"

// Copies file from targetFileName path to copyFileName path.
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

// Close a file with error handling
func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		slog.Error("Error closing file", "error", err)
	}
}

func asciiHostnameToPascalCase(input string) string {
	return stringCaseHostname(input, false)
}

func stringCaseHostname(input string, ignoreFirst bool) string {
	const capitalDifference byte = 'a' - 'A'
	words := strings.FieldsFunc(strings.ToLower(input), isHostnameSeparator)

	for i, word := range words {
		if ignoreFirst && i == 0 {
			continue
		}
		wordParts := []byte(word)
		wordParts[0] = wordParts[0] - capitalDifference
		words[i] = string(wordParts)
	}

	return strings.Join(words, "")
}

func asciiHostnameToCamelCase(input string) string {
	return stringCaseHostname(input, true)
}

// Returns true for ASCII hostname separators
func isHostnameSeparator(r rune) bool {
	return r == '-' || r == '_'
}
