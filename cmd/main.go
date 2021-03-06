package main

import (
	"fmt"
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/lib"
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/workflow"
	"github.com/jessevdk/go-flags"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Options struct {
	Debug bool `short:"d" long:"debug" description:"enable debug level for logs"`
}
type RunCommand struct {
	File string `required:"yes" short:"f" long:"file" description:"workflow yaml file"`
}
type ValidateCommand struct {
	File string `required:"yes" short:"f" long:"file" description:"workflow yaml file"`
}
type SplitCommand struct {
	File string `required:"yes" short:"f" long:"file" description:"workflow yaml file"`
}
type VersionCommand struct{}

var options Options
var parser = flags.NewParser(&options, flags.Default)
var runCommand RunCommand
var validateCommand ValidateCommand
var splitCommand SplitCommand
var versionCommand VersionCommand

var logger *zap.Logger

func init() {
	var err error

	_, err = parser.AddCommand("run", "run a workflow", "run a workflow yaml file", &runCommand)
	if err != nil {
		log.Fatalln(fmt.Errorf("error adding command line commands: %w", err))
		return
	}
	_, err = parser.AddCommand("validate", "validate a workflow yaml file",
		"validate a workflow yaml file content", &validateCommand)
	if err != nil {
		log.Fatalln(fmt.Errorf("error adding command line commands: %w", err))
		return
	}
	_, err = parser.AddCommand("split", "split a workflow yaml file",
		"split a workflow yaml file content and output new workflow yaml(s) if nodes are assigned to each task",
		&splitCommand)
	if err != nil {
		log.Fatalln(fmt.Errorf("error adding command line commands: %w", err))
		return
	}
	_, err = parser.AddCommand("version", "display version", "display version information",
		&versionCommand)
	if err != nil {
		log.Fatalln(fmt.Errorf("error adding command line commands: %w", err))
		return
	}

	logger, err = lib.NewZapLogger(options.Debug)
	if err != nil {
		log.Fatalln(fmt.Errorf("error creating new zap logger: %w", err))
		return
	}
}

func main() {
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
}

func (rc *RunCommand) Execute(_ []string) error {
	logger.Info("running the workflow yaml file", zap.String("path", rc.File))
	w, err := workflow.NewWorkflows(nil, logger).NewWorkflow(rc.File)
	if err != nil {
		logger.Fatal("error creating new workflow", zap.Error(err))
		return err
	}
	if err = w.Run(); err != nil {
		logger.Fatal("error running workflow: %w", zap.Error(err))
		return err
	}
	logger.Info("successfully ran the workflow yaml file", zap.String("path", rc.File))
	return nil
}

func (vc *ValidateCommand) Execute(_ []string) error {
	logger.Info("validating the workflow yaml file", zap.String("path", vc.File))
	w, err := workflow.NewWorkflows(nil, logger).NewWorkflow(vc.File)
	if err != nil {
		logger.Fatal("error creating new workflow", zap.Error(err))
		return err
	}
	if err = w.Validate(); err != nil {
		logger.Fatal("error validating workflow: %w", zap.Error(err))
		return err
	}
	logger.Info("successfully validated the workflow yaml file", zap.String("path", vc.File))
	return nil
}

func (sc *SplitCommand) Execute(_ []string) error {
	logger.Info("splitting the workflow yaml file content", zap.String("path", sc.File))
	w, err := workflow.NewWorkflows(nil, logger).NewWorkflow(sc.File)
	if err != nil {
		logger.Fatal("error creating new workflow", zap.Error(err))
		return err
	}
	splitWorkflow := w.SplitNodes()
	if len(splitWorkflow) == 0 {
		logger.Info("failed split up the workflow yaml file content", zap.String("path", sc.File))
	} else {
		for i, w := range splitWorkflow {
			fmt.Printf(">>> printing workflow: %d >>>\n", i+1)
			bytes, _ := yaml.Marshal(w)
			fmt.Println(string(bytes))
		}
		logger.Info("successfully split up the workflow yaml file content", zap.String("path", sc.File))
	}
	return nil
}

func (_ *VersionCommand) Execute(_ []string) error {
	fmt.Println("roiergasias workflow engine: v0.1.2")
	return nil
}
