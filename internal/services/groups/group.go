package groups

import (
	"algobot/internal/domain/models"
	"algobot/internal/lib/logger/sl"
	"algobot/internal/lib/sort"
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrNotValidCookie = errors.New("not a valid cookie")
	ErrNoGroups       = errors.New("groups not found")
)

type GroupGetter interface {
	Groups(uid int64) ([]models.Group, error)
}

type Group struct {
	log          *slog.Logger
	getter       GroupGetter
	groupFetcher GroupFetcher
	domainSetter DomainSetter
}

func NewGroup(log *slog.Logger, getter GroupGetter, setter DomainSetter, groupFetcher GroupFetcher) *Group {
	return &Group{log: log, getter: getter, domainSetter: setter, groupFetcher: groupFetcher}
}

func (g *Group) Groups(uid int64, traceID interface{}) ([]models.Group, error) {
	const op = "services.Group.Groups"
	log := g.log.With(
		slog.String("op", op),
		slog.Any("trace_id", traceID),
	)

	groups, err := g.getter.Groups(uid)
	if err != nil {
		log.Warn("error while get groups", sl.Err(err))
		return []models.Group{}, fmt.Errorf("%s error while get groups: %w", op, err)
	}

	sort.GroupsByDate(groups)

	return groups, nil
}
