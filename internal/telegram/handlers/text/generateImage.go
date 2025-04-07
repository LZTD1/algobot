package text

import (
	"algobot/internal/lib/logger/sl"
	"fmt"
	"gopkg.in/telebot.v4"
	"log/slog"
	"strings"
)

type GeneratorImage interface {
	GenerateImage(uid int64, promt string, traceID interface{}) (string, error)
}

func GenerateImage(generator GeneratorImage, log *slog.Logger) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		const op = "text.GenerateImage"
		uid := ctx.Sender().ID
		traceID := ctx.Get("trace_id")
		promt := getPromtFromMessage(ctx.Message().Text)
		log = log.With(
			slog.String("op", op),
			slog.Any("trace_id", traceID),
		)

		msg, err := ctx.Bot().Send(telebot.ChatID(uid), "⚙️ Генерирую изображение ...")
		if err != nil {
			log.Warn("failed to send prepare msg", sl.Err(err))
			return fmt.Errorf("%s failed to send prepare msg: %w", op, err)
		}

		imgURL, err := generator.GenerateImage(uid, promt, traceID)
		if err != nil {
			log.Warn("failed to generate image", sl.Err(err))
			ctx.Bot().Edit(msg, "⚠️ К сожалению, я не смог сгенерировать изображение, попробуйте снова чуть позже")
			return fmt.Errorf("%s failed to generate image: %w", op, err)
		}

		_, err = ctx.Bot().Edit(msg, &telebot.Photo{
			File: telebot.FromURL(imgURL),
		})

		return err
	}
}

func getPromtFromMessage(text string) string {
	return strings.TrimSpace(strings.TrimLeft(text, "/image"))
}
