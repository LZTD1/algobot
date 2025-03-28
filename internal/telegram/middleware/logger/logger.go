package logger

import (
	"fmt"
	tele "gopkg.in/telebot.v4"
	"log/slog"
)

func New(log *slog.Logger) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		log = log.With(
			slog.String("component", "middleware/logger"),
		)

		return func(c tele.Context) error {
			if msg := c.Message(); msg != nil {
				log.Info("incoming message",
					slog.Int64("from_id", c.Sender().ID),
					slog.String("from", c.Sender().Username),
					slog.String("full_name", fmt.Sprintf("%s %s", c.Sender().FirstName, c.Sender().LastName)),
					slog.String("message", msg.Text),
				)
			}
			if cb := c.Callback(); cb != nil {
				log.Info("incoming callback",
					slog.Int64("from_id", c.Sender().ID),
					slog.String("from", c.Sender().Username),
					slog.String("full_name", fmt.Sprintf("%s %s", c.Sender().FirstName, c.Sender().LastName)),
					slog.String("message", cb.Data),
				)
			}

			return next(c)
		}
	}
}
