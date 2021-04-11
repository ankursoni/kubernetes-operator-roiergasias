package steps

type Step struct {
}

type StepWorkflow interface {
	Run()
}

func NewStep(stepType string, stepArguments []interface{}) (step StepWorkflow) {
	switch stepType {
	case "print":
		step = NewPrintStep(stepArguments)
		return
	default:
		return
	}
}
