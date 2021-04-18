package lib

import (
	"os"
	"strings"
)

func SetEnvironmentVariables(inputList map[string]interface{}) {
	for k := range inputList {
		os.Setenv(k, inputList[k].(string))
	}
}

func ResolveEnvironmentVariables(inputList []string) (outputList []string) {
	outputList = inputList
	environmentVariablePrefix := "env:"
	for k := range inputList {
		if strings.HasPrefix(inputList[k], environmentVariablePrefix) {
			environmentVariableValue := os.Getenv(strings.TrimPrefix(inputList[k], environmentVariablePrefix))
			outputList[k] = environmentVariableValue
		}
	}
	return
}
