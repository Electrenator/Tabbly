package browser

import (
	"fmt"
	"log/slog"
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
				slog.Info(fmt.Sprintf("Found %d - %s", pid, name))
				return true
			}
		}
	}

	return false
}

func (browser *AbstractBrowser) getWindowData() []WindowInfo {
	panic(fmt.Sprintf(
		"Browser '%s' is unimplemented!",
		reflect.TypeOf(browser).String(),
	))
}
