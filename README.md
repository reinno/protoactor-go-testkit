# TestKit for ProtoActor-Go
Provide a AKKA TestKit like for [Protoactor](https://github.com/AsynkronIT/protoactor-go).

[![Build Status](https://travis-ci.org/reinno/protoactor-go-testkit.svg?branch=master)](https://travis-ci.org/reinno/protoactor-go-testkit)
[![Go Report Card](https://goreportcard.com/badge/github.com/reinno/protoactor-go-testkit)](https://goreportcard.com/report/github.com/reinno/protoactor-go-testkit)

# Example
``` go
func helloActorProps() *actor.Props {
	return actor.FromFunc(func(context actor.Context) {
		switch context.Message() {
		case hello:
			context.Respond("world")
		}
	})
}

func TestTestProbeReceiveMsg(t *testing.T) {
    helloActor := actor.Spawn(helloActorProps())
	tp := NewTestProbe(t)

	tp.Request(helloActor, "hello")
	tp.ExpectMsg("world")
	tp.ExpectNoMsg()

	tp.StopGraceful()
}
```