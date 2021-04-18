package steps

import (
	"github.com/ankursoni/kubernetes-operator-roiergasias/lib"
)

type Step struct {
	StepArgumentList   []interface{}
	SetEnvironmentList map[string]interface{}
}

type StepWorkflow interface {
	Run()
}

func NewStep(stepType string, stepArguments []interface{}, otherStepArguments map[string]interface{}) (step StepWorkflow) {
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
		return
	default:
		return
	}
}

func (step *Step) Run() {
	lib.SetEnvironmentVariables(step.SetEnvironmentList)
}
