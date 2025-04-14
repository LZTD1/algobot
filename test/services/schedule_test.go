package test

import (
	"algobot/internal/domain/scheduler"
	"algobot/internal/services/schedule"
	mocks "algobot/test/mocks/services"
	"go.uber.org/mock/gomock"
	"gopkg.in/telebot.v4"
	"testing"
	"time"
)

func TestShedule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ch := make(chan scheduler.Message, 2)
	sender := mocks.NewMockSender(ctrl)

	sch := schedule.NewSchedule(ch, sender)

	ch <- scheduler.Message{
		To:      123,
		From:    "From",
		Theme:   "Theme",
		Link:    "Link",
		Text:    "Some Text",
		LinkURL: "",
	}
	ch <- scheduler.Message{
		To:      333,
		From:    "From2",
		Theme:   "Theme2",
		Link:    "Link2",
		Text:    "Some Text2",
		LinkURL: "Link2",
	}
	sender.EXPECT().Send(telebot.ChatID(123),
		"🔔 Новое сообщение\n\nОт: From\nТема: Theme\nСсылка: Link\n\n```Сообщение:\nSome Text\n```",
		telebot.ModeMarkdown,
		telebot.NoPreview)
	sender.EXPECT().Send(telebot.ChatID(333),
		&telebot.Photo{File: telebot.FromURL("Link2"), Caption: "🔔 Новое сообщение\n\nОт: From2\nТема: Theme2\nСсылка: Link2\n\n"},
		telebot.ModeMarkdown,
		telebot.NoPreview)

	go sch.Process()
	time.Sleep(100 * time.Millisecond)
	close(ch)

}
