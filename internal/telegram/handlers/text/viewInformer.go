package text

import (
	"algobot/internal/domain/models"
	"algobot/internal/lib/logger/sl"
	"algobot/internal/lib/serdes"
	"gopkg.in/telebot.v4"
	"log/slog"
	"strings"
)

type FetchView interface {
	GetGroupView(uid int64, groupID int) models.GroupView
}

type Serializator interface {
	GetType(encoded string) (serdes.SerType, error)
}

func ViewInformer(ser Serializator, log *slog.Logger) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		const op = "text.GenerateImage"

		uid := ctx.Sender().ID
		traceID := ctx.Get("trace_id")
		data := getData(ctx.Message().Text)

		log = log.With(
			slog.String("op", op),
			slog.Any("trace_id", traceID),
		)

		serType, err := ser.GetType(data)
		if err != nil {
			log.Warn("can't get ser type", sl.Err(err))
			return ctx.Send("⚠️ Ошибка при расшифровке запроса!")
		}

		switch serType {
		case serdes.UserType:
			return nil
		case serdes.GroupType:
			return nil
		default:
			return ctx.Send("⚠️ Не удалось определить обработчик")
		}
	}
}

func getData(text string) string {
	return strings.TrimSpace(strings.TrimLeft(text, "/start"))
}
