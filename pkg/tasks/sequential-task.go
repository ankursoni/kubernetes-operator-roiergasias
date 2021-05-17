package tasks

import (
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/steps"
)

var _ TaskWorkflow = &SequentialTask{}

type SequentialTask struct {
	Task,
	Steps []steps.StepWorkflow
}

func NewSequentialTask(steps []steps.StepWorkflow) *SequentialTask {
	return &SequentialTask{
		Steps: steps,
	}
}

func (task *SequentialTask) Run() {
	for k := range task.Steps {
		task.Steps[k].Run()
	}
}
