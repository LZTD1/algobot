package defaultState

import (
	"errors"
	"gopkg.in/telebot.v4"
	"strings"
	"tgbot/internal/config"
	appError "tgbot/internal/error"
	"tgbot/internal/helpers"
	"tgbot/internal/service"
	"time"
)

type AbsentKids struct {
	s service.Service
}

func NewAbsentKids(s service.Service) *AbsentKids {
	return &AbsentKids{s: s}
}

func (a AbsentKids) CanHandle(ctx telebot.Context) bool {
	if strings.HasPrefix(ctx.Message().Text, "/abs") {
		return true
	}
	return false
}

func (a AbsentKids) Process(ctx telebot.Context) error {
	if ctx.Message().Payload == "" {
		return ctx.Send("Формат сообщения - '/abs 2025-01-12 15:32'\nВыдаст статистику за 2025г. 12 Января, 15ч 32м")
	}
	t, err := time.Parse("2006-01-02 15:04", ctx.Message().Payload)
	if err != nil {
		return ctx.Send("Формат сообщения - '/abs 2025-01-12 15:32'\nВыдаст статистику за 2025г. 12 Января, 15ч 32м")
	}
	uid := ctx.Message().Sender.ID

	g, e := a.s.CurrentGroup(uid, t)
	if e != nil {
		if errors.Is(e, appError.ErrHasNone) {
			return ctx.Send(config.CurrentGroupDontFind)
		}

		return helpers.LogError(e, ctx, "Произошла непредвиденная ошибка при попытке получить текущую группу")
	}
	actual, err := a.s.ActualInformation(uid, t, g.GroupID)
	if err != nil {
		if errors.Is(err, appError.ErrNotValid) {
			return ctx.Send(config.CookieNotSetException)
		}

		return helpers.LogError(e, ctx, "Произошла непредвиденная ошибка при попытке подгрузить информацию о группе")
	}

	allKids, err := a.s.AllKidsNames(uid, g.GroupID)
	if err != nil {
		if errors.Is(err, appError.ErrNotValid) {
			return ctx.Send(config.CookieNotSetException)
		}

		return helpers.LogError(e, ctx, "Произошла непредвиденная ошибка при попытке подгрузить имена детей")
	}

	return ctx.Send(msgMissingKids(g, actual, allKids), getMissingKidsKeyboard(g, actual), telebot.ModeMarkdown)
}
