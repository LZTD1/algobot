package stater

import (
	"algobot/internal/lib/fsm"
	router "github.com/LZTD1/telebot-router"
	"gopkg.in/telebot.v4"
)

type Stater interface {
	State(uid int64) fsm.State
}

func New(stater Stater, onState fsm.State) func(next router.RouteHandler) router.RouteHandler {
	return func(next router.RouteHandler) router.RouteHandler {
		return router.HandlerFunc(func(ctx telebot.Context) error {
			if stater.State(ctx.Sender().ID) == onState {
				return next.ServeContext(ctx)
			}

			return nil
		})
	}
}
