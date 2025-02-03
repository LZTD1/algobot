package schedulers

import (
	"fmt"
	"gopkg.in/telebot.v4"
	"strings"
	"tgbot/internal/models"
	"tgbot/internal/service"
)

type Message struct {
	b   telebot.API
	svc service.Service
}

func NewMessage(b telebot.API, svc service.Service) *Message {
	return &Message{b: b, svc: svc}
}

func (m Message) Schedule() {
	uids, _ := m.svc.UserUidsByNotif(true)
	for _, uid := range uids {
		allMessages, _ := m.svc.NewMessageByUID(uid)
		for _, msg := range allMessages {
			m.b.Send(RecipientUser{uid}, getMsg(msg))
		}
	}
}

func getMsg(msg models.Message) string {
	sb := strings.Builder{}
	sb.WriteString("🔔 Новое сообщение\n\n")
	sb.WriteString(fmt.Sprintf("От: %s\n", msg.From))
	sb.WriteString(fmt.Sprintf("Тема: %s\n", msg.Theme))
	sb.WriteString(fmt.Sprintf("Ссылка: %s\n\n", msg.Link))
	sb.WriteString(fmt.Sprintf("<%s>\n", strings.Repeat("=", 5)))
	sb.WriteString(msg.Content)
	sb.WriteString(fmt.Sprintf("\n<%s>", strings.Repeat("=", 5)))

	return sb.String()
}

type RecipientUser struct {
	recipient string
}

func (r RecipientUser) Recipient() string {
	return r.recipient
}
