package models

import "time"

type Group struct {
	GroupID    int
	Title      string
	TimeLesson time.Time
}
