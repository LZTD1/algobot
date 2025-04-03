package auth

import (
	"algobot/internal/lib/logger/sl"
	"fmt"
	"gopkg.in/telebot.v4"
	"log/slog"
)

type Auther interface {
	IsRegistered(uid int64) (bool, error)
	Register(uid int64) error
}

func New(auth Auther, log *slog.Logger) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		log = log.With(
			slog.String("component", "middleware/auth"),
		)

		return func(ctx telebot.Context) error {
			traceID := ctx.Get("trace_id")
			log.With("trace_id", traceID)

			uid := ctx.Sender().ID

			isReg, err := auth.IsRegistered(uid)
			if err != nil {
				log.Warn("error while checking if user exists", sl.Err(err))
				return fmt.Errorf("error while checking if user exists: %w", err)
			}
			if !isReg {
				if err := auth.Register(uid); err != nil {
					log.Warn("error while register user", sl.Err(err))
					return fmt.Errorf("error while register user: %w", err)
				}

				log.Info("user is registered", slog.Int64("uid", uid))
			}

			return next(ctx)
		}
	}
}
