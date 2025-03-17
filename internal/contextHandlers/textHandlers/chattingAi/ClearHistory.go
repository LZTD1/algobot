package chattingAi

import (
	"gopkg.in/telebot.v4"
	"tgbot/internal/config"
	"tgbot/internal/helpers"
	"tgbot/internal/service"
)

type ClearHistory struct {
	ai service.AIService
}

func NewClearHistory(ai service.AIService) *ClearHistory {
	return &ClearHistory{ai: ai}
}

func (c ClearHistory) CanHandle(ctx telebot.Context) bool {
	if ctx.Message().Text == config.ClearHistoryBtn.Text {
		return true
	}
	return false
}

func (c ClearHistory) Process(ctx telebot.Context) error {
	err := c.ai.ClearAllHistory(ctx.Sender().ID)
	if err != nil {
		return helpers.LogError(err, ctx, "Ошибка при отчистке памяти!")
	}

	return ctx.Send("Успешно отчищено!")
}
