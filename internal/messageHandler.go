package internal

import (
	"gopkg.in/telebot.v4"
	"tgbot/internal/config"
	handlers "tgbot/internal/contextHandlers"
	"tgbot/internal/contextHandlers/textHandlers"
	"tgbot/internal/service"
)

type MessageHandler struct {
	h []handlers.ContextHandler
	s service.Service
}

func NewMessageHandler(s service.Service) *MessageHandler {
	h := []handlers.ContextHandler{
		&textHandlers.Start{},
		textHandlers.NewMissingKids(s),
		textHandlers.NewSettings(s),
		textHandlers.NewMyGroups(s),
	}

	return &MessageHandler{h: h, s: s}
}

func (m *MessageHandler) Process(ctx telebot.Context) handlers.Response {
	uid := ctx.Message().Sender.ID

	if response := m.handleUserRegistration(uid); response != nil {
		return *response
	}

	for _, h := range m.h {
		if h.CanHandle(ctx) {
			return h.Process(ctx)
		}
	}
	return handlers.Response{Message: config.Incorrect, Keyboard: config.StartKeyboard}
}

func (m *MessageHandler) handleUserRegistration(uid int64) *handlers.Response {
	if m.s.IsUserRegistered(uid) == false {
		m.s.RegisterUser(uid)
		return &handlers.Response{Message: config.HelloWorld, Keyboard: config.StartKeyboard}
	}
	return nil
}
