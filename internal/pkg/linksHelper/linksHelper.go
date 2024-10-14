package linksHelper

import (
	"bufio"
	"fmt"
	"okieoth/schemaguesser/internal/pkg/meta"
	"okieoth/schemaguesser/internal/pkg/utils"
	"os"
	"slices"
	"strings"
)

type ColDependency struct {
	Db             string
	Collection     string
	AttributeNames []string
}

type CollectionDependencies struct {
	CollectionInfo meta.MetaInfo
	Dependencies   []meta.MetaInfo
}

func NewColDependency(dbName string, collName string) ColDependency {
	return ColDependency{
		Db:             dbName,
		Collection:     collName,
		AttributeNames: make([]string, 0),
	}
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

func FoundKeyValue(keyValueDir string, dbName string, collName string, valueToFind string) (ColDependency, error) {

	ret := NewColDependency(dbName, collName)
	file, err := OpenKeyValuesFile(keyValueDir, dbName, collName)
	if err != nil {
		return ret, fmt.Errorf("error while open key-values file: dir=%s, db=%s, colName=%s", keyValueDir, dbName, collName)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			continue
		}
		// TODO
		// if (parts[1] == valueToFind) && (!slices.Contains(ret, parts[0])) {

		// 	ret = append(ret, parts[0])
		// }
	}

	if err := scanner.Err(); err != nil {
		return ret, fmt.Errorf("error while reading the file: %v", err)
	}

	return ret, nil
}
