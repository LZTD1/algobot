package text

import (
	"algobot/internal/lib/logger/sl"
	"fmt"
	"gopkg.in/telebot.v4"
	"log/slog"
)

type Chatter interface {
	ChatAI(uid int64, message string, traceID interface{}) (string, error)
}

func ChatAI(chatter Chatter, log *slog.Logger) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		const op = "text.ChatAI"

		uid := ctx.Sender().ID
		message := ctx.Message().Text
		traceID := ctx.Get("trace_id")
		log = log.With(
			slog.String("op", op),
			slog.Any("trace_id", traceID),
		)

		msg, err := ctx.Bot().Reply(ctx.Message(), "⚙️ Думаю что ответить ...")
		if err != nil {
			log.Warn("failed to send prepare msg", sl.Err(err))
			return fmt.Errorf("%s failed to send prepare msg: %w", op, err)
		}

		resp, err := chatter.ChatAI(uid, message, traceID)
		if err != nil {
			log.Warn("failed to ChatAI", sl.Err(err))
			ctx.Bot().Edit(msg, "⚠️ К сожалению, я не смог ответить на ваше сообщение, попробуйте снова чуть позже")
			return fmt.Errorf("%s failed to chatting ai: %w", op, err)
		}

		_, err = ctx.Bot().Edit(msg, resp, telebot.ModeMarkdown)
		return err
	}
}
