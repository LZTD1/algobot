package textHandlers

import (
	"gopkg.in/telebot.v4"
	"tgbot/internal/config"
	"tgbot/internal/contextHandlers"
)

type Start struct {
}

func (h *Start) CanHandle(ctx telebot.Context) bool {
	if ctx.Message().Text == "/start" {
		return true
	}
	return false
}
func (h *Start) Process(ctx telebot.Context) contextHandlers.Response {
	return contextHandlers.Response{
		Message:  config.StartText,
		Keyboard: config.StartKeyboard,
	}
}
