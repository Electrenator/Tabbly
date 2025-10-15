package internal_storage_test

import (
	"os"
	"testing"
	"time"

	internal_setting "github.com/electrenator/tabbly/internal/setting"
	internal_storage "github.com/electrenator/tabbly/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestAsciiHostnameToCamelCase(test *testing.T) {
	cases := map[string]string{
		"hostname":                             "hostname",
		"HOSTNAME":                             "hostname",
		"Hostname":                             "hostname",
		"hostname-with-separator":              "hostnameWithSeparator",
		"HOSNAME_MORE_SCREEEEEMINGGG-SEPRATOR": "hosnameMoreScreeeeemingggSeprator",
		"Some-Intr_hostWITH,inValIdSEpERaToRS": "someIntrHostwith,invalidseperators",
	}

	for caseInput, caseExpectedOutput := range cases {
		output := internal_storage.AsciiHostnameToCamelCase(caseInput)
		assert.Equal(test, caseExpectedOutput, output)
	}
}

func TestCloseFile(test *testing.T) {
	const FILE_PATH = "/tmp/tabblyTestFile"
	test.Cleanup(cleanFile(test, FILE_PATH))

	file, err := os.OpenFile(FILE_PATH, os.O_CREATE|os.O_RDWR, internal_setting.DefaultFilePerms)
	if err != nil {
		test.Error(err)
		test.FailNow()
	}
	internal_storage.CloseFile(file)

	_, err = file.Write([]byte(":3"))
	assert.ErrorIs(test, err, os.ErrClosed)
}

func TestCopyFile(test *testing.T) {
	const FROM_FILE = "/tmp/tabblyFromFile"
	const TO_FILE = "/tmp/tabblyToFile"
	var content = "Cheese\n\t" + time.Now().String()

	test.Cleanup(cleanFile(test, FROM_FILE))
	test.Cleanup(cleanFile(test, TO_FILE))

	err := prepareFile(FROM_FILE, content)
	if err != nil {
		test.Error(err)
		test.FailNow()
	}

	err = internal_storage.CopyFile(FROM_FILE, TO_FILE)
	if err != nil {
		test.Error(err)
	}

	fileData, err := os.ReadFile(TO_FILE)
	if err != nil {
		test.Error(err)
	}

	assert.EqualValues(test, content, fileData)
}

func prepareFile(fileName string, contents string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, internal_setting.DefaultFilePerms)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(contents))
	if err == nil {
		err = file.Sync()
		if err == nil {
			err = file.Close()
		}
	}
	return err
}

func cleanFile(test *testing.T, filepath string) func() {
	cleanupFunction := func() {
		err := os.Remove(filepath)
		if err != nil {
			test.Error(err)
		}
	}
	return cleanupFunction
}
