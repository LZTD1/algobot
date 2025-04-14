package models

import "time"

type GroupView struct {
	GroupID        int
	GroupTitle     string
	GroupContent   string
	NextLessonTime string
	LessonsTotal   int
	LessonsPassed  int
	ActiveKids     []GroupKid
	NotActiveKids  []GroupKid
}

type GroupKid struct {
	ID        int
	FullName  string
	LastGroup KidGroup
}
type KidGroup struct {
	ID        int
	StartTime time.Time
	EndTime   time.Time
}
