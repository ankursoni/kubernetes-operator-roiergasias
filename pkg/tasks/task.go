package tasks

import (
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/steps"
)

type Task struct {
	Node string
}

type TaskWorkflow interface {
	Run()
}

func NewTask(taskType string, stepData map[string]interface{}, node string) (task TaskWorkflow) {
	var keys []string
	for k := range stepData {
		keys = append(keys, k)
	}
	stepType := keys[0]
	var stepArguments []interface{}
	otherStepArguments := make(map[string]interface{})
	for i := 0; i < len(keys); i++ {
		if i == 0 {
			stepArguments = stepData[keys[i]].([]interface{})
		} else {
			otherStepData := stepData[keys[i]].([]interface{})
			for otherStepArgumentType := range otherStepData {
				otherStepArguments[keys[i]] = otherStepData[otherStepArgumentType]
			}
		}
	}

	switch taskType {
	case "sequential":
		var sequentialSteps []steps.StepWorkflow
		step := steps.NewStep(stepType, stepArguments, otherStepArguments)
		sequentialSteps = append(sequentialSteps, step)
		task = NewSequentialTask(sequentialSteps, node)
		return
	default:
		return
	}
}
