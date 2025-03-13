package storage

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Electrenator/Tabbly/src-go/browser"
	internal_status "github.com/Electrenator/Tabbly/src-go/internal/status"
	"github.com/Electrenator/Tabbly/src-go/util"
	_ "github.com/mattn/go-sqlite3"
)

const dbStructureDirectory = "database"

var Databasefiles embed.FS
var ApplicationSettings *util.Settings

// Keeps track that the parent dirs have been verified to exist. Only run 1 time
// per application start
var checkedStorageParents = false

// Keeps track if the DB version has been checked since application start. Only
// allow a DB version 1 time per start
var checkedDbVersion = false

func SaveToDb(browserInfoList []browser.BrowserInfo) {

	walkFilesFunc(Databasefiles, false)

	db, err := connectToDb()
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		os.Exit(internal_status.UNSPECIFIED_PRIMARY_FUNCTION_ERROR)
	}

	var version string
	db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)
	fmt.Println("DB version:", version)

	db.Close()
}

func getDbFileName() string {
	const dbFileName = "test"
	const dbFileExt = ".sqlite"

	if ApplicationSettings.IsDevelopmentBuild {
		return filepath.Join(ApplicationSettings.DataPath, dbFileName+"-dev"+dbFileExt)
	}
	return dbFileExt + dbFileExt
}

// Makes a SQLite DB connection if possible. Initializes the database if it
// does not exist yet and (ToDo) runs migrations when the DB version is not
// up to date with the current application version
func connectToDb() (*sql.DB, error) {
	if !checkedStorageParents {
		// This path checks if the config directory & Tabbly DB exists, if not it creates it
		err := os.MkdirAll(ApplicationSettings.DataPath, util.DefaultDirPerms)

		if err == nil {
			var file *os.File
			file, err = os.OpenFile(getDbFileName(), os.O_CREATE, util.DefaultFilePerms)

			if err != nil {
				file.Close()
			}
		}
		if err != nil {
			slog.Error("Unable to create DB file!", "error", err)
			os.Exit(internal_status.FILE_CREATION_ERROR)
		}
	}

	db, err := sql.Open("sqlite3", "file:"+getDbFileName())
	if err != nil {
		return nil, err
	}

	if !checkedDbVersion {
		latestMigrationVersion := getLatestDbVersion()
		currentDbVersion := getCurrentDbVersion(db)

		if currentDbVersion < latestMigrationVersion {
			slog.Warn("Database out of date. Trying to migrate!",
				"currentVersion", currentDbVersion,
				"latestApplicationDb", latestMigrationVersion,
			)

			err := errors.New("migration unimplemented")
			if err != nil {
				slog.Error("Unable to migrate...", "error", err)
				os.Exit(internal_status.DB_MIGRATION_ERROR)
			}
		}

		fmt.Printf("Latest version; %d\n", getLatestDbVersion())
		fmt.Printf("Current DB version; %d\n", getCurrentDbVersion(db))
	}
	return db, nil
}

func walkFilesFunc(file fs.FS, function bool) {
	fs.WalkDir(file, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return err
		}
		fmt.Println(path, err)
		return err
	})
}

// Get the latest migration version from `database/`. These files are
// formatted like `123-filename` with the version before `-`. The last
// of these files will be counted as the latest migration version.  files
// not in that format will be ignored with this.
//
// If no version files are found -1 is returned. Should never happen when
// it's correctly build.
func getLatestDbVersion() int {
	files, err := Databasefiles.ReadDir(dbStructureDirectory)
	if err != nil {
		slog.Error("Unable to read db migration directory",
			"error", err,
		)
		return -1
	}
	for i := len(files) - 1; i >= 0; i-- {
		filename := files[i].Name()
		if version := getVersionFromFile(filename); version >= 0 {
			return version
		}
	}
	return -1
}

func getVersionFromFile(filename string) int {
	nameParts := strings.Split(filename, "-")

	if len(nameParts) <= 0 {
		return -1
	}
	version, err := strconv.Atoi(nameParts[0])
	if err != nil {
		return -1
	}

	return version
}

func getCurrentDbVersion(db *sql.DB) int {
	var version int
	err := db.QueryRow("SELECT `version` FROM `Database`").Scan(&version)
	if err != nil {
		return -1
	}
	return version
}
