package defaultState

import (
	"algobot/internal_old/config"
	"gopkg.in/telebot.v4"
)

type Start struct {
}

func (h *Start) CanHandle(ctx telebot.Context) bool {
	if ctx.Message().Text == "/start" {
		return true
	}
	return false
}
func (h *Start) Process(ctx telebot.Context) error {
	return ctx.Send(config.StartText, config.StartKeyboard)
}
