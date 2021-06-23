package steps

import (
	"go.uber.org/zap"
)

type Step struct {
	StepArgumentList []interface{}
	Logger           *zap.Logger
}

type IStepWorkflow interface {
	Run() error
}

func NewStep(stepType string, stepArguments []interface{}, logger *zap.Logger) (
	step IStepWorkflow) {
	logger.Debug("creating new step", zap.String("step type", stepType),
		zap.Any("step arguments", stepArguments))
	newStep := &Step{Logger: logger}
	newStep.StepArgumentList = stepArguments

	switch stepType {
	case "environment":
		step = newStep.NewEnvironmentStep()
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
