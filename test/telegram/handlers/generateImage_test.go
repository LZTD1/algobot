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

func TestGen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gen := mocks.NewMockGeneratorImage(ctrl)
	log := mocks2.NewMockLogger()
	mctx := mocks3.NewMockContext(ctrl)
	mapi := mocks3.NewMockAPI(ctrl)

	handler := text.GenerateImage(gen, log)

	mctx.EXPECT().Bot().Return(mapi).AnyTimes()
	mctx.EXPECT().Get(gomock.Any()).Return("").AnyTimes()

	t.Run("happy path", func(t *testing.T) {
		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			mctx.EXPECT().Message().Return(&tele.Message{Text: "/image firefox"}).Times(1),

			mctx.EXPECT().Message().Return(&tele.Message{ID: 1}).Times(1),
			mapi.EXPECT().Reply(&tele.Message{ID: 1}, "⚙️ Генерирую изображение ...").Return(&tele.Message{ID: 1}, nil).Times(1),

			gen.EXPECT().GenerateImage(int64(1), "firefox", "").Return("https", nil),
			mapi.EXPECT().Edit(&tele.Message{ID: 1}, &tele.Photo{
				File: tele.FromURL("https"),
			}).Return(nil, nil).Times(1),
		)
		err := handler(mctx)
		assert.NoError(t, err)
	})
	t.Run("first send ret err", func(t *testing.T) {
		errExp := errors.New("err")

		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			mctx.EXPECT().Message().Return(&tele.Message{Text: "/image firefox"}).Times(1),

			mctx.EXPECT().Message().Return(&tele.Message{ID: 1}).Times(1),
			mapi.EXPECT().Reply(&tele.Message{ID: 1}, "⚙️ Генерирую изображение ...").Return(&tele.Message{ID: 1}, errExp).Times(1),
		)
		err := handler(mctx)
		assert.ErrorIs(t, err, errExp)
	})
	t.Run("GenerateImage send err", func(t *testing.T) {
		errExp := errors.New("err")

		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			mctx.EXPECT().Message().Return(&tele.Message{Text: "/image firefox"}).Times(1),

			mctx.EXPECT().Message().Return(&tele.Message{ID: 1}).Times(1),
			mapi.EXPECT().Reply(&tele.Message{ID: 1}, "⚙️ Генерирую изображение ...").Return(&tele.Message{ID: 1}, nil).Times(1),

			gen.EXPECT().GenerateImage(int64(1), "firefox", "").Return("https", errExp),
			mapi.EXPECT().Edit(&tele.Message{ID: 1}, "⚠️ К сожалению, я не смог сгенерировать изображение, попробуйте снова чуть позже").Return(&tele.Message{ID: 1}, nil).Times(1),
		)

		err := handler(mctx)
		assert.ErrorIs(t, err, errExp)
	})
}
