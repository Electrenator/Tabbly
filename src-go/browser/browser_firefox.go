package browser

import "fmt"

type FirefoxBrowser struct {
	*AbstractBrowser
	test float32
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
			storageLocations: nil,
		},
		0.00001,
	}
}

func (browser *FirefoxBrowser) GetInfo() BrowserInfo {
	return BrowserInfo{
		browser.typicalName,
		browser.isActive(),
		browser.getWindowData(),
		fmt.Sprintf("%f", browser.test),
	}
}

func (browser *FirefoxBrowser) getWindowData() []WindowInfo {
	return []WindowInfo{}
}
