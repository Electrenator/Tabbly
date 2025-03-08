package browser

type TesterBrowser struct {
	*AbstractBrowser
}

// Testing browser → to be deleted
func GetTesterBrowser() Browser {
	return &TesterBrowser{
		&AbstractBrowser{
			typicalName:      "browser",
			processAliases:   nil,
			storageLocations: nil,
		},
	}
}

func (browser *TesterBrowser) GetInfo() BrowserInfo {
	return BrowserInfo{
		browser.typicalName,
		BROWSER_STATE_UNKNOWN,
		browser.getWindowData(),
	}
}

func (browser *TesterBrowser) getWindowData() []WindowInfo {
	return []WindowInfo{}
}
