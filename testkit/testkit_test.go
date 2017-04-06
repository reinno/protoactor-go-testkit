package testkit

import (
	"testing"
	"github.com/AsynkronIT/protoactor-go/actor"
)

func TestTestProbeReceiveMsg(t *testing.T) {
	hello := "hello"
	world := "world"

	props := actor.FromFunc(func(context actor.Context) {
		switch context.Message() {
		case hello:
			context.Respond(world)
		}
	})
	rec := actor.Spawn(props)

	tp := NewTestProbe(t)
	tp.Request(rec, hello)
	tp.ExpectMsg(world)
	tp.ExpectNoMsg()
}
