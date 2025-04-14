package text

import (
	"algobot/internal/domain/telegram/keyboards"
	"algobot/internal/lib/logger/sl"
	"fmt"
	"gopkg.in/telebot.v4"
	"log/slog"
	"strings"
)

type UserInformer interface {
	Cookies(uid int64) (string, error)
	Notification(uid int64) (bool, error)
}

func NewSettings(uInformer UserInformer, log *slog.Logger) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		const op = "text.NewSettings"

		uid := ctx.Sender().ID
		log = log.With(
			slog.String("op", op),
			slog.Any("trace_id", ctx.Get("trace_id")),
		)

		cookies, err := uInformer.Cookies(uid)
		if err != nil {
			log.Warn("error while get cookies", sl.Err(err))
			return fmt.Errorf("%s: error while get cookies %w", op, err)
		}

		notification, err := uInformer.Notification(uid)
		if err != nil {
			log.Warn("error while get notification", sl.Err(err))
			return fmt.Errorf("%s: error while get notification %w", op, err)
		}

		return ctx.Send(GetSettingsMessage(cookies, notification), keyboards.Settings())
	}
}

func GetSettingsMessage(cookies string, notification bool) string {
	sb := strings.Builder{}
	sb.WriteString("🔧 Ваши настройки:\n")
	sb.WriteString("\nКуки: ")
	if cookies != "" {
		sb.WriteString("✅")
	} else {
		sb.WriteString("✖️")
	}
	sb.WriteString("\nУведомление от чата:")
	if notification {
		sb.WriteString("✅")
	} else {
		sb.WriteString("✖️")
	}
	return sb.String()
}
