package schedule

import (
	"algobot/internal/domain/scheduler"
	"fmt"
	"gopkg.in/telebot.v4"
	"strings"
)

type Sender interface {
	Send(to telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error)
}

type Schedule struct {
	ch     chan scheduler.Message
	sender Sender
}

func NewSchedule(ch chan scheduler.Message, sender Sender) *Schedule {
	return &Schedule{ch: ch, sender: sender}
}

func (s *Schedule) Process() {
	for msg := range s.ch {
		if msg.LinkURL != "" {
			p := &telebot.Photo{File: telebot.FromURL(msg.LinkURL), Caption: getMsg(msg)}
			s.sender.Send(
				telebot.ChatID(msg.To),
				p,
				telebot.ModeMarkdown,
				telebot.NoPreview,
			)
			continue
		}
		s.sender.Send(
			telebot.ChatID(msg.To),
			getMsg(msg),
			telebot.ModeMarkdown,
			telebot.NoPreview,
		)
	}
}

func getMsg(msg scheduler.Message) string {
	sb := strings.Builder{}
	sb.WriteString("🔔 Новое сообщение\n\n")
	sb.WriteString(fmt.Sprintf("От: %s\n", msg.From))
	sb.WriteString(fmt.Sprintf("Тема: %s\n", msg.Theme))
	sb.WriteString(fmt.Sprintf("Ссылка: %s\n\n", msg.Link))
	if msg.LinkURL == "" {
		sb.WriteString("```Сообщение:\n")
		sb.WriteString(msg.Text)
		sb.WriteString("\n```")
	}
	return sb.String()
}
