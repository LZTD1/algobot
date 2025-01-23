package contextHandlers

import (
	"gopkg.in/telebot.v4"
	"tgbot/internal/config"
	"tgbot/internal/contextHandlers/defaultHandler"
	"tgbot/internal/contextHandlers/textHandlers"
	"tgbot/internal/service"
)

type OnText struct {
	h []defaultHandler.ContextHandler
	s service.Service
}

func NewOnText(s service.Service) *OnText {
	h := []defaultHandler.ContextHandler{
		&textHandlers.Start{},
		textHandlers.NewMissingKids(s),
		textHandlers.NewSettings(s),
		textHandlers.NewMyGroups(s),
	}

	return &OnText{h: h, s: s}
}

func (m *OnText) Process(ctx telebot.Context) defaultHandler.Response {
	uid := ctx.Message().Sender.ID

	if response := m.handleUserRegistration(uid); response != nil {
		return *response
	}

	for _, h := range m.h {
		if h.CanHandle(ctx) {
			return h.Process(ctx)
		}
	}
	return defaultHandler.Response{Message: config.Incorrect, Keyboard: config.StartKeyboard}
}

func (m *OnText) handleUserRegistration(uid int64) *defaultHandler.Response {
	if m.s.IsUserRegistered(uid) == false {
		m.s.RegisterUser(uid)
		return &defaultHandler.Response{Message: config.HelloWorld, Keyboard: config.StartKeyboard}
	}
	return nil
}
