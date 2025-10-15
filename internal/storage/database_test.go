package internal_storage_test

import (
	"database/sql"
	"testing"

	internal_storage "github.com/electrenator/tabbly/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestGetCurrentDbSchemaVersion(test *testing.T) {
	db, err := sql.Open("sqlite3", "file:testData/v0.sqlite")
	if err != nil {
		test.Error(err)
	}
	version := internal_storage.GetCurrentDbSchemaVersion(db)

	assert.Equal(test, 0, version)
}

func TestGetCurrentDbSchemaVersionUninitialized(test *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		test.Error(err)
	}
	version := internal_storage.GetCurrentDbSchemaVersion(db)

	assert.Equal(test, -1, version)
}

func TestGetVersionFromFile(test *testing.T) {
	cases := map[string]int{
		"000-something.sql": 0,
		"010-newStuff!":     10,
		".901-invalid":      -1,
		"":                  -1,
		"901-aaa-4444":      901,
		"9-2-391-call":      9,
		"123456":            123456,
		"123456.sql":        -1,
		"nan-123456":        -1,
		"123456-solution":   123456,
	}

	for caseInput, caseExpectedOutput := range cases {
		output := internal_storage.GetVersionFromFile(caseInput)
		assert.Equal(test, caseExpectedOutput, output, "With input \"%s\"", caseInput)
	}
}
