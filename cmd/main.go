package main

import (
	"fmt"
	"os"

	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/workflow"
)

func main() {
	fmt.Printf("Running the workflow with input yaml: %s\n", os.Args[1])
	workflow := workflow.NewWorkflow(os.Args[1])
	workflow.Run()
}
