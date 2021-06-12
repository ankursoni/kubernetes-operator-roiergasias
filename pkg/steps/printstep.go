package steps

import (
	"fmt"

	"github.com/ankursoni/kubernetes-operator-roiergasias/lib"
)

var _ IStepWorkflow = &PrintStep{}

type PrintStep struct {
	MessageList []string
	Step
}

func (step *Step) NewPrintStep() (printStep *PrintStep) {
	messageList := []string{}
	for k := range step.StepArgumentList {
		messageList = append(messageList, step.StepArgumentList[k].(string))
	}
	messageList = lib.ResolveEnvironmentVariables(messageList)
	printStep = &PrintStep{
		MessageList: messageList,
		Step:        *step,
	}
	return
}

func (printStep *PrintStep) Run() {
	for message := range printStep.MessageList {
		fmt.Println(printStep.MessageList[message])
	}
	printStep.Step.Run()
}
