package test

import (
	"algobot/internal/domain/telegram/keyboards"
	"algobot/internal/lib/fsm"
	"algobot/internal/telegram/handlers/callback"
	mocks2 "algobot/test/mocks/telegram"
	mocks "algobot/test/mocks/telegram/handlers"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gopkg.in/telebot.v4"
	"testing"
)

func TestChangeCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stater := mocks.NewMockStateChanger(ctrl)
	mctx := mocks2.NewMockContext(ctrl)

	gomock.InOrder(
		mctx.EXPECT().Sender().Return(&telebot.User{ID: 1}).Times(1),
		stater.EXPECT().SetState(int64(1), fsm.SendingCookie).Times(1),
		mctx.EXPECT().Send("Отправьте мне свои cookie 🍪\nИнструкция: https://telegra.ph/Kak-dobavit-v-bota-svoi-Cookie-02-05", keyboards.RejectKeyboard()).Return(nil).Times(1),
	)

	err := callback.NewChangeCookie(stater)(mctx)
	assert.NoError(t, err)
}
