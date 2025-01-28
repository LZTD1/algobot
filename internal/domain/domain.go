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
	cookie        string
	userAgent     string
	notifications bool
	groups        []Group
}

func (u User) Cookie() string {
	return u.cookie
}

func (u User) UserAgent() string {
	return u.userAgent
}

func (u User) Notifications() bool {
	return u.notifications
}

func (u User) Groups() []Group {
	return u.groups
}

type Domain interface {
	User(uid int64) (User, error)
	Cookie(uid int64) (string, error)
	SetCookie(uid int64, cookie string)
	SetUserAgent(uid int64, agent string)
	Groups(uid int64) ([]Group, error)
	SetGroups(uid int64, groups []Group)
	Notification(uid int64) bool
	SetNotification(uid int64, value bool)
	RegisterUser(uid int64)
}
