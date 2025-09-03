package browser

import (
	"fmt"
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

		groupWindowMap := browser.getSavedGroupWindowMap(browserData)
		fmt.Printf("Recieved map: %+v\n", groupWindowMap)

		resultWindows := make([]WindowInfo, len(browserData.Windows))

		for i, window := range browserData.Windows {
			windowTabCount := len(window.Tabs)
			fmt.Println(window.Title)
			fmt.Printf("workspace id %s\n", *window.WorkspaceId)

			// Also count saved groups if available
			fmt.Printf("Recieved gruops: %+v\n", groupWindowMap[window.LastSessionWindowId])
			if savedTabList, exists := groupWindowMap[window.LastSessionWindowId]; exists {
				fmt.Printf("%+v", len(savedTabList))
				for _, savedTabs := range savedTabList {
					fmt.Printf("%+v", len(*savedTabs))
					windowTabCount += len(*savedTabs)
				}
				delete(groupWindowMap, window.LastSessionWindowId)
			}
			resultWindows[i] = WindowInfo{windowTabCount}
		}

		windowData = append(windowData, resultWindows...)
	}

	return windowData
}

func (browser *FirefoxBrowser) getSavedGroupWindowMap(browserData *mozillaLz4.MozillaRecoveryFormat) map[string][]*[]mozillaLz4.ClosedTab {
	groupWindowMap := map[string][]*[]mozillaLz4.ClosedTab{}

	for i := 0; i < len(browserData.SavedGroups); i++ {
		savedGroup := &browserData.SavedGroups[i]

		if !*savedGroup.Saved {
			continue
		}
		if tabList, exists := groupWindowMap[savedGroup.SourceWindowId.String()]; !exists {
			fmt.Println("TABLIST Reused")
			groupWindowMap[savedGroup.SourceWindowId.String()] = append(tabList, &savedGroup.Tabs)
			fmt.Printf("%+v", groupWindowMap[savedGroup.SourceWindowId.String()])
		} else {
			tabList := make([]*[]mozillaLz4.ClosedTab, 1)
			fmt.Println("TABLIST MADE")
			groupWindowMap[savedGroup.SourceWindowId.String()] = append(tabList, &savedGroup.Tabs)
		}
	}
	fmt.Printf("GeneratedGroups: %+v\n", groupWindowMap)
	return groupWindowMap
}
