package testkit

import (
	"github.com/AsynkronIT/protoactor-go/actor"
)

type (
	testActor struct{}

	setAutoPilot struct {
		ap AutoPilot
	}

	IgnoreFunc func(interface{}) bool
	setIgnore  struct {
		fn IgnoreFunc
	}

	testActorPlugin struct {
		msgQueue  chan realMessage
		autopilot AutoPilot
		ignore    IgnoreFunc
	}
)

func IgnoreNone(interface{}) bool { return false }

func (ta *testActor) Receive(context actor.Context) {}

func (ta *testActorPlugin) Receive(context actor.Context, next actor.ActorFunc) {
	next(context)

	switch msg := context.Message().(type) {
	case *actor.Started:

	case *actor.Stopping:
	L:
		for {
			select {
			case m := <-ta.msgQueue:
				deadLetter, _ := actor.ProcessRegistry.Get(nil)
				deadLetter.SendUserMessage(context.Sender(), &actor.MessageEnvelope{nil, m.msg, m.sender})
			default:
				break L
			}
		}

	case setAutoPilot:
		ta.autopilot = msg.ap

	case setIgnore:
		ta.ignore = msg.fn

	default:
		if !ta.ignore(msg) {
			switch ap := ta.autopilot.Run(context.Sender(), msg); ap {
			case KeepRunning:
			default:
				ta.autopilot = ap
			}

			ta.msgQueue <- realMessage{msg, context.Sender()}
		}
	}
}

func newTestActorProps(msgQueue chan realMessage) *actor.Props {
	return actor.FromInstance(&testActor{}).
		WithMiddleware(use(&testActorPlugin{msgQueue, NoAutoPilot, IgnoreNone}))
}

func newTestActorPropsWithProps(msgQueue chan realMessage, props *actor.Props) *actor.Props {
	return props.WithMiddleware(use(
		&testActorPlugin{msgQueue,
			NoAutoPilot,
			IgnoreNone}))
}
