package mocks

import (
	"tgbot/internal/clients"
)

type MockWebClient struct {
}

func (m MockWebClient) GetKidsNamesByGroup(cookie, group string) (*clients.GroupResponse, *clients.ClientError) {
	//TODO implement me
	panic("implement me")
}

func (m MockWebClient) GetKidsStatsByGroup(cookie, group string) (*clients.KidsStats, *clients.ClientError) {
	//TODO implement me
	panic("implement me")
}

func (m MockWebClient) OpenLession(cookie, group, lession string) *clients.ClientError {
	//TODO implement me
	panic("implement me")
}

func (m MockWebClient) CloseLession(cookie, group, lession string) *clients.ClientError {
	//TODO implement me
	panic("implement me")
}

func (m MockWebClient) GetKidsMessages(cookie string) (*clients.KidsMessages, *clients.ClientError) {
	//TODO implement me
	panic("implement me")
}

func (m MockWebClient) GetAllGroupsByUser(cookie string) ([]clients.AllGroupsUser, *clients.ClientError) {
	//TODO implement me
	panic("implement me")
}
