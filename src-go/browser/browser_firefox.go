package browser

import "fmt"

type FirefoxBrowser struct {
	*AbstractBrowser
	test float32
}

func GetFirefoxBrowser() Browser {
	return &FirefoxBrowser{
		&AbstractBrowser{
			typicalName:              "Firefox",
			possibleApplicationNames: nil,
			storageLocations:         nil,
		},
		0.00001,
	}
}

func (browser *FirefoxBrowser) GetInfo() BrowserInfo {
	other := fmt.Sprintf("%f", browser.test)
	fmt.Println(other)
	return BrowserInfo{
		browser.typicalName,
		browser.isActive(),
		browser.getWindowData(),
		other,
	}
}
