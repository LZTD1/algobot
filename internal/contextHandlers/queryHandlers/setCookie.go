package queryHandlers

import (
	"gopkg.in/telebot.v4"
	"tgbot/internal/contextHandlers"
)

type SetCookie struct {
}

func (s SetCookie) CanHandle() string {
	return "set_cookie"
}

func (s SetCookie) Process(ctx telebot.Context) contextHandlers.Response {
	//TODO implement me
	panic("implement me")
}
