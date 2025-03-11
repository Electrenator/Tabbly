package main

import (
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"github.com/Electrenator/Tabbly/src-go/browser"
	"github.com/Electrenator/Tabbly/src-go/storage"
	"github.com/spf13/pflag"
)

type Settings struct {
	verbose        bool
	dryRun         bool
	updateInterval uint16
}

func main() {
	settings := initSettings()

	if !settings.verbose {
		slog.SetLogLoggerLevel(slog.LevelWarn)
	}

	slog.Info(fmt.Sprintf("Application settings: %+v\n", settings))
	for {
		availableBrowsers := browser.GetAvailableBrowsers()

		if settings.verbose {
			browser.LogAllBrowserStates()
		}

		stats := make([]browser.BrowserInfo, len(availableBrowsers))

		for _, browserInst := range availableBrowsers {
			if state := browserInst.GetState(); state == browser.BROWSER_OPEN || state == browser.BROWSER_STATE_UNKNOWN {
				slog.Info(browserInst.GetName())
				stats = append(stats, browser.BrowserInfo{
					Name:    browserInst.GetName(),
					IsOpen:  state,
					Windows: browserInst.GetherWindowData(),
				})
			}
		}

		storage.SaveCsv(stats)

		runtime.GC() // Can run GC if where going to sleep anyways
		time.Sleep(time.Second * time.Duration(settings.updateInterval))
	}
}

func initSettings() Settings {
	verboseFlag := pflag.BoolP("verbose", "v", false, "Verbose logging output")
	dryRunFlag := pflag.Bool("dryrun", false, "Verbose logging output")
	intervalFlag := pflag.Uint16("interval", 60, "Time between tab checks in seconds")

	pflag.Parse()
	settings := Settings{
		verbose:        *verboseFlag,
		dryRun:         *dryRunFlag,
		updateInterval: *intervalFlag,
	}
	return settings
}
