package browser

import (
	"fmt"
	"log/slog"
)

type BrowserInfo struct {
	Name    string
	IsOpen  BrowserState
	Windows []WindowInfo
}

type WindowInfo struct {
	TabCount int
}

type BrowserState int

const (
	BROWSER_CLOSED BrowserState = iota
	BROWSER_OPEN
	BROWSER_STATE_UNKNOWN = -1
)

var availableBrowsers = []Browser{
	GetTesterBrowser(),
	GetFirefoxBrowser(),
	GetFirefoxDeveloperBrowser(),
}

func GetAvailableBrowsers() []Browser {
	return availableBrowsers
}

func LogAllBrowserStates() {
	combinedInfo := make([]BrowserInfo, len(availableBrowsers))

	for i := range availableBrowsers {
		combinedInfo[i] = availableBrowsers[i].GatherInfo()
	}
	slog.Info(fmt.Sprintf("Browsers: %+v\n", combinedInfo))
}
