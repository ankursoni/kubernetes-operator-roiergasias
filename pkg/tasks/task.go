package tasks

import (
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/steps"
	"go.uber.org/zap"
)

type ITasks interface {
	NewTask(string, map[string]interface{}, string) ITaskWorkflow
}

type Tasks struct {
	SequentialTasks ISequentialTasks
	Logger          *zap.Logger
}

var _ ITasks = &Tasks{}

func NewTasks(logger *zap.Logger) (tasks ITasks) {
	logger.Debug("creating new tasks")
	tasks = &Tasks{SequentialTasks: NewSequentialTasks(), Logger: logger}
	logger.Debug("successfully created new tasks")
	return
}

type ITaskWorkflow interface {
	Run() error
}

type Task struct {
	Node string
}

func (t *Tasks) NewTask(taskType string, stepData map[string]interface{}, node string) (task ITaskWorkflow) {
	logger := t.Logger
	logger.Debug("creating new task using arguments", zap.String("task type", taskType),
		zap.Any("step data", stepData), zap.String("node", node))
	keys := []string{}
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
		sequentialSteps := []steps.IStepWorkflow{}
		logger.Debug("creating new step for sequential task", zap.String("step type", stepType),
			zap.Any("step arguments", stepArguments), zap.Any("other step arguments", otherStepArguments))
		step := steps.NewStep(stepType, stepArguments, otherStepArguments)
		logger.Debug("successfully created new step for sequential task")

		sequentialSteps = append(sequentialSteps, step)
		logger.Debug("creating new sequential task", zap.Any("sequential steps", sequentialSteps),
			zap.String("node", node))
		task = t.SequentialTasks.NewSequentialTask(sequentialSteps, node)
		logger.Debug("successfully created new sequential task")
	}
	logger.Debug("successfully created new task")
	return
}
