package testkit

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"reflect"
	"testing"
	"time"
)

type (
	TestProbe interface {
		ExpectMsg(obj interface{}) interface{}
		ExpectMsgInTime(max time.Duration, obj interface{}) interface{}
		ExpectNoMsg()
		ExpectNoMsgInTime(max time.Duration)
		ExpectMsgType(t reflect.Type) interface{}
		ExpectMsgTypeInTime(max time.Duration, t reflect.Type) interface{}

		Request(actor *actor.PID, msg interface{})
		SetAutoPilot(ap AutoPilot)
		Pid() *actor.PID
	}
)

func newTestProbe(t *testing.T, testActorFunc func(chan RealMessage) *actor.PID) TestProbe {
	msgQueue := make(chan RealMessage, MaxMsgQueueNum)

	return &TestBase{
		t:              t,
		testActor:      testActorFunc(msgQueue),
		msgQueue:       msgQueue,
		defaultTimeout: DefaultTimeout,
	}
}

func NewTestProbe(t *testing.T) TestProbe {
	return newTestProbe(t,
		func(msgQueue chan RealMessage) *actor.PID {
			return actor.Spawn(NewTestActorProps(msgQueue))
		})
}

func NewTestProbeNamed(t *testing.T, name string) TestProbe {
	return newTestProbe(t,
		func(msgQueue chan RealMessage) *actor.PID {
			return actor.SpawnNamed(NewTestActorProps(msgQueue), name)
		})
}
