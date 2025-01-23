package defaultHandler

import "gopkg.in/telebot.v4"

type ContextHandler interface {
	CanHandle(ctx telebot.Context) bool
	Process(ctx telebot.Context) error
}

type ActionType string

type Response struct {
	Action   ActionType
	Message  string
	Keyboard *telebot.ReplyMarkup
}
