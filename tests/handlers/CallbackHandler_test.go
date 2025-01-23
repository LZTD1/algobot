package test

import (
	"fmt"
	"testing"
	"tgbot/internal/config"
	"tgbot/internal/contextHandlers"
	"tgbot/internal/stateMachine"
	"tgbot/tests/mocks"
)

func TestCallback(t *testing.T) {
	t.Run("Set cookie", func(t *testing.T) {
		ms := mocks.NewMockService(map[int64]bool{
			12: true,
		})
		mockState := mocks.MockStateMachine{}

		mockContext := mocks.MockContext{}
		mockContext.SetUserMessage(12, "set_cookie")

		queryHandler := contextHandlers.NewOnCallback(ms, &mockState)

		queryHandler.Process(&mockContext)

		assertContextOptsLen(t, mockContext.SentMessages[0], 1)
		assertMessages(t, mockContext.SentMessages[0], config.SendingCookie)
		assertKeyboards(t, mockContext.SentMessages[0], config.RejectKeyboard)

		assertMockStatement(t, mockState, stateMachine.SendingCookie, 3)
	})
	t.Run("Change notification", func(t *testing.T) {
		ms := mocks.NewMockService(map[int64]bool{
			12: true,
		})

		mockState := mocks.MockStateMachine{}
		mockContext := mocks.MockContext{}
		mockContext.SetUserMessage(12, "change_notification")

		queryHandler := contextHandlers.NewOnCallback(ms, &mockState)

		if ms.StubNotification != false {
			t.Fatalf("Wanted notif false, got true")
		}
		queryHandler.Process(&mockContext)
		if ms.StubNotification != true {
			t.Fatalf("Wanted notif true, got false")
		}

		assertContextOptsLen(t, mockContext.SentMessages[0], 1)
		assertMessages(t, mockContext.SentMessages[0], fmt.Sprintf(
			"%s\n\n%s%s\n%s%s",
			config.Settings,
			config.Cookie,
			config.NotSetParam,
			config.ChatNotifications,
			config.SetParam,
		))
		assertKeyboards(t, mockContext.SentMessages[0], config.SettingsKeyboard)

	})
	t.Run("Refresh groups", func(t *testing.T) {
		ms := mocks.NewMockService(map[int64]bool{
			12: true,
		})

		mockState := mocks.MockStateMachine{}
		mockContext := mocks.MockContext{}
		mockContext.SetUserMessage(12, "refresh_groups")

		queryHandler := contextHandlers.NewOnCallback(ms, &mockState)
		got := queryHandler.Process(&mockContext)
		fmt.Println(got)
	})
}

func assertMockStatement(t *testing.T, mockState mocks.MockStateMachine, wantedState stateMachine.Statement, wantedLen int) {
	if mockState.Current != wantedState {
		t.Errorf("Wanted %+v, got %+v", wantedState, mockState.Current)
	}
	if len(mockState.Calls) != wantedLen {
		t.Errorf("Wanted len %d, got %d", wantedLen, len(mockState.Calls))
	}
}
