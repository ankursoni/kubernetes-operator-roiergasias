package workflow

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/tasks"
	"gopkg.in/yaml.v3"
)

type IWorkflows interface {
	NewWorkflow(string) (*Workflow, error)
	NewWorkflowFromText(string) (*Workflow, error)
}

type Workflows struct {
	Tasks tasks.ITasks
}

var _ IWorkflows = &Workflows{}

func NewWorkflows(t tasks.ITasks) (workflows IWorkflows) {
	if t == nil {
		workflows = &Workflows{Tasks: tasks.NewTasks()}
	} else {
		workflows = &Workflows{Tasks: t}
	}
	return
}

func (w Workflows) NewWorkflow(filePath string) (workflow *Workflow, err error) {
	bytes, _ := ioutil.ReadFile(filePath)
	err = yaml.Unmarshal(bytes, &workflow)
	workflow.Tasks = w.Tasks
	return
}

func (w Workflows) NewWorkflowFromText(text string) (workflow *Workflow, err error) {
	err = yaml.Unmarshal([]byte(text), &workflow)
	workflow.Tasks = w.Tasks
	return
}

type Workflow struct {
	Version         string                   `yaml:"version,omitempty"`
	EnvironmentList []map[string]string      `yaml:"environment,omitempty"`
	TaskList        []map[string]interface{} `yaml:"task,omitempty"`
	Node            string                   `yaml:",omitempty"`
	Tasks           tasks.ITasks             `yaml:",omitempty"`
}

func (w *Workflow) Run() error {
	// setup workflow environment
	for j := range w.EnvironmentList {
		environmentVariableList := w.EnvironmentList[j]
		for k := range environmentVariableList {
			if err := os.Setenv(k, environmentVariableList[k]); err != nil {
				return fmt.Errorf("error setting up workflow environment: %w", err)
			}
		}
	}

	// setup workflow tasks and steps
	for i := range w.TaskList {
		taskData := w.TaskList[i]
		var node string
		var taskType string
		var stepList []interface{}
		for j := range taskData {
			switch j {
			case "node":
				node = taskData[j].(string)
			default:
				taskType = j
				stepList = taskData[j].([]interface{})
			}
		}
		for k := range stepList {
			step := stepList[k].(map[string]interface{})
			task := w.Tasks.NewTask(taskType, step, node)
			if err := task.Run(); err != nil {
				return fmt.Errorf("error running task: %w", err)
			}
		}
	}
	return nil
}

func (w *Workflow) SplitNodes() (newWorkflowList []Workflow) {
	additionalEnvironmentList := []map[string]string{}
	for i := range w.TaskList {
		taskData := w.TaskList[i]
		var node string
		var stepList []interface{}
		for j := range taskData {
			switch j {
			case "node":
				node = taskData[j].(string)
			default:
				stepList = taskData[j].([]interface{})
			}
		}

		if node != "" {
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
					if l == "set-environment" {
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
		} else {
			newWorkflowList = nil
			return
		}
	}
	return
}
