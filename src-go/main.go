package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Electrenator/Tabbly/src-go/browser"
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
		slog.Info(fmt.Sprintf("Thingy: %+v\n", browser.GetUsageInfo()))
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
