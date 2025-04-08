package test

import (
	"algobot/internal/telegram/handlers/callback"
	mocks2 "algobot/test/mocks"
	mocks3 "algobot/test/mocks/telegram"
	mocks "algobot/test/mocks/telegram/handlers"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	tele "gopkg.in/telebot.v4"
	"testing"
)

func TestRefreshGroups(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	refresher := mocks.NewMockGroupRefresher(ctrl)
	log := mocks2.NewMockLogger()
	mctx := mocks3.NewMockContext(ctrl)

	handler := callback.RefreshGroup(refresher, log)

	mctx.EXPECT().Get(gomock.Any()).Return("").AnyTimes()
	mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).AnyTimes()
	t.Run("happy path", func(t *testing.T) {
		gomock.InOrder(
			refresher.EXPECT().RefreshGroup(int64(1), "").Return(nil),
			mctx.EXPECT().Edit("Успешно обновлено!"),
		)
		err := handler(mctx)
		assert.NoError(t, err)
	})
	t.Run("refresher return err", func(t *testing.T) {
		errExp := errors.New("exp")

		gomock.InOrder(
			refresher.EXPECT().RefreshGroup(int64(1), "").Return(errExp),
		)
		err := handler(mctx)
		assert.ErrorIs(t, err, errExp)
	})
}
