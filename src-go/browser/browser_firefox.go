package browser

type FirefoxBrowser struct {
	*AbstractBrowser
}

func GetFirefoxBrowser() Browser {
	return &FirefoxBrowser{
		&AbstractBrowser{
			typicalName: "Firefox",
			processAliases: []string{
				// Firefox GNU/Linux used on:
				// - Ubuntu …?
				// - Fedora …?
				"GeckoMain",
				// Firefox GNU/Linux used on:
				// - Manjaro 25.0.0 (Zetar)
				"firefox",
			},
			storageLocations: []string{
				// Firefox GNU/Linux seen on:
				// - Ubuntu ?
				// - Fedora ?
				// - Manjaro 25.0.0 (Zetar)
				"~/.mozilla/firefox*/*.default*/sessionstore-backups/recovery.jsonlz4",
			},
		},
	}
}

func (browser *FirefoxBrowser) GetInfo() BrowserInfo {
	return BrowserInfo{
		browser.typicalName,
		browser.isActive(),
		browser.getWindowData(),
	}
}

func (browser *FirefoxBrowser) getWindowData() []WindowInfo {
	_ = browser.getSessionStorageLocation()
	// Go some magic~
	return []WindowInfo{}
}
