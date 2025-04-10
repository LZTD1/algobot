package text

import (
	"algobot/internal/domain/models"
	"algobot/internal/lib/logger/sl"
	"algobot/internal/services/groups"
	"errors"
	"fmt"
	"gopkg.in/telebot.v4"
	"log/slog"
	"strings"
	"time"
)

type ActualGroup interface {
	CurrentGroup(uid int64, time time.Time, traceID interface{}) (models.CurrentGroup, error)
}

func NewMissingKids(log *slog.Logger, actualGroup ActualGroup) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		const op = "text.NewMissingKids"
		traceID := ctx.Get("trace_id")
		uid := ctx.Sender().ID

		log := log.With(
			slog.String("op", op),
			slog.Any("trace_id", traceID),
		)

		group, err := actualGroup.CurrentGroup(uid, time.Now(), traceID)
		if err != nil {
			if errors.Is(err, groups.ErrNoGroups) {
				return ctx.Send("В данный момент, никакой группы не найдено!")
			}
			if errors.Is(err, groups.ErrNotValidCookie) {
				return ctx.Send("Вам необходимо установить свои cookie!")
			}
			log.Warn("error while fetching CurrentGroup", sl.Err(err))
			return fmt.Errorf("%s error while fetching CurrentGroup: %w", op, err)
		}

		return ctx.Send(getMsg(group), telebot.ModeMarkdown)
	}
}

func getMsg(gr models.CurrentGroup) string {
	miss := strings.Builder{}
	miss.WriteString("\n```Отсутствующие\n")
	for _, kid := range gr.MissingKids {
		miss.WriteString(kid.Fullname)
		if kid.Count > 1 {
			miss.WriteString(fmt.Sprintf(" (Уже %d занятие)", kid.Count))
		}
		miss.WriteString("\n")
	}
	miss.WriteString("```")

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Группа: %s", gr.Title))
	sb.WriteString(fmt.Sprintf("\nЛекция: %s\n", gr.Lesson))
	sb.WriteString(fmt.Sprintf("\nОбщее число детей: %d", len(gr.Kids)))
	sb.WriteString(fmt.Sprintf("\nОтсутствуют: %d\n", len(gr.MissingKids)))
	sb.WriteString(miss.String())

	return sb.String()
}
