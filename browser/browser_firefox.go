package browser

import (
	"log/slog"

	"gitlab.com/electrenator/mozilla-lz4-decoder/mozillaLz4"
)

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
				"~/.mozilla/firefox*/*.default-release/sessionstore-backups/recovery.jsonlz4",
				// Found on:
				// - Windows 11
				"%APPDATA%/Mozilla/Firefox/Profiles/*.default-release-*/sessionstore-backups/recovery.jsonlz4",
			},
		},
	}
}

func (browser *FirefoxBrowser) GetherWindowData() []WindowInfo {
	var windowData []WindowInfo

	for _, path := range browser.getSessionStorageLocation() {
		browserData, err := mozillaLz4.Read(path)
		if err != nil {
			// Already logged in above function!
			slog.Error(
				"Error while parsing session file",
				"path", path,
				"error", err,
			)
			continue
		}

		resultWindows := make([]WindowInfo, len(browserData.Windows))

		for i, window := range browserData.Windows {
			resultWindows[i] = WindowInfo{len(window.Tabs)}
		}
		windowData = append(windowData, resultWindows...)

		if true { // TODO; cli option for enabling this (disabled by default)
			savedGroupTabCount := browser.getSavedGroupTotalTabCount(browserData.SavedGroups)
			if savedGroupTabCount > 0 {
				windowData = append(windowData, WindowInfo{savedGroupTabCount})
			}
		}
	}

	return windowData
}

func (browser *FirefoxBrowser) getSavedGroupTotalTabCount(groupList []mozillaLz4.ClosedTabGroup) int {
	var count int
	for _, group := range groupList {
		if group.Saved == nil && !*group.Saved {
			continue
		}
		if !*group.Saved {
			continue
		}
		count += len(group.Tabs)
	}
	return count
}
