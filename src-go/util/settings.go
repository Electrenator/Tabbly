package util

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	internal_status "github.com/Electrenator/Tabbly/src-go/internal/status"
	"github.com/spf13/pflag"
)

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
	verboseFlag := pflag.BoolP("verbose", "v", false, "Verbose logging output")
	dryRunFlag := pflag.Bool("dryrun", false, "Disable file writing")
	intervalFlag := pflag.Uint16("interval", 60, "Time between tab checks in seconds")
	applicationStorageLocation := getApplicationStorageLocation()

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

func getApplicationStorageLocation() string {
	const applicationName = "tabbly"
	var storageDirectory string
	var err error

	switch runtime.GOOS {
	case "windows", "darwin", "ios", "plan9":
		var baseDirectory string
		baseDirectory, err = os.UserConfigDir()
		storageDirectory = filepath.Join(baseDirectory, applicationName)
	default: // Unix
		var baseDirectory string
		baseDirectory, err = os.UserHomeDir()
		storageDirectory = filepath.Join(baseDirectory, "."+applicationName)
	}

	if err != nil {
		slog.Error("Error getting application directory", "error", err)
		os.Exit(internal_status.UNSPECIFIED_PRIMARY_FUNCTION_ERROR)
	}

	return storageDirectory
}
