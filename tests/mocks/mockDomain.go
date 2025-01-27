package mocks

import "tgbot/internal/domain"

type MockDomain struct {
}

func (m MockDomain) User(uid int64) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockDomain) Cookie(uid int64) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockDomain) SetCookie(uid int64, cookie string) {
	//TODO implement me
	panic("implement me")
}

func (m MockDomain) SetUserAgent(uid int64, agent string) {
	//TODO implement me
	panic("implement me")
}

func (m MockDomain) Groups(uid int64) ([]domain.Group, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockDomain) SetGroups(uid int64, groups []domain.Group) {
	//TODO implement me
	panic("implement me")
}

func (m MockDomain) Notification(uid int64) bool {
	//TODO implement me
	panic("implement me")
}

func (m MockDomain) RegisterUser(uid int64) {
	//TODO implement me
	panic("implement me")
}
