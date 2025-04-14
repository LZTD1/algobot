package sort

import (
	"algobot/internal/domain/models"
	"sort"
)

func GroupsByDate(groups []models.Group) {
	sort.Slice(groups, func(i, j int) bool {
		dayI, dayJ := groups[i].TimeLesson.Weekday(), groups[j].TimeLesson.Weekday()
		if dayI != dayJ {
			return dayI > dayJ
		}
		return groups[i].TimeLesson.Hour() < groups[j].TimeLesson.Hour() ||
			(groups[i].TimeLesson.Hour() == groups[j].TimeLesson.Hour() && groups[i].TimeLesson.Minute() < groups[j].TimeLesson.Minute())
	})
}
