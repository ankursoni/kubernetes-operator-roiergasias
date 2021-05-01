package main

import (
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/workflow"
)

func main() {
	workflow := workflow.NewWorkflow("./cmd/hello-world/hello-world.yaml")
	workflow.Run()
}
