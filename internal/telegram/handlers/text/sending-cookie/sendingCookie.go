package sendingCookie

import (
	"gopkg.in/telebot.v4"
	"log/slog"
)

type SendingCookie struct {
	log *slog.Logger
}

func New(log *slog.Logger) *SendingCookie {
	return &SendingCookie{}
}

func (d SendingCookie) Handle(c telebot.Context) error {
	panic("implement me")
}
