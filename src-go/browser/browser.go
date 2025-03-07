package browser

type Browser interface {
	GetInfo() BrowserInfo
	isActive() bool
	getWindowData() []WindowInfo
}

type AbstractBrowser struct {
	Browser
	typicalName              string
	possibleApplicationNames []string
	storageLocations         []string
}

func (browser *AbstractBrowser) GetInfo() {
	panic("Unimplemented!")
}

func (browser *AbstractBrowser) isActive() bool {
	// panic("Unimplemented")
	return false
}

func (browser *AbstractBrowser) getWindowData() []WindowInfo {
	// var data []WindowInfo

	// ToDo
	return []WindowInfo{}
}
