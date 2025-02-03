package mocks

import (
	"errors"
	"fmt"
	"strconv"
	"tgbot/internal/models"
	"time"
)

type MockService struct {
	m                map[int64]bool
	cookie           string
	StubNotification bool
	gr               *models.Group
	grs              []models.Group
	SettedCookie     []string
	Actual           models.ActualInformation
	AllNames         models.AllKids
	grErr            error
	Calls            []string
}

func (n *MockService) ActualInformation(uid int64, t time.Time, groupId int) (models.ActualInformation, error) {
	return n.Actual, nil
}

func (n *MockService) AllKidsNames(uid int64, groupId int) (models.AllKids, error) {
	return n.AllNames, nil
}

func NewMockService(m map[int64]bool) *MockService {
	return &MockService{m: m}
}

func (n *MockService) CloseLesson(uid int64, lessonId int, groupId int) error {
	n.Calls = append(n.Calls, fmt.Sprintf("CloseLesson(%d, %d, %d)", uid, lessonId, groupId))
	return nil
}

func (n *MockService) OpenLesson(uid int64, lessonId int, groupId int) error {
	n.Calls = append(n.Calls, fmt.Sprintf("OpenLesson(%d, %d, %d)", uid, lessonId, groupId))
	return nil
}

func (n *MockService) GetAllCredentials(uid int64, groupId int) (map[string]string, error) {
	n.Calls = append(n.Calls, fmt.Sprintf("GetAllCredentials(%d, %d)", uid, groupId))
	return map[string]string{
		"Ваня": "van:12",
	}, nil
}

func (n *MockService) SetMockCookie(s string) {
	n.cookie = s
}
func (n *MockService) RefreshGroups(uid int64) error {
	return nil
}
func (n *MockService) SetCurrentGroup(group *models.Group) {
	n.gr = group
}
func (n *MockService) SetGroups(groups []models.Group) {
	n.grs = groups
}

func (n *MockService) CurrentGroup(uid int64, t time.Time) (models.Group, error) {
	if n.gr == nil {
		return models.Group{}, errors.New("no gr")
	}
	return *n.gr, nil
}

func (n *MockService) Groups(uid int64) ([]models.Group, error) {
	if n.grErr != nil {
		return nil, n.grErr
	}
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

func (n *MockService) SetGroupsErr(err error) {
	n.grErr = err
}
