package storage

type User struct {
	cookie    string
	userAgent string
	groups    []string
}

type Storage interface {
	User(uid int64) (User, error)
	Cookie(uid int64) (string, error)
	SetCookie(uid int64, cookie string)
	SetUserAgent(uid int64, agent string)
	Groups(uid int64) []string
	SetGroups(uid int64, groups []string)
	Notification(uid int64) bool
	RegisterUser(uid int64)
}
