package lib

import (
	"os"
	"regexp"
	"strings"
)

func SetEnvironmentVariables(inputList map[string]interface{}) {
	for k := range inputList {
		os.Setenv(k, inputList[k].(string))
	}
}

func ResolveEnvironmentVariables(inputList []string) (outputList []string) {
	outputList = inputList
	regExp, _ := regexp.Compile("{{env:[a-zA-Z0-9-_]+}}")
	for k := range inputList {
		matches := regExp.FindAllString(inputList[k], -1)
		for m := range matches {
			environmentVariableValue := os.Getenv(strings.TrimPrefix(strings.Trim(matches[m], "{}"), "env:"))
			outputList[k] = strings.ReplaceAll(inputList[k], matches[m], environmentVariableValue)
		}
	}
	return
}
