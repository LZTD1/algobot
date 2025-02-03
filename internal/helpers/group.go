package helpers

import (
	"sort"
	appError "tgbot/internal/error"
	"tgbot/internal/models"
	"time"
)

// GetSortedGroups получение списка групп отсортированных по дням
// И внутри дней по часам
func GetSortedGroups(groups []models.Group) []models.Group {
	sort.Slice(groups, func(i, j int) bool {
		dayI, dayJ := groups[i].TimeLesson.Weekday(), groups[j].TimeLesson.Weekday()
		if dayI != dayJ {
			return dayI > dayJ
		}
		return groups[i].TimeLesson.Hour() < groups[j].TimeLesson.Hour() ||
			(groups[i].TimeLesson.Hour() == groups[j].TimeLesson.Hour() && groups[i].TimeLesson.Minute() < groups[j].TimeLesson.Minute())
	})
	return groups
}

// GetCurrentGroup , логика выдачи групп:
// За 30 минут до начала и во время группы выдывать конкретную группу
func GetCurrentGroup(t time.Time, g []models.Group) (models.Group, error) {
	current, err := GetGroupsByDay(t, g)
	if err != nil {
		return models.Group{}, err
	}
	for _, group := range current {
		if inDiapazon(-30, 90, t, group.TimeLesson) {
			return group, nil
		}
	}

	return models.Group{}, appError.ErrHasNone
}

// GetGroupsByDay получение групп по текущему дню
func GetGroupsByDay(t time.Time, g []models.Group) ([]models.Group, error) {
	var filtered []models.Group

	for _, group := range g {
		if t.Weekday() == group.TimeLesson.Weekday() {
			filtered = append(filtered, group)
		}
	}
	if len(filtered) == 0 {
		return nil, appError.ErrHasNone
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
