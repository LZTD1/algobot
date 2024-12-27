package tests

import (
	"gopkg.in/telebot.v4"
	"reflect"
	"testing"
	"tgbot"
	"tgbot/handlers"
	"tgbot/storage"
)

type MockContext struct {
	data string
	telebot.Context
}

func (m *MockContext) setData(data string) {
	m.data = data
}

type MockStorage struct {
}

func (s MockStorage) User(uid int64) storage.User {
	panic("implement me")
}

func (s MockStorage) SetCookie(uid int64, cookie string) {
	panic("implement me")
}

func (s MockStorage) SetUserAgent(uid int64, agent string) {
	panic("implement me")
}

func (s MockStorage) Groups(uid int64) []string {
	panic("implement me")
}

func (s MockStorage) SetGroups(uid int64, groups []string) {
	panic("implement me")
}

func TestDefaultHandler(t *testing.T) {
	mc := MockContext{}
	mock := MockStorage{}
	dh := handlers.NewMessageHandler(mock)

	t.Run("If user is not register", func(t *testing.T) {
		got := dh.Process(mc)
		want := handlers.HandlerResponse{
			Text: tgbot.HelloWorld,
		}

		if reflect.DeepEqual(got, want) != true {
			t.Fatalf("Want %+v, got %+v", want, got)
		}
	})
	t.Run("If user register", func(t *testing.T) {
		got := dh.Process(mc)
		want := handlers.HandlerResponse{
			Text: tgbot.HelloWorld,
		}

		if reflect.DeepEqual(got, want) != true {
			t.Fatalf("Want %+v, got %+v", want, got)
		}
	})
}
