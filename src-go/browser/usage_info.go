package browser

type BrowserInfo struct {
	name    string
	isOpen  BrowserState
	windows []WindowInfo
}

type WindowInfo struct {
	tabCount int
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

func GetUsageInfo() []BrowserInfo {
	info := make([]BrowserInfo, len(availableBrowsers))

	for i := range availableBrowsers {
		info[i] = availableBrowsers[i].GetInfo()
	}

	return info
}
