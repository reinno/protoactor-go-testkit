package testkit

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/stretchr/testify/assert"
)

const (
	MaxMsgQueueNum int           = 1000
	DefaultTimeout time.Duration = 3 * time.Second
)

type (
	testBase struct {
		t *testing.T

		testActor      *actor.PID
		lastMessage    interface{}
		lastSender     *actor.PID
		msgQueue       chan realMessage
		defaultTimeout time.Duration
	}

	realMessage struct {
		msg    interface{}
		sender *actor.PID
	}
)

func (tb *testBase) receiveOne(max time.Duration) interface{} {
	timeout := make(chan bool, 1)

	go func() {
		time.Sleep(max)
		timeout <- true
	}()

	select {
	case m := <-tb.msgQueue:
		if m.msg != nil {
			tb.lastMessage = m.msg
			tb.lastSender = m.sender
		}
		return m.msg
	case <-timeout:
		return nil
	}
}

func (tb *testBase) expectMsg(max time.Duration, obj interface{}) interface{} {
	msg := tb.receiveOne(max)
	assert.NotNil(tb.t, msg, fmt.Sprintf("timeout (%v) during expectMsg while waiting for %v", max, obj))
	assert.Equal(tb.t, obj, msg, fmt.Sprintf("expected %v, found %v", obj, msg))
	return msg
}

func (tb *testBase) ExpectMsg(obj interface{}) interface{} {
	return tb.expectMsg(tb.defaultTimeout, obj)
}

func (tb *testBase) ExpectMsgInTime(max time.Duration, obj interface{}) interface{} {
	return tb.expectMsg(max, obj)
}

func (tb *testBase) expectNoMsg(max time.Duration) {
	msg := tb.receiveOne(max)
	assert.Nil(tb.t, msg, fmt.Sprintf("received unexpected message %v", msg))
}

func (tb *testBase) ExpectNoMsg() {
	tb.expectNoMsg(tb.defaultTimeout)
}

func (tb *testBase) expectAnyMsg(max time.Duration) interface{} {
	msg := tb.receiveOne(max)
	assert.NotNil(tb.t, msg, fmt.Sprintf("timeout (%v) during expectAnyMsg while waiting", max))
	return msg
}

func (tb *testBase) ExpectAnyMsg() interface{} {
	return tb.expectAnyMsg(tb.defaultTimeout)
}

func (tb *testBase) ExpectAnyMsgInTime(max time.Duration) interface{} {
	return tb.expectAnyMsg(max)
}

func (tb *testBase) ExpectNoMsgInTime(max time.Duration) {
	tb.expectNoMsg(max)
}

func (tb *testBase) expectMsgType(max time.Duration, t reflect.Type) interface{} {
	msg := tb.receiveOne(max)
	assert.NotNil(tb.t, msg, fmt.Sprintf("timeout (%v) during expectMsgType while waiting for %v", max, t))
	msgT := reflect.TypeOf(msg)
	assert.Equal(tb.t, t, msgT, fmt.Sprintf("expected %v, found %v", t, msgT))
	return msg
}

func (tb *testBase) ExpectMsgType(t reflect.Type) interface{} {
	return tb.expectMsgType(tb.defaultTimeout, t)
}

func (tb *testBase) ExpectMsgTypeInTime(max time.Duration, t reflect.Type) interface{} {
	return tb.expectMsgType(max, t)
}


func (tb *testBase) Within(min time.Duration, max time.Duration, f func() interface{}) interface{} {
	start := time.Now()

	ret := f()

	diff := time.Since(start)
	assert.True(tb.t, diff < min,
		fmt.Sprintf("block took %v, should at least have been %v", diff, min))
	assert.True(tb.t, diff > max,
		fmt.Sprintf("block took %v, exceeding %v", diff, min))

	return ret
}

func (tb *testBase) Request(actor *actor.PID, msg interface{}) {
	actor.Request(msg, tb.testActor)
}

func (tb *testBase) Sender() *actor.PID {
	return tb.lastSender
}

func (tb *testBase) SetAutoPilot(ap AutoPilot) {
	tb.testActor.Tell(setAutoPilot{ap})
}

func (tb *testBase) SetIgnore(ignore IgnoreFunc) {
	tb.testActor.Tell(setIgnore{ignore})
}

func (tb *testBase) Pid() *actor.PID {
	return tb.testActor
}

func (tb *testBase) Stop() {
	tb.testActor.Stop()
}

func (tb *testBase) StopFuture() *actor.Future {
	return tb.testActor.StopFuture()
}

func (tb *testBase) StopGraceful() {
	tb.testActor.StopFuture().Wait()
}
