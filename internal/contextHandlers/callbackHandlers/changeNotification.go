package callbackHandlers

import (
	"gopkg.in/telebot.v4"
	"tgbot/internal/contextHandlers/defaultHandler"
	"tgbot/internal/contextHandlers/textHandlers"
	"tgbot/internal/service"
)

type ChangeNotification struct {
	svc      service.Service
	settings defaultHandler.ContextHandler
}

func NewChangeNotification(svc service.Service) *ChangeNotification {
	return &ChangeNotification{svc: svc, settings: textHandlers.NewSettings(svc)}
}

func (c ChangeNotification) CanHandle(ctx telebot.Context) bool {
	if ctx.Callback().Data == "change_notification" {
		return true
	}
	return false
}

func (c ChangeNotification) Process(ctx telebot.Context) defaultHandler.Response {
	uid := ctx.Callback().Sender.ID
	notify := c.svc.Notification(uid)
	c.svc.SetNotification(uid, !notify)

	settings := c.settings.Process(ctx)
	return defaultHandler.Response{
		Action:   defaultHandler.EditMessage,
		Message:  settings.Message,
		Keyboard: settings.Keyboard,
	}
}
