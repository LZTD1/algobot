package schedulers

import (
	"algobot/internal_old/models"
	"algobot/internal_old/service"
	"fmt"
	"gopkg.in/telebot.v4"
	"log"
	"strconv"
	"strings"
)

type Message struct {
	b   telebot.API
	svc service.Service
}

func NewMessage(b telebot.API, svc service.Service) *Message {
	return &Message{b: b, svc: svc}
}

func (m Message) Schedule() {
	users, err := m.svc.UsersByNotif(true)
	if err != nil {
		log.Println(err)
	}
	for _, user := range users {
		allMessages, err := m.svc.NewMessageByUID(user.UID)
		if err != nil {
			log.Println(err)
		}
		for _, msg := range allMessages {
			if msg.Type == "img" {
				p := &telebot.Photo{File: telebot.FromURL(msg.Content), Caption: getMsg(msg)}
				m.b.Send(RecipientUser{strconv.FormatInt(user.UID, 10)}, p, telebot.ModeMarkdown, telebot.NoPreview)
			} else {
				m.b.Send(RecipientUser{strconv.FormatInt(user.UID, 10)}, getMsg(msg), telebot.ModeMarkdown, telebot.NoPreview)
			}
		}
	}
}

func getMsg(msg models.Message) string {
	sb := strings.Builder{}
	sb.WriteString("🔔 Новое сообщение\n\n")
	sb.WriteString(fmt.Sprintf("От: %s\n", msg.From))
	sb.WriteString(fmt.Sprintf("Тема: %s\n", msg.Theme))
	sb.WriteString(fmt.Sprintf("Ссылка: %s\n\n", msg.Link))
	if msg.Type != "img" {
		sb.WriteString("```Сообщение:\n")
		sb.WriteString(msg.Content)
		sb.WriteString("\n```")
	}
	return sb.String()
}

type RecipientUser struct {
	ID string
}

func (r RecipientUser) Recipient() string {
	return r.ID
}
