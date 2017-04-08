package testkit

import (
	"github.com/AsynkronIT/protoactor-go/actor"
)

type (
	TestActor struct {
		msgQueue  chan RealMessage
		autopilot AutoPilot
	}

	SetAutoPilot struct {
		ap AutoPilot
	}
)

func (ta *TestActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:

	case *actor.Stopping:
		for m := range ta.msgQueue {
			deadLetter, _ := actor.ProcessRegistry.Get(nil)
			deadLetter.SendUserMessage(context.Sender(), m.msg, m.sender)
		}

	case SetAutoPilot:
		ta.autopilot = msg.ap

	default:
		//fmt.Printf("recieved msg: %v", msg)
		switch ap := ta.autopilot.Run(context.Sender(), msg); ap {
		case KeepRunning:
		default:
			ta.autopilot = ap
		}

		ta.msgQueue <- RealMessage{msg, context.Sender()}
	}
}

func NewTestActorProps(msgQueue chan RealMessage) *actor.Props {
	return actor.FromInstance(&TestActor{msgQueue, NoAutoPilot})
}
