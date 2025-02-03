package test

import (
	"fmt"
	"gopkg.in/telebot.v4"
	"reflect"
	"strings"
	"testing"
	"tgbot/internal/config"
	"tgbot/internal/contextHandlers"
	appError "tgbot/internal/error"
	"tgbot/internal/models"
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
				gr := models.Group{
					GroupID:    1,
					Title:      "Title",
					TimeLesson: getDayByTime(28, 10, 0),
				}
				ms := mocks.NewMockService(map[int64]bool{
					12: true,
				})
				ms.Actual = models.ActualInformation{
					LessonTitle: "LTitle",
					LessonId:    0,
					MissingKids: []int{1, 2},
				}
				ms.AllNames = map[int]string{
					1: "vasya",
					2: "petya",
					3: "kirill",
				}
				ms.SetCurrentGroup(&gr)

				mockContext := mocks.MockContext{}

				mockState := mocks.MockStateMachine{}
				mockState.SetStatement(12, stateMachine.Default)
				messageHandler := contextHandlers.NewOnText(ms, &mockState)

				mockContext.SetUserMessageWithTime(12, "Получить отсутсвующих", getUnixByDay(28, 9, 40))

				messageHandler.Handle(&mockContext)
				assertContextOptsLen(t, mockContext.SentMessages[0], 1)
				assertMessages(t, mockContext.SentMessages[0], "Группа по курсу: Title\nЛекция: LTitle\n\nОбщее число детей: 3\nОтсутствуют: 2\n\n```Отсутствующие\nvasya\npetya\n```")
				// TODO assertKeyboards
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
				assertContextOptsLen(t, mockContext.SentMessages[0], 1)
				containtsMessages(t, mockContext.SentMessages[0], config.CurrentGroupDontFind)
			})
		})
		t.Run("Send my groups", func(t *testing.T) {
			t.Run("If user have groups", func(t *testing.T) {
				g := []models.Group{
					{
						GroupID:    1,
						Title:      "Гр 1",
						TimeLesson: getDayByTime(28, 10, 0),
					},
					{
						GroupID:    3,
						Title:      "Гр 3",
						TimeLesson: getDayByTime(27, 14, 0),
					},
					{
						GroupID:    2,
						Title:      "Гр 2",
						TimeLesson: getDayByTime(21, 12, 0),
					},
					{
						GroupID:    4,
						Title:      "Гр 4",
						TimeLesson: getDayByTime(27, 10, 0),
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
				ms.SetGroupsErr(appError.ErrHasNone)

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

func containtsMessages(t *testing.T, excepted mocks.SentMessage, wanted string) {
	t.Helper()

	if strings.Contains(excepted.What.(string), wanted) {

		t.Fatalf("Wanted %s, got %s", wanted, excepted.What.(string))
	}
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
