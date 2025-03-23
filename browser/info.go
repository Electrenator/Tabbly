package browser

type BrowserInfo struct {
	Name    string
	IsOpen  BrowserState
	Windows []WindowInfo
}

// Get the total number of tabs within this browser.
func (info *BrowserInfo) TotalTabs() int {
	sum := 0
	for _, tabs := range info.Windows {
		sum += tabs.TabCount
	}
	return sum
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
	GetFirefoxBrowser(),
	GetFirefoxDeveloperBrowser(),
	GetChromiumBrowser(),
}

func GetAvailableBrowsers() []Browser {
	return availableBrowsers
}

// Creates summary list containing all tabs open in every open window from
// all given.
func CreateSummary(info []BrowserInfo) []int {
	var windowTabs []int

	for _, browser := range info {
		for _, window := range browser.Windows {
			windowTabs = append(windowTabs, window.TabCount)
		}
	}

	return windowTabs
}
