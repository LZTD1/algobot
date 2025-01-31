package helpers

import (
	"errors"
	"sort"
	"tgbot/internal/domain"
	appError "tgbot/internal/error"
	"time"
)

// GetSortedGroups , получение списка групп отсортированных по дням
// И внутри дней по часам
func GetSortedGroups(groups []domain.Group) []domain.Group {
	sort.Slice(groups, func(i, j int) bool {
		dayI, dayJ := groups[i].Time.Weekday(), groups[j].Time.Weekday()
		if dayI != dayJ {
			return dayI > dayJ
		}
		return groups[i].Time.Hour() < groups[j].Time.Hour() ||
			(groups[i].Time.Hour() == groups[j].Time.Hour() && groups[i].Time.Minute() < groups[j].Time.Minute())
	})
	return groups
}

// GetCurrentGroup , логика выдачи групп:
// За 30 минут до начала и во время группы выдывать конкретную группу
func GetCurrentGroup(t time.Time, g []domain.Group) (domain.Group, error) {
	current, err := GetGroupsByDay(t, g)
	if err != nil {
		return domain.Group{}, err
	}
	for _, group := range current {
		if inDiapazon(-30, 90, t, group.Time) {
			return group, nil
		}
	}

	return domain.Group{}, appError.ErrHasNone
}

// GetGroupsByDay получение групп по текущему дню
func GetGroupsByDay(t time.Time, g []domain.Group) ([]domain.Group, error) {
	var filtered []domain.Group

	for _, group := range g {
		if t.Weekday() == group.Time.Weekday() {
			filtered = append(filtered, group)
		}
	}
	if len(filtered) == 0 {
		return nil, errors.New("no groups found")
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
