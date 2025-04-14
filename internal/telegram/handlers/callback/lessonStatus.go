package callback

import (
	"algobot/internal/lib/logger/sl"
	"algobot/internal/services/backoffice"
	"fmt"
	"gopkg.in/telebot.v4"
	"log/slog"
	"strings"
)

type LessonStatuser interface {
	SetLessonStatus(uid int64, groupID string, lessonID string, status backoffice.LessonStatus, traceID interface{}) error
}

func LessonStatus(ls LessonStatuser, status backoffice.LessonStatus, log *slog.Logger) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		const op = "callback.LessonStatus"

		traceID := ctx.Get("trace_id")
		log := log.With(
			slog.String("op", op),
			slog.Any("trace_id", traceID),
		)
		uid := ctx.Sender().ID

		var data []string
		switch status {
		case backoffice.CloseLesson:
			data = strings.Split(strings.TrimPrefix(ctx.Callback().Data, "\fclose_lesson_"), "_")
		case backoffice.OpenLesson:
			data = strings.Split(strings.TrimPrefix(ctx.Callback().Data, "\fopen_lesson_"), "_")
		}
		if len(data) != 2 {
			log.Warn("data is not correct", slog.Any("data", data))
			return ctx.Send("⚠️ Ошибка при анализе данных от кнопки")
		}

		if err := ls.SetLessonStatus(uid, data[0], data[1], status, traceID); err != nil {
			log.Warn("error while refreshing group", sl.Err(err))
			return fmt.Errorf("%s error while refreshing group: %w", op, err)
		}

		return ctx.Send("Статус переключен!")
	}
}
