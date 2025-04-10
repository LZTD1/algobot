package models

type CurrentGroup struct {
	GroupID     int
	Title       string
	Lesson      string
	LessonID    int
	Kids        []string
	MissingKids []MissingKid
}
type MissingKid struct {
	Fullname string
	Count    int
}
