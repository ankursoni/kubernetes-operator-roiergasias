package lib

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

type TaskType string

const (
	SequentialTaskType TaskType = "sequential"
)

type StepType string

const (
	EnvironmentStepType StepType = "environment"
	PrintStepType       StepType = "print"
	ExecuteStepType     StepType = "execute"
)

type Attribute string

const (
	NodeAttribute Attribute = "node"
)

var (
	regex = "{{env:[a-zA-Z0-9-_]+}}"
)

func SetEnvironmentVariables(inputList map[string]string) (err error) {
	for k := range inputList {
		if envErr := os.Setenv(k, inputList[k]); envErr != nil {
			err = fmt.Errorf("error setting environment: %w", envErr)
			return
		}
	}
	return
}

func ResolveEnvironmentVariablesInList(inputList []string) (outputList []string) {
	outputList = inputList
	regExp, _ := regexp.Compile(regex)
	for k := range inputList {
		matches := regExp.FindAllString(inputList[k], -1)
		for m := range matches {
			environmentVariableValue := os.Getenv(strings.TrimPrefix(strings.Trim(matches[m], "{}"), "env:"))
			outputList[k] = strings.ReplaceAll(inputList[k], matches[m], environmentVariableValue)
		}
	}
	return
}

func ResolveEnvironmentVariablesInMap(inputMap map[string]string) (outputMap map[string]string) {
	outputMap = inputMap
	regExp, _ := regexp.Compile(regex)
	for k := range outputMap {
		matches := regExp.FindAllString(outputMap[k], -1)
		for m := range matches {
			environmentVariableValue := os.Getenv(strings.TrimPrefix(strings.Trim(matches[m], "{}"), "env:"))
			outputMap[k] = strings.ReplaceAll(inputMap[k], matches[m], environmentVariableValue)
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
