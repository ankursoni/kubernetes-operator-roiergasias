package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/workflow"
)

func main() {
	fmt.Printf("Running the workflow with input yaml: %s\n", os.Args[1])
	w, err := workflow.NewWorkflows(nil).NewWorkflow(os.Args[1])
	if err != nil {
		log.Fatalln(fmt.Errorf("error creating new workflow: %w", err))
		return
	}
	if err = w.Run(); err != nil {
		log.Fatalln(fmt.Errorf("error running workflow: %w", err))
		return
	}

	//splitList := w.SplitNodes()
	//for _, split := range splitList {
	//	bytes, _ := yaml.Marshal(split)
	//	fmt.Println(string(bytes))
	//}
	//if splitList == nil || len(splitList) == 0 {
	//	bytes, _ := yaml.Marshal(w)
	//	fmt.Println(string(bytes))
	//}
}
