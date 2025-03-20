package storage

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/Electrenator/Tabbly/src-go/browser"
	internal_status "github.com/Electrenator/Tabbly/src-go/internal/status"
	"github.com/Electrenator/Tabbly/src-go/util"
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
func SaveToCsv(browserList []browser.BrowserInfo) {
	const separator = ';'
	summary := browser.CreateSummary(browserList)

	if len(summary) == previousTabCount {
		return
	}
	slog.Info("Tab change detected → Continuing to log")
	previousTabCount = len(summary)

	csvPath := filepath.Join(ApplicationSettings.DataPath, getCsvName())
	var file *os.File

	if _, err := os.Stat(csvPath); err != nil {
		slog.Info(fmt.Sprintf("Creating new file: %s", csvPath))
		file = initCsvFile(csvPath, separator)
	}

	if file == nil {
		tmpFile, err := os.OpenFile(csvPath, os.O_WRONLY|os.O_APPEND, 0)
		if err != nil {
			slog.Error("Can't open csv for writing", "error", err)
			os.Exit(internal_status.FILE_OPEN_ERROR)
		}
		file = tmpFile
	}
	defer file.Close()

	_, err := file.WriteString(fmt.Sprintf(
		"%d%c%d%c%d%c%+v\n",
		time.Now().Unix(), separator,
		len(summary), separator,
		util.SumSlice(summary), separator,
		summary,
	))

	if err != nil {
		slog.Error("Error writing csv", "error", err)
	}
	file.Close()
}

func getCsvName() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}

	if ApplicationSettings.IsDevelopmentBuild {
		hostname += developmentFileAddition
	}

	return "tabUsage" + asciiHostnameToPascalCase(hostname) + ".csv"
}

func initCsvFile(path string, sperator rune) *os.File {
	err := os.MkdirAll(filepath.Dir(path), util.DefaultDirPerms)
	if err != nil {
		slog.Error("Can't create log folder", "error", err)
		os.Exit(internal_status.FILE_CREATION_ERROR)
	}

	// Write flag too since this connection is being kept, so the connection can be re-used
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, util.DefaultFilePerms)
	if err != nil {
		slog.Error("Can't create log file", "error", err)
		os.Exit(internal_status.FILE_CREATION_ERROR)
	}

	file.WriteString(fmt.Sprintf(
		"'UNIX timestamp'%c'Total window count'%c'Total tab count'%c'List of total tabs per window'\n",
		sperator, sperator, sperator))
	return file
}
