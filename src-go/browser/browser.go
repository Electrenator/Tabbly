package browser

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"reflect"

	"github.com/Electrenator/Tabbly/src-go/util"
	"github.com/shirou/gopsutil/v4/process"
)

type Browser interface {
	GetInfo() BrowserInfo
	isActive() bool
	getWindowData() []WindowInfo
}

type AbstractBrowser struct {
	Browser
	typicalName      string
	processAliases   []string
	storageLocations []string
}

func (browser *AbstractBrowser) isActive() bool {
	pids, err := process.Pids()

	if err != nil {
		slog.Error(
			"Can't find system processes? -",
			"error", err,
		)
		return false
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
				return true
			}
		}
	}

	return false
}

func (browser *AbstractBrowser) getWindowData() []WindowInfo {
	panic(fmt.Sprintf(
		"Browser '%s' (%s) is unimplemented!",
		reflect.TypeOf(browser).String(),
		browser.typicalName,
	))
}

func (browser *AbstractBrowser) getSessionStorageLocation() []string {
	var expandedPaths []string

	for _, pattern := range browser.storageLocations {
		pattern, err := util.ExpandHomeDirectory(pattern)
		if err != nil {
			slog.Error(
				"Can't get current user for '~' expansion",
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
