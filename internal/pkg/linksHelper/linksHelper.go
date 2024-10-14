package linksHelper

import (
	"bufio"
	"fmt"
	"okieoth/schemaguesser/internal/pkg/utils"
	"os"
	"slices"
	"strings"
)

// Stores the reference details for one found reference
type AttribRefDetails struct {
	// database where the reference was found
	Db string
	// collection where the reference was found
	Collection string
	// list of attribute string that are referenced
	Attributes []string
}

// Stores the references of a single attribute
type AttribRef struct {
	AttribStr string
	// The different references to this attributes
	References []AttribRefDetails
}

// Stores all the references of a collection to other collections, over all databases
type ColRefs struct {
	Db         string
	Collection string
	AttribRefs []AttribRef
}

// This function read a key values file, extract the unique key values and return them
// as map, where the key value is key of the map and ...
func GetKeyValues(keyValueDir string, dbName string, collName string) (map[string][]string, error) {
	ret := make(map[string][]string, 0)
	file, err := OpenKeyValuesFile(keyValueDir, dbName, collName)
	if err != nil {
		return ret, fmt.Errorf("error while open key-values file: dir=%s, db=%s, colName=%s", keyValueDir, dbName, collName)
	}

	// 1. Read the file content per line
	// 2. split the line content by the string ': '
	// 3. check if the second part is already in 'ret' as key, if not insert it and put the first part of the match in the value array of 'ret

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		value := parts[1]

		if _, exists := ret[value]; !exists {
			ret[value] = []string{}
			ret[value] = append(ret[value], key)
		}
		if !slices.Contains(ret[value], key) {
			ret[value] = append(ret[value], key)
		}
	}

	if err := scanner.Err(); err != nil {
		return ret, fmt.Errorf("error while reading the file: %v", err)
	}

	return ret, nil
}

func OpenKeyValuesFile(keyValueDir string, dbName string, colName string) (*os.File, error) {
	filePath := utils.GetFileName(keyValueDir, "key-values.txt", dbName, colName)
	return os.Open(filePath)
}

func FoundKeyValue(keyValueDir string, destDbName string, destCollName string, valueToFind string, sourceAttribsWithValue []string, sourceDbName string, sourceCollName string, chIn chan<- ColRefs) error {
	file, err := OpenKeyValuesFile(keyValueDir, destDbName, destCollName)
	if err != nil {
		return fmt.Errorf("error while open key-values file: dir=%s, db=%s, colName=%s", keyValueDir, destDbName, destCollName)
	}
	defer file.Close()

	foundAttribs := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			continue
		}
		if (parts[1] == valueToFind) && (!slices.Contains(foundAttribs, parts[0])) {
			foundAttribs = append(foundAttribs, parts[0])
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error while reading the file: %v", err)
	}

	colRefs := new(ColRefs)
	colRefs.Db = sourceDbName
	colRefs.Collection = sourceCollName

	if len(foundAttribs) > 0 {
		for _, e := range sourceAttribsWithValue {
			ref := new(AttribRef)
			ref.AttribStr = e
			details := new(AttribRefDetails)
			details.Db = destDbName
			details.Collection = destCollName
			details.Attributes = foundAttribs
			ref.References = append(ref.References, *details)
			colRefs.AttribRefs = append(colRefs.AttribRefs, *ref)
		}
	}

	return nil
}
