package textHandlers

import (
	"gopkg.in/telebot.v4"
	"strings"
	"tgbot/internal/config"
	"tgbot/internal/contextHandlers"
	"tgbot/internal/service"
)

type Settings struct {
	svc service.Service
}

func NewSettings(svc service.Service) *Settings {
	return &Settings{svc: svc}
}

func (s *Settings) CanHandle(ctx telebot.Context) bool {
	if ctx.Message().Text == "Настройки" {
		return true
	}
	return false
}
func (s *Settings) Process(ctx telebot.Context) contextHandlers.Response {
	uid := ctx.Message().Sender.ID

	c, err := s.svc.Cookie(uid)
	if err != nil {
		c = ""
	}
	n := s.svc.Notification(uid)

	return contextHandlers.Response{
		Message:  getMessageSettings(c, n),
		Keyboard: config.SettingsKeyboard,
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
