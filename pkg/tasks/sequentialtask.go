package tasks

import (
	"fmt"
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

func (task *SequentialTask) Run() (err error) {
	for k := range task.Steps {
		if taskErr := task.Steps[k].Run(); taskErr != nil {
			err = fmt.Errorf("error running task: %w", taskErr)
			return
		}
	}
	return
}
