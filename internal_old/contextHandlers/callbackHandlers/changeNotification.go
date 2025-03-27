package callbackHandlers

import (
	"algobot/internal_old/contextHandlers/defaultHandler"
	"algobot/internal_old/contextHandlers/textHandlers/defaultState"
	"algobot/internal_old/helpers"
	"algobot/internal_old/service"
	"gopkg.in/telebot.v4"
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
		return helpers.LogError(err, ctx, "Ошибка при получении нотификаций!")
	}
	err = c.svc.SetNotification(uid, !notify)
	if err != nil {
		return helpers.LogError(err, ctx, "Ошибка при установлении нотификаций!")
	}

	return ctx.Edit("Настройки уведомлений были изменены!")
}
