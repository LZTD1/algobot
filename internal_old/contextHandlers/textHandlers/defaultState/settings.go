package defaultState

import (
	"algobot/internal_old/config"
	"algobot/internal_old/helpers"
	"algobot/internal_old/service"
	"gopkg.in/telebot.v4"
	"strings"
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
func (s *Settings) Process(ctx telebot.Context) error {
	uid := ctx.Message().Sender.ID

	c, err := s.svc.Cookie(uid)
	if err != nil {
		return helpers.LogError(err, ctx, "Ошибка при формировании настроек (получение cookie) !")
	}
	n, err := s.svc.Notification(uid)
	if err != nil {
		return helpers.LogError(err, ctx, "Ошибка при формировании настроек (получение нотификаций) !")
	}

	return ctx.Send(GetMessageSettings(c, n), config.SettingsKeyboard)
}

func GetMessageSettings(c string, n bool) string {
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
