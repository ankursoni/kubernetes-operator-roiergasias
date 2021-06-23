package steps

import (
	"fmt"
	"github.com/ankursoni/kubernetes-operator-roiergasias/pkg/lib"
	"go.uber.org/zap"
)

var _ IStepWorkflow = &EnvironmentStep{}

type EnvironmentStep struct {
	EnvironmentList []map[string]string
	Step
}

func (step *Step) NewEnvironmentStep() (environmentStep *EnvironmentStep) {
	logger := step.Logger
	logger.Debug("creating new environment step")
	environmentList := []map[string]string{}
	for i := range step.StepArgumentList {
		stepArgument := step.StepArgumentList[i].(map[string]interface{})
		environmentMap := map[string]string{}
		for j := range stepArgument {
			environmentMap[j] = stepArgument[j].(string)
			environmentMap = lib.ResolveEnvironmentVariablesInMap(environmentMap)
		}
		environmentList = append(environmentList, environmentMap)
	}
	environmentStep = &EnvironmentStep{
		EnvironmentList: environmentList,
		Step:            *step,
	}
	logger.Debug("successfully created new environment step")
	return
}

func (environmentStep *EnvironmentStep) Run() (err error) {
	logger := environmentStep.Logger
	logger.Debug("started running environment step", zap.Any("environment list", environmentStep.EnvironmentList))
	for _, environmentMap := range environmentStep.EnvironmentList {
		if environmentErr := lib.SetEnvironmentVariables(environmentMap); environmentErr != nil {
			err = fmt.Errorf("error environment step: %w", environmentErr)
			logger.Error(err.Error(), zap.Error(err))
			return
		}
	}
	logger.Debug("successfully ran environment step")
	return
}
