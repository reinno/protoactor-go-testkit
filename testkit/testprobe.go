package testkit

import (
	"reflect"
	"testing"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type (
	TestProbe interface {
		// Same as `expectMsg(testkit.DefaultTimeout, obj)`.
		ExpectMsg(obj interface{}) interface{}

		// Receive one message from the test actor and assert that it equals the
		// given object. Wait time is bounded by the given duration, with an
		// AssertionFailure being thrown in case of timeout.
		ExpectMsgInTime(max time.Duration, obj interface{}) interface{}

		// Same as `expectNoMsgInTime(testkit.DefaultTimeout)`.
		ExpectNoMsg()

		// Assert that no message is received for the specified time.
		ExpectNoMsgInTime(max time.Duration)

		// Same as `expectAnyInTime(testkit.DefaultTimeout)`.
		ExpectAnyMsg() interface{}

		// Receive one message from the test actor. Wait time is bounded by the given duration, with an
		// AssertionFailure being thrown in case of timeout.
		ExpectAnyMsgInTime(max time.Duration) interface{}

		// Same as `expectMsgType(testkit.DefaultTimeout, obj)`.
		ExpectMsgType(t reflect.Type) interface{}

		// Receive one message from the test actor and assert that it conforms to the
		// given type (after erasure). Wait time is bounded by the given duration,
		// with an AssertionFailure being thrown in case of timeout.
		ExpectMsgTypeInTime(max time.Duration, t reflect.Type) interface{}

		Within(min time.Duration, max time.Duration, f func() interface{}) interface{}

		// Request sends a message to the given PID and also provides probe's test actor PID as sender.
		Request(actor *actor.PID, msg interface{})

		// Get sender of last received message.
		Sender() *actor.PID

		// Install an AutoPilot to drive the testActor:
		//   the AutoPilot will be run for each received message and can be used to send or forward messages, etc.
		// Each invocation must return the AutoPilot for the next round.
		SetAutoPilot(ap AutoPilot)

		// Ignore all messages in the test actor for which the given function returns true.
		SetIgnore(ignore IgnoreFunc)

		// PID of the test actor.
		Pid() *actor.PID

		// Stop Test Probe
		Stop()

		// Stop Test Probe with end callback
		StopFuture() *actor.Future

		// Stop Test Probe and wait it end
		// Same as `StopFuture().Wait()`.
		StopGraceful()
	}
)

func newTestProbe(t *testing.T, testActorFunc func(chan realMessage) (*actor.PID, error)) (TestProbe, error) {
	msgQueue := make(chan realMessage, MaxMsgQueueNum)

	pid, err := testActorFunc(msgQueue)
	if err != nil {
		return nil, err
	}

	return &testBase{
		t:              t,
		testActor:      pid,
		msgQueue:       msgQueue,
		defaultTimeout: DefaultTimeout}, nil
}

func NewTestProbe(t *testing.T) TestProbe {
	tp, _ := newTestProbe(t,
		func(msgQueue chan realMessage) (*actor.PID, error) {
			return actor.Spawn(newTestActorProps(msgQueue)), nil
		})
	return tp
}

func NewTestProbeWithProps(t *testing.T, props *actor.Props) TestProbe {
	tp, _ := newTestProbe(t,
		func(msgQueue chan realMessage) (*actor.PID, error) {
			return actor.Spawn(newTestActorPropsWithProps(msgQueue, props)), nil
		})
	return tp
}

func NewTestProbeNamed(t *testing.T, name string) (TestProbe, error) {
	return newTestProbe(t,
		func(msgQueue chan realMessage) (*actor.PID, error) {
			return actor.SpawnNamed(newTestActorProps(msgQueue), name)
		})
}

func NewTestProbeNamedWithProps(t *testing.T, name string, props *actor.Props) (TestProbe, error) {
	return newTestProbe(t,
		func(msgQueue chan realMessage) (*actor.PID, error) {
			return actor.SpawnNamed(newTestActorPropsWithProps(msgQueue, props), name)
		})
}
