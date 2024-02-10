package utils

import (
	"fmt"
	"os"
	"strings"
)

func GetStrVar(varName string, defaultValue string) string {
	value, exists := os.LookupEnv(varName)
	if !exists {
		return defaultValue
	}
	return value
}

func ReplaceWithEnvVar(s string, varName string, defaultValue string) string {
	txtToSearch := fmt.Sprintf("{%s}", varName)
	if strings.Contains(s, txtToSearch) {
		varValue := GetStrVar(varName, defaultValue)
		return strings.ReplaceAll(s, txtToSearch, varValue)
	} else {
		return s
	}
}
