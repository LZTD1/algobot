package defaultState

import (
	"algobot/internal_old/config"
	appError "algobot/internal_old/error"
	"algobot/internal_old/helpers"
	"algobot/internal_old/models"
	"algobot/internal_old/service"
	"errors"
	"fmt"
	"gopkg.in/telebot.v4"
	"strings"
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
	t := ctx.Message().Time()
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

	return ctx.Send(msgMissingKids(g, actual, allKids), getMissingKidsKeyboard(g, actual), telebot.ModeMarkdown)
}

func getMissingKidsKeyboard(g models.Group, actual models.ActualInformation) *telebot.ReplyMarkup {
	markup := telebot.ReplyMarkup{ResizeKeyboard: true}
	markup.Inline(
		markup.Row(markup.Data(config.CloseLessonBtn, fmt.Sprintf("close_lesson_%d_%d", g.GroupID, actual.LessonId)), markup.Data(config.OpenLessonBtn, fmt.Sprintf("open_lesson_%d_%d", g.GroupID, actual.LessonId))),
		markup.Row(markup.Data(config.GetCredsBtn, fmt.Sprintf("get_creds_%d", g.GroupID))),
	)

	return &markup
}

func msgMissingKids(g models.Group, actual models.ActualInformation, kids models.AllKids) string {
	miss := strings.Builder{}
	miss.WriteString("\n```Отсутствующие\n")
	missingCount := 0
	for _, kid := range actual.MissingKids {
		if v, ok := kids[kid.Id]; ok == true {
			missingCount++
			miss.WriteString(fmt.Sprintf("%s", v.FullName))
			if kid.Count > 1 {
				miss.WriteString(fmt.Sprintf(" (Уже %d занятие)", kid.Count))
			}
			miss.WriteString("\n")
		}
	}
	miss.WriteString("```")

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s%s", config.GroupName, g.Title))
	sb.WriteString(fmt.Sprintf("\n%s%s\n", config.Lection, actual.LessonTitle))
	sb.WriteString(fmt.Sprintf("\n%s%d", config.TotalKids, len(kids)))
	sb.WriteString(fmt.Sprintf("\n%s%d\n", config.MissingKids, missingCount))
	sb.WriteString(miss.String())

	return sb.String()
}
