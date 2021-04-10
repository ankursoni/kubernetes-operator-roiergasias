package workflow

import (
	"io/ioutil"

	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/tasks"
	"gopkg.in/yaml.v3"
)

type Workflow struct {
	Version string                   `yaml:"version"`
	Tasks   []map[string]interface{} `yaml:"tasks"`
}

func NewWorkflow(filePath string) (workflow Workflow) {
	bytes, _ := ioutil.ReadFile(filePath)
	yaml.Unmarshal(bytes, &workflow)
	return
}

func (workflow *Workflow) Run() {
	for i := range workflow.Tasks {
		tasksList := workflow.Tasks[i]
		for j := range tasksList {
			stepsList := tasksList[j].([]interface{})
			for k := range stepsList {
				step := stepsList[k].(map[string]interface{})
				task := tasks.NewTask(j, step)
				task.Run()
			}
		}
	}
}
