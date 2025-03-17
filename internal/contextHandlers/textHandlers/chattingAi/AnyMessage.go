package chattingAi

import (
	"gopkg.in/telebot.v4"
	"strconv"
	"tgbot/internal/config"
	"tgbot/internal/helpers"
	"tgbot/internal/schedulers"
	"tgbot/internal/service"
)

type AnyMessage struct {
	s service.AIService
}

func NewAnyMessage(s service.AIService) *AnyMessage {
	return &AnyMessage{s: s}
}

func (a AnyMessage) CanHandle(ctx telebot.Context) bool {
	if ctx.Message().Text != config.BackBtn.Text && ctx.Message().Text != config.ClearHistoryBtn.Text {
		return true
	}
	return false
}

func (a AnyMessage) Process(ctx telebot.Context) error {
	m, _ := ctx.Bot().Send(schedulers.RecipientUser{
		ID: strconv.FormatInt(ctx.Sender().ID, 10),
	}, "⚙️ Думаю как ответить ....")
	suggest, err := a.s.GetSuggestion(ctx.Sender().ID, ctx.Message().Text)
	if err != nil {
		return helpers.LogError(err, ctx, "Ошибка при получении ответа от бота!")
	}
	_, err = ctx.Bot().Edit(m, suggest, telebot.ModeMarkdown)
	return err
}
