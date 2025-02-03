package defaultState

import (
	"errors"
	"fmt"
	"gopkg.in/telebot.v4"
	"strings"
	"tgbot/internal/config"
	appError "tgbot/internal/error"
	"tgbot/internal/helpers"
	"tgbot/internal/models"
	"tgbot/internal/service"
	"time"
)

type MissingKids struct {
	s service.Service
}

func NewMissingKids(s service.Service) *MissingKids {
	return &MissingKids{s}
}

func (m *MissingKids) CanHandle(ctx telebot.Context) bool {
	if ctx.Message().Text == "Получить отсутсвующих" {
		return true
	}
	return false
}

func (m *MissingKids) Process(ctx telebot.Context) error {
	t := time.Date(2025, 2, 8, 9, 40, 0, 0, time.UTC)
	uid := ctx.Message().Sender.ID

	g, e := m.s.CurrentGroup(uid, t)
	if e != nil {
		if errors.Is(e, appError.ErrHasNone) {
			return ctx.Send(config.CurrentGroupDontFind)
		}

		return helpers.LogError(e, ctx, "Произошла непредвиденная ошибка при попытке получить текущую группу")
	}

	actual, err := m.s.ActualInformation(uid, t, g.GroupID)
	if err != nil {
		if errors.Is(err, appError.ErrNotValid) {
			return ctx.Send(config.CookieNotSetException)
		}

		return helpers.LogError(e, ctx, "Произошла непредвиденная ошибка при попытке подгрузить информацию о группе")
	}

	allKids, err := m.s.AllKidsNames(uid, g.GroupID)
	if err != nil {
		if errors.Is(err, appError.ErrNotValid) {
			return ctx.Send(config.CookieNotSetException)
		}

		return helpers.LogError(e, ctx, "Произошла непредвиденная ошибка при попытке подгрузить имена детей")
	}

	return ctx.Send(msg(g, actual, allKids), telebot.ModeMarkdown)
}

func msg(g models.Group, actual models.ActualInformation, kids models.AllKids) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s%s", config.GroupName, g.Title))
	sb.WriteString(fmt.Sprintf("\n%s%s\n", config.Lection, actual.LessonTitle))
	sb.WriteString(fmt.Sprintf("\n%s%d", config.TotalKids, len(kids)))
	sb.WriteString(fmt.Sprintf("\n%s%d\n", config.MissingKids, len(actual.MissingKids)))
	sb.WriteString("\n```Отсутствующие\n")
	for _, kid := range actual.MissingKids {
		sb.WriteString(fmt.Sprintf("%s\n", kids[kid]))
	}
	sb.WriteString("```")

	return sb.String()
}
