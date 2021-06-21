package tasks

import (
	"fmt"
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/steps"
	"go.uber.org/zap"
)

type ISequentialTasks interface {
	NewSequentialTask([]steps.IStepWorkflow, string) ITaskWorkflow
}

type SequentialTasks struct {
	Logger *zap.Logger
}

func NewSequentialTasks(logger *zap.Logger) (sequentialTasks ISequentialTasks) {
	logger.Debug("creating new sequential tasks")
	sequentialTasks = &SequentialTasks{Logger: logger}
	logger.Debug("successfully created new sequential tasks")
	return
}

func (st *SequentialTasks) NewSequentialTask(steps []steps.IStepWorkflow, node string) (sequentialTask ITaskWorkflow) {
	logger := st.Logger
	logger.Debug("creating new sequential task", zap.Any("steps", steps), zap.String("node", node))
	sequentialTask = &SequentialTask{
		Steps: steps,
		Task:  Task{Node: node, Logger: logger},
	}
	logger.Debug("successfully created new sequential task")
	return
}

type SequentialTask struct {
	Steps []steps.IStepWorkflow
	Task
}

var _ ITaskWorkflow = &SequentialTask{}

func (sequentialTask *SequentialTask) Run() (err error) {
	logger := sequentialTask.Logger
	logger.Debug("started running sequential task", zap.Any("steps", sequentialTask.Steps))
	for k := range sequentialTask.Steps {
		if taskErr := sequentialTask.Steps[k].Run(); taskErr != nil {
			err = fmt.Errorf("error running sequential task: %w", taskErr)
			logger.Error(err.Error(), zap.Error(err))
			return
		}
	}
	logger.Debug("successfully ran sequential task")
	return
}
