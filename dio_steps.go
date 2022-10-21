package dionysus

import (
	"github.com/gowins/dionysus/step"
)

func (d *Dio) RegUserFirstPreRunStep(instanceStep step.InstanceStep) error {
	return d.persistentPreRunE.RegActionSteps(instanceStep.StepName, 1+step.SystemPrioritySteps, instanceStep.Func)
}

func (d *Dio) RegUserSecondPreRunStep(instanceStep step.InstanceStep) error {
	return d.persistentPreRunE.RegActionSteps(instanceStep.StepName, 2+step.SystemPrioritySteps, instanceStep.Func)
}

func (d *Dio) RegUserThirdPreRunStep(instanceStep step.InstanceStep) error {
	return d.persistentPreRunE.RegActionSteps(instanceStep.StepName, 3+step.SystemPrioritySteps, instanceStep.Func)
}

func (d *Dio) RegUserFourthPreRunStep(instanceStep step.InstanceStep) error {
	return d.persistentPreRunE.RegActionSteps(instanceStep.StepName, 4+step.SystemPrioritySteps, instanceStep.Func)
}

func (d *Dio) RegUserFifthPreRunStep(instanceStep step.InstanceStep) error {
	return d.persistentPreRunE.RegActionSteps(instanceStep.StepName, 5+step.SystemPrioritySteps, instanceStep.Func)
}

func (d *Dio) RegUserSixthPreRunStep(instanceStep step.InstanceStep) error {
	return d.persistentPreRunE.RegActionSteps(instanceStep.StepName, 6+step.SystemPrioritySteps, instanceStep.Func)
}

func (d *Dio) RegUserSeventhPreRunStep(instanceStep step.InstanceStep) error {
	return d.persistentPreRunE.RegActionSteps(instanceStep.StepName, 7+step.SystemPrioritySteps, instanceStep.Func)
}

func (d *Dio) RegUserEighthPreRunStep(instanceStep step.InstanceStep) error {
	return d.persistentPreRunE.RegActionSteps(instanceStep.StepName, 8+step.SystemPrioritySteps, instanceStep.Func)
}

func (d *Dio) RegUserNinethPreRunStep(instanceStep step.InstanceStep) error {
	return d.persistentPreRunE.RegActionSteps(instanceStep.StepName, 9+step.SystemPrioritySteps, instanceStep.Func)
}

func (d *Dio) RegUserTenthPreRunStep(instanceStep step.InstanceStep) error {
	return d.persistentPreRunE.RegActionSteps(instanceStep.StepName, 10+step.SystemPrioritySteps, instanceStep.Func)
}

func (d *Dio) RegUserFirstPostRunStep(instanceStep step.InstanceStep) error {
	return d.persistentPostRunE.RegActionSteps(instanceStep.StepName, 1+step.SystemPrioritySteps, instanceStep.Func)
}

func (d *Dio) RegUserSecondPostRunStep(instanceStep step.InstanceStep) error {
	return d.persistentPostRunE.RegActionSteps(instanceStep.StepName, 2+step.SystemPrioritySteps, instanceStep.Func)
}

func (d *Dio) RegUserThirdPostRunStep(instanceStep step.InstanceStep) error {
	return d.persistentPostRunE.RegActionSteps(instanceStep.StepName, 3+step.SystemPrioritySteps, instanceStep.Func)
}

func (d *Dio) RegUserFourthPostRunStep(instanceStep step.InstanceStep) error {
	return d.persistentPostRunE.RegActionSteps(instanceStep.StepName, 4+step.SystemPrioritySteps, instanceStep.Func)
}

func (d *Dio) RegUserFifthPostRunStep(instanceStep step.InstanceStep) error {
	return d.persistentPostRunE.RegActionSteps(instanceStep.StepName, 5+step.SystemPrioritySteps, instanceStep.Func)
}

func (d *Dio) RegUserSixthPostRunStep(instanceStep step.InstanceStep) error {
	return d.persistentPostRunE.RegActionSteps(instanceStep.StepName, 6+step.SystemPrioritySteps, instanceStep.Func)
}

func (d *Dio) RegUserSeventhPostRunStep(instanceStep step.InstanceStep) error {
	return d.persistentPostRunE.RegActionSteps(instanceStep.StepName, 7+step.SystemPrioritySteps, instanceStep.Func)
}

func (d *Dio) RegUserEighthPostRunStep(instanceStep step.InstanceStep) error {
	return d.persistentPostRunE.RegActionSteps(instanceStep.StepName, 8+step.SystemPrioritySteps, instanceStep.Func)
}

func (d *Dio) RegUserNinethPostRunStep(instanceStep step.InstanceStep) error {
	return d.persistentPostRunE.RegActionSteps(instanceStep.StepName, 9+step.SystemPrioritySteps, instanceStep.Func)
}

func (d *Dio) RegUserTenthPostRunStep(instanceStep step.InstanceStep) error {
	return d.persistentPostRunE.RegActionSteps(instanceStep.StepName, 10+step.SystemPrioritySteps, instanceStep.Func)
}

// PreRunStepsAppend append step will exec after step with priority which define by func PreRunRegWithPriority
func (d *Dio) PreRunStepsAppend(instanceSteps ...step.InstanceStep) error {
	for _, instanceStep := range instanceSteps {
		if err := d.persistentPreRunE.ActionStepsAppend(instanceStep.StepName, instanceStep.Func); err != nil {
			return err
		}
	}
	return nil
}

// PostRunStepsAppend append step will exec after step with priority which define by func PostRunRegWithPriority
func (d *Dio) PostRunStepsAppend(instanceSteps ...step.InstanceStep) error {
	for _, instanceStep := range instanceSteps {
		if err := d.persistentPostRunE.ActionStepsAppend(instanceStep.StepName, instanceStep.Func); err != nil {
			return err
		}
	}
	return nil
}
