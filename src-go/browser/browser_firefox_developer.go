package browser

type FirefoxDeveloperBrowser struct {
	*FirefoxBrowser
}

func GetFirefoxDeveloperBrowser() Browser {
	return &FirefoxDeveloperBrowser{
		&FirefoxBrowser{
			&AbstractBrowser{
				typicalName:    "Firefox Developer Edition",
				processAliases: nil, // Can't distinguish it from base Firefox it seems
				storageLocations: []string{
					// Found on:
					// - Manjaro 25.0.0 (Zetar)
					"~/.mozilla/firefox*/*.dev-edition-default/sessionstore-backups/recovery.jsonlz4",
				},
			},
		},
	}
}

func (browser *FirefoxDeveloperBrowser) GatherInfo() BrowserInfo {
	return BrowserInfo{
		browser.typicalName,
		BROWSER_STATE_UNKNOWN,
		browser.GetherWindowData(),
	}
}
