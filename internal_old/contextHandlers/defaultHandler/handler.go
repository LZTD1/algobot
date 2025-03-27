package defaultHandler

import "gopkg.in/telebot.v4"

type ContextHandler interface {
	CanHandle(ctx telebot.Context) bool
	Process(ctx telebot.Context) error
}
