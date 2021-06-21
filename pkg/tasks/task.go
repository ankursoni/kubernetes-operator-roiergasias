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

func NewTasks(logger *zap.Logger) (tasks ITasks) {
	logger.Debug("creating new tasks")
	tasks = &Tasks{SequentialTasks: NewSequentialTasks(logger), Logger: logger}
	logger.Debug("successfully created new tasks")
	return
}

type ITaskWorkflow interface {
	Run() error
}

type Task struct {
	Node   string
	Logger *zap.Logger
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
		step := steps.NewStep(stepType, stepArguments, otherStepArguments, logger)
		if step == nil {
			return nil
		}
		task = t.SequentialTasks.NewSequentialTask([]steps.IStepWorkflow{step}, node)
	default:
		return nil
	}
	logger.Debug("successfully created new task")
	return
}
