package test

import (
	"algobot/internal/domain/telegram/keyboards"
	"algobot/internal/lib/fsm"
	"algobot/internal/telegram/handlers/text"
	"algobot/test/mocks"
	mocks3 "algobot/test/mocks/telegram"
	mocks2 "algobot/test/mocks/telegram/handlers"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	tele "gopkg.in/telebot.v4"
	"testing"
)

func TestSendingCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocks.NewMockLogger()
	stater := mocks2.NewMockCookieStater(ctrl)
	setter := mocks2.NewMockCookieSetter(ctrl)
	mctx := mocks3.NewMockContext(ctrl)

	handler := text.NewSendingCookie(log, setter, stater)

	mctx.EXPECT().Get(gomock.Any()).Return("").AnyTimes()
	t.Run("happy path", func(t *testing.T) {
		cookieText := "ckc"

		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			mctx.EXPECT().Message().Return(&tele.Message{Text: cookieText}).Times(1),
			setter.EXPECT().SetCookie(int64(1), cookieText).Return(nil).Times(1),
			stater.EXPECT().SetState(int64(1), fsm.Default).Times(1),
			mctx.EXPECT().Send("Cookie успешно установлены", keyboards.Start()).Return(nil).Times(1),
		)

		err := handler(mctx)
		assert.NoError(t, err)
	})
	t.Run("SetCookie return err", func(t *testing.T) {
		cookieText := "ckc"
		errCookie := errors.New("err")

		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			mctx.EXPECT().Message().Return(&tele.Message{Text: cookieText}).Times(1),
			setter.EXPECT().SetCookie(int64(1), cookieText).Return(errCookie).Times(1),
		)

		err := handler(mctx)
		assert.ErrorIs(t, err, errCookie)
	})
}
