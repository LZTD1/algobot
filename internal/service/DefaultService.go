package service

import (
	"tgbot/internal/clients"
	"tgbot/internal/domain"
	"tgbot/internal/helpers"
	"time"
)

type DefaultService struct {
	domain    domain.Domain
	webClient clients.WebClient
}

func NewDefaultService(domain domain.Domain, webClient clients.WebClient) *DefaultService {
	return &DefaultService{domain: domain, webClient: webClient}
}

func (d DefaultService) CurrentGroup(uid int64, t time.Time) (domain.Group, error) {
	allGroups, err := d.domain.Groups(uid)
	if err != nil {
		return domain.Group{}, err
	}

	group, err := helpers.GetCurrentGroup(t, allGroups)
	if err != nil {
		return domain.Group{}, err
	}

	kids, err := d.MissingKids(uid, t, group.Id)
	if err != nil {
		return domain.Group{}, err
	}

	group.MissingKids = kids

	return group, nil
}

func (d DefaultService) Groups(uid int64) ([]domain.Group, error) {
	//TODO implement me
	panic("implement me")
}

func (d DefaultService) MissingKids(uid int64, t time.Time, g int) ([]string, error) {
	cookie, err := d.domain.Cookie(uid)
	if err != nil {
		return nil, err
	}

	names, err := d.webClient.GetKidsNamesByGroup(cookie, string(g))
	if err != nil {
		return nil, err
	}

	stats, err := d.webClient.GetKidsStatsByGroup(cookie, string(g))
	if err != nil {
		return nil, err
	}

	for i, datum := range stats.Data {
		_ = i
		_ = datum
	}

	_ = names
	panic("implement me")
}

func (d DefaultService) Cookie(uid int64) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (d DefaultService) SetCookie(uid int64, cookie string) {
	//TODO implement me
	panic("implement me")
}

func (d DefaultService) Notification(uid int64) bool {
	//TODO implement me
	panic("implement me")
}

func (d DefaultService) SetNotification(uid int64, notification bool) {
	//TODO implement me
	panic("implement me")
}

func (d DefaultService) IsUserRegistered(uid int64) bool {
	//TODO implement me
	panic("implement me")
}

func (d DefaultService) RegisterUser(uid int64) {
	//TODO implement me
	panic("implement me")
}

func (d DefaultService) RefreshGroups(uid int64) error {
	//TODO implement me
	panic("implement me")
}
