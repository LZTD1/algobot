package callbackHandlers

import (
	"algobot/internal_old/config"
	"algobot/internal_old/service"
	"algobot/internal_old/stateMachine"
	"gopkg.in/telebot.v4"
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

func (s *SetCookie) Process(ctx telebot.Context) error {
	s.state.SetStatement(ctx.Callback().Sender.ID, stateMachine.SendingCookie)

	return ctx.Send(config.SendingCookie, config.RejectKeyboard)
}
