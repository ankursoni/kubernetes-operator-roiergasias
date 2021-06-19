package steps

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/lib"
)

var _ IStepWorkflow = &ExecuteStep{}

type ExecuteStep struct {
	CommandList []string
	Step
}

func (step *Step) NewExecuteStep() (executeStep *ExecuteStep) {
	commandList := []string{}
	for k := range step.StepArgumentList {
		commandList = append(commandList, step.StepArgumentList[k].(string))
	}
	commandList = lib.ResolveEnvironmentVariables(commandList)
	executeStep = &ExecuteStep{
		CommandList: commandList,
		Step:        *step,
	}
	return
}

func (executeStep *ExecuteStep) Run() (err error) {
	for command := range executeStep.CommandList {
		cmdWithArgs := strings.SplitAfter(executeStep.CommandList[command], " ")
		for k := range cmdWithArgs {
			cmdWithArgs[k] = strings.TrimSpace(cmdWithArgs[k])
		}
		cmd := exec.Command(cmdWithArgs[0])
		cmd.Args = cmdWithArgs
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout

		if cmdErr := cmd.Run(); cmdErr != nil {
			err = fmt.Errorf("error executing step: %w", cmdErr)
			return
		}
	}
	err = executeStep.Step.Run()
	return
}
