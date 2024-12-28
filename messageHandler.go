package tgbot

import (
	"gopkg.in/telebot.v4"
	"tgbot/config"
	"tgbot/handlers"
	"tgbot/storage"
)

type MessageHandler struct {
	handlers map[string]handlers.Handler
	store    storage.Storage
}

func NewMessageHandler(s storage.Storage) *MessageHandler {
	h := []handlers.Handler{
		&handlers.SettingsHandler{},
		handlers.NewSettingsHandler(s),
	}

	return &MessageHandler{handlers: getHandlerMap(h), store: s}
}

func (m *MessageHandler) Process(ctx telebot.Context) handlers.Response {
	uid := ctx.Message().Sender.ID
	msg := ctx.Message().Text

	if response := m.handleUserRegistration(uid); response != nil {
		return *response
	}

	v, ok := m.handlers[msg]
	if !ok {
		return handlers.Response{Message: config.Incorrect}
	}
	return v.Process(ctx)
}

func (m *MessageHandler) handleUserRegistration(uid int64) *handlers.Response {
	_, err := m.store.User(uid)
	if err != nil {
		m.store.RegisterUser(uid)
		return &handlers.Response{Message: config.HelloWorld}
	}
	return nil
}

func getHandlerMap(h []handlers.Handler) map[string]handlers.Handler {
	m := make(map[string]handlers.Handler)
	for _, v := range h {
		m[v.Message()] = v
	}
	return m
}
