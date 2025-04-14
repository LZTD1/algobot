package callback

import (
	"algobot/internal/lib/logger/sl"
	"fmt"
	"gopkg.in/telebot.v4"
	"log/slog"
)

type NotificationChanger interface {
	SetNotification(uid int64, isEnable bool) error
	Notification(uid int64) (bool, error)
}

func NewChangeNotification(n NotificationChanger, log *slog.Logger) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		const op = "callback.NewChangeNotification"
		log := log.With(
			slog.String("op", op),
			slog.Any("trace_id", ctx.Get("trace_id")),
		)

		uid := ctx.Sender().ID

		nstat, err := n.Notification(uid)
		if err != nil {
			log.Warn("error while getting notification", sl.Err(err))
			return fmt.Errorf("%s get notif: %w", op, err)
		}

		if err := n.SetNotification(uid, !nstat); err != nil {
			log.Warn("error while set notification", sl.Err(err))
			return fmt.Errorf("%s set notif: %w", op, err)
		}

		return ctx.Edit("Уведомления переключены")
	}
}
