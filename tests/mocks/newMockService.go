package mocks

import (
	"errors"
	"strconv"
	"tgbot/internal/domain"
	"time"
)

type MockService struct {
	m                map[int64]bool
	cookie           string
	StubNotification bool
	gr               *domain.Group
	grs              []domain.Group
	SettedCookie     []string
}

func NewMockService(m map[int64]bool) *MockService {
	return &MockService{m: m}
}

func (n *MockService) SetMockCookie(s string) {
	n.cookie = s
}
func (n *MockService) RefreshGroups(uid int64) error {
	return nil
}
func (n *MockService) SetCurrentGroup(group *domain.Group) {
	n.gr = group
}
func (n *MockService) SetGroups(groups []domain.Group) {
	n.grs = groups
}

func (n *MockService) CurrentGroup(uid int64, t time.Time) (domain.Group, error) {
	if n.gr == nil {
		return domain.Group{}, errors.New("no gr")
	}
	return *n.gr, nil
}

func (n *MockService) Groups(uid int64) ([]domain.Group, error) {
	if n.grs == nil {
		return nil, errors.New("no gr")
	}
	return n.grs, nil
}

func (n *MockService) MissingKids(uid int64, t time.Time, g int) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (n *MockService) Cookie(uid int64) (string, error) {
	return n.cookie, nil
}

func (n *MockService) SetCookie(uid int64, cookie string) error {
	n.SettedCookie = []string{strconv.FormatInt(uid, 10), cookie}
	return nil
}
func (n *MockService) SetNotification(uid int64, notification bool) error {
	n.StubNotification = notification
	return nil
}

func (n *MockService) Notification(uid int64) (bool, error) {
	return n.StubNotification, nil
}

func (n *MockService) IsUserRegistered(uid int64) (bool, error) {
	v, ok := n.m[uid]
	if ok != true {
		return false, nil
	}
	return v, nil
}

func (n *MockService) RegisterUser(uid int64) error {
	n.m[uid] = true
	return nil
}
