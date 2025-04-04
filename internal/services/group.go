package services

import (
	"algobot/internal/domain/models"
	"algobot/internal/lib/logger/sl"
	"algobot/internal/lib/sort"
	"fmt"
	"log/slog"
)

type GroupGetter interface {
	Groups(uid int64) ([]models.Group, error)
}

type Group struct {
	log    *slog.Logger
	getter GroupGetter
}

func NewGroup(log *slog.Logger, getter GroupGetter) *Group {
	return &Group{log: log, getter: getter}
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
