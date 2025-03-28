package test

import (
	"algobot/internal/lib/fsm"
	"algobot/internal/telegram/dispatcher/text"
	"algobot/test_v2/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_Dispatcher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h1 := mocks.NewMockHandler(ctrl)
	h2 := mocks.NewMockHandler(ctrl)
	h3 := mocks.NewMockHandler(ctrl)

	log := mocks.NewMockLogger()

	dispather := text.NewDispatcher(log)
	dispather.Register(fsm.Default, h1)
	dispather.Register(fsm.SendingCookie, h2)
	dispather.Register(fsm.ChattingAI, h3)

	handler := dispather.GetHandlers(fsm.Default)
	assert.Same(t, h1, handler)
	handler = dispather.GetHandlers(fsm.SendingCookie)
	assert.Same(t, h2, handler)
	handler = dispather.GetHandlers(fsm.ChattingAI)
	assert.Same(t, h3, handler)

}
