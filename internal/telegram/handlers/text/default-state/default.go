package defaultstate

import (
	"gopkg.in/telebot.v4"
	"log/slog"
)

type DefaultState struct {
	log *slog.Logger
}

func New(log *slog.Logger) *DefaultState {
	return &DefaultState{
		log: log,
	}
}

func (d *DefaultState) Handle(c telebot.Context) error {
	panic("implement me")
}
