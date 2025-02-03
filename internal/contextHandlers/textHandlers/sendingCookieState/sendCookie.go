package sendingCookieState

import (
	"gopkg.in/telebot.v4"
	"tgbot/internal/config"
	"tgbot/internal/helpers"
	"tgbot/internal/service"
	"tgbot/internal/stateMachine"
)

type SendingCookieAction struct {
	state   stateMachine.StateMachine
	service service.Service
}

func NewSendingCookieAction(state stateMachine.StateMachine, service service.Service) *SendingCookieAction {
	return &SendingCookieAction{state: state, service: service}
}

func (s SendingCookieAction) CanHandle(ctx telebot.Context) bool {
	if ctx.Message().Text != "Отменить действие" {
		return true
	}
	return false
}

func (s SendingCookieAction) Process(ctx telebot.Context) error {
	uid := ctx.Sender().ID
	cookie := ctx.Message().Text

	err := s.service.SetCookie(uid, cookie)
	if err != nil {
		return helpers.LogError(err, ctx, "Ошибка при установке Cookie!")
	}
	s.state.SetStatement(uid, stateMachine.Default)

	return ctx.Send(config.CookieSet, config.StartKeyboard)
}
