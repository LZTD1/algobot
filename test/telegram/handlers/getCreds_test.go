package test

import (
	"algobot/internal/domain/models"
	"algobot/internal/telegram/handlers/callback"
	"algobot/test/mocks"
	mocks3 "algobot/test/mocks/telegram"
	mocks2 "algobot/test/mocks/telegram/handlers"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gopkg.in/telebot.v4"
	"testing"
)

func TestGetterCreds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocks.NewMockLogger()
	creds := mocks2.NewMockGetterCreds(ctrl)
	mctx := mocks3.NewMockContext(ctrl)

	handler := callback.GetCreds(creds, log)

	mctx.EXPECT().Get(gomock.Any()).Return("").AnyTimes()
	mctx.EXPECT().Sender().Return(&telebot.User{ID: 1}).AnyTimes()

	t.Run("happy path", func(t *testing.T) {
		gomock.InOrder(
			mctx.EXPECT().Callback().Return(&telebot.Callback{Data: "\fget_creds_123"}),
			creds.EXPECT().Creds(int64(1), "123", "").Return([]models.Credential{
				{
					Fullname: "f",
					Login:    "l",
					Password: "p",
				},
			}, nil).Times(1),
			mctx.EXPECT().Send("<i>f</i> - l : p\n", telebot.ModeHTML).Return(nil).Times(1),
		)
		err := handler(mctx)
		assert.NoError(t, err)
	})
	t.Run("creds return err", func(t *testing.T) {
		errExp := errors.New("err")

		gomock.InOrder(
			mctx.EXPECT().Callback().Return(&telebot.Callback{Data: "\fget_creds_123"}),
			creds.EXPECT().Creds(int64(1), "123", "").Return(nil, errExp).Times(1),
		)
		err := handler(mctx)
		assert.ErrorIs(t, err, errExp)
	})
}
