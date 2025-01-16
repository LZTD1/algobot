package contextHandlers

import (
	"fmt"
	"gopkg.in/telebot.v4"
	"strings"
	"tgbot/internal/config"
	"tgbot/internal/domain"
	"tgbot/internal/service"
)

type MissingKids struct {
	s service.Service
}

func NewMissingKids(s service.Service) *MissingKids {
	return &MissingKids{s}
}

func (m *MissingKids) Message() string {
	return "Получить отсутсвующих"
}
func (m *MissingKids) Process(ctx telebot.Context) Response {
	g, e := m.s.CurrentGroup(ctx.Message().Sender.ID, ctx.Message().Time())
	if e != nil {
		return Response{
			Message: e.Error(),
		}
	}
	return Response{Message: message(g)}
}

func message(g domain.Group) string {
	return fmt.Sprintf(
		"%s%s\n%s%s\n%s%d\n%s%d\n%s",
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
