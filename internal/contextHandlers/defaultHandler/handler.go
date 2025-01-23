package defaultHandler

import "gopkg.in/telebot.v4"

type ContextHandler interface {
	CanHandle(ctx telebot.Context) bool
	Process(ctx telebot.Context) Response
}

type ActionType string

var (
	SendMessage ActionType = "sendMessage"
	EditMessage ActionType = "editMessage"
)

type Response struct {
	Action   ActionType
	Message  string
	Keyboard *telebot.ReplyMarkup
}
