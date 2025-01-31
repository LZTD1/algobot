package test

import (
	"fmt"
	"gopkg.in/telebot.v4"
	"reflect"
	"strings"
	"testing"
	"tgbot/internal/config"
	"tgbot/internal/contextHandlers"
	"tgbot/internal/domain"
	"tgbot/internal/stateMachine"
	"tgbot/tests/mocks"
	"time"
)

func TestDefaultHandler(t *testing.T) {
	t.Run("If user is not register", func(t *testing.T) {
		ms := mocks.NewMockService(make(map[int64]bool))

		mockContext := mocks.MockContext{}

		mockState := mocks.MockStateMachine{}
		mockState.SetStatement(12, stateMachine.Default)

		messageHandler := contextHandlers.NewOnText(ms, &mockState)

		mockContext.SetUserMessage(12, "hello world!")

		messageHandler.Handle(&mockContext)
		assertContextOptsLen(t, mockContext.SentMessages[0], 1)
		assertMessages(t, mockContext.SentMessages[0], config.Incorrect)
		assertKeyboards(t, mockContext.SentMessages[0], config.StartKeyboard)

	})
	t.Run("If user register", func(t *testing.T) {
		t.Run("Send any bullshit", func(t *testing.T) {
			mockContext := mocks.MockContext{}

			ms := mocks.NewMockService(map[int64]bool{
				12: true,
			})

			mockState := mocks.MockStateMachine{}
			mockState.SetStatement(12, stateMachine.Default)
			messageHandler := contextHandlers.NewOnText(ms, &mockState)

			mockContext.SetUserMessage(12, "aezakmi")

			messageHandler.Handle(&mockContext)

			assertContextOptsLen(t, mockContext.SentMessages[0], 1)
			assertMessages(t, mockContext.SentMessages[0], config.Incorrect)
			assertKeyboards(t, mockContext.SentMessages[0], config.StartKeyboard)
		})
		t.Run("Send settings", func(t *testing.T) {

			t.Run("Cookie set, notif off", func(t *testing.T) {
				mockContext := mocks.MockContext{}

				ms := mocks.NewMockService(map[int64]bool{
					12: true,
				})

				mockState := mocks.MockStateMachine{}
				mockState.SetStatement(12, stateMachine.Default)
				messageHandler := contextHandlers.NewOnText(ms, &mockState)

				ms.SetMockCookie("Cookie")
				mockContext.SetUserMessage(12, "Настройки")

				messageHandler.Handle(&mockContext)
				assertContextOptsLen(t, mockContext.SentMessages[0], 1)
				assertMessages(t, mockContext.SentMessages[0], fmt.Sprintf(
					"%s\n\n%s%s\n%s%s",
					config.Settings,
					config.Cookie,
					config.SetParam,
					config.ChatNotifications,
					config.NotSetParam,
				))
				assertKeyboards(t, mockContext.SentMessages[0], config.SettingsKeyboard)
			})
			t.Run("Cookie unset, notif off", func(t *testing.T) {
				mockContext := mocks.MockContext{}

				ms := mocks.NewMockService(map[int64]bool{
					12: true,
				})

				mockState := mocks.MockStateMachine{}
				mockState.SetStatement(12, stateMachine.Default)
				messageHandler := contextHandlers.NewOnText(ms, &mockState)

				ms.SetMockCookie("")
				mockContext.SetUserMessage(12, "Настройки")

				messageHandler.Handle(&mockContext)
				assertContextOptsLen(t, mockContext.SentMessages[0], 1)
				assertMessages(t, mockContext.SentMessages[0], fmt.Sprintf(
					"%s\n\n%s%s\n%s%s",
					config.Settings,
					config.Cookie,
					config.NotSetParam,
					config.ChatNotifications,
					config.NotSetParam,
				))
				assertKeyboards(t, mockContext.SentMessages[0], config.SettingsKeyboard)
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

				mockContext := mocks.MockContext{}

				mockState := mocks.MockStateMachine{}
				mockState.SetStatement(12, stateMachine.Default)
				messageHandler := contextHandlers.NewOnText(ms, &mockState)

				mockContext.SetUserMessageWithTime(12, "Получить отсутсвующих", getUnixByDay(28, 9, 40))

				messageHandler.Handle(&mockContext)
				assertContextOptsLen(t, mockContext.SentMessages[0], 0)
				assertMessages(t, mockContext.SentMessages[0], fmt.Sprintf(
					"%s%s\n%s%s\n\n%s%d\n%s%d\n\n```Отсутсвующие\n%s\n```",
					config.GroupName,
					gr.Name,
					config.Lection,
					gr.Lesson,
					config.TotalKids,
					gr.AllKids,
					config.MissingKids,
					len(gr.MissingKids),
					strings.Join(gr.MissingKids, "\n"),
				))
				// TODO assertKeyboards(t, mockContext.SentMessages[0], config.StartKeyboard)
			})
			t.Run("Group Non exists", func(t *testing.T) {
				ms := mocks.NewMockService(map[int64]bool{
					12: true,
				})
				ms.SetCurrentGroup(nil)

				mockContext := mocks.MockContext{}
				mockState := mocks.MockStateMachine{}
				mockState.SetStatement(12, stateMachine.Default)
				messageHandler := contextHandlers.NewOnText(ms, &mockState)
				mockContext.SetUserMessageWithTime(12, "Получить отсутсвующих", getUnixByDay(28, 22, 40))

				messageHandler.Handle(&mockContext)
				assertContextOptsLen(t, mockContext.SentMessages[0], 0)
				assertMessages(t, mockContext.SentMessages[0], config.CurrentGroupDontFind)
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

				mockContext := mocks.MockContext{}
				mockState := mocks.MockStateMachine{}
				mockState.SetStatement(12, stateMachine.Default)
				messageHandler := contextHandlers.NewOnText(ms, &mockState)

				mockContext.SetUserMessage(12, "Мои группы")

				messageHandler.Handle(&mockContext)

				assertContextOptsLen(t, mockContext.SentMessages[0], 1)
				assertMessages(t, mockContext.SentMessages[0], fmt.Sprintf(
					"%s4\n\n%s\n\n%s",
					config.MyGroups,
					"1. Гр 4 🕐 сб 10:00\n2. Гр 3 🕐 сб 14:00",
					"1. Гр 1 🕐 вс 10:00\n2. Гр 2 🕐 вс 12:00",
				))
				assertKeyboards(t, mockContext.SentMessages[0], config.MyGroupsKeyboard)
			})
			t.Run("If user dont have groups", func(t *testing.T) {
				ms := mocks.NewMockService(map[int64]bool{
					12: true,
				})
				ms.SetGroups(nil)

				mockContext := mocks.MockContext{}

				mockState := mocks.MockStateMachine{}
				mockState.SetStatement(12, stateMachine.Default)
				messageHandler := contextHandlers.NewOnText(ms, &mockState)
				mockContext.SetUserMessage(12, "Мои группы")

				messageHandler.Handle(&mockContext)
				assertContextOptsLen(t, mockContext.SentMessages[0], 1)
				assertMessages(t, mockContext.SentMessages[0], config.UserDontHaveGroup)
				assertKeyboards(t, mockContext.SentMessages[0], config.MyGroupsKeyboard)
			})
		})
	})

}

func assertMessages(t *testing.T, got mocks.SentMessage, wantedText string) {
	t.Helper()

	if got.What.(string) != wantedText {
		t.Errorf("Wanted [%s], but got [%s]", wantedText, got.What.(string))
	}
}
func assertKeyboards(t *testing.T, got mocks.SentMessage, wantedMarkup *telebot.ReplyMarkup) {
	t.Helper()

	var gotMarkup *telebot.ReplyMarkup
	for _, opt := range got.Opts {
		if markup, ok := opt.(*telebot.ReplyMarkup); ok {
			gotMarkup = markup
			break
		}
	}

	if !reflect.DeepEqual(gotMarkup, wantedMarkup) {
		t.Errorf("Wanted keyboard [%+v],\n but got [%+v]", wantedMarkup, gotMarkup)
	}
}

func assertContextOptsLen(t *testing.T, sent mocks.SentMessage, i int) {
	t.Helper()

	if len(sent.Opts) != i {
		t.Errorf("%+v\n", sent)
		t.Errorf("Wanted context len = %d, got, %d", i, len(sent.Opts))
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
