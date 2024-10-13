package linksHelper

import (
	"slices"
	"testing"
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
	keyValues, err := GetKeyValues(keyValueDir, "odd", "cmd")
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
