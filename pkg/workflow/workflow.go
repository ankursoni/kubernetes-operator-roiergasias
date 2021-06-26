package workflow

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"reflect"
	"strconv"

	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/lib"
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/tasks"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type IWorkflows interface {
	NewWorkflow(string) (*Workflow, error)
	NewWorkflowFromText(string) (*Workflow, error)
}

type Workflows struct {
	Tasks  tasks.ITasks
	Logger *zap.Logger
}

func NewWorkflows(t tasks.ITasks, logger *zap.Logger) (workflows IWorkflows) {
	if logger == nil {
		newLogger, err := lib.NewZapLogger(false)
		if err != nil {
			log.Fatalln(fmt.Errorf("error creating new zap logger: %w", err))
			return
		} else {
			logger = newLogger
		}
	}
	logger.Debug("creating new workflows using dependencies", zap.Any("ITasks", t))
	if t == nil {
		workflows = &Workflows{Tasks: tasks.NewTasks(logger), Logger: logger}
	} else {
		workflows = &Workflows{Tasks: t, Logger: logger}
	}
	logger.Debug("successfully created new workflows", zap.Any("Workflows", workflows))
	return
}

func (w *Workflows) NewWorkflow(filePath string) (workflow *Workflow, err error) {
	logger := w.Logger
	logger.Debug("reading yaml file", zap.String("path", filePath))
	bytes, readErr := ioutil.ReadFile(filePath)
	if readErr != nil {
		err = fmt.Errorf("error reading yaml file: %w", readErr)
		logger.Error(err.Error(), zap.Error(err))
		return
	}
	logger.Debug("parsing yaml text to workflow", zap.String("text", string(bytes)))
	parseErr := yaml.Unmarshal(bytes, &workflow)
	if parseErr != nil {
		err = fmt.Errorf("error parsing yaml text: %w", parseErr)
		logger.Error(err.Error(), zap.Error(err))
		return
	}
	workflow.Tasks = w.Tasks
	workflow.Logger = w.Logger
	logger.Debug("successfully parsed yaml text to workflow", zap.Any("workflow", workflow))
	return
}

func (w *Workflows) NewWorkflowFromText(text string) (workflow *Workflow, err error) {
	logger := w.Logger
	logger.Debug("parsing yaml text to workflow", zap.String("path", text))
	parseErr := yaml.Unmarshal([]byte(text), &workflow)
	if parseErr != nil {
		err = fmt.Errorf("error parsing yaml text: %w", parseErr)
		logger.Error(err.Error(), zap.Error(err))
		return
	}
	workflow.Tasks = w.Tasks
	workflow.Logger = w.Logger
	logger.Debug("successfully parsed yaml text to workflow", zap.Any("workflow", workflow))
	return
}

type Workflow struct {
	Version         string                   `yaml:"version,omitempty"`
	EnvironmentList []map[string]string      `yaml:"environment,omitempty"`
	TaskList        []map[string]interface{} `yaml:"task,omitempty"`
	Node            string                   `yaml:",omitempty"`
	Tasks           tasks.ITasks             `yaml:",omitempty"`
	Logger          *zap.Logger              `yaml:",omitempty"`
}

func (w *Workflow) Run() (err error) {
	logger := w.Logger
	if validationErr := w.Validate(); validationErr != nil {
		err = validationErr
		return
	}
	if len(w.EnvironmentList) > 0 {
		logger.Debug("setting up workflow environment", zap.Any("environment list", w.EnvironmentList))
		for i := range w.EnvironmentList {
			environmentMap := lib.ResolveEnvironmentVariablesInMap(w.EnvironmentList[i])
			if environmentErr := lib.SetEnvironmentVariables(environmentMap); environmentErr != nil {
				err = fmt.Errorf("error setting up workflow environment: %w", err)
				logger.Error(err.Error(), zap.Error(err))
				return
			}
		}
		logger.Debug("successfully set up workflow environment")
	}
	logger.Debug("setting up workflow tasks, steps and then run", zap.Any("task list", w.TaskList))
	for _, taskData := range w.TaskList {
		var node string
		var taskType string
		var stepList []interface{}
		for j := range taskData {
			switch j {
			case string(lib.NodeAttribute):
				node = taskData[j].(string)
			default:
				taskType = j
				stepList = taskData[j].([]interface{})
			}
		}
		for k := range stepList {
			step := stepList[k].(map[string]interface{})
			task := w.Tasks.NewTask(taskType, step, node)
			if task == nil {
				err = fmt.Errorf("error creating new task type as invalid task type or step type")
				logger.Error(err.Error(), zap.Error(err))
				return
			}
			if runErr := task.Run(); runErr != nil {
				err = fmt.Errorf("error running task: %w", runErr)
				logger.Error(err.Error(), zap.Error(err))
				return
			}
		}
	}
	logger.Debug("successfully set up workflow tasks, steps and then ran")
	return
}

