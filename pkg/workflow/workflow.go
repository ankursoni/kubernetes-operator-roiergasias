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
}

func NewWorkflow(filePath string) (workflow Workflow) {
	bytes, _ := ioutil.ReadFile(filePath)
	yaml.Unmarshal(bytes, &workflow)
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
		taskList := workflow.TaskList[i]
		for j := range taskList {
			stepList := taskList[j].([]interface{})
			for k := range stepList {
				step := stepList[k].(map[string]interface{})
				task := tasks.NewTask(j, step)
				task.Run()
			}
		}
	}
}
