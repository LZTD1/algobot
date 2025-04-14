package models

type User struct {
	ID               int
	Uid              int64
	Cookie           string
	LastNotification string
	Notification     int
}
