package steps

import (
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/lib"
)

type Step struct {
	StepArgumentList   []interface{}
	SetEnvironmentList map[string]interface{}
}

type IStepWorkflow interface {
	Run() error
}

func NewStep(stepType string, stepArguments []interface{}, otherStepArguments map[string]interface{}) (step IStepWorkflow) {
	newStep := &Step{}
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
	}
	return
}

func (step *Step) Run() (err error) {
	err = lib.SetEnvironmentVariables(step.SetEnvironmentList)
	return
}
