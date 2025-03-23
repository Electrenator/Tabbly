package storage

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Electrenator/Tabbly/src-go/browser"
	internal_status "github.com/Electrenator/Tabbly/src-go/internal/status"
	"github.com/Electrenator/Tabbly/src-go/util"
	_ "github.com/mattn/go-sqlite3"
)

const dbStructureDirectory = "database"

var Databasefiles embed.FS

// Keeps track if the DB version has been checked since application start. Only
// allow a DB version 1 time per start
var checkedDbSchemaVersion = false

type TimedBrowserInfo struct {
	Timestamp   int64
	BrowserInfo browser.BrowserInfo
}

// Save a single entry with multiple browsers to the database at the current timestamp.
func SaveToDb(browserInfoList []browser.BrowserInfo) {
	db, err := connectToDb()
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		os.Exit(internal_status.DB_CONNECT_ERROR)
	}
	defer db.Close()

	db.Exec("BEGIN TRANSACTION")
	err = saveToDbUsingConnection(time.Now().Unix(), browserInfoList, db)

	if err != nil {
		dbErrorRollback(db, err)
	}
	db.Exec("COMMIT")
}

// Save multiple info entries from different timestamps to the database from **1**
// browser per entry.
func SaveMultipleToDb(timeBrowserInfoMap *[]TimedBrowserInfo) error {
	db, err := connectToDb()
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		os.Exit(internal_status.DB_CONNECT_ERROR)
	}
	defer db.Close()

	db.Exec("BEGIN TRANSACTION")

	for _, timedBrowserInfo := range *timeBrowserInfoMap {
		err := saveToDbUsingConnection(
			timedBrowserInfo.Timestamp,
			[]browser.BrowserInfo{timedBrowserInfo.BrowserInfo},
			db,
		)
		if err != nil {
			dbErrorRollback(db, err)
			return err
		}
	}
	db.Exec("COMMIT")
	return nil
}

// Save entry to the database at a specified timestamp. This can re-use an existing DB connection.
//
// Todo; Make smarter. Only insert changed browsers. Also add an empty browser entry when the browser closed
func saveToDbUsingConnection(time int64, browserInfoList []browser.BrowserInfo, db *sql.DB) error {

	for _, browser := range browserInfoList {
		var browserId int64
		err := db.QueryRow("SELECT `id` FROM `Browser` WHERE `name` == ?", browser.Name).Scan(&browserId)

		if err == sql.ErrNoRows {
			result, insertErr := db.Exec("INSERT INTO `Browser` (`name`) VALUES (?)", browser.Name)

			if insertErr != nil {
				err = insertErr
			} else {
				browserId, err = result.LastInsertId()
			}
		}
		if err != nil {
			dbErrorRollback(db, err)
			return err
		}

		result, err := db.Exec(
			"INSERT INTO `Entry` (`timestamp`, `browserId`) VALUES (?, ?)",
			time, browserId,
		)
		if err != nil {
			dbErrorRollback(db, err)
			return err
		}

		entryId, err := result.LastInsertId()
		if err != nil {
			dbErrorRollback(db, err)
			return err
		}

		for _, window := range browser.Windows {
			_, err = db.Exec(
				"INSERT INTO `Window` (`entryId`, `openTabs`) VALUES (?, ?)",
				entryId, window.TabCount,
			)
			if err != nil {
				dbErrorRollback(db, err)
				return err
			}
		}
	}
	return nil
}

// Returns the DB file path in the following format where hostname will be
// replaced with the current devices name
//
//	[DataPath]/tabs-hostname[-dev].sqlite
func getDbFilePath() string {
	const dbName = "tabs"
	const dbFileExt = ".sqlite"

	if ApplicationSettings.PreferredDbSavePath != "" {
		return ApplicationSettings.PreferredDbSavePath
	}
	hostname, err := os.Hostname()

	if err != nil {
		hostname = ""
	}

	hostname = asciiHostnameToCamelCase(hostname)
	tmpFileName := fmt.Sprintf("%s-%s", dbName, hostname)

	if ApplicationSettings.IsDevelopmentBuild {
		tmpFileName += developmentFileAddition
	}
	return filepath.Join(ApplicationSettings.DataPath, tmpFileName+dbFileExt)
}

// Makes a SQLite DB connection if possible. Initializes the database if it
// does not exist yet and runs migrations when the DB version is not up to
// date with the current application version.
func connectToDb() (*sql.DB, error) {
	if err := os.MkdirAll(ApplicationSettings.DataPath, util.DefaultDirPerms); err != nil {
		slog.Error("Unable to application data directory!", "error", err)
		os.Exit(internal_status.FILE_CREATION_ERROR)
	}

	if file, err := os.OpenFile(getDbFilePath(), os.O_CREATE, util.DefaultFilePerms); err != nil {
		slog.Error("Unable to create DB file!", "error", err)
		os.Exit(internal_status.FILE_CREATION_ERROR)
	} else {
		file.Close()
	}

	db, err := sql.Open("sqlite3", "file:"+getDbFilePath())
	if err != nil {
		return nil, err
	}

	if !checkedDbSchemaVersion {
		latestMigrationVersion := getAvailableMigrationVersion()
		currentDbSchemaVersion := getCurrentDbSchemaVersion(db)

		var version string
		db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)
		slog.Info(fmt.Sprintln("SQLite version:", version))

		if currentDbSchemaVersion < latestMigrationVersion {
			slog.Warn("Database out of date. Going to migrate!",
				"schemaVersion", currentDbSchemaVersion,
				"latestSchema", latestMigrationVersion,
			)

			err := migrateDatabase(db, currentDbSchemaVersion)
			if err != nil {
				db.Close()
				slog.Error("Unable to migrate...", "error", err)
				os.Exit(internal_status.DB_MIGRATION_ERROR)
			}
		} else {
			slog.Info(fmt.Sprintf("Database on schema on version %d", currentDbSchemaVersion))
		}
		checkedDbSchemaVersion = true
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
		err := copyFile(getDbFilePath(), fmt.Sprintf("%s.v%d.bck", getDbFilePath(), fromVersion))
		if err != nil {
			return fmt.Errorf("unable to create db backup: %s", err.Error())
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
			dbErrorRollback(db, err)
			return err
		}
		affectedRows, _ := result.RowsAffected()
		affectedRowCount += affectedRows
	}
	slog.Info("Migration success!",
		"rowChanges", affectedRowCount,
		"currentSchemaVersion", getCurrentDbSchemaVersion(db),
	)

	return err
}

func dbErrorRollback(db *sql.DB, err error) {
	slog.Error("Error interacting with database, aborting interaction!", "error", err)
	db.Exec("ROLLBACK")
}
