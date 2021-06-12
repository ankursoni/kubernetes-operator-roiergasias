package tasks

import (
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/steps"
)

type ISequentialTasks interface {
	NewSequentialTask([]steps.IStepWorkflow, string) ITaskWorkflow
}

type SequentialTasks struct{}

var _ ISequentialTasks = &SequentialTasks{}

func NewSequentialTasks() (sequentialTasks ISequentialTasks) {
	sequentialTasks = &SequentialTasks{}
	return
}

func (_ *SequentialTasks) NewSequentialTask(steps []steps.IStepWorkflow, node string) ITaskWorkflow {
	return &SequentialTask{
		Steps: steps,
		Task:  Task{Node: node},
	}
}

type SequentialTask struct {
	Steps []steps.IStepWorkflow
	Task
}

var _ ITaskWorkflow = &SequentialTask{}

func (task *SequentialTask) Run() error {
	for k := range task.Steps {
		task.Steps[k].Run()
	}
	return nil
}
