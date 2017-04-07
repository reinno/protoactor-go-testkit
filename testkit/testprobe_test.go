package testkit

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"reflect"
	"testing"
	"time"
)

type World struct{ say string }

const (
	hello = "hello"
	world = "world"
	hey   = "hey"
)

func helloActorPropsWithSleep(d time.Duration) actor.Props {
	return actor.FromFunc(func(context actor.Context) {
		switch context.Message() {
		case hello:
			time.Sleep(d)
			context.Respond(world)

		case hey:
			context.Respond(&World{"haha"})
		}
	})
}

func helloActorProps() actor.Props {
	return helloActorPropsWithSleep(0)
}

func TestTestProbeReceiveMsg(t *testing.T) {
	tp := NewTestProbe(t)
	tp.Request(actor.Spawn(helloActorProps()), hello)
	tp.ExpectMsg(world)
	tp.ExpectNoMsg()
}

func TestTestProbeReceiveMsgInTime(t *testing.T) {
	tp := NewTestProbe(t)
	tp.Request(
		actor.Spawn(helloActorPropsWithSleep(5*time.Second)),
		hello)
	tp.ExpectMsgInTime(6*time.Second, world)
}

func TestTestProbeReceiveMsgInTime0(t *testing.T) {
	tp := NewTestProbe(t)
	tp.Request(actor.Spawn(helloActorProps()), hello)
	time.Sleep(1 * time.Millisecond)
	tp.ExpectMsgInTime(0*time.Second, world)
}

func TestTestProbeReceiveMsgType(t *testing.T) {
	tp := NewTestProbe(t)
	tp.Request(actor.Spawn(helloActorProps()), hey)
	tp.ExpectMsgType(reflect.TypeOf(&World{"wowo"}))
}

func TestTestProbeAutoPilot(t *testing.T) {
	tp1 := NewTestProbe(t)
	tp2 := NewTestProbe(t)

	ap := &struct{ AutoPilotHelper }{}
	ap.runMethod = func(sender *actor.PID, msg interface{}) AutoPilot {
		switch msg {
		case hey:
			sender.Tell(&World{"haha"})
		}
		return KeepRunning
	}

	tp1.SetAutoPilot(ap)
	tp2.Request(tp1.Pid(), hey)
	tp2.ExpectMsgType(reflect.TypeOf(&World{"wowo"}))
}
