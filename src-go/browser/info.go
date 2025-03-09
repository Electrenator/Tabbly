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
	GetChromiumBrowser(),
}

func GetAvailableBrowsers() []Browser {
	return availableBrowsers
}

func LogAllBrowserStates() {
	combinedInfo := make([]BrowserInfo, len(availableBrowsers))

	for i, browser := range availableBrowsers {
		combinedInfo[i] = BrowserInfo{
			browser.GetName(),
			browser.GetState(),
			browser.GetherWindowData(),
		}
	}
	slog.Info(fmt.Sprintf("Browsers: %+v\n", combinedInfo))
}
