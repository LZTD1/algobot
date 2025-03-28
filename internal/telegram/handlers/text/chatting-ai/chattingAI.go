package chattingAI

import (
	"gopkg.in/telebot.v4"
	"log/slog"
)

type ChattingAI struct {
	log *slog.Logger
}

func New(log *slog.Logger) *ChattingAI {
	return &ChattingAI{}
}

func (d ChattingAI) Handle(c telebot.Context) error {
	panic("implement me")
}
