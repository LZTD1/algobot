package groups

import (
	"algobot/internal/domain/models"
	"algobot/internal/lib/logger/sl"
	"fmt"
	"log/slog"
)

type DomainSetter interface {
	SetGroups(uid int64, groups []models.Group) error
	Cookies(uid int64) (string, error)
}

type GroupFetcher interface {
	Group(cookie string) ([]models.Group, error)
}

func (g *Group) RefreshGroup(uid int64, traceID interface{}) error {
	const op = "services.Group.RefreshGroup"

	log := g.log.With(
		slog.String("op", op),
		slog.Any("trace_id", traceID),
	)

	cookie, err := g.domain.Cookies(uid)
	if err != nil {
		log.Warn("failed to get cookies", sl.Err(err))
		return fmt.Errorf("%s failed to get cookies: %w", op, err)
	}
	if cookie == "" {
		return fmt.Errorf("%s cookie is empty: %w", op, ErrNotValidCookie)
	}

	groups, err := g.groupFetcher.Group(cookie)
	if err != nil {
		log.Warn("failed to fetch groups", sl.Err(err))
		return fmt.Errorf("%s failed to fetch groups: %w", op, err)
	}

	if len(groups) == 0 {
		return fmt.Errorf("%s no groups found: %w", op, ErrNoGroups)
	}

	if err := g.domain.SetGroups(uid, groups); err != nil {
		log.Warn("failed to set groups", sl.Err(err))
		return fmt.Errorf("%s failed to set groups: %w", op, err)
	}

	return nil
}
