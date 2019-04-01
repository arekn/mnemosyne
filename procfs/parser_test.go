package procfs

import (
	"github.com/arekn/mnemosyne"
	"strings"
	"testing"
)

func TestParseMemInfo_should_parseLinesWithoutError(t *testing.T) {
	// given
	const validFile = `Name:	cat
Groups:	4 24 27 30 46 118 128 999 1000 
VmSize:	    8828 kB
`
	// when
	result, e := main.ParseProcFile(strings.NewReader(validFile))

	// then
	if len(result) != 3 {
		t.Errorf("result size was %v, but expected to be 6", len(result))
	}
	if e != nil {
		t.Error(e)
	}
}

func TestParseMemInfo_should_returnProperKeyValuePair(t *testing.T) {
	// given
	const expectedKey = "catName"
	const expectedValue = "dumpling"
	const testFile = expectedKey + ":" + expectedValue

	// when
	result, e := main.ParseProcFile(strings.NewReader(testFile))

	// then
	if result[expectedKey] != expectedValue {
		t.Errorf("returned value for key %s was %s, but expected was %s", expectedKey, result[expectedKey], expectedValue)
	}
	if e != nil {
		t.Error(e)
	}
}

func TestParseMemInfo_should_returnErrorWhenSeparatorIsMissing(t *testing.T) {
	// given
	const testFile = "some test file"

	// when
	result, e := main.ParseProcFile(strings.NewReader(testFile))

	// then
	if result != nil {
		t.Errorf("returned value should be nil")
	}
	if e == nil {
		t.Error("error is missing")
	}
}
