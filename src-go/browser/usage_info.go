package browser

type BrowserInfo struct {
	name    string
	isOpen  bool
	windows []WindowInfo
}

type WindowInfo struct {
	openTabs int
}

var supportedBrowsers = []Browser{
	GetTesterBrowser(),
	GetFirefoxBrowser(),
}

func GetUsageInfo() []BrowserInfo {
	info := make([]BrowserInfo, len(supportedBrowsers))

	for i := range supportedBrowsers {
		info[i] = supportedBrowsers[i].GetInfo()
	}

	return info
}
