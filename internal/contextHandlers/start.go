package contextHandlers

import (
	"gopkg.in/telebot.v4"
	"tgbot/internal/config"
)

type Start struct {
}

func (h *Start) Message() string {
	return "/start"
}
func (h *Start) Process(ctx telebot.Context) Response {
	return Response{
		Message:  config.StartText,
		Keyboard: config.StartKeyboard,
	}
}
