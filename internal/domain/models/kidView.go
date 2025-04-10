package models

import "time"

type ExtraInfo string

var NotAccessible ExtraInfo = "not_accessible"

type KidView struct {
	Extra ExtraInfo
	Kid   Kid
}
type Kid struct {
	FullName   string         `json:"fullName"`
	ParentName string         `json:"parentName"`
	Email      string         `json:"email"`
	Phone      string         `json:"phone"`
	Age        int            `json:"age"`
	BirthDate  time.Time      `json:"birthDate"`
	Username   string         `json:"username"`
	Password   string         `json:"password"`
	Groups     []KidViewGroup `json:"groups"`
}

type KidViewGroup struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Status    int       `json:"status"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}
