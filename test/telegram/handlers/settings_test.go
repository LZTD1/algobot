package test

import (
	"algobot/internal/domain/telegram/keyboards"
	"algobot/internal/telegram/handlers/text"
	mocks3 "algobot/test/mocks"
	mocks2 "algobot/test/mocks/telegram"
	mocks "algobot/test/mocks/telegram/handlers"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	tele "gopkg.in/telebot.v4"
	"testing"
)

func TestSettings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uinformer := mocks.NewMockUserInformer(ctrl)
	mctx := mocks2.NewMockContext(ctrl)
	log := mocks3.NewMockLogger()

	t.Run("Happy path", func(t *testing.T) {
		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			mctx.EXPECT().Get(gomock.Any()).Return("a-1").Times(1),
			uinformer.EXPECT().Cookies(int64(1)).Return("", nil).Times(1),
			uinformer.EXPECT().Notification(int64(1)).Return(true, nil).Times(1),
			mctx.EXPECT().Send(text.GetSettingsMessage("", true), keyboards.Settings()).Return(nil).Times(1),
		)

		settings := text.NewSettings(uinformer, log)
		err := settings(mctx)
		assert.NoError(t, err)
	})
	t.Run("Cookie returns err", func(t *testing.T) {
		errCookie := errors.New("err")

		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			mctx.EXPECT().Get(gomock.Any()).Return("a-1").Times(1),
			uinformer.EXPECT().Cookies(int64(1)).Return("", errCookie).Times(1),
		)

		settings := text.NewSettings(uinformer, log)
		err := settings(mctx)
		assert.ErrorIs(t, err, errCookie)
	})
	t.Run("Notification returns err", func(t *testing.T) {
		errNotification := errors.New("err")

		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			mctx.EXPECT().Get(gomock.Any()).Return("a-1").Times(1),
			uinformer.EXPECT().Cookies(int64(1)).Return("", nil).Times(1),
			uinformer.EXPECT().Notification(int64(1)).Return(false, errNotification).Times(1),
		)

		settings := text.NewSettings(uinformer, log)
		err := settings(mctx)
		assert.ErrorIs(t, err, errNotification)
	})
}
