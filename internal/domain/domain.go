package domain

import "time"

type Group struct {
	Id          int
	Name        string
	Lesson      string
	Time        time.Time
	AllKids     int
	MissingKids []string
}

type User struct {
	cookie    string
	userAgent string
	groups    []Group
}

type Domain interface {
	User(uid int64) (User, error)
	Cookie(uid int64) (string, error)
	SetCookie(uid int64, cookie string)
	SetUserAgent(uid int64, agent string)
	Groups(uid int64) ([]Group, error)
	SetGroups(uid int64, groups []Group)
	Notification(uid int64) bool
	RegisterUser(uid int64)
}
