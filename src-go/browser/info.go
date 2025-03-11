package browser

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

func CreateSummary(info []BrowserInfo) []int {
	var windowTabs []int

	for _, browser := range info {
		for _, window := range browser.Windows {
			windowTabs = append(windowTabs, window.TabCount)
		}
	}

	return windowTabs
}
