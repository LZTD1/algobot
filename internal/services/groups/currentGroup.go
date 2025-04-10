package groups

import (
	"algobot/internal/domain/models"
	"algobot/internal/lib/logger/sl"
	"fmt"
	"log/slog"
	"time"
)

func (g *Group) CurrentGroup(uid int64, time time.Time, traceID interface{}) (models.CurrentGroup, error) {
	const op = "services.groups.CurrentGroup"
	log := g.log.With(
		slog.String("op", op),
		slog.Any("trace_id", traceID),
	)

	cookie, err := g.domain.Cookies(uid)
	if err != nil {
		log.Warn("failed to get cookies", sl.Err(err))
		return models.CurrentGroup{}, fmt.Errorf("%s failed to get cookies: %w", op, err)
	}
	if cookie == "" {
		return models.CurrentGroup{}, fmt.Errorf("%s cookie is empty: %w", op, ErrNotValidCookie)
	}

	panic("todo implementation")
}
