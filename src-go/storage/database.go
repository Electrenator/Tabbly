package storage

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
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

// Keeps track if the DB version has been checked since application start. Only
// allow a DB version 1 time per start
var checkedDbVersion = false

func SaveToDb(browserInfoList []browser.BrowserInfo) {
	db, err := connectToDb()
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		os.Exit(internal_status.UNSPECIFIED_PRIMARY_FUNCTION_ERROR)
	}
	defer db.Close()

	var version string
	db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)
	fmt.Println("DB version:", version)
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
	if err := os.MkdirAll(ApplicationSettings.DataPath, util.DefaultDirPerms); err != nil {
		slog.Error("Unable to application data directory!", "error", err)
		os.Exit(internal_status.FILE_CREATION_ERROR)
	}

	if file, err := os.OpenFile(getDbFileName(), os.O_CREATE, util.DefaultFilePerms); err != nil {
		slog.Error("Unable to create DB file!", "error", err)
		os.Exit(internal_status.FILE_CREATION_ERROR)
	} else {
		file.Close()
	}

	db, err := sql.Open("sqlite3", "file:"+getDbFileName())
	if err != nil {
		return nil, err
	}

	if !checkedDbVersion {
		latestMigrationVersion := getAvailableMigrationVersion()
		currentDbVersion := getCurrentDbSchemaVersion(db)

		if currentDbVersion < latestMigrationVersion {
			slog.Warn("Database out of date. Trying to migrate!",
				"currentVersion", currentDbVersion,
				"latestApplicationDb", latestMigrationVersion,
			)

			err := migrateDatabase(db, currentDbVersion)
			if err != nil {
				db.Close()
				slog.Error("Unable to migrate...", "error", err)
				os.Exit(internal_status.DB_MIGRATION_ERROR)
			}
		} else {
			slog.Info(fmt.Sprintf("Database on schema on version %d", currentDbVersion))
		}
	}
	return db, nil
}

// Get the latest migration version from `database/`. These files are
// formatted like `123-filename` with the version before `-`. The last
// of these files will be counted as the latest migration version.  files
// not in that format will be ignored with this.
//
// If no version files are found -1 is returned. Should never happen when
// it's correctly build.
func getAvailableMigrationVersion() int {
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

func getCurrentDbSchemaVersion(db *sql.DB) int {
	var version int
	err := db.QueryRow("SELECT `version` FROM `Database`").Scan(&version)
	if err != nil {
		return -1
	}
	return version
}

func migrateDatabase(db *sql.DB, fromVersion int) error {
	if fromVersion >= 0 {
		err := util.CopyFile(getDbFileName(), fmt.Sprintf("%s.v%d.bck", getDbFileName(), fromVersion))
		if err != nil {
			return errors.New(fmt.Sprintf("unable to create db backup: %s", err.Error()))
		}
	}

	migrationFiles, err := Databasefiles.ReadDir(dbStructureDirectory)
	if err != nil {
		return err
	}

	var affectedRowCount int64 = 0

	for _, file := range migrationFiles {
		fileVersion := getVersionFromFile(file.Name())

		if fromVersion >= fileVersion {
			continue
		}

		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") || fileVersion < 0 {
			slog.Warn("Non-SQL or migration file found in migration folder, skipping.",
				"filePath", file.Name(),
				"isDirectory", file.IsDir(),
			)
			continue
		}

		slog.Info(fmt.Sprintf("Running migration to version %d", fileVersion))

		fileData, err := Databasefiles.ReadFile(filepath.Join(dbStructureDirectory, file.Name()))
		if err != nil {
			return err
		}

		result, err := db.Exec(string(fileData))
		if err != nil {
			db.Exec("ROLLBACK")
			return err
		}
		affectedRows, _ := result.RowsAffected()
		affectedRowCount += affectedRows
	}
	slog.Info("Migration success!", "rowChanges", affectedRowCount)

	return err
}
