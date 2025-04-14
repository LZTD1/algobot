package test

import (
	"algobot/internal/telegram/handlers/callback"
	mocks3 "algobot/test/mocks"
	mocks "algobot/test/mocks/telegram"
	mocks2 "algobot/test/mocks/telegram/handlers"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	tele "gopkg.in/telebot.v4"
	"testing"
)

func TestNotification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mctx := mocks.NewMockContext(ctrl)
	notif := mocks2.NewMockNotificationChanger(ctrl)
	log := mocks3.NewMockLogger()

	handler := callback.NewChangeNotification(notif, log)
	mctx.EXPECT().Get("trace_id").Return("a-1").AnyTimes()

	t.Run("happy path", func(t *testing.T) {
		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			notif.EXPECT().Notification(int64(1)).Return(true, nil).Times(1),
			notif.EXPECT().SetNotification(int64(1), false).Return(nil).Times(1),
			mctx.EXPECT().Edit("Уведомления переключены").Return(nil).Times(1),
		)
		err := handler(mctx)
		assert.NoError(t, err)
	})
	t.Run("Notification err", func(t *testing.T) {
		expErr := errors.New("test error")

		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			notif.EXPECT().Notification(int64(1)).Return(false, expErr).Times(1),
		)
		err := handler(mctx)
		assert.ErrorIs(t, err, expErr)
	})
	t.Run("SetNotification err", func(t *testing.T) {
		expErr := errors.New("test error")

		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			notif.EXPECT().Notification(int64(1)).Return(true, nil).Times(1),
			notif.EXPECT().SetNotification(int64(1), false).Return(expErr).Times(1),
		)
		err := handler(mctx)
		assert.ErrorIs(t, err, expErr)
	})
}
