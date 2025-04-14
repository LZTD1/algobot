package test

import (
	"algobot/internal/app/scheduler"
	"algobot/internal/config"
	"algobot/internal/domain/models"
	scheduler2 "algobot/internal/domain/scheduler"
	"algobot/test/mocks"
	mocks2 "algobot/test/mocks/app"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestSchedule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocks.NewMockLogger()
	domain := mocks2.NewMockDomain(ctrl)
	bo := mocks2.NewMockBackoffice(ctrl)
	cfg := &config.Config{Backoffice: config.Backoffice{
		MessageTimer: 5 * time.Second,
	}}

	app := scheduler.New(log, cfg, domain, bo)

	go func() {
		msg := <-app.Chan()
		assert.Equal(t, assetMSG, msg)
		app.Stop()
	}()

	domain.EXPECT().UsersByNotification(1).Return([]models.User{
		{
			ID:               1,
			Uid:              1,
			Cookie:           "Cookie",
			LastNotification: "LastNotification",
			Notification:     1,
		},
	}, nil).Times(1)
	bo.EXPECT().MessagesUser(int64(1), "LastNotification").Return([]scheduler2.Message{assetMSG}, nil).Times(1)
	domain.EXPECT().ChaneNotifDate(int64(1), "newTime").Return(nil).Times(1)

	app.GetMessage()
}

var assetMSG = scheduler2.Message{
	To:      1,
	From:    "From",
	Theme:   "Theme",
	Link:    "Link",
	Text:    "Text",
	LinkURL: "",
	Time:    "newTime",
}
