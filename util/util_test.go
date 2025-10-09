package util_test

import (
	"os/user"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/electrenator/tabbly/util"
	"github.com/stretchr/testify/assert"
)

func TestStringContains(test *testing.T) {
	const TEST_STRING = "Bean salad 190213"
	const NOT_IN_STRING = "Zero"
	stringParts := strings.SplitSeq(TEST_STRING, "")

	for part := range stringParts {
		if !(util.StringContains(TEST_STRING, strings.ToUpper(part)) &&
			util.StringContains(TEST_STRING, strings.ToLower(part))) {
			test.Errorf("%s is marked not to be in %s but should be!", part, TEST_STRING)
			test.FailNow()
		}
	}
	if util.StringContains(TEST_STRING, strings.ToUpper(NOT_IN_STRING)) &&
		util.StringContains(TEST_STRING, strings.ToLower(NOT_IN_STRING)) {
		test.Errorf("%s should not be marked to be in %s", NOT_IN_STRING, TEST_STRING)
		test.FailNow()
	}
}

func TestSumSlice(test *testing.T) {
	const SUM = 7654321
	array := []int{1, 20, 300, 4000, 50000, 600000, 7000000}
	sliceSum := util.SumSlice(array)

	if sliceSum != SUM {
		test.Errorf("Sum slice was %d and should be %d", sliceSum, SUM)
	}
}

func TestSumSliceNegative(test *testing.T) {
	if util.SumSlice[int](nil) != -1 {
		test.Error("Empty slice should give -1 result")
		test.Fail()
	}
}

func pathExpansionConfigDir() string {
	switch runtime.GOOS {
	case "windows":
		return ""
	case "darwin", "ios":
		return "Library/Application Support"
	case "plan9":
		return "lib"
	default:
		return ".config"
	}
}

func TestExpandUserDirectory(test *testing.T) {
	configDir := pathExpansionConfigDir()
	usr, err := user.Current()
	if err != nil {
		test.Error(err)
	}

	test.Setenv("AppData", usr.HomeDir)
	cases := map[string]string{
		"~/directory":                    filepath.Join(usr.HomeDir, "directory"),
		"some/weird/~/directory":         "some/weird/~/directory",
		"%APPDATA%/directory":            filepath.Join(usr.HomeDir, configDir, "directory"),
		"some/weird/%APPDATA%/directory": "some/weird/%APPDATA%/directory",
	}

	for caseInput, caseOutput := range cases {
		path, err := util.ExpandUserDirectory(caseInput)
		if err != nil {
			test.Error(err)
		}

		if strings.Compare(caseOutput, path) != 0 {
			test.Errorf("Path expansion of '%s' should be '%s' but expanded to: %s",
				caseInput, caseOutput, path)
		}
	}
}

func TestSameSlice(test *testing.T) {
	var longValue []int = []int{98, 76, 54}
	var shortValue []int = []int{12, 34}
	assert := assert.New(test)

	assert.False(util.SameSlice(longValue, shortValue))
	assert.False(util.SameSlice(shortValue, longValue))
	assert.True(util.SameSlice(shortValue, shortValue))
	assert.True(util.SameSlice(longValue, longValue))

	// And the inverse
	inverseValue := slices.Clone(longValue)
	slices.Reverse(inverseValue)
	assert.False(util.SameSlice(longValue, inverseValue))
}

func TestSameSliceNill(test *testing.T) {
	var value []int = []int{12, 34}
	var noValue []int = nil
	assert := assert.New(test)

	assert.False(util.SameSlice(value, noValue))
	assert.False(util.SameSlice(noValue, value))
	assert.True(util.SameSlice(noValue, noValue))
}
