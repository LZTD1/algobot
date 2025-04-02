package test

import (
	"algobot/internal/domain/telegram/keyboards"
	"algobot/internal/lib/fsm"
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

func TestAI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocks3.NewMockLogger()
	mctx := mocks2.NewMockContext(ctrl)
	ai := mocks.NewMockAIInformer(ctrl)
	stater := mocks.NewMockAIStater(ctrl)

	h := text.NewAI(ai, log, stater)
	mctx.EXPECT().Get("trace_id").Return("a-1").AnyTimes()

	t.Run("happy path", func(t *testing.T) {
		aiRet := text.AIInfo{
			TextModel:  "1",
			ImageModel: "1",
		}

		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			ai.EXPECT().GetAIInfo().Return(aiRet, nil).Times(1),
			stater.EXPECT().SetState(int64(1), fsm.ChattingAI).Times(1),
			mctx.EXPECT().Send(text.GetAIMessage(aiRet), keyboards.RejectKeyboard()).Times(1),
		)
		err := h(mctx)
		assert.NoError(t, err)
	})
	t.Run("GetAIInfo return err", func(t *testing.T) {
		aiRet := text.AIInfo{}
		aiErr := errors.New("GetAIInfo error")

		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			ai.EXPECT().GetAIInfo().Return(aiRet, aiErr).Times(1),
			mctx.EXPECT().Send("Упс, AI сейчас не работает!").Times(1),
		)
		err := h(mctx)
		assert.NoError(t, err)
	})
}
