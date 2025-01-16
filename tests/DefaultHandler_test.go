package tests

import (
	"fmt"
	"gopkg.in/telebot.v4"
	"reflect"
	"strings"
	"testing"
	"tgbot/internal"
	"tgbot/internal/config"
	"tgbot/internal/contextHandlers"
	"tgbot/internal/domain"
	"time"
)

type MockContext struct {
	userId      int64
	userMessage string
	unixTime    int64
	telebot.Context
}

func (m *MockContext) Message() *telebot.Message {
	return &telebot.Message{
		Sender: &telebot.User{
			ID: m.userId,
		},
		Text:     m.userMessage,
		Unixtime: m.unixTime,
	}
}

func (m *MockContext) setUserMessage(uid int64, msg string) {
	m.userId = uid
	m.userMessage = msg
}

func (m *MockContext) setUserMessageWithTime(uid int64, msg string, unix int64) {
	m.userId = uid
	m.userMessage = msg
	m.unixTime = unix
}

func TestDefaultHandler(t *testing.T) {
	t.Run("If user is not register", func(t *testing.T) {
		mockContext := MockContext{}
		mockStorage := NewMockStorage()
		mockService := NewMockService(mockStorage)

		messageHandler := internal.NewMessageHandler(mockService)

		mockContext.setUserMessage(12, "hello world!")
		mockStorage.setUserInfo(12, false)

		got := messageHandler.Process(&mockContext)
		want := contextHandlers.Response{
			Message:  config.HelloWorld,
			Keyboard: config.StartKeyboard,
		}

		if len(mockStorage.userMap) != 1 {
			t.Fatalf("want 1 register user, but got %d", len(mockStorage.userMap))
		}
		if len(got.Keyboard.ReplyKeyboard) != 2 {
			t.Errorf("want 2 rows, but got %d", len(got.Keyboard.ReplyKeyboard))
		}
		assertObjects(t, got, want)

		got = messageHandler.Process(&mockContext)
		want = contextHandlers.Response{
			Message:  config.Incorrect,
			Keyboard: config.StartKeyboard,
		}
		if len(got.Keyboard.ReplyKeyboard) != 2 {
			t.Errorf("want 2 rows, but got %d", len(got.Keyboard.ReplyKeyboard))
		}
		assertObjects(t, got, want)
	})
	t.Run("If user register", func(t *testing.T) {
		t.Run("Send any bullshit", func(t *testing.T) {
			mockContext := MockContext{}
			mockStorage := NewMockStorage()
			mockService := NewMockService(mockStorage)
			messageHandler := internal.NewMessageHandler(mockService)

			mockContext.setUserMessage(12, "aezakmi")
			mockStorage.setUserInfo(12, true)

			got := messageHandler.Process(&mockContext)
			want := contextHandlers.Response{
				Message:  config.Incorrect,
				Keyboard: config.StartKeyboard,
			}
			if len(got.Keyboard.ReplyKeyboard) != 2 {
				t.Errorf("want 2 rows, but got %d", len(got.Keyboard.ReplyKeyboard))
			}
			assertObjects(t, got, want)
		})
		t.Run("Send settings", func(t *testing.T) {
			mockContext := MockContext{}
			mockStorage := NewMockStorage()
			mockService := NewMockService(mockStorage)
			messageHandler := internal.NewMessageHandler(mockService)

			t.Run("Cookie set, notif off", func(t *testing.T) {
				want := contextHandlers.Response{
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
				want := contextHandlers.Response{
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
		t.Run("Send get missing kids", func(t *testing.T) {
			t.Run("Group exists", func(t *testing.T) {
				g := []domain.Group{
					{
						Id:     1,
						Name:   "Лекция 1",
						Lesson: "Lession 1",
						Time:   getDayByTime(28, 10, 0),
					},
					{
						Id:     2,
						Name:   "Лекция 2",
						Lesson: "Lession 2",
						Time:   getDayByTime(21, 12, 0),
					},
					{
						Id:     3,
						Name:   "Лекция 3",
						Lesson: "Lession 3",
						Time:   getDayByTime(27, 10, 0),
					},
				}

				mockContext := MockContext{}
				mockStorage := NewMockStorage()
				mockService := NewMockService(mockStorage)
				messageHandler := internal.NewMessageHandler(mockService)

				mockStorage.setUserInfo(12, true)
				mockStorage.SetGroups(12, g)
				mockContext.setUserMessageWithTime(12, "Получить отсутсвующих", getUnixByDay(28, 9, 40))

				got := messageHandler.Process(&mockContext)
				want := contextHandlers.Response{
					Message: fmt.Sprintf(
						"%s%s\n%s%s\n%s%d\n%s%d\n%s",
						config.GroupName,
						g[0].Name,
						config.Lection,
						g[0].Lesson,
						config.TotalKids,
						g[0].AllKids,
						config.MissingKids,
						len(g[0].MissingKids),
						strings.Join(g[0].MissingKids, "\n"),
					),
				}
				assertObjects(t, got, want)
			})
			t.Run("Group Non exists", func(t *testing.T) {
				g := []domain.Group{
					{
						Id:     1,
						Name:   "Лекция 1",
						Lesson: "Lession 1",
						Time:   getDayByTime(28, 10, 0),
					},
					{
						Id:     2,
						Name:   "Лекция 2",
						Lesson: "Lession 2",
						Time:   getDayByTime(21, 12, 0),
					},
					{
						Id:     3,
						Name:   "Лекция 3",
						Lesson: "Lession 3",
						Time:   getDayByTime(27, 10, 0),
					},
				}

				mockContext := MockContext{}
				mockStorage := NewMockStorage()
				mockService := NewMockService(mockStorage)
				messageHandler := internal.NewMessageHandler(mockService)

				mockStorage.setUserInfo(12, true)
				mockStorage.SetGroups(12, g)
				mockContext.setUserMessageWithTime(12, "Получить отсутсвующих", getUnixByDay(28, 22, 40))

				got := messageHandler.Process(&mockContext)
				want := contextHandlers.Response{
					Message: "group not found",
				}
				assertObjects(t, got, want)
			})
		})
		t.Run("Send my groups", func(t *testing.T) {
			g := []domain.Group{
				{
					Id:     1,
					Name:   "Гр 1",
					Lesson: "Lession 1",
					Time:   getDayByTime(28, 10, 0),
				},
				{
					Id:     3,
					Name:   "Гр 3",
					Lesson: "Lession 4",
					Time:   getDayByTime(27, 14, 0),
				},
				{
					Id:     2,
					Name:   "Гр 2",
					Lesson: "Lession 2",
					Time:   getDayByTime(21, 12, 0),
				},

				{
					Id:     4,
					Name:   "Гр 4",
					Lesson: "Lession 3",
					Time:   getDayByTime(27, 10, 0),
				},
			}

			mockContext := MockContext{}
			mockStorage := NewMockStorage()
			mockService := NewMockService(mockStorage)
			messageHandler := internal.NewMessageHandler(mockService)

			mockStorage.setUserInfo(12, true)
			mockStorage.SetGroups(12, g)
			mockContext.setUserMessage(12, "Мои группы")

			want := contextHandlers.Response{
				Message: fmt.Sprintf(
					"%s4\n\n%s\n\n%s",
					config.MyGroups,
					"1. Гр 4 🕐 сб 10:00\n2. Гр 3 🕐 сб 14:00",
					"1. Гр 1 🕐 вс 10:00\n2. Гр 2 🕐 вс 12:00",
				),
			}
			got := messageHandler.Process(&mockContext)
			assertObjects(t, got, want)
		})
	})

}

func assertObjects(t *testing.T, got contextHandlers.Response, want contextHandlers.Response) {
	t.Helper()
	if reflect.DeepEqual(got, want) != true {
		t.Errorf("Expected: [%s]\nActual: [%s]\n", got.Message, want.Message)
		t.Errorf("Length Expected: %d, Length Actual: %d\n", len(got.Message), len(want.Message))
	}
}

// getDayByTime 28, 21 вск ||  27, 20 сб
func getDayByTime(day, hour, min int) time.Time {
	return time.Date(2025, 9, day, hour, min, 0, 0, time.UTC)
}

// getUnixByDay поскольку телеграм переводит из Unix время в мое время на операционной системе, нужно предварительно при переводе вычесть разницу времен что бы перевод был корректен
func getUnixByDay(day, hour, min int) int64 {
	utcTime := time.Date(2025, 9, day, hour, min, 0, 0, time.UTC)
	_, offset := time.Now().Zone()

	unixTime := utcTime.Add(-time.Duration(offset) * time.Second).Unix()
	return unixTime
}
