package workflow

import (
	"io/ioutil"
	"os"

	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/tasks"
	"gopkg.in/yaml.v3"
)

type Workflow struct {
	Version         string                   `yaml:"version"`
	EnvironmentList []map[string]string      `yaml:"environment"`
	TaskList        []map[string]interface{} `yaml:"task"`
	Node            string
}

func NewWorkflow(filePath string) (workflow *Workflow) {
	bytes, _ := ioutil.ReadFile(filePath)
	yaml.Unmarshal(bytes, &workflow)
	return
}

func NewWorkflowFromText(text string) (workflow *Workflow) {
	yaml.Unmarshal([]byte(text), &workflow)
	return
}

func (workflow *Workflow) Run() {
	// setup workflow environment
	for j := range workflow.EnvironmentList {
		environmentVariableList := workflow.EnvironmentList[j]
		for k := range environmentVariableList {
			os.Setenv(k, environmentVariableList[k])
		}
	}

	// setup workflow tasks and steps
	for i := range workflow.TaskList {
		taskData := workflow.TaskList[i]
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
			task := tasks.NewTask(taskType, step, node)
			task.Run()
		}
	}
}

func (workflow *Workflow) SplitNodes() (newWorkflowList []Workflow) {
	// setup workflow tasks and steps
	for i := range workflow.TaskList {
		taskData := workflow.TaskList[i]
		node := ""
		for j := range taskData {
			switch j {
			case "node":
				node = taskData[j].(string)
			}
		}

		if node != "" {
			newWf := &Workflow{
				Version:         workflow.Version,
				EnvironmentList: workflow.EnvironmentList,
				TaskList:        []map[string]interface{}{workflow.TaskList[i]},
				Node:            node,
			}
			newWorkflowList = append(newWorkflowList, *newWf)
		} else {
			newWorkflowList = nil
			return
		}
	}
	return
}
