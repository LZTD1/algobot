package contextHandlers

import "gopkg.in/telebot.v4"

type ContextHandler interface {
	Message() string
	Process(ctx telebot.Context) Response
}

type Response struct {
	Message  string
	Keyboard *telebot.ReplyMarkup
}
