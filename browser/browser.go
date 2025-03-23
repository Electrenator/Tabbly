package browser

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/electrenator/tabbly/util"
	"github.com/shirou/gopsutil/v4/process"
)

type Browser interface {
	GetName() string
	GetState() BrowserState
	GetherWindowData() []WindowInfo
}

type AbstractBrowser struct {
	Browser
	typicalName      string
	processAliases   []string
	storageLocations []string
}

func (browser *AbstractBrowser) GetState() BrowserState {
	if len(browser.processAliases) == 0 {
		return BROWSER_STATE_UNKNOWN
	}
	pids, err := process.Pids()

	if err != nil {
		slog.Error(
			"Can't find system processes? -",
			"error", err,
		)
		return BROWSER_CLOSED
	}

	for _, pid := range pids {
		proc, err := process.NewProcess(pid)
		if err != nil {
			continue
		}

		name, err := proc.Name()
		if err != nil {
			continue
		}

		for _, needle := range browser.processAliases {
			if util.StringContains(name, needle) {
				slog.Info(fmt.Sprintf("Found %s as process %d", name, pid))
				return BROWSER_OPEN
			}
		}
	}

	return BROWSER_CLOSED
}

func (browser *AbstractBrowser) GetName() string {
	return browser.typicalName
}

func (browser *AbstractBrowser) getSessionStorageLocation() []string {
	var expandedPaths []string

	for _, pattern := range browser.storageLocations {
		pattern, err := util.ExpandHomeDirectory(pattern)
		if err != nil {
			slog.Error(
				"Can't get current user for home directory expansion",
				"error", err,
			)
		}
		foundPaths, err := filepath.Glob(pattern)
		if err != nil {
			slog.Error(
				"Can't expand glob pattern",
				"pattern", pattern,
				"error", err,
			)
		}
		expandedPaths = append(expandedPaths, foundPaths...)
	}
	if len(expandedPaths) > 0 {
		slog.Info(fmt.Sprintf(
			"(%s) Detected browser storage: %+v",
			browser.typicalName,
			expandedPaths,
		))
	}
	return expandedPaths
}
