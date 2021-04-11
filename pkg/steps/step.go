package steps

import (
	"os"
	"strings"
)

type Step struct {
}

type StepWorkflow interface {
	Run()
}

func NewStep(stepType string, stepArgumentList []interface{}) (step StepWorkflow) {
	switch stepType {
	case "print":
		messageList := []string{}
		for k := range stepArgumentList {
			messageList = append(messageList, stepArgumentList[k].(string))
		}
		messageList = ResolveEnvironmentVariables(messageList)
		step = NewPrintStep(messageList)
		return
	default:
		return
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
