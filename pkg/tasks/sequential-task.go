package tasks

import (
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/steps"
)

var _ TaskWorkflow = &SequentialTask{}

type SequentialTask struct {
	Steps []steps.StepWorkflow
	Task
}

func NewSequentialTask(steps []steps.StepWorkflow, node string) *SequentialTask {
	return &SequentialTask{
		Steps: steps,
		Task:  Task{Node: node},
	}
}

func (task *SequentialTask) Run() {
	for k := range task.Steps {
		task.Steps[k].Run()
	}
}
