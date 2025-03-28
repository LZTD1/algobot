package trace

import (
	"fmt"
	"github.com/google/uuid"
	tele "gopkg.in/telebot.v4"
	"log/slog"
)

// New
// Generates a unique trace_id for each request, which can be accessed from the context like this:
// var context tele.Context
// traceID := context.Get("trace_id")
func New(log *slog.Logger) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		log = log.With(
			slog.String("component", "middleware/trace"),
		)

		return func(c tele.Context) error {
			newUUID, err := uuid.NewUUID()
			if err != nil {
				log.Warn("failed to generate UUID")
			}
			traceID := fmt.Sprintf("%d/%s/%s", c.Sender().ID, c.Sender().Username, newUUID.String())

			c.Set("trace_id", traceID)

			return next(c)
		}
	}
}
