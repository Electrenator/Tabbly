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
					// Found on:
					// - Manjaro 25.0.0 (Zetar)
					"~/.mozilla/firefox*/*.dev-edition-default-*/sessionstore-backups/recovery.jsonlz4",
					// Found on:
					// - Windows 11
					"%APPDATA%/Mozilla/Firefox/Profiles/*.dev-edition-default/sessionstore-backups/recovery.jsonlz4",
				},
			},
		},
	}
}
