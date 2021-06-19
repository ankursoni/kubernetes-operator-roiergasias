package lib

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

func SetEnvironmentVariables(inputList map[string]interface{}) (err error) {
	for k := range inputList {
		if envErr := os.Setenv(k, inputList[k].(string)); envErr != nil {
			err = fmt.Errorf("error setting environment: %w", envErr)
			return
		}
	}
	return
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

func NewZapLogger(debug bool) (logger *zap.Logger, err error) {
	if debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	return
}