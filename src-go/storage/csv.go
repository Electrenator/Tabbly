package storage

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Electrenator/Tabbly/src-go/browser"
	internal_status "github.com/Electrenator/Tabbly/src-go/internal/status"
	"github.com/Electrenator/Tabbly/src-go/util"
)

const separator = ';'
const (
	legacyCsvTimestamp = iota
	legacyCsvWindowCount
	legacyCsvTabCount
	legacyCsvTabsPerWindow
)

var legacyCsvHeader = fmt.Sprintf(
	"'UNIX timestamp'%c'Total window count'%c'Total tab count'%c'List of total tabs per window'",
	separator, separator, separator,
)
var standardLegacyBrowser = browser.GetFirefoxBrowser().GetName()

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
		file = initCsvFile(csvPath)
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

func initCsvFile(path string) *os.File {
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

	file.WriteString(legacyCsvHeader + "\n")
	return file
}

func ImportLegacyCsv(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		slog.Error("Unable to open file", "error", err)
		os.Exit(internal_status.FILE_OPEN_ERROR)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !isLegacyCsv(scanner) {
		slog.Warn("Given file is not a legacy csv save file")
		os.Exit(internal_status.UNSPECIFIED_PRIMARY_FUNCTION_ERROR)
	}

	var lastTimestamp int64
	var lastTabsPerWindow []browser.WindowInfo
	windowsTimeMap := make(map[int64]*browser.BrowserInfo)

	for scanner.Scan() {
		timestamp, browserInfo, err := legacyCsvLineToBrowserInfo(scanner.Bytes())

		if err != nil {
			slog.Error(err.Error())
			os.Exit(internal_status.UNSPECIFIED_PRIMARY_FUNCTION_ERROR)
		}

		if *timestamp == lastTimestamp || util.SameSlice(browserInfo.Windows, lastTabsPerWindow) {
			continue
		}

		windowsTimeMap[*timestamp] = browserInfo
		lastTimestamp = *timestamp
		lastTabsPerWindow = browserInfo.Windows
	}

	if err := scanner.Err(); err != nil {
		slog.Error("Issue reading file", "error", err)
		os.Exit(internal_status.UNSPECIFIED_PRIMARY_FUNCTION_ERROR)
	}

	// Todo : save to DB
}

func isLegacyCsv(scanner *bufio.Scanner) bool {
	scanner.Scan()
	line := scanner.Bytes()
	return util.SameSlice(line, []byte(legacyCsvHeader))
}

func legacyCsvLineToBrowserInfo(line []byte) (*int64, *browser.BrowserInfo, error) {
	var windowCount, tabCount int
	lineParts := bytes.Split(line, []byte{separator})
	timestampFloat, err := strconv.ParseFloat(string(lineParts[legacyCsvTimestamp]), 64)

	if err != nil {
		return nil, nil, fmt.Errorf("unable to read timestamp of '%s': %w",
			lineParts[legacyCsvTimestamp], err)
	}

	timestamp := int64(timestampFloat)
	windowCount, err = strconv.Atoi(string(lineParts[legacyCsvWindowCount]))

	if err == nil && windowCount == 0 {
		return &timestamp, &browser.BrowserInfo{
			Name:    standardLegacyBrowser,
			IsOpen:  0,
			Windows: nil,
		}, nil
	}

	tabCount, err = strconv.Atoi(string(lineParts[legacyCsvTabCount]))

	if err != nil {
		return nil, nil, err
	}

	if len(lineParts[legacyCsvTabsPerWindow]) == 0 {
		// State with old csv saves → has no tabs per window
		return &timestamp, &browser.BrowserInfo{
			Name:    standardLegacyBrowser,
			IsOpen:  1,
			Windows: generateLegacyTabsPerWindow(windowCount, tabCount),
		}, nil
	}
	windowInfo, err := parseLegacyTabsPerWindow(lineParts[legacyCsvTabsPerWindow])

	if err != nil {
		return nil, nil, err
	}
	return &timestamp, &browser.BrowserInfo{
		Name:    standardLegacyBrowser,
		IsOpen:  1,
		Windows: windowInfo,
	}, nil
}

// Legacy csv had a point where it didn't have tabs per window. This function adds extra 1 size
// windows, so that they are still taken into account within the history even though there not accurate.
func generateLegacyTabsPerWindow(windowCount int, tabCount int) []browser.WindowInfo {
	if windowCount <= 0 {
		return nil
	}
	windowInfo := []browser.WindowInfo{{TabCount: tabCount - (windowCount - 1)}}

	for range windowCount - 1 {
		windowInfo = append(windowInfo, browser.WindowInfo{TabCount: 1})
	}
	return windowInfo
}

// Parses the json array of tabs per window from the legacy csv. Data looks like;
//
//	[1, 2, 3]
func parseLegacyTabsPerWindow(rawJson []byte) ([]browser.WindowInfo, error) {
	var tabsPerWindow []int

	err := json.Unmarshal(rawJson, &tabsPerWindow)

	if err != nil {
		return nil, err
	}
	windowInfo := make([]browser.WindowInfo, len(tabsPerWindow))

	for i, tabCount := range tabsPerWindow {
		windowInfo[i] = browser.WindowInfo{TabCount: tabCount}
	}
	return windowInfo, nil
}
