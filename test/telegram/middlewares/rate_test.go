package test

import (
	"algobot/internal/config"
	rate2 "algobot/internal/telegram/middleware/rate"
	"algobot/test/mocks"
	mocks2 "algobot/test/mocks/telegram"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gopkg.in/telebot.v4"
	"testing"
	"time"
)

func TestRate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocks.NewMockLogger()
	expected := 3

	actual := 0

	rate := rate2.New(log, config.RateLimit{
		FillPeriod:  1 * time.Second,
		BucketLimit: expected,
	})

	mctx := mocks2.NewMockContext(ctrl)
	mctx.EXPECT().Get(gomock.Any()).Return(nil).AnyTimes()
	mctx.EXPECT().Sender().Return(&telebot.User{ID: 1}).AnyTimes()

	mctx.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
	hfunc := func(ctx telebot.Context) error {
		actual++
		return nil
	}
	for i := 0; i < expected+1; i++ {
		err := rate(hfunc)(mctx)
		assert.NoError(t, err)
	}
	assert.Equal(t, expected, actual)
}
