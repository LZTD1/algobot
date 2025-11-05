package text

import (
	"algobot/internal/domain/telegram/keyboards"
	"algobot/internal/lib/logger/sl"
	"algobot/internal/services/groups"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"gopkg.in/telebot.v4"
)

func NewAbsentKids(actualGroup ActualGroup, log *slog.Logger) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		const op = "text.NewAbsentKids"

		uid := ctx.Sender().ID
		traceID := ctx.Get("trace_id")
		data := getDate(ctx.Message().Text)

		log := log.With(
			slog.String("op", op),
			slog.Any("trace_id", traceID),
		)

		date, err := time.Parse("2006-01-02 15:04", data)
		if err != nil {
			t, err2 := time.Parse("15:04", data)
			if err2 != nil {
				return ctx.Reply("Не удалось распарсить дату, пожалуйста, введите дату в формате YYYY-MM-DD HH:MM\nИли дату за сегодняшний день в формате HH:MM")
			}

			now := time.Now()
			date = time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				t.Hour(),
				t.Minute(),
				0, 0,
				now.Location(),
			)
		}

		group, err := actualGroup.CurrentGroup(uid, date, traceID)
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

		return ctx.Reply(GetMissingMessage(group), keyboards.MissingKids(group.GroupID, group.LessonID), telebot.ModeMarkdown)
	}
}
func getDate(text string) string {
	return strings.TrimSpace(strings.TrimLeft(text, "/abs"))
}
