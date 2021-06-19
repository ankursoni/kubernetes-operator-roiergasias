package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/lib"
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/workflow"
	"go.uber.org/zap"
)

func main() {
	debug := flag.Bool("debug", false, "to indicate debug mode")
	flag.Parse()
	logger, err := lib.NewZapLogger(*debug)
	if err != nil {
		log.Fatalln(fmt.Errorf("error creating new zap logger: %w", err))
		return
	}
	//defer func(logger *zap.Logger) {
	//	err := logger.Sync()
	//	if err != nil {
	//		log.Println(fmt.Errorf("error syncing zap logger: %w", err))
	//		return
	//	}
	//}(logger)

	fmt.Printf("Running the workflow with input yaml: %s\n", os.Args[len(os.Args)-1])
	w, err := workflow.NewWorkflows(nil, logger).NewWorkflow(os.Args[len(os.Args)-1])
	if err != nil {
		logger.Fatal("error creating new workflow", zap.Error(err))
		return
	}
	if err = w.Run(); err != nil {
		logger.Fatal("error running workflow: %w", zap.Error(err))
		return
	}
}
