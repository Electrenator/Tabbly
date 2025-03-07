package main

import (
	"fmt"
	"log/slog"

	"github.com/Electrenator/Tabbly/src-go/browser"
	"github.com/spf13/pflag"
)

type Settings struct {
	verbose bool
	dryRun  bool
}

func main() {
	settings := initSettings()

	if !settings.verbose {
		slog.SetLogLoggerLevel(slog.LevelWarn)
	}

	slog.Info(fmt.Sprintf("Application settings: %+v\n", settings))
	slog.Info(fmt.Sprintf("Thingy: %+v\n", browser.GetUsageInfo()))
}

func initSettings() Settings {
	verboseFlag := pflag.BoolP("verbose", "v", false, "Verbose logging output")
	dryRunFlag := pflag.Bool("dryrun", false, "Verbose logging output")

	pflag.Parse()
	settings := Settings{
		verbose: *verboseFlag,
		dryRun:  *dryRunFlag,
	}
	return settings
}
