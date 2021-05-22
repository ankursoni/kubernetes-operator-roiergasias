package steps

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ankursoni/kubernetes-operator-roiergasias/lib"
)

var _ StepWorkflow = &ExecuteStep{}

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

func (executeStep *ExecuteStep) Run() {
	for command := range executeStep.CommandList {
		cmdWithArgs := strings.SplitAfter(executeStep.CommandList[command], " ")
		for k := range cmdWithArgs {
			cmdWithArgs[k] = strings.TrimSpace(cmdWithArgs[k])
		}
		cmd := exec.Command(cmdWithArgs[0])
		cmd.Args = cmdWithArgs
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout

		if err := cmd.Run(); err != nil {
			fmt.Println("Error: ", err)
		}
	}
	executeStep.Step.Run()
}
