package test

import (
	"algobot/internal/telegram/middleware/auth"
	"algobot/test/mocks"
	mocks2 "algobot/test/mocks/telegram"
	mocks3 "algobot/test/mocks/telegram/middleware"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	tele "gopkg.in/telebot.v4"
	"testing"
)

func TestAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocks.NewMockLogger()
	mctx := mocks2.NewMockContext(ctrl)
	auther := mocks3.NewMockAuther(ctrl)

	handler := auth.New(auther, log)

	mctx.EXPECT().Get(gomock.Any()).Return("").AnyTimes()
	t.Run("happy path is not registered", func(t *testing.T) {
		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			auther.EXPECT().IsRegistered(int64(1)).Return(false, nil).Times(1),
			auther.EXPECT().Register(int64(1)).Return(nil).Times(1),
		)
		err := handler(func(context tele.Context) error {
			return nil
		})(mctx)
		assert.NoError(t, err)
	})
	t.Run("happy path is registered", func(t *testing.T) {
		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			auther.EXPECT().IsRegistered(int64(1)).Return(true, nil).Times(1),
		)
		err := handler(func(context tele.Context) error {
			return nil
		})(mctx)
		assert.NoError(t, err)
	})
	t.Run("if isReg return err", func(t *testing.T) {
		errExpected := errors.New("err")

		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			auther.EXPECT().IsRegistered(int64(1)).Return(false, errExpected).Times(1),
		)
		err := handler(func(context tele.Context) error {
			return nil
		})(mctx)
		assert.ErrorIs(t, err, errExpected)
	})
	t.Run("if isReg return err", func(t *testing.T) {
		errExpected := errors.New("err")

		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			auther.EXPECT().IsRegistered(int64(1)).Return(false, nil).Times(1),
			auther.EXPECT().Register(int64(1)).Return(errExpected).Times(1),
		)
		err := handler(func(context tele.Context) error {
			return nil
		})(mctx)
		assert.ErrorIs(t, err, errExpected)
	})
}
