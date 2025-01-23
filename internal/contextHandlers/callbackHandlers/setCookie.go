package callbackHandlers

import (
	"gopkg.in/telebot.v4"
	"tgbot/internal/config"
	"tgbot/internal/contextHandlers/defaultHandler"
	"tgbot/internal/service"
	"tgbot/internal/stateMachine"
)

type SetCookie struct {
	svc   service.Service
	state stateMachine.StateMachine
}

func NewSetCookie(svc service.Service, state stateMachine.StateMachine) *SetCookie {
	return &SetCookie{svc: svc, state: state}
}

func (s *SetCookie) CanHandle(ctx telebot.Context) bool {
	if ctx.Callback().Data == "set_cookie" {
		return true
	}

	return false
}

func (s *SetCookie) Process(ctx telebot.Context) defaultHandler.Response {
	s.state.SetStatement(ctx.Callback().Sender.ID, stateMachine.SendingCookie)
	return defaultHandler.Response{
		Message:  config.SendingCookie,
		Keyboard: config.RejectKeyboard,
	}
}
