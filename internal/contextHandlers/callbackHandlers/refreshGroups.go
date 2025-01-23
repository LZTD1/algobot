package callbackHandlers

import (
	"gopkg.in/telebot.v4"
	"tgbot/internal/contextHandlers/defaultHandler"
	"tgbot/internal/service"
)

type RefreshGroups struct {
	svc service.Service
}

func NewRefreshGroups(svc service.Service) *RefreshGroups {
	return &RefreshGroups{svc: svc}
}

func (r RefreshGroups) CanHandle(ctx telebot.Context) bool {
	if ctx.Callback().Data == "refresh_groups" {
		return true
	}
	return false
}

func (r RefreshGroups) Process(ctx telebot.Context) defaultHandler.Response {
	return defaultHandler.Response{
		Message: "",
	}
}
