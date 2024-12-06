package linksHelper

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

var keyValueDir string = "../../../resources/key_values"

func TestOpenKeyValuesFile(t *testing.T) {
	expectedFile := "../../../resources/key_values/odd_cmd.key-values.txt"
	file, err := OpenKeyValuesFile(keyValueDir, "odd", "cmd")

	if err != nil {
		t.Errorf("Received error while opening key-values file: %s", expectedFile)
		return
	}

	fileInfo, err := file.Stat()
	if err != nil {
		t.Errorf("Received error while reading file stats: %s", expectedFile)
		return
	}

	if fileInfo.IsDir() {
		t.Errorf("Key values file seems to be a directory: %s", expectedFile)
		return
	}

	_, err = OpenKeyValuesFile(keyValueDir, "dummy_odd", "cmd")

	expectedFile2 := "../../../resources/key_values/dummy_odd_cmd.key-values.txt"
	if err == nil {
		t.Errorf("Received error while opening key-values file: %s", expectedFile2)
		return
	}
}

func TestGetKeyValues(t *testing.T) {
	keyValues, err := GetKeyValues(keyValueDir, "odd", "cmd", []string{})
	if err != nil {
		t.Errorf("Received error while reading key-values: %v", err)
		return
	}

	expectedMapSize := 47
	mapSize := len(keyValues)
	if mapSize != expectedMapSize {
		t.Errorf("Retrieved map has the wrong size: got: %d, expected: %d", mapSize, expectedMapSize)
	}

	for k, v := range keyValues {
		slices.Sort(v)
		var lastE string
		for _, e := range v {
			if lastE == "" {
				lastE = e
			} else {
				if lastE == e {
					t.Errorf("Found duplicated attrib entries for: key=%s, value=%s", k, v)
				}
			}
		}
	}
}

func TestAreSrcAndDestAttribsTheSame(t *testing.T) {
	tests := []struct {
		name                   string
		sourceAttribsWithValue []string
		destAttrib             string
		expectedResult         bool
	}{
		{
			name:                   "Matching attributes",
			sourceAttribsWithValue: []string{"attrib1", "attrib2", "Attrib3"},
			destAttrib:             "Attrib2",
			expectedResult:         true,
		},
		{
			name:                   "No matching attributes",
			sourceAttribsWithValue: []string{"attrib1", "attrib2", "Attrib3"},
			destAttrib:             "Attrib4",
			expectedResult:         false,
		},
		{
			name:                   "Matching after harmonization",
			sourceAttribsWithValue: []string{"AttribA", "AttribB"},
			destAttrib:             "attribA", // Assuming HarmonizeName will normalize case
			expectedResult:         true,
		},
		{
			name:                   "Empty source attributes",
			sourceAttribsWithValue: []string{},
			destAttrib:             "Attrib1",
			expectedResult:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := srcAndDestAttribsAreTheSame(tt.sourceAttribsWithValue, tt.destAttrib)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestHarmonizeLinkAttribName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloWorld", "helloworld"},
		{"Hello World!", "hello_world_"},
		{"123-ABC_xyz", "abc_xyz"},
		{"123-ABC_xyz-", "123-abc_xyz-"},
	}

	for _, test := range tests {
		result := harmonizeLinkAttribName(test.input)
		assert.Equal(t, test.expected, result, "Expected %s but got %s", test.expected, result)
	}
}
