package tests

import (
	"tgbot/internal/domain"
	"tgbot/internal/helpers"
	"time"
)

type MockService struct {
	s domain.Domain
}

func NewMockService(s domain.Domain) *MockService {
	return &MockService{s}
}

func (m MockService) IsUserRegistered(uid int64) bool {
	_, e := m.s.User(uid)
	if e != nil {
		return false
	}
	return true
}

func (m MockService) RegisterUser(uid int64) {
	m.s.RegisterUser(uid)
}

func (m MockService) Groups(uid int64) ([]domain.Group, error) {
	all, err := m.s.Groups(uid)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (m MockService) CurrentGroup(uid int64, t time.Time) (domain.Group, error) {
	all, err := m.s.Groups(uid)
	if err != nil {
		return domain.Group{}, err
	}
	g, e := helpers.GetCurrentGroup(t, all)
	if e != nil {
		return domain.Group{}, e
	}
	return g, nil
}

func (m MockService) MissingKids(t time.Time, g int) ([]string, error) {
	//TODO implement me
	panic("implement me MissingKids")
}

func (m MockService) Cookie(uid int64) (string, error) {
	return m.s.Cookie(uid)
}

func (m MockService) Notification(uid int64) bool {
	return m.s.Notification(uid)
}
