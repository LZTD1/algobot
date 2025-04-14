package test

import (
	"algobot/internal/telegram/handlers/text"
	mocks2 "algobot/test/mocks"
	mocks3 "algobot/test/mocks/telegram"
	mocks "algobot/test/mocks/telegram/handlers"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	tele "gopkg.in/telebot.v4"
	"testing"
)

func TestReset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reseter := mocks.NewMockReseter(ctrl)
	log := mocks2.NewMockLogger()
	mctx := mocks3.NewMockContext(ctrl)

	handler := text.NewReset(reseter, log)

	mctx.EXPECT().Get(gomock.Any()).Return("").AnyTimes()
	t.Run("happy path", func(t *testing.T) {
		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			reseter.EXPECT().ResetHistory(int64(1), "").Return(nil).Times(1),
			mctx.EXPECT().Send("История успешно отчищена").Return(nil).Times(1),
		)

		err := handler(mctx)
		assert.NoError(t, err)
	})
	t.Run("ResetHistory returns err", func(t *testing.T) {
		errExp := errors.New("err")

		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			reseter.EXPECT().ResetHistory(int64(1), "").Return(errExp).Times(1),
		)

		err := handler(mctx)
		assert.ErrorIs(t, err, errExp)
	})
}
