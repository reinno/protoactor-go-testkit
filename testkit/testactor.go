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

	case SetAutoPilot:
		ta.autopilot = msg.ap

	default:
		switch ap := ta.autopilot.Run(context.Sender(), msg); ap {
		case KeepRunning:
		default:
			ta.autopilot = ap
		}

		ta.msgQueue <- RealMessage{msg, context.Sender()}
	}
}

func NewTestActorProps(msgQueue chan RealMessage) actor.Props {
	return actor.FromInstance(&TestActor{msgQueue, NoAutoPilot})
}
