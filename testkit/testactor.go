package testkit

import (
	"github.com/AsynkronIT/protoactor-go/actor"
)

type (
	TestActor struct {
		msgQueue chan RealMessage
	}
)

func (ta *TestActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
	default:
		ta.msgQueue <- RealMessage{msg, context.Sender()}
	}
}