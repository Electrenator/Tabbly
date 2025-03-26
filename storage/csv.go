package storage

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/electrenator/tabbly/browser"
	internal_status "github.com/electrenator/tabbly/internal/status"
	"github.com/electrenator/tabbly/util"
)

const separator = ';'
const (
	legacyCsvTimestamp = iota
	legacyCsvWindowCount
	legacyCsvTabCount
	legacyCsvTabsPerWindow
)

// Import the legacy CSV from V0.1 into the new database storage format.
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
	var windowsTimeMap []TimedBrowserInfo
	var totalWindowCount int

	for scanner.Scan() {
		timestamp, browserInfo, err := legacyCsvLineToBrowserInfo(scanner.Bytes())

		if err != nil {
			slog.Error(err.Error())
			os.Exit(internal_status.UNSPECIFIED_PRIMARY_FUNCTION_ERROR)
		}

		// Is data not already in read before?
		if *timestamp == lastTimestamp || util.SameSlice(browserInfo.Windows, lastTabsPerWindow) ||
			// While at it, do a little sanity check for the tab count
			browserInfo.TotalTabs()-len(browserInfo.Windows) < 0 {
			continue
		}

		windowsTimeMap = append(windowsTimeMap, TimedBrowserInfo{
			Timestamp:   *timestamp,
			BrowserInfo: *browserInfo,
		})
		lastTimestamp = *timestamp
		lastTabsPerWindow = browserInfo.Windows
		totalWindowCount += len(browserInfo.Windows)
	}

	if err := scanner.Err(); err != nil {
		slog.Error("Issue reading file", "error", err)
		os.Exit(internal_status.UNSPECIFIED_PRIMARY_FUNCTION_ERROR)
	}
	slog.Info("Scanned legacy save!",
		"totalDeduplicatedEntries", len(windowsTimeMap),
		"totalWindowRecords", totalWindowCount,
	)

	if err := SaveMultipleToDb(&windowsTimeMap); err != nil {
		// Already logged in DB rollback
		os.Exit(internal_status.DB_CONNECT_ERROR)
	}
}

// Checks if the file being scanned is in the legacy CSV format used by V0.1.
func isLegacyCsv(scanner *bufio.Scanner) bool {
	scanner.Scan()
	line := scanner.Bytes()
	return util.SameSlice(line, []byte(fmt.Sprintf(
		"'UNIX timestamp'%c'Total window count'%c'Total tab count'%c'List of total tabs per window'",
		separator, separator, separator,
	)))
}

// Parses read line from the V0.1 CSV format into the corresponding
// BrowserInfo struct. Will also give back the timestamp associated, at least
// when it doesn't error.
func legacyCsvLineToBrowserInfo(line []byte) (*int64, *browser.BrowserInfo, error) {
	var windowCount, tabCount int

	standardLegacyBrowser := browser.GetFirefoxBrowser().GetName()
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

// Legacy csv had a point where it didn't have tabs per window. This function
// adds extra 1 size windows, so that they are still taken into account within
// the history even though there not accurate.
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
