package service

import (
	"tgbot/internal/domain"
	"time"
)

type Service interface {
	CurrentGroup(uid int64, t time.Time) (domain.Group, error)
	Groups(uid int64) ([]domain.Group, error)
	Cookie(uid int64) (string, error)
	SetCookie(uid int64, cookie string) error
	Notification(uid int64) (bool, error)
	SetNotification(uid int64, notification bool) error
	IsUserRegistered(uid int64) (bool, error)
	RegisterUser(uid int64) error
	RefreshGroups(uid int64) error
}
