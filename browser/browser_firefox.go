package browser

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/electrenator/tabbly/util"
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

func (browser *FirefoxBrowser) GetherWindowData() []WindowInfo {
	var windowData []WindowInfo

	for _, path := range browser.getSessionStorageLocation() {
		jsonValues, err := parseSessionBackup(path)
		if err != nil {
			// Already logged in above function!
			continue
		}

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

func parseSessionBackup(path string) (map[string]any, error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		slog.Error(
			"Can't read file",
			"path", path,
			"error", err,
		)
		return nil, err
	}

	data, err := jsonlz4.Uncompress(fileData)
	if err != nil {
		slog.Error(
			"Can't decompress file for reading",
			"path", path,
			"error", err,
		)
		return nil, err
	}
	slog.Info(fmt.Sprintf(
		"Decompressed %s of size %s (%s uncompressed)",
		path,
		humanize.IBytes(uint64(len(fileData))),
		humanize.IBytes(uint64(len(data))),
	))

	jsonValues := make(map[string]any)
	// Following operation somehow massively increases application memory usage
	// What's happening in that Unmarshal that makes the application go from ±11
	// to 75MiB (depending on session file) at sleep back in main?
	if err := json.Unmarshal(data, &jsonValues); err != nil {
		slog.Error(
			"Can't json decode file",
			"file", path,
			"error", err,
		)
		return nil, err
	}
	return jsonValues, nil
}
