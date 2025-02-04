package callbackHandlers

import (
	"errors"
	"gopkg.in/telebot.v4"
	"tgbot/internal/config"
	appError "tgbot/internal/error"
	"tgbot/internal/helpers"
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
		if errors.Is(err, appError.ErrHasNone) {
			return ctx.Edit("Я не смог найти ни одной группы, может быть дело в cookie?")
		}
		return helpers.LogError(err, ctx, "Ошибка при обновлении групп!")
	}

	return ctx.Edit(config.UpdateEnd)
}
