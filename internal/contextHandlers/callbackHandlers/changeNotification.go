package callbackHandlers

import (
	"gopkg.in/telebot.v4"
	"tgbot/internal/contextHandlers/defaultHandler"
	"tgbot/internal/contextHandlers/textHandlers/defaultState"
	"tgbot/internal/helpers"
	"tgbot/internal/service"
)

type ChangeNotification struct {
	svc      service.Service
	settings defaultHandler.ContextHandler
}

func NewChangeNotification(svc service.Service) *ChangeNotification {
	return &ChangeNotification{svc: svc, settings: defaultState.NewSettings(svc)}
}

func (c ChangeNotification) CanHandle(ctx telebot.Context) bool {
	if ctx.Callback().Data == "change_notification" {
		return true
	}
	return false
}

func (c ChangeNotification) Process(ctx telebot.Context) error {
	uid := ctx.Callback().Sender.ID
	notify, err := c.svc.Notification(uid)
	if err != nil {
		t := helpers.LogWithRandomToken(err)
		return ctx.Send(t + " | Ошибка при получении нотификаций! ")
	}
	err = c.svc.SetNotification(uid, !notify)
	if err != nil {
		t := helpers.LogWithRandomToken(err)
		return ctx.Send(t + " | Ошибка при установлении нотификаций! ")
	}

	return c.settings.Process(ctx)
}
