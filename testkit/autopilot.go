package testkit

import "github.com/AsynkronIT/protoactor-go/actor"

var (
	NoAutoPilot = &noAutoPilot{}
	KeepRunning = &keepRunning{}
)

type (
	AutoPilot interface {
		Run(sender *actor.PID, msg interface{}) AutoPilot
	}

	AutoPilotHelper struct {
		runMethod func(sender *actor.PID, msg interface{}) AutoPilot
	}

	noAutoPilot struct{}
	keepRunning struct{}
)

func (ap *noAutoPilot) Run(sender *actor.PID, msg interface{}) AutoPilot {
	return ap
}

func (ap *keepRunning) Run(sender *actor.PID, msg interface{}) AutoPilot {
	panic("must not call")
}

func (aph *AutoPilotHelper) Run(sender *actor.PID, msg interface{}) AutoPilot {
	return aph.runMethod(sender, msg)
}
