package tasks

import (
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/steps"
)

type ISequentialTasks interface {
	NewSequentialTask([]steps.StepWorkflow, string) *SequentialTask
}

type SequentialTasks struct{}

var _ ISequentialTasks = &SequentialTasks{}

func (_ *SequentialTasks) NewSequentialTask(steps []steps.StepWorkflow, node string) *SequentialTask {
	return &SequentialTask{
		Steps: steps,
		Task:  Task{Node: node},
	}
}

type SequentialTask struct {
	Steps []steps.StepWorkflow
	Task
}

var _ ITaskWorkflow = &SequentialTask{}

func (task *SequentialTask) Run() error {
	for k := range task.Steps {
		task.Steps[k].Run()
	}
	return nil
}
