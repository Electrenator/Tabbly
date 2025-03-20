package main

import (
	"embed"
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"github.com/Electrenator/Tabbly/src-go/browser"
	"github.com/Electrenator/Tabbly/src-go/storage"
	"github.com/Electrenator/Tabbly/src-go/util"
)

//go:embed database/*
var databasefiles embed.FS

func main() {
	settings := util.InitSettings()
	storage.Databasefiles = databasefiles
	storage.ApplicationSettings = &settings

	if !settings.Verbose {
		slog.SetLogLoggerLevel(slog.LevelWarn)
	}

	slog.Info(fmt.Sprintf("Application settings: %+v\n", settings))
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
			}
		}

		if settings.Verbose {
			slog.Info(fmt.Sprintf("Browsers: %+v\n", stats))
		}

		if !settings.DryRun {
			storage.SaveToCsv(stats)
			storage.SaveToDb(stats)
		}

		runtime.GC() // Can run GC if where going to sleep anyways
		time.Sleep(time.Second * time.Duration(settings.UpdateInterval))
	}

}
