package defaultState

import (
	"errors"
	"fmt"
	"gopkg.in/telebot.v4"
	"strings"
	"tgbot/internal/config"
	"tgbot/internal/domain"
	appError "tgbot/internal/error"
	"tgbot/internal/helpers"
	"tgbot/internal/service"
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
	g, e := m.s.CurrentGroup(ctx.Message().Sender.ID, ctx.Message().Time())
	if e != nil {
		if errors.Is(e, appError.ErrHasNone) {
			return ctx.Send(config.CurrentGroupDontFind)
		}
		// TODO зарефачить эту логику
		t := helpers.LogWithRandomToken(e)
		return ctx.Send(t + " | Ошибка при получении групп!")
	}
	return ctx.Send(message(g), telebot.ModeMarkdown)
}

func message(g domain.Group) string {
	return fmt.Sprintf(
		"%s%s\n%s%s\n\n%s%d\n%s%d\n\n```Отсутсвующие\n%s\n```",
		config.GroupName,
		g.Name,
		config.Lection,
		g.Lesson,
		config.TotalKids,
		g.AllKids,
		config.MissingKids,
		len(g.MissingKids),
		strings.Join(g.MissingKids, "\n"),
	)
}
