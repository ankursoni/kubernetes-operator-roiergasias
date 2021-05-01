package main

import (
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/workflow"
)

func main() {
	workflow := workflow.NewWorkflow("./cmd/machine-learning/machine-learning.yaml")
	workflow.Run()
}
