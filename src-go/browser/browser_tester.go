package browser

type TesterBrowser struct {
	*AbstractBrowser
}

// Testing browser → to be deleted
func GetTesterBrowser() Browser {
	return &TesterBrowser{
		&AbstractBrowser{
			typicalName:              "browser",
			possibleApplicationNames: nil,
			storageLocations:         nil,
		},
	}
}

func (browser *TesterBrowser) GetInfo() BrowserInfo {
	var info BrowserInfo
	info.name = browser.typicalName
	info.isOpen = browser.isActive()
	info.windows = browser.getWindowData()
	return info
}
