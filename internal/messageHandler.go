package internal

import (
	"gopkg.in/telebot.v4"
	"tgbot/internal/config"
	handlers "tgbot/internal/contextHandlers"
	"tgbot/internal/service"
)

type MessageHandler struct {
	h map[string]handlers.ContextHandler
	s service.Service
}

func NewMessageHandler(s service.Service) *MessageHandler {
	h := []handlers.ContextHandler{
		&handlers.Start{},
		handlers.NewMissingKids(s),
		handlers.NewSettings(s),
		handlers.NewMyGroups(s),
	}

	return &MessageHandler{h: getHandlerMap(h), s: s}
}

func (m *MessageHandler) Process(ctx telebot.Context) handlers.Response {
	uid := ctx.Message().Sender.ID
	msg := ctx.Message().Text

	if response := m.handleUserRegistration(uid); response != nil {
		return *response
	}

	v, ok := m.h[msg]
	if !ok {
		return handlers.Response{Message: config.Incorrect, Keyboard: config.StartKeyboard}
	}
	return v.Process(ctx)
}

func (m *MessageHandler) handleUserRegistration(uid int64) *handlers.Response {
	if m.s.IsUserRegistered(uid) == false {
		m.s.RegisterUser(uid)
		return &handlers.Response{Message: config.HelloWorld, Keyboard: config.StartKeyboard}
	}
	return nil
}

func getHandlerMap(h []handlers.ContextHandler) map[string]handlers.ContextHandler {
	m := make(map[string]handlers.ContextHandler)
	for _, v := range h {
		m[v.Message()] = v
	}
	return m
}
