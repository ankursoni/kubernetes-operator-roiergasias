package steps

import (
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/lib"
	"go.uber.org/zap"
)

type Step struct {
	StepArgumentList   []interface{}
	SetEnvironmentList map[string]interface{}
	Logger             *zap.Logger
}

type IStepWorkflow interface {
	Run() error
}

func NewStep(stepType string, stepArguments []interface{}, otherStepArguments map[string]interface{}, logger *zap.Logger) (
	step IStepWorkflow) {
	logger.Debug("creating new step", zap.String("step type", stepType),
		zap.Any("step arguments", stepArguments), zap.Any("other step arguments", otherStepArguments))
	newStep := &Step{Logger: logger}
	newStep.StepArgumentList = stepArguments
	for k := range otherStepArguments {
		switch k {
		case "set-environment":
			newStep.SetEnvironmentList = otherStepArguments[k].(map[string]interface{})
		}
	}

	switch stepType {
	case "print":
		step = newStep.NewPrintStep()
	case "execute":
		step = newStep.NewExecuteStep()
	default:
		return nil
	}
	logger.Debug("successfully created new step")
	return
}

func (step *Step) Run() (err error) {
	logger := step.Logger
	logger.Debug("started running base step's set environment list",
		zap.Any("environment list", step.SetEnvironmentList))
	err = lib.SetEnvironmentVariables(step.SetEnvironmentList)
	logger.Debug("successfully ran base step's set environment list")
	return
}
