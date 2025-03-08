package browser

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/Electrenator/Tabbly/src-go/util"
	"github.com/dustin/go-humanize"
	"github.com/giulianopz/go-dejsonlz4/jsonlz4"
	"github.com/itchyny/gojq"
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
			},
		},
	}
}

func (browser *FirefoxBrowser) GetInfo() BrowserInfo {
	var windowData []WindowInfo
	active := browser.isActive()

	if active == BROWSER_OPEN {
		windowData = browser.getWindowData()
	}
	return BrowserInfo{
		browser.typicalName,
		active,
		windowData,
	}
}

func (browser *FirefoxBrowser) getWindowData() []WindowInfo {
	var windowData []WindowInfo

	for _, path := range browser.getSessionStorageLocation() {
		fileData, err := os.ReadFile(path)
		if err != nil {
			slog.Error(
				"Can't read file",
				"path", path,
				"error", err,
			)
			continue
		}

		data, err := jsonlz4.Uncompress(fileData)
		if err != nil {
			slog.Error(
				"Can't decompress file for reading",
				"path", path,
				"error", err,
			)
			continue
		}
		slog.Info(fmt.Sprintf(
			"Decompressed %s of size %s (%s uncompressed)",
			path,
			humanize.IBytes(uint64(len(fileData))),
			humanize.IBytes(uint64(len(data))),
		))

		jsonValues := make(map[string]interface{})
		json.Unmarshal(data, &jsonValues)

		query, err := gojq.Parse(".windows | map(.tabs | length)")
		if err != nil {
			slog.Error("Unable to parse jq", "error", err)
			continue
		}

		resultIterator := query.Run(jsonValues)
		for {
			resultItem, ok := resultIterator.Next()
			if !ok {
				break
			}

			resultValues := util.ConvertSlice[int](resultItem.([]any))
			resultWindows := make([]WindowInfo, len(resultValues))

			for i, tab := range resultValues {
				resultWindows[i] = WindowInfo{tab}
			}
			windowData = append(windowData, resultWindows...)
		}
	}

	return windowData
}
