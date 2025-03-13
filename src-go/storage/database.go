package storage

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/Electrenator/Tabbly/src-go/browser"
	internal_status "github.com/Electrenator/Tabbly/src-go/internal/status"
	"github.com/mattn/go-sqlite3"
)

const dbFileDirectory = "database"

var Databasefiles embed.FS

func SaveToDb(browserInfoList []browser.BrowserInfo) {

	walkFilesFunc(Databasefiles, false)
	fmt.Printf("Latest version; %d\n", getLatestDbVersion())

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

	return dbFileExt + dbFileExt
}

// Makes a SQLite DB connection if possible. Initializes the database if it
// does not exist yet and (ToDo) runs migrations when the DB version is not
// up to date with the current application version
func connectToDb() (*sql.DB, error) {
	// todo; create parent directories
	db, err := sql.Open("sqlite3", "file:"+getDbFileName())
	if err != nil {
		return nil, err
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
	files, err := Databasefiles.ReadDir(dbFileDirectory)
	if err != nil {
		slog.Error("Unable to read db migration directory",
			"error", err,
		)
		return -1
	}
	for i := len(files) - 1; i >= 0; i-- {
		filename := files[i].Name()
		nameParts := strings.Split(filename, "-")

		if len(nameParts) <= 0 {
			continue
		}
		if version, err := strconv.Atoi(nameParts[0]); err != nil {
			continue
		} else {
			return version
		}
	}
	return -1
}
