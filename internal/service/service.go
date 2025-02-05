package service

import (
	"tgbot/internal/models"
	"time"
)

type Service interface {
	CurrentGroup(uid int64, t time.Time) (models.Group, error)
	Groups(uid int64) ([]models.Group, error)
	Cookie(uid int64) (string, error)
	SetCookie(uid int64, cookie string) error
	Notification(uid int64) (bool, error)
	SetNotification(uid int64, notification bool) error
	IsUserRegistered(uid int64) (bool, error)
	RegisterUser(uid int64) error
	RefreshGroups(uid int64) error
	ActualInformation(uid int64, t time.Time, groupId int) (models.ActualInformation, error)
	AllKidsNames(uid int64, groupId int) (models.AllKids, error)
	OpenLesson(uid int64, lessonId int, groupId int) error
	CloseLesson(uid int64, lessonId int, groupId int) error
	AllCredentials(uid int64, groupId int) (map[string]string, error)
	UsersByNotif(status bool) ([]models.ScheduleData, error)
	NewMessageByUID(uid int64) ([]models.Message, error)
	FullGroupInfo(uid int64, groupId int) (models.FullGroupInfo, error)
	FullKidInfo(uid int64, kidID int) (models.FullKidInfo, error)
}
