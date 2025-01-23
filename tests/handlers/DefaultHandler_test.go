package handlers

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
	"tgbot/tests/mocks"
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
		ms := mocks.NewMockService(make(map[int64]bool))

		mockContext := MockContext{}

		messageHandler := internal.NewMessageHandler(ms)

		mockContext.setUserMessage(12, "hello world!")

		got := messageHandler.Process(&mockContext)
		want := contextHandlers.Response{
			Message:  config.HelloWorld,
			Keyboard: config.StartKeyboard,
		}

		assertMessages(t, got, want)
		assertKeyboards(t, got, want)

		got = messageHandler.Process(&mockContext)
		want = contextHandlers.Response{
			Message:  config.Incorrect,
			Keyboard: config.StartKeyboard,
		}
		assertMessages(t, got, want)
		assertKeyboards(t, got, want)
	})
	t.Run("If user register", func(t *testing.T) {
		t.Run("Send any bullshit", func(t *testing.T) {
			mockContext := MockContext{}

			ms := mocks.NewMockService(map[int64]bool{
				12: true,
			})

			messageHandler := internal.NewMessageHandler(ms)

			mockContext.setUserMessage(12, "aezakmi")

			got := messageHandler.Process(&mockContext)
			want := contextHandlers.Response{
				Message:  config.Incorrect,
				Keyboard: config.StartKeyboard,
			}

			assertMessages(t, got, want)
			assertKeyboards(t, got, want)
		})
		t.Run("Send settings", func(t *testing.T) {
			mockContext := MockContext{}

			ms := mocks.NewMockService(map[int64]bool{
				12: true,
			})

			messageHandler := internal.NewMessageHandler(ms)

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
					Keyboard: config.SettingsKeyboard,
				}

				ms.SetCookie("Cookie")
				mockContext.setUserMessage(12, "Настройки")

				got := messageHandler.Process(&mockContext)
				assertMessages(t, got, want)
				assertKeyboards(t, got, want)
			})
			t.Run("Cookie unset, notif off", func(t *testing.T) {
				want := contextHandlers.Response{
					Message: fmt.Sprintf(
						"%s\n\n%s%s\n%s%s",
						config.Settings,
						config.Cookie,
						config.NotSetParam,
						config.ChatNotifications,
						config.NotSetParam,
					),
					Keyboard: config.SettingsKeyboard,
				}

				ms.SetCookie("")
				mockContext.setUserMessage(12, "Настройки")

				got := messageHandler.Process(&mockContext)
				assertMessages(t, got, want)
				assertKeyboards(t, got, want)
			})
		})
		t.Run("Send get missing kids", func(t *testing.T) {
			t.Run("Group exists", func(t *testing.T) {
				gr := domain.Group{
					Id:          1,
					Name:        "Лекция 1",
					Lesson:      "Lession 1",
					Time:        getDayByTime(28, 10, 0),
					AllKids:     10,
					MissingKids: []string{"Name1", "Name2"},
				}
				ms := mocks.NewMockService(map[int64]bool{
					12: true,
				})
				ms.SetCurrentGroup(&gr)

				mockContext := MockContext{}

				messageHandler := internal.NewMessageHandler(ms)

				mockContext.setUserMessageWithTime(12, "Получить отсутсвующих", getUnixByDay(28, 9, 40))

				got := messageHandler.Process(&mockContext)
				want := contextHandlers.Response{
					Message: fmt.Sprintf(
						"%s%s\n%s%s\n%s%d\n%s%d\n%s",
						config.GroupName,
						gr.Name,
						config.Lection,
						gr.Lesson,
						config.TotalKids,
						gr.AllKids,
						config.MissingKids,
						len(gr.MissingKids),
						strings.Join(gr.MissingKids, "\n"),
					),
				}
				assertMessages(t, got, want)
			})
			t.Run("Group Non exists", func(t *testing.T) {
				ms := mocks.NewMockService(map[int64]bool{
					12: true,
				})
				ms.SetCurrentGroup(nil)

				mockContext := MockContext{}
				messageHandler := internal.NewMessageHandler(ms)
				mockContext.setUserMessageWithTime(12, "Получить отсутсвующих", getUnixByDay(28, 22, 40))

				got := messageHandler.Process(&mockContext)
				want := contextHandlers.Response{
					Message: config.CurrentGroupDontFind,
				}
				assertMessages(t, got, want)
			})
		})
		t.Run("Send my groups", func(t *testing.T) {
			t.Run("If user have groups", func(t *testing.T) {
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
				ms := mocks.NewMockService(map[int64]bool{
					12: true,
				})
				ms.SetGroups(g)

				mockContext := MockContext{}
				messageHandler := internal.NewMessageHandler(ms)

				mockContext.setUserMessage(12, "Мои группы")

				want := contextHandlers.Response{
					Message: fmt.Sprintf(
						"%s4\n\n%s\n\n%s",
						config.MyGroups,
						"1. Гр 4 🕐 сб 10:00\n2. Гр 3 🕐 сб 14:00",
						"1. Гр 1 🕐 вс 10:00\n2. Гр 2 🕐 вс 12:00",
					),
					Keyboard: config.MyGroupsKeyboard,
				}
				got := messageHandler.Process(&mockContext)
				assertMessages(t, got, want)
				assertKeyboards(t, got, want)
			})
			t.Run("If user dont have groups", func(t *testing.T) {
				ms := mocks.NewMockService(map[int64]bool{
					12: true,
				})
				ms.SetGroups(nil)

				mockContext := MockContext{}

				messageHandler := internal.NewMessageHandler(ms)
				mockContext.setUserMessage(12, "Мои группы")

				want := contextHandlers.Response{
					Message:  config.UserDontHaveGroup,
					Keyboard: config.MyGroupsKeyboard,
				}
				got := messageHandler.Process(&mockContext)
				assertMessages(t, got, want)
				assertKeyboards(t, got, want)
			})
		})
	})

}

func assertKeyboards(t *testing.T, got contextHandlers.Response, want contextHandlers.Response) {
	t.Helper()

	if reflect.DeepEqual(got.Keyboard, want.Keyboard) != true {
		t.Errorf("Wanted: [%v]\nGot: [%v]\n", want.Keyboard, got.Keyboard)
	}
}

func assertMessages(t *testing.T, got contextHandlers.Response, want contextHandlers.Response) {
	t.Helper()
	if reflect.DeepEqual(got.Message, want.Message) != true {
		t.Errorf("Wanted: [%s]\nGot: [%s]\n", want.Message, got.Message)
		t.Errorf("Length Wanted: %d, Length Got: %d\n", len(want.Message), len(got.Message))
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
