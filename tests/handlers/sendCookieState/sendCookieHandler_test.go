package test

import (
	"algobot/internal_old/config"
	"algobot/internal_old/contextHandlers"
	"algobot/internal_old/stateMachine"
	"algobot/tests/mocks"
	"github.com/golang/mock/gomock"
	"gopkg.in/telebot.v4"
	"reflect"
	"testing"
)

func TestSending(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAI := mocks.NewMockAIService(ctrl)

	t.Run("Send reject action", func(t *testing.T) {
		ms := mocks.NewMockService(make(map[int64]bool))

		mockContext := mocks.MockContext{}

		mockState := mocks.MockStateMachine{}
		mockState.SetStatement(12, stateMachine.SendingCookie)

		messageHandler := contextHandlers.NewOnText(ms, &mockState, mockAI)

		mockContext.SetUserMessage(12, "Отменить действие")

		messageHandler.Handle(&mockContext)

		if mockState.Current != stateMachine.Default {
			t.Fatalf("Wanted default got %+v\n", mockState.Current)
		}
		assertContextOptsLen(t, mockContext.SentMessages[0], 1)
		assertMessages(t, mockContext.SentMessages[0], config.StartText)
		assertKeyboards(t, mockContext.SentMessages[0], config.StartKeyboard)
	})
	t.Run("Send cookie", func(t *testing.T) {
		ms := mocks.NewMockService(make(map[int64]bool))

		mockContext := mocks.MockContext{}

		mockState := mocks.MockStateMachine{}
		mockState.SetStatement(12, stateMachine.SendingCookie)

		messageHandler := contextHandlers.NewOnText(ms, &mockState, mockAI)

		mockContext.SetUserMessage(12, "aezakmi")

		messageHandler.Handle(&mockContext)

		if mockState.Current != stateMachine.Default {
			t.Fatalf("Wanted default got %+v\n", mockState.Current)
		}
		if ms.SettedCookie[0] != "12" || ms.SettedCookie[1] != "aezakmi" {
			t.Fatalf("Wanted setted cookie, got %+v\n", ms.SettedCookie)
		}
		assertContextOptsLen(t, mockContext.SentMessages[0], 1)
		assertMessages(t, mockContext.SentMessages[0], config.CookieSet)
		assertKeyboards(t, mockContext.SentMessages[0], config.StartKeyboard)
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
