package callback

import (
	"algobot/internal/domain/models"
	"algobot/internal/lib/logger/sl"
	"fmt"
	"gopkg.in/telebot.v4"
	"log/slog"
	"strings"
)

type GetterCreds interface {
	Creds(uid int64, groupID string, traceID interface{}) ([]models.Credential, error)
}

func GetCreds(creds GetterCreds, log *slog.Logger) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		const op = "callback.GetCreds"

		traceID := ctx.Get("trace_id")
		log := log.With(
			slog.String("op", op),
			slog.Any("trace_id", traceID),
		)
		uid := ctx.Sender().ID

		groupID := strings.TrimPrefix(ctx.Callback().Data, "\fget_creds_")

		c, err := creds.Creds(uid, groupID, traceID)
		if err != nil {
			log.Warn("error while get creds", sl.Err(err))
			return fmt.Errorf("%s error while get creds: %w", op, err)
		}

		return ctx.Send(getCredsMsg(c), telebot.ModeHTML)
	}
}

func getCredsMsg(c []models.Credential) string {
	sb := strings.Builder{}
	for _, credential := range c {
		sb.WriteString("<i>")
		sb.WriteString(credential.Fullname)
		sb.WriteString("</i>")
		sb.WriteString(" - ")
		sb.WriteString(credential.Login)
		sb.WriteString(" : ")
		sb.WriteString(credential.Password)
		sb.WriteString("\n")
	}
	return sb.String()
}
