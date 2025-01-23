package service

import (
	"tgbot/internal/domain"
	"time"
)

type Service interface {
	CurrentGroup(uid int64, t time.Time) (domain.Group, error)
	Groups(uid int64) ([]domain.Group, error)
	MissingKids(t time.Time, g int) ([]string, error)
	Cookie(uid int64) (string, error)
	Notification(uid int64) bool
	IsUserRegistered(uid int64) bool
	RegisterUser(uid int64)
}
