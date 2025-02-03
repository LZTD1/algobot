package domain

import (
	"time"
)

type Group struct {
	GroupID    int
	Title      string
	TimeLesson time.Time
}

type User struct {
	Cookie        string
	UserAgent     string
	Notifications bool
	Groups        []Group
}

type Domain interface {
	User(uid int64) (User, error)
	Cookie(uid int64) (string, error)
	SetCookie(uid int64, cookie string) error
	SetUserAgent(uid int64, agent string) error
	Groups(uid int64) ([]Group, error)
	SetGroups(uid int64, groups []Group) error
	Notification(uid int64) (bool, error)
	SetNotification(uid int64, value bool) error
	RegisterUser(uid int64) error
}
