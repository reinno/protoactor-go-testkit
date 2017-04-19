package testkit

import (
	"github.com/AsynkronIT/protoactor-go/actor"
)

type pluginEx interface {
	Receive(actor.Context, actor.ActorFunc)
}

func use(plugin pluginEx) func(next actor.ActorFunc) actor.ActorFunc {
	return func(next actor.ActorFunc) actor.ActorFunc {
		return func(context actor.Context) {
			plugin.Receive(context, next)
		}
	}
}