func (w *Workflow) Validate() (err error) {
	logger := w.Logger
	logger.Debug("validating workflow", zap.Any("workflow", w))
	if version, verErr := strconv.ParseFloat(w.Version, 32); verErr != nil || math.Round(version*10)/10 != 0.1 {
		err = fmt.Errorf("validation error: invalid version or unsupported version (not 0.1)")
		logger.Error(err.Error(), zap.Error(err))
		return
	}
	if len(w.TaskList) == 0 {
		err = fmt.Errorf("validation error: no task list found")
		logger.Error(err.Error(), zap.Error(err))
		return
	}
	for i, taskData := range w.TaskList {
		var taskType string
		var stepList []interface{}
		for j := range taskData {
			switch j {
			case string(lib.NodeAttribute):
			case string(lib.SequentialTaskType):
				taskType = j
				if reflect.TypeOf(taskData[j]) != reflect.TypeOf([]interface{}{}) {
					err = fmt.Errorf("validation error: invalid step list in task number %s with task type %s",
						j, taskType)
					logger.Error(err.Error(), zap.Error(err))
					return
				}
				stepList = taskData[j].([]interface{})
			default:
				err = fmt.Errorf("validation error: invalid task type %s", j)
				logger.Error(err.Error(), zap.Error(err))
				return
			}
		}
		if len(stepList) == 0 {
			err = fmt.Errorf("validation error: no step list found for task number %d with task type %s", i+1, taskType)
			logger.Error(err.Error(), zap.Error(err))
			return
		}
		for k := range stepList {
			step := stepList[k].(map[string]interface{})
			for l := range step {
				switch l {
				case string(lib.EnvironmentStepType):
				case string(lib.PrintStepType):
				case string(lib.ExecuteStepType):
				default:
					err = fmt.Errorf("validation error: invalid step type %s in task number %d with task type %s",
						l, i+1, taskType)
					logger.Error(err.Error(), zap.Error(err))
					return
				}
			}
		}
	}
	logger.Debug("successfully validated workflow")
	return
}

func (w *Workflow) SplitNodes() (newWorkflowList []Workflow) {
	logger := w.Logger
	if validationErr := w.Validate(); validationErr != nil {
		return
	}
	logger.Debug("splitting up nodes", zap.Any("task list", w.TaskList))
	additionalEnvironmentList := []map[string]string{}
	for i, taskData := range w.TaskList {
		var node string
		var stepList []interface{}
		for j := range taskData {
			switch j {
			case string(lib.NodeAttribute):
				node = taskData[j].(string)
			default:
				stepList = taskData[j].([]interface{})
			}
		}

		if node == "" {
			newWorkflowList = nil
			logger.Debug("failed split up nodes as node value not found", zap.Any("task data", taskData))
			return
		} else {
			newWf := &Workflow{
				Version:         w.Version,
				EnvironmentList: w.EnvironmentList,
				TaskList:        []map[string]interface{}{w.TaskList[i]},
				Node:            node,
			}
			taskAdditionalEnvironmentList := []map[string]string{}
			for k := range stepList {
				step := stepList[k].(map[string]interface{})
				for l := range step {
					if l == string(lib.EnvironmentStepType) {
						stepEnvironmentData := step[l].([]interface{})
						for m := range stepEnvironmentData {
							stepEnvironmentDataList := stepEnvironmentData[m].(map[string]interface{})
							for n := range stepEnvironmentDataList {
								taskAdditionalEnvironmentList = append(taskAdditionalEnvironmentList,
									map[string]string{n: stepEnvironmentDataList[n].(string)})
							}
						}
					}
				}
			}
			for o := range additionalEnvironmentList {
				newWf.EnvironmentList = append(newWf.EnvironmentList, additionalEnvironmentList[o])
			}
			for p := range taskAdditionalEnvironmentList {
				additionalEnvironmentList = append(additionalEnvironmentList, taskAdditionalEnvironmentList[p])
			}
			newWorkflowList = append(newWorkflowList, *newWf)
		}
	}
	logger.Debug("successfully split up nodes", zap.Any("workflow list", newWorkflowList))
	return
}
