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

func newTestProbe(t *testing.T, testActorFunc func(chan RealMessage) (*actor.PID, error)) (TestProbe, error) {
	msgQueue := make(chan RealMessage, MaxMsgQueueNum)

	pid, err := testActorFunc(msgQueue)
	if err != nil {
		return nil, err
	}

	return &TestBase{
			t:              t,
			testActor:      pid,
			msgQueue:       msgQueue,
			defaultTimeout: DefaultTimeout}, nil
}

func NewTestProbe(t *testing.T) TestProbe {
	tp, _ := newTestProbe(t,
		func(msgQueue chan RealMessage) (*actor.PID, error) {
			return actor.Spawn(NewTestActorProps(msgQueue)), nil
		})
	return tp
}

func NewTestProbeNamed(t *testing.T, name string) (TestProbe, error) {
	return newTestProbe(t,
		func(msgQueue chan RealMessage) (*actor.PID, error) {
			return actor.SpawnNamed(NewTestActorProps(msgQueue), name)
		})
}
