package util

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	internal_status "github.com/electrenator/tabbly/internal/status"
	"github.com/spf13/pflag"
)

type Settings struct {
	Verbose             bool
	IsDevelopmentBuild  bool
	UpdateInterval      uint16
	DataPath            string
	LegacyFileForImport string
	PreferredDbSavePath string
	ShowVersion         bool
	CountSavedGroup     bool // If we should count the tabs from a saved tab group as another window
}

const DefaultDirPerms = 0755
const DefaultFilePerms = 0660

// Reference here, so others can use it. Will prevent circular dependances.
var AppSettings *Settings

// Define CLI arguments and build settings from received arguments.
func InitSettings() *Settings {
	applicationStorageLocation := getApplicationStorageLocation()
	verboseFlag := pflag.BoolP("verbose", "v", false, "Verbose logging output")
	intervalFlag := pflag.Uint16("interval", 60, "Time between tab checks in seconds")
	legacyFileToImport := pflag.String("import-legacy", "", "Legacy file to import into "+
		"application database. Not recommended to import into already existing database "+
		"files given it doesn't sort imported entries",
	)
	dbSaveLocation := pflag.String("db-location", "",
		"Override where the db will be saved. Handy in combination with '--import-legacy'",
	)
	showVersion := pflag.Bool("version", false, "Print the application version then exit")
	pflag.Parse()

	if *dbSaveLocation != "" {
		absoluteSavePath, err := filepath.Abs(*dbSaveLocation)
		if err != nil {
			slog.Error("Unable to expand db location", "path", *dbSaveLocation, "error", err)
			os.Exit(internal_status.UNSPECIFIED_PRIMARY_FUNCTION_ERROR)
		}
		dbSaveLocation = &absoluteSavePath
	}

	settings := Settings{
		Verbose:             *verboseFlag,
		IsDevelopmentBuild:  isDevelopmentBuild,
		UpdateInterval:      *intervalFlag,
		DataPath:            applicationStorageLocation,
		LegacyFileForImport: *legacyFileToImport,
		PreferredDbSavePath: *dbSaveLocation,
		ShowVersion:         *showVersion,
		CountSavedGroup:     *countSavedGroups,
	}
	AppSettings = &settings
	return &settings
}

// Finds the default storage location for this application. This is an OS
// specific path.
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
