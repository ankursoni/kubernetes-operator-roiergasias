package steps

import (
	"fmt"
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/lib"
	"go.uber.org/zap"
)

var _ IStepWorkflow = &PrintStep{}

type PrintStep struct {
	MessageList []string
	Step
}

func (step *Step) NewPrintStep() (printStep *PrintStep) {
	logger := step.Logger
	logger.Debug("creating new print step")
	messageList := []string{}
	for k := range step.StepArgumentList {
		messageList = append(messageList, step.StepArgumentList[k].(string))
	}
	messageList = lib.ResolveEnvironmentVariables(messageList)
	printStep = &PrintStep{
		MessageList: messageList,
		Step:        *step,
	}
	logger.Debug("successfully created new print step")
	return
}

func (printStep *PrintStep) Run() (err error) {
	logger := printStep.Logger
	logger.Debug("started running print step", zap.Any("message list", printStep.MessageList))
	for message := range printStep.MessageList {
		if _, printErr := fmt.Println(printStep.MessageList[message]); printErr != nil {
			err = fmt.Errorf("error printing step: %w", printErr)
			logger.Error(err.Error(), zap.Error(err))
			return
		}
	}
	err = printStep.Step.Run()
	logger.Debug("successfully ran print step")
	return
}
