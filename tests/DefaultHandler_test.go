package tests

import (
	"errors"
	"fmt"
	"gopkg.in/telebot.v4"
	"reflect"
	"testing"
	"tgbot"
	"tgbot/config"
	"tgbot/handlers"
	"tgbot/storage"
)

type MockContext struct {
	userId      int64
	userMessage string
	telebot.Context
}

func (m *MockContext) Message() *telebot.Message {
	return &telebot.Message{
		Sender: &telebot.User{
			ID: m.userId,
		},
		Text: m.userMessage,
	}
}

func (m *MockContext) setUserMessage(uid int64, msg string) {
	m.userId = uid
	m.userMessage = msg
}

func TestDefaultHandler(t *testing.T) {
	t.Run("If user is not register", func(t *testing.T) {
		mockContext := MockContext{}
		mockStorage := NewMockStorage()
		messageHandler := tgbot.NewMessageHandler(mockStorage)

		mockContext.setUserMessage(12, "hello world!")
		mockStorage.setUserInfo(12, false)

		got := messageHandler.Process(&mockContext)
		want := handlers.Response{
			Message: config.HelloWorld,
		}

		if len(mockStorage.userMap) != 1 {
			t.Fatalf("want 1 register user, but got %d", len(mockStorage.userMap))
		}
		assertObjects(t, got, want)
	})
	t.Run("If user send any bullshit", func(t *testing.T) {
		mockContext := MockContext{}
		mockStorage := NewMockStorage()
		messageHandler := tgbot.NewMessageHandler(mockStorage)

		mockContext.setUserMessage(12, "aezakmi")
		mockStorage.setUserInfo(12, true)

		got := messageHandler.Process(&mockContext)
		want := handlers.Response{
			Message: config.Incorrect,
		}

		assertObjects(t, got, want)
	})
	t.Run("If user register, and send setting", func(t *testing.T) {
		mockContext := MockContext{}
		mockStorage := NewMockStorage()
		messageHandler := tgbot.NewMessageHandler(mockStorage)

		t.Run("Cookie set, notif off", func(t *testing.T) {
			want := handlers.Response{
				Message: fmt.Sprintf(
					"%s\n\n%s%s\n%s%s",
					config.Settings,
					config.Cookie,
					config.SetParam,
					config.ChatNotifications,
					config.NotSetParam,
				),
			}

			mockStorage.setUserInfo(12, true)
			mockStorage.SetCookie(12, "cookie")

			mockContext.setUserMessage(12, "Настройки")

			got := messageHandler.Process(&mockContext)
			assertObjects(t, got, want)
		})
		t.Run("Cookie unset, notif on", func(t *testing.T) {
			want := handlers.Response{
				Message: fmt.Sprintf(
					"%s\n\n%s%s\n%s%s",
					config.Settings,
					config.Cookie,
					config.NotSetParam,
					config.ChatNotifications,
					config.SetParam,
				),
			}

			mockStorage.setUserInfo(12, true)
			mockStorage.SetNotification(12, true)

			mockContext.setUserMessage(12, "Настройки")

			got := messageHandler.Process(&mockContext)
			assertObjects(t, got, want)
		})
	})

}

func assertObjects(t *testing.T, got handlers.Response, want handlers.Response) {
	if reflect.DeepEqual(got, want) != true {
		t.Fatalf("Want %+v, got %v", want, got)
	}
}

type UserInfo struct {
	IsRegistered bool
	Cookie       string
	Notification bool
}

type MockStorage struct {
	userMap map[int64]UserInfo
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		userMap: make(map[int64]UserInfo),
	}
}

func (s MockStorage) User(uid int64) (storage.User, error) {
	val, ok := s.userMap[uid]
	if ok {
		if !val.IsRegistered {
			return storage.User{}, errors.New("not found")
		}
		return storage.User{}, nil
	}
	return storage.User{}, errors.New("not found")
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
	panic("implement me")
}

func (s MockStorage) Groups(uid int64) []string {
	panic("implement me")
}

func (s MockStorage) SetGroups(uid int64, groups []string) {
	panic("implement me")
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
