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

func (browser *TesterBrowser) GetherWindowData() []WindowInfo {
	return []WindowInfo{}
}
