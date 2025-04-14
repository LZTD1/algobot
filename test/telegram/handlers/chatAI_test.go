package test

import (
	"algobot/internal/telegram/handlers/text"
	mocks2 "algobot/test/mocks"
	mocks "algobot/test/mocks/telegram"
	mocks3 "algobot/test/mocks/telegram/handlers"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	tele "gopkg.in/telebot.v4"
	"testing"
)

func TestChatAI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mctx := mocks.NewMockContext(ctrl)
	chatter := mocks3.NewMockChatter(ctrl)
	log := mocks2.NewMockLogger()
	mapi := mocks.NewMockAPI(ctrl)

	mctx.EXPECT().Bot().Return(mapi).AnyTimes()
	handler := text.ChatAI(chatter, log)
	t.Run("happy path", func(t *testing.T) {
		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			mctx.EXPECT().Message().Return(&tele.Message{Text: "some"}).Times(1),
			mctx.EXPECT().Get(gomock.Any()).Return("").Times(1),

			mctx.EXPECT().Message().Return(&tele.Message{ID: 1}).Times(1),
			mapi.EXPECT().Reply(&tele.Message{ID: 1}, "⚙️ Думаю что ответить ...").Return(&tele.Message{ID: 2}, nil).Times(1),

			chatter.EXPECT().ChatAI(int64(1), "some", "").Return("msg", nil).Times(1),

			mapi.EXPECT().Edit(&tele.Message{ID: 2}, "msg", tele.ModeMarkdown).Return(&tele.Message{ID: 2}, nil).Times(1),
		)
		err := handler(mctx)
		assert.NoError(t, err)
	})
	t.Run("ChatAI returns err", func(t *testing.T) {
		errExp := errors.New("")
		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			mctx.EXPECT().Message().Return(&tele.Message{Text: "some"}).Times(1),
			mctx.EXPECT().Get(gomock.Any()).Return("").Times(1),

			mctx.EXPECT().Message().Return(&tele.Message{ID: 1}).Times(1),
			mapi.EXPECT().Reply(&tele.Message{ID: 1}, "⚙️ Думаю что ответить ...").Return(&tele.Message{ID: 2}, nil).Times(1),

			chatter.EXPECT().ChatAI(int64(1), "some", "").Return("", errExp).Times(1),

			mapi.EXPECT().Edit(&tele.Message{ID: 2}, "⚠️ К сожалению, я не смог ответить на ваше сообщение, попробуйте снова чуть позже").Return(&tele.Message{ID: 2}, nil).Times(1),
		)
		err := handler(mctx)
		assert.ErrorIs(t, err, errExp)
	})
}
