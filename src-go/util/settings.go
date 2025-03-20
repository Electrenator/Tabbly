package util

import "github.com/spf13/pflag"

type Settings struct {
	Verbose            bool
	DryRun             bool
	IsDevelopmentBuild bool
	UpdateInterval     uint16
	DataPath           string
}

const DefaultDirPerms = 0755
const DefaultFilePerms = 0660

func InitSettings() Settings {
	const applicationStorageLocation = "."
	verboseFlag := pflag.BoolP("verbose", "v", false, "Verbose logging output")
	dryRunFlag := pflag.Bool("dryrun", false, "Verbose logging output")
	intervalFlag := pflag.Uint16("interval", 60, "Time between tab checks in seconds")

	pflag.Parse()
	settings := Settings{
		Verbose:            *verboseFlag,
		DryRun:             *dryRunFlag,
		IsDevelopmentBuild: isDevelopmentBuild,
		UpdateInterval:     *intervalFlag,
		DataPath:           applicationStorageLocation,
	}
	return settings
}
