package models

import (
	"tgbot/internal/domain"
	"time"
)

type ActualInformation struct {
	LessonTitle string
	LessonId    int
	MissingKids []int
}

type AllKids map[int]string

type Group struct {
	GroupID    int
	Title      string
	TimeLesson time.Time
}

func GroupMap(domains []domain.Group) []Group {
	gr := make([]Group, len(domains))
	for i, group := range domains {
		gr[i] = mapGroup(group)
	}
	return gr
}

func mapGroup(group domain.Group) Group {
	return Group{
		GroupID:    group.GroupID,
		Title:      group.Title,
		TimeLesson: group.TimeLesson,
	}
}
