package steps

import (
	"fmt"
	"go.uber.org/zap"
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
	logger := step.Logger
	logger.Debug("creating new execute step")
	commandList := []string{}
	for k := range step.StepArgumentList {
		commandList = append(commandList, step.StepArgumentList[k].(string))
	}
	commandList = lib.ResolveEnvironmentVariables(commandList)
	executeStep = &ExecuteStep{
		CommandList: commandList,
		Step:        *step,
	}
	logger.Debug("successfully created new execute step")
	return
}

func (executeStep *ExecuteStep) Run() (err error) {
	logger := executeStep.Logger
	logger.Debug("started running execute step", zap.Any("command list", executeStep.CommandList))
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
	logger.Debug("successfully ran execute step")
	return
}
