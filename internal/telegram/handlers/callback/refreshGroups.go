package callback

import (
	"algobot/internal/lib/logger/sl"
	"algobot/internal/services/groups"
	"errors"
	"fmt"
	"log/slog"

	"gopkg.in/telebot.v4"
)

type GroupRefresher interface {
	RefreshGroup(uid int64, traceID interface{}) error
}

func RefreshGroup(refresher GroupRefresher, log *slog.Logger) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		const op = "callback.RefreshGroup"

		traceID := ctx.Get("trace_id")
		log := log.With(
			slog.String("op", op),
			slog.Any("trace_id", traceID),
		)
		uid := ctx.Sender().ID

		if err := ctx.Edit("⚙️ Обновляю группы..."); err != nil {
			log.Warn("error while editing message", sl.Err(err))
			return fmt.Errorf("%s error while editing message: %w", op, err)
		}

		if err := refresher.RefreshGroup(uid, traceID); err != nil {
			if errors.Is(err, groups.ErrNoGroups) {
				return ctx.Edit("Не найдено ни 1 группы!\nВозможно, проблема в актуальности cookie!")
			}
			if errors.Is(err, groups.ErrNotValidCookie) {
				return ctx.Edit("У вас не валидные cookie!")
			}
			log.Warn("error while refreshing group", sl.Err(err))
			return fmt.Errorf("%s error while refreshing group: %w", op, err)
		}

		return ctx.Edit("Успешно обновлено!")
	}
}
