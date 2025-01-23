package contextHandlers

import (
	"gopkg.in/telebot.v4"
	"tgbot/internal/config"
	"tgbot/internal/contextHandlers/callbackHandlers"
	"tgbot/internal/contextHandlers/defaultHandler"
	"tgbot/internal/service"
	"tgbot/internal/stateMachine"
)

type OnCallback struct {
	h     []defaultHandler.ContextHandler
	s     service.Service
	state stateMachine.StateMachine
}

func NewOnCallback(s service.Service, machine stateMachine.StateMachine) *OnCallback {
	h := []defaultHandler.ContextHandler{
		callbackHandlers.NewSetCookie(s, machine),
		callbackHandlers.NewChangeNotification(s),
		callbackHandlers.NewRefreshGroups(s),
	}

	return &OnCallback{h: h, s: s}
}

func (h *OnCallback) Process(ctx telebot.Context) error {
	uid := ctx.Callback().Sender.ID

	if response := h.handleUserRegistration(uid); response != nil {
		return ctx.Send(response.Message, response.Keyboard)
	}

	for _, handlers := range h.h {
		if handlers.CanHandle(ctx) {
			return handlers.Process(ctx)
		}
	}
	return ctx.Send(config.Incorrect, config.StartKeyboard)
}

func (h *OnCallback) handleUserRegistration(uid int64) *defaultHandler.Response {
	if h.s.IsUserRegistered(uid) == false {
		h.s.RegisterUser(uid)
		return &defaultHandler.Response{Message: config.HelloWorld, Keyboard: config.StartKeyboard}
	}
	return nil
}
