package main

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/electrenator/tabbly/browser"
	internal_status "github.com/electrenator/tabbly/internal/status"
	"github.com/electrenator/tabbly/storage"
	"github.com/electrenator/tabbly/util"
)

//go:embed database/*
var databasefiles embed.FS

func main() {
	settings := util.InitSettings()
	storage.Databasefiles = databasefiles

	if !settings.Verbose {
		slog.SetLogLoggerLevel(slog.LevelWarn)
	}
	slog.Info(fmt.Sprintf("Application settings: %+v\n", settings))

	if settings.ShowVersion {
		fmt.Println(getApplicationVersion())
		os.Exit(0)
	} else {
		slog.Info(getApplicationVersion())
	}

	if settings.LegacyFileForImport != "" {
		storage.ImportLegacyCsv(settings.LegacyFileForImport)
		os.Exit(internal_status.OK)
	}

	for {
		availableBrowsers := browser.GetAvailableBrowsers()
		stats := []browser.BrowserInfo{}

		for _, browserInst := range availableBrowsers {
			if state := browserInst.GetState(); state == browser.BROWSER_OPEN ||
				state == browser.BROWSER_STATE_UNKNOWN {
				stats = append(stats, browser.BrowserInfo{
					Name:    browserInst.GetName(),
					IsOpen:  state,
					Windows: browserInst.GetherWindowData(),
				})
			} else {
				stats = append(stats, browser.BrowserInfo{
					Name:    browserInst.GetName(),
					IsOpen:  state,
					Windows: nil,
				})
			}
		}

		if settings.Verbose {
			slog.Info(fmt.Sprintf("Browsers: %+v\n", stats))
		}
		storage.SaveToDb(stats)

		runtime.GC() // Can run GC if where going to sleep anyways
		runtime.GC() // Also a second time to remove more marked objects from memory
		time.Sleep(time.Second * time.Duration(settings.UpdateInterval))
	}

}

// Formats a version information string to be printed or logged.
func getApplicationVersion() string {
	applicationName := "Tabbly"

	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		return fmt.Sprintf(
			"%s (%s)\n ├ Version: %s\n └ Go version: %s",
			applicationName,
			buildInfo.Path,
			buildInfo.Main.Version,
			buildInfo.GoVersion,
		)

	}
	return fmt.Sprintf("%s version unknown\n └ Go version: %s",
		applicationName,
		runtime.Version(),
	)
}
