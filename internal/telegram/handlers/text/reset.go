package text

import (
	"algobot/internal/lib/logger/sl"
	"fmt"
	"gopkg.in/telebot.v4"
	"log/slog"
)

type Reseter interface {
	ResetHistory(uid int64, traceID interface{}) error
}

func NewReset(reseter Reseter, log *slog.Logger) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		const op = "text.NewReset"

		uid := ctx.Sender().ID
		traceID := ctx.Get("trace_id")
		log = log.With(
			slog.String("op", op),
			slog.Any("trace_id", traceID),
		)

		if err := reseter.ResetHistory(uid, traceID); err != nil {
			log.Warn("failed to reset history", sl.Err(err))
			return fmt.Errorf("%s failed to reset history: %w", op, err)
		}

		return ctx.Send("История успешно отчищена")
	}
}
