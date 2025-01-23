package mocks

import (
	"errors"
	"tgbot/internal/domain"
	"time"
)

type MockService struct {
	m      map[int64]bool
	cookie string
	gr     *domain.Group
	grs    []domain.Group
}

func NewMockService(m map[int64]bool) *MockService {
	return &MockService{m: m}
}

func (n *MockService) SetCookie(s string) {
	n.cookie = s
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

func (n *MockService) MissingKids(t time.Time, g int) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (n *MockService) Cookie(uid int64) (string, error) {
	return n.cookie, nil
}

func (n *MockService) Notification(uid int64) bool {
	return false
}

func (n *MockService) IsUserRegistered(uid int64) bool {
	v, ok := n.m[uid]
	if ok != true {
		return false
	}
	return v
}

func (n *MockService) RegisterUser(uid int64) {
	n.m[uid] = true
}
