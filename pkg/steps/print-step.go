package steps

import (
	"fmt"
)

var _ StepWorkflow = &PrintStep{}

type PrintStep struct {
	Step,
	MessageList []string
}

func NewPrintStep(messageList []string) *PrintStep {
	return &PrintStep{
		MessageList: messageList,
	}
}

func (step *PrintStep) Run() {
	for message := range step.MessageList {
		fmt.Println(step.MessageList[message])
	}
}
