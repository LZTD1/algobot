package mocks

import (
	"tgbot/internal/domain"
	"time"
)

type MockDomain struct {
	MockGroups []domain.Group
	errCookie  error
	errNotif   error
	DataNotif  string
}

func (m *MockDomain) LastNotificationDate(uid int64) (string, error) {
	return "14 дек. `24, 18:36", nil
}

func (m *MockDomain) SetLastNotificationDate(uid int64, data string) error {
	m.DataNotif = data
	return nil
}

func (m *MockDomain) GetUsersByNotification(notifications int) ([]domain.User, error) {
	return []domain.User{
		{
			UID:           1,
			Cookie:        "2",
			UserAgent:     "2",
			Notifications: false,
			Groups:        nil,
		},
	}, nil
}

func (m *MockDomain) User(uid int64) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockDomain) Cookie(uid int64) (string, error) {
	if m.errCookie == nil {
		return "cookie", nil
	}
	return "", m.errCookie
}

func (m *MockDomain) SetCookie(uid int64, cookie string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockDomain) SetUserAgent(uid int64, agent string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockDomain) Groups(uid int64) ([]domain.Group, error) {
	return []domain.Group{
		{
			GroupID:    1,
			Title:      "test1",
			TimeLesson: time.Date(2024, 10, 6, 14, 55, 55, 3, time.UTC),
		},
		{
			GroupID:    1,
			Title:      "test2",
			TimeLesson: time.Date(2024, 11, 3, 14, 55, 55, 3, time.UTC),
		},
	}, nil
}

func (m *MockDomain) SetGroups(uid int64, groups []domain.Group) error {
	m.MockGroups = groups
	return nil
}

func (m *MockDomain) Notification(uid int64) (bool, error) {
	if m.errNotif == nil {
		return true, nil
	}
	return false, m.errNotif
}
func (m *MockDomain) SetNotification(uid int64, value bool) error {
	return nil
}
func (m *MockDomain) RegisterUser(uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockDomain) SetErrorCookie(err error) {
	m.errCookie = err
}

func (m *MockDomain) SetErrorNotif(err error) {
	m.errNotif = err
}
