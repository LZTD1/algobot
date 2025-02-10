package models

import (
	"tgbot/internal/clients"
	"tgbot/internal/domain"
	"time"
)

type MissingKid struct {
	Id    int
	Count int
}

type ActualInformation struct {
	LessonTitle string
	LessonId    int
	MissingKids []MissingKid
}

type KidData struct {
	FullName string
	Login    string
	Password string
}
type ScheduleData struct {
	UID    int64
	Cookie string
}

type Message struct {
	Id      string
	From    string
	Theme   string
	Type    string
	Link    string
	Content string
}

type AllKids map[int]KidData

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

type FullGroupInfo struct {
	GroupID        int
	GroupTitle     string
	GroupContent   string
	NextLessonTime string
	LessonsTotal   int
	LessonsPassed  int
	ActiveKids     []clients.Student
	NotActiveKids  []clients.Student
}

type ExtraInfo string

var NotAccessible ExtraInfo = "not_accessible"

type FullKidInfo struct {
	Extra ExtraInfo
	Kid   clients.Student
}
