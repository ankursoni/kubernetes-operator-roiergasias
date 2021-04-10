package steps

import (
	"fmt"
)

var _ StepWorkflow = &PrintStep{}

type PrintStep struct {
	Step,
	Messages []string
}

func NewPrintStep(stepArguments []interface{}) *PrintStep {
	messages := []string{}
	for k := range stepArguments {
		messages = append(messages, stepArguments[k].(string))
	}
	return &PrintStep{
		Messages: messages,
	}
}

func (step *PrintStep) Run() {
	for message := range step.Messages {
		fmt.Println(step.Messages[message])
	}
}
