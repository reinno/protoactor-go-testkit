package testkit

import (
	"reflect"
	"testing"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/eventstream"
	"github.com/stretchr/testify/assert"
)

type World struct{ say string }

const (
	hello = "hello"
	world = "world"
	hey   = "hey"
)

func helloActorPropsWithSleep(d time.Duration) *actor.Props {
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

func helloActorProps() *actor.Props {
	return helloActorPropsWithSleep(0)
}

func TestTestProbeReceiveMsg(t *testing.T) {
	tp := NewTestProbe(t)
	tp.Request(actor.Spawn(helloActorProps()), hello)
	msg := tp.ExpectMsg(world)
	assert.Equal(t, msg, world)
	tp.ExpectNoMsg()
	tp.StopGraceful()
}

func TestTestProbeReceiveMsgInTime(t *testing.T) {
	tp := NewTestProbe(t)
	tp.Request(
		actor.Spawn(helloActorPropsWithSleep(5*time.Second)),
		hello)
	tp.ExpectMsgInTime(6*time.Second, world)
	tp.StopGraceful()
}

func TestTestProbeReceiveMsgInTime0(t *testing.T) {
	tp := NewTestProbe(t)
	tp.Request(actor.Spawn(helloActorProps()), hello)
	time.Sleep(1 * time.Millisecond)
	tp.ExpectMsgInTime(0*time.Second, world)
	tp.StopGraceful()
}

func TestTestProbeReceiveAny(t *testing.T) {
	tp := NewTestProbe(t)
	tp.Request(actor.Spawn(helloActorProps()), hello)
	msg := tp.ExpectAnyMsg()
	assert.Equal(t, msg, world)
	tp.ExpectNoMsg()
	tp.StopGraceful()
}

func TestTestProbeReceiveMsgType(t *testing.T) {
	tp := NewTestProbe(t)
	tp.Request(actor.Spawn(helloActorProps()), hey)
	tp.ExpectMsgType(reflect.TypeOf(&World{"wowo"}))
	tp.StopGraceful()
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
	tp1.ExpectMsg(hey)
	tp2.ExpectMsgType(reflect.TypeOf(&World{"wowo"}))
	tp1.StopGraceful()
	tp2.StopGraceful()
}

func TestTestProbeSender(t *testing.T) {
	tp := NewTestProbe(t)

	sender := actor.Spawn(actor.FromFunc(func(context actor.Context) {
		switch context.Message().(type) {
		case *actor.Started:
			context.Request(tp.Pid(), hey)
		}
	}))

	tp.ExpectMsg(hey)
	assert.Equal(t, tp.Sender(), sender)

	tp.StopGraceful()
}

func TestTestProbeUnExpectMsg(t *testing.T) {
	tp := NewTestProbe(t)

	sender := actor.Spawn(actor.FromFunc(func(context actor.Context) {
		switch context.Message().(type) {
		case *actor.Started:
			context.Request(tp.Pid(), hey)
			context.Request(tp.Pid(), hey)
		}
	}))

	deadLetterReceived := 0
	sub := eventstream.Subscribe(func(msg interface{}) {
		if deadLetter, ok := msg.(*actor.DeadLetterEvent); ok {
			assert.Equal(t, deadLetter.Sender, sender)
			assert.Equal(t, deadLetter.Message, hey)
			deadLetterReceived += 1
		}
	})
	defer eventstream.Unsubscribe(sub)

	//tp.ExpectMsg(hey)
	time.Sleep(time.Millisecond)
	tp.StopGraceful()
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, 2, deadLetterReceived)
}

func TestTestProbeIgnoreMsg(t *testing.T) {
	tp := NewTestProbe(t)

	sender := actor.Spawn(actor.FromFunc(func(context actor.Context) {
		switch context.Message().(type) {
		case *actor.Started:
			context.Request(tp.Pid(), hey)
		}
	}))

	tp.SetIgnore(func(m interface{}) bool {
		switch m {
		case hey:
			return true
		default:
			return false
		}
	})
	deadLetterReceived := 0
	sub := eventstream.Subscribe(func(msg interface{}) {
		if deadLetter, ok := msg.(*actor.DeadLetterEvent); ok {
			assert.Equal(t, deadLetter.Sender, sender)
			assert.Equal(t, deadLetter.Message, hey)
			deadLetterReceived += 1
		}
	})
	defer eventstream.Unsubscribe(sub)

	//tp.ExpectMsg(hey)
	time.Sleep(time.Millisecond)
	tp.StopGraceful()
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, 0, deadLetterReceived)
}

func TestTestProbePropsIgnoreMsg(t *testing.T) {
	tp := NewTestProbe(t)

	props := actor.FromFunc(func(context actor.Context) {
		switch context.Message() {
		case hello:
			context.Request(tp.Pid(), hey)
		}
	})
	sender := NewTestProbeWithProps(t, props)
	sender.SetIgnore(func(m interface{}) bool {
		switch m {
		case hello:
			return true
		default:
			return false
		}
	})

	sender.Pid().Tell(hello)

	deadLetterReceived := 0
	sub := eventstream.Subscribe(func(msg interface{}) {
		if deadLetter, ok := msg.(*actor.DeadLetterEvent); ok {
			assert.Equal(t, deadLetter.Message, hey)
			assert.Equal(t, deadLetter.Sender, sender.Pid())
			deadLetterReceived += 1
		}
	})
	defer eventstream.Unsubscribe(sub)

	//tp.ExpectMsg(hey)
	time.Sleep(time.Millisecond)
	tp.StopGraceful()
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, 1, deadLetterReceived)
}

func TestTestProbeStop(t *testing.T) {
	tp1, err := NewTestProbeNamed(t, "a")
	assert.Nil(t, err)
	assert.NotNil(t, tp1)

	tp2, err := NewTestProbeNamed(t, "a")
	assert.NotNil(t, err)
	assert.Nil(t, tp2)

	tp1.StopFuture().Wait()
	tp1.StopGraceful()
	tp3, err := NewTestProbeNamed(t, "a")
	assert.Nil(t, err)
	assert.NotNil(t, tp3)
	tp3.StopGraceful()
}
