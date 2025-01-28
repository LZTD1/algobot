package mocks

import (
	"tgbot/internal/domain"
	"time"
)

type MockDomain struct {
	MockGroups []domain.Group
}

func (m *MockDomain) User(uid int64) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockDomain) Cookie(uid int64) (string, error) {
	return "", nil
}

func (m *MockDomain) SetCookie(uid int64, cookie string) {
	//TODO implement me
	panic("implement me")
}

func (m *MockDomain) SetUserAgent(uid int64, agent string) {
	//TODO implement me
	panic("implement me")
}

func (m *MockDomain) Groups(uid int64) ([]domain.Group, error) {
	return []domain.Group{
		{
			Id:          1,
			Name:        "test1",
			Lesson:      "",
			Time:        time.Date(2024, 10, 6, 14, 55, 55, 3, time.UTC),
			AllKids:     0,
			MissingKids: nil,
		},
		{
			Id:          1,
			Name:        "test2",
			Lesson:      "",
			Time:        time.Date(2024, 11, 3, 14, 55, 55, 3, time.UTC),
			AllKids:     0,
			MissingKids: nil,
		},
	}, nil
}

func (m *MockDomain) SetGroups(uid int64, groups []domain.Group) {
	m.MockGroups = groups
}

func (m *MockDomain) Notification(uid int64) bool {
	//TODO implement me
	panic("implement me")
}
func (m *MockDomain) SetNotification(uid int64, value bool) {

}
func (m *MockDomain) RegisterUser(uid int64) {
	//TODO implement me
	panic("implement me")
}
