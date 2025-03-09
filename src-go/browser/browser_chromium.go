package browser

type ChromiumBrowser struct {
	*AbstractBrowser
}

// ToDo; Wip → implement this!
func GetChromiumBrowser() Browser {
	return &ChromiumBrowser{
		&AbstractBrowser{
			typicalName: "Chromium",
			processAliases: []string{
				// Found on:
				// - Manjaro 25.0.0 (Zetar)
				"chromium",
			},
			storageLocations: []string{
				// Found on:
				// - Manjaro 25.0.0 (Zetar)
				"~/.config/chromium/Default/Sessions/",
			},
		},
	}
}

// Yet unimplemented
func (browser *ChromiumBrowser) GetherWindowData() []WindowInfo {
	//browser.getSessionStorageLocation()
	// Do some magic~
	// Following things could be handy;
	// - https://github.com/JRBANCEL/Chromagnon/wiki/Reverse-Engineering-SSNS-Format
	// - https://github.com/lemnos/chrome-session-dump/tree/master
	// - https://github.com/phacoxcll/SNSS_Reader/tree/master
	return []WindowInfo{}
}
