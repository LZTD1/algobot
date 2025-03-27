package test

import (
	"algobot/internal_old/config"
	"algobot/internal_old/contextHandlers"
	"algobot/internal_old/stateMachine"
	"algobot/tests/mocks"
	"fmt"
	"testing"
)

func TestCallback(t *testing.T) {
	t.Run("Set cookie", func(t *testing.T) {
		ms := mocks.NewMockService(map[int64]bool{
			12: true,
		})
		mockState := mocks.MockStateMachine{}
		mockState.SetStatement(12, stateMachine.Default)

		mockContext := mocks.MockContext{}
		mockContext.SetUserMessage(12, "set_cookie")

		queryHandler := contextHandlers.NewOnCallback(ms, &mockState)

		queryHandler.Handle(&mockContext)

		assertContextOptsLen(t, mockContext.SentMessages[0], 1)
		assertMessages(t, mockContext.SentMessages[0], config.SendingCookie)
		assertKeyboards(t, mockContext.SentMessages[0], config.RejectKeyboard)

		assertMockStatement(t, mockState, stateMachine.SendingCookie, 8)
	})
	t.Run("Change notification", func(t *testing.T) {
		ms := mocks.NewMockService(map[int64]bool{
			12: true,
		})

		mockState := mocks.MockStateMachine{}
		mockState.SetStatement(12, stateMachine.Default)
		mockContext := mocks.MockContext{}
		mockContext.SetUserMessage(12, "change_notification")

		queryHandler := contextHandlers.NewOnCallback(ms, &mockState)

		if ms.StubNotification != false {
			t.Fatalf("Wanted notif false, got true")
		}
		queryHandler.Handle(&mockContext)
		if ms.StubNotification != true {
			t.Fatalf("Wanted notif true, got false")
		}

		assertContextOptsLen(t, mockContext.SentMessages[0], 0)
		assertMessages(t, mockContext.SentMessages[0], "Настройки уведомлений были изменены!")

	})
	t.Run("Refresh groups without error", func(t *testing.T) {
		ms := mocks.NewMockService(map[int64]bool{
			12: true,
		})

		mockState := mocks.MockStateMachine{}
		mockState.SetStatement(12, stateMachine.Default)
		mockContext := mocks.MockContext{}
		mockContext.SetUserMessage(12, "refresh_groups")

		queryHandler := contextHandlers.NewOnCallback(ms, &mockState)
		queryHandler.Handle(&mockContext)

		if len(mockContext.SentMessages) != 2 {
			t.Fatalf("Wanted 2, got %d", len(mockContext.SentMessages))
		}
		assertMessages(t, mockContext.SentMessages[0], config.UpdateStarted)
		assertMessages(t, mockContext.SentMessages[1], config.UpdateEnd)
	})
	t.Run("Close lesson", func(t *testing.T) {
		ms := mocks.NewMockService(map[int64]bool{
			12: true,
		})
		mockState := mocks.MockStateMachine{}
		mockState.SetStatement(12, stateMachine.Default)

		mockContext := mocks.MockContext{}
		mockContext.SetUserMessage(12, "close_lesson_1_1")

		queryHandler := contextHandlers.NewOnCallback(ms, &mockState)

		queryHandler.Handle(&mockContext)

		assertContextOptsLen(t, mockContext.SentMessages[0], 0)
		sprintf := fmt.Sprintf("CloseLesson(%d, %d, %d)", 12, 1, 1)
		if ms.Calls[0] != sprintf {
			t.Errorf("Wanted %s, got %s", sprintf, ms.Calls[0])
		}
		assertMessages(t, mockContext.SentMessages[0], config.SuccessfulChangeStatus)
	})
	t.Run("Open lesson", func(t *testing.T) {
		ms := mocks.NewMockService(map[int64]bool{
			12: true,
		})
		mockState := mocks.MockStateMachine{}
		mockState.SetStatement(12, stateMachine.Default)

		mockContext := mocks.MockContext{}
		mockContext.SetUserMessage(12, "open_lesson_1_1")

		queryHandler := contextHandlers.NewOnCallback(ms, &mockState)

		queryHandler.Handle(&mockContext)

		assertContextOptsLen(t, mockContext.SentMessages[0], 0)
		sprintf := fmt.Sprintf("OpenLesson(%d, %d, %d)", 12, 1, 1)
		if ms.Calls[0] != sprintf {
			t.Errorf("Wanted %s, got %s", sprintf, ms.Calls[0])
		}
		assertMessages(t, mockContext.SentMessages[0], config.SuccessfulChangeStatus)
	})
	t.Run("Get creds", func(t *testing.T) {
		ms := mocks.NewMockService(map[int64]bool{
			12: true,
		})
		mockState := mocks.MockStateMachine{}
		mockState.SetStatement(12, stateMachine.Default)

		mockContext := mocks.MockContext{}
		mockContext.SetUserMessage(12, "get_creds_1")

		queryHandler := contextHandlers.NewOnCallback(ms, &mockState)

		queryHandler.Handle(&mockContext)

		assertContextOptsLen(t, mockContext.SentMessages[0], 0)
		sprintf := fmt.Sprintf("AllCredentials(%d, %d)", 12, 1)
		if ms.Calls[0] != sprintf {
			t.Errorf("Wanted %s, got %s", sprintf, ms.Calls[0])
		}
		assertMessages(t, mockContext.SentMessages[0], fmt.Sprintf("Учетные записи детей:\n\nВаня [van:12]"))
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
