package tasks

import (
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/steps"
)

type Task struct {
}

type TaskWorkflow interface {
	Run()
}

func NewTask(taskType string, stepsData map[string]interface{}) (task TaskWorkflow) {
	switch taskType {
	case "sequential":
		var sequentialSteps []steps.StepWorkflow
		for stepType := range stepsData {
			step := steps.NewStep(stepType, stepsData[stepType].([]interface{}))
			sequentialSteps = append(sequentialSteps, step)
		}
		task = NewSequentialTask(sequentialSteps)
		return
	default:
		return
	}
}
