package text

import (
	"algobot/internal/domain/telegram/keyboards"
	"algobot/internal/lib/fsm"
	"algobot/internal/lib/logger/sl"
	"fmt"
	"gopkg.in/telebot.v4"
	"log/slog"
)

type CookieSetter interface {
	SetCookie(uid int64, cookie string) error
}

type CookieStater interface {
	SetState(uid int64, state fsm.State)
}

func NewSendingCookie(log *slog.Logger, cookieSetter CookieSetter, stater CookieStater) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		const op = "text.NewSendingCookie"

		log = log.With(
			slog.String("op", op),
			slog.Any("trace_id", ctx.Get("trace_id")),
		)

		uid := ctx.Sender().ID
		cookie := ctx.Message().Text

		if err := cookieSetter.SetCookie(uid, cookie); err != nil {
			log.Warn("error while setting cookie", sl.Err(err))
			return fmt.Errorf("%s error while setting cookie: %w", op, err)
		}

		stater.SetState(uid, fsm.Default)
		return ctx.Send("Cookie успешно установлены", keyboards.Start())
	}
}
