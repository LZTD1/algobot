package groups

import (
	"algobot/internal/domain/models"
	"algobot/internal/lib/logger/sl"
	"algobot/internal/lib/sort"
	"errors"
	"fmt"
	"log/slog"
	"time"
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
	domain       DomainSetter
	kidsStats    KidStats
}

func NewGroup(
	log *slog.Logger,
	getter GroupGetter,
	groupFetcher GroupFetcher,
	domain DomainSetter,
	kidsStats KidStats,
) *Group {
	return &Group{
		log:          log,
		getter:       getter,
		groupFetcher: groupFetcher,
		domain:       domain,
		kidsStats:    kidsStats,
	}
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

func CurrentGroup(t time.Time, g []models.Group) (models.Group, error) {
	current, err := GroupsByDay(t, g)
	if err != nil {
		return models.Group{}, err
	}
	for _, group := range current {
		if inDiapazon(-30, 90, t, group.TimeLesson) {
			return group, nil
		}
	}

	return models.Group{}, ErrNoGroups
}
func GroupsByDay(t time.Time, g []models.Group) ([]models.Group, error) {
	var filtered []models.Group

	for _, group := range g {
		if t.Weekday() == group.TimeLesson.Weekday() {
			filtered = append(filtered, group)
		}
	}
	if len(filtered) == 0 {
		return nil, ErrNoGroups
	}

	return filtered, nil
}
func inDiapazon(start, end int, now, group time.Time) bool {
	s := group.Add(time.Duration(start) * time.Minute)
	e := group.Add(time.Duration(end) * time.Minute)

	startTime := time.Date(1970, 1, 1, s.Hour(), s.Minute(), 0, 0, time.UTC)
	endTime := time.Date(1970, 1, 1, e.Hour(), e.Minute(), 0, 0, time.UTC)
	currentTime := time.Date(1970, 1, 1, now.Hour(), now.Minute(), 0, 0, time.UTC)

	if currentTime.After(startTime) && currentTime.Before(endTime) {
		return true
	}
	return false
}
