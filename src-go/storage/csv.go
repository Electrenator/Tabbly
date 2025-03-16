package storage

import (
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Electrenator/Tabbly/src-go/browser"
)

var previousTabCount int = -1

// Function combining al usage and saving that to CSV in an V0.1 compatible way.
// Which entails that most data gets lost apart from the tab and window data which gets
// saved in a duplicate way.
// ToDo;
//
//	This should be changed changed into a better format sometime before entering main! Mainly
//	only written this, so it has parity to the current python version
//
// This function also has wayyy to many choice decisions in it and should only be responsible for saving!
func SaveCsv(browserList []browser.BrowserInfo) {
	const separator = ';'
	var windowTabs []int
	tabSum := 0

	for _, browser := range browserList {
		for _, window := range browser.Windows {
			windowTabs = append(windowTabs, window.TabCount)
			tabSum += window.TabCount
		}
	}

	if tabSum == previousTabCount {
		return
	}
	slog.Info("Tab change detected → Continuing to log")
	previousTabCount = tabSum

	csvPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	csvPath = filepath.Join(csvPath, "log/", getCsvName())

	if _, err := os.Stat(csvPath); err != nil {
		slog.Info(fmt.Sprintf("Creating new file: %s", csvPath))
		initCsvFile(csvPath, separator)
	}

	file, err := os.OpenFile(csvPath, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		slog.Error("Can't open csv for writing", "error", err)
		os.Exit(1)
	}

	file.WriteString(fmt.Sprintf(
		"%d%c%d%c%d%c%+v\n",
		time.Now().Unix(), separator,
		len(windowTabs), separator,
		tabSum, separator,
		windowTabs,
	))
	file.Close()
}

func getCsvName() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}

	return "tabUsage" + asciiHostnameToPascalCase(hostname) + "V2.csv"
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

func initCsvFile(path string, sperator rune) {
	err := os.MkdirAll(filepath.Dir(path), fs.ModePerm)
	if err != nil {
		slog.Error("Can't create log folder", "error", err)
		os.Exit(1)
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, fs.ModePerm)
	if err != nil {
		slog.Error("Can't create log file", "error", err)
		os.Exit(1)
	}

	file.WriteString(fmt.Sprintf(
		"'UNIX timestamp'%c'Total window count'%c'Total tab count'%c'List of total tabs per window'\n",
		sperator, sperator, sperator))
	file.Close()
}
