package stater

import (
	"algobot/internal/lib/fsm"
	"fmt"
	router "github.com/LZTD1/telebot-router"
	"gopkg.in/telebot.v4"
)

type Stater interface {
	State(uid int64) fsm.State
}

func New(stater Stater, onState fsm.State) func(next router.RouteHandler) router.RouteHandler {
	return func(next router.RouteHandler) router.RouteHandler {
		return router.HandlerFunc(func(ctx telebot.Context) error {
			fmt.Printf("username: %s ", ctx.Sender().Username)
			fmt.Printf("message: %s ", ctx.Message().Text)
			fmt.Printf("onState: %d ", onState)
			fmt.Printf("stater.State: %d \n", stater.State(ctx.Sender().ID))
			if stater.State(ctx.Sender().ID) == onState {
				return next.ServeContext(ctx)
			}

			return nil
		})
	}
}
