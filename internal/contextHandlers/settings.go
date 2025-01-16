package contextHandlers

import (
	"gopkg.in/telebot.v4"
	"strings"
	"tgbot/internal/config"
	"tgbot/internal/service"
)

type Settings struct {
	svc service.Service
}

func NewSettings(svc service.Service) *Settings {
	return &Settings{svc: svc}
}

func (s *Settings) Message() string {
	return "Настройки"
}
func (s *Settings) Process(ctx telebot.Context) Response {
	uid := ctx.Message().Sender.ID

	c, err := s.svc.Cookie(uid)
	if err != nil {
		c = ""
	}
	n := s.svc.Notification(uid)

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
