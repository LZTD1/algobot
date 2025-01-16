package tests

import (
	"errors"
	"tgbot/internal/domain"
)

type UserInfo struct {
	IsRegistered bool
	Cookie       string
	Notification bool
	Groups       []domain.Group
}

type MockStorage struct {
	userMap map[int64]UserInfo
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		userMap: make(map[int64]UserInfo),
	}
}

func (s MockStorage) User(uid int64) (domain.User, error) {
	val, ok := s.userMap[uid]
	if ok {
		if !val.IsRegistered {
			return domain.User{}, errors.New("not found")
		}
		return domain.User{}, nil
	}
	return domain.User{}, errors.New("not found")
}

func (s MockStorage) Cookie(uid int64) (string, error) {
	val, ok := s.userMap[uid]
	if !ok || val.Cookie == "" {
		return "", errors.New("not found")
	}
	return val.Cookie, nil
}

func (s MockStorage) SetCookie(uid int64, cookie string) {
	val, ok := s.userMap[uid]
	if ok {
		val.Cookie = cookie
		s.userMap[uid] = val
	}
}
func (s MockStorage) RegisterUser(uid int64) {
	s.userMap[uid] = UserInfo{
		IsRegistered: true,
		Cookie:       "",
		Notification: false,
	}
}
func (s MockStorage) SetUserAgent(uid int64, agent string) {
	panic("implement me SetUserAgent")
}

func (s MockStorage) Groups(uid int64) ([]domain.Group, error) {
	if val, ok := s.userMap[uid]; ok {
		g := val.Groups
		if len(g) == 0 {
			return nil, errors.New("groups dont sets")
		}
		return g, nil
	}
	return nil, errors.New("user not found")
}

func (s MockStorage) SetGroups(uid int64, groups []domain.Group) {
	val := s.userMap[uid]
	val.Groups = groups
	s.userMap[uid] = val
}
func (s MockStorage) Notification(uid int64) bool {
	return s.userMap[uid].Notification
}
func (s MockStorage) SetNotification(uid int64, state bool) {
	val, ok := s.userMap[uid]
	if ok {
		val.Notification = state
		s.userMap[uid] = val
	}
}
func (s MockStorage) setUserInfo(uid int64, isRegister bool) {
	s.userMap[uid] = UserInfo{
		IsRegistered: isRegister,
		Cookie:       "",
		Notification: false,
	}
}
