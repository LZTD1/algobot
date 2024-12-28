package handlers

import (
	"gopkg.in/telebot.v4"
	"strings"
	"tgbot/config"
	"tgbot/storage"
)

type SettingsHandler struct {
	storage storage.Storage
}

func NewSettingsHandler(storage storage.Storage) *SettingsHandler {
	return &SettingsHandler{storage}
}

func (s *SettingsHandler) Message() string {
	return "Настройки"
}
func (s *SettingsHandler) Process(ctx telebot.Context) Response {
	uid := ctx.Message().Sender.ID

	c, err := s.storage.Cookie(uid)
	if err != nil {
		c = ""
	}
	n := s.storage.Notification(uid)

	return Response{
		Message: getMessageSettings(c, n),
	}
}

func getMessageSettings(c string, n bool) string {
	msg := strings.Builder{}
	msg.WriteString(config.Settings)
	msg.WriteString("\n\n")
	msg.WriteString(config.Cookie)
	if c != "" {
		msg.WriteString(config.SetParam)
	} else {
		msg.WriteString(config.NotSetParam)
	}
	msg.WriteString("\n")
	msg.WriteString(config.ChatNotifications)
	if n {
		msg.WriteString(config.SetParam)
	} else {
		msg.WriteString(config.NotSetParam)
	}

	return msg.String()
}
