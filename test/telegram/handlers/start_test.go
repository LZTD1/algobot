package test

import (
	"algobot/internal/domain/telegram/keyboards"
	"algobot/internal/lib/fsm"
	"algobot/internal/telegram/handlers/text"
	mocks2 "algobot/test/mocks/telegram"
	mocks "algobot/test/mocks/telegram/handlers"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gopkg.in/telebot.v4"
	"testing"
)

func TestStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stater := mocks.NewMockSetStater(ctrl)
	mctx := mocks2.NewMockContext(ctrl)

	gomock.InOrder(
		mctx.EXPECT().Sender().Return(&telebot.User{ID: 1}).Times(1),
		stater.EXPECT().SetState(int64(1), fsm.Default).Times(1),
		mctx.EXPECT().Send("Открыто главное меню:", keyboards.Start()).Return(nil).Times(1),
	)

	err := text.NewStart(stater)(mctx)
	assert.NoError(t, err)
}
