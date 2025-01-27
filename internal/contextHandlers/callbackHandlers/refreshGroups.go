package callbackHandlers

import (
	"gopkg.in/telebot.v4"
	"tgbot/internal/config"
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

func (r RefreshGroups) Process(ctx telebot.Context) error {
	err := ctx.Edit(config.UpdateStarted)
	if err != nil {
		return err
	}

	err = r.svc.RefreshGroups(ctx.Sender().ID)
	if err != nil {
		ctx.Edit(err.Error())
	}

	err = ctx.Edit(config.UpdateEnd)
	return err
}
