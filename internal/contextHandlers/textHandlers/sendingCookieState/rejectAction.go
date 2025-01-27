package sendingCookieState

import (
	"gopkg.in/telebot.v4"
	"tgbot/internal/config"
	"tgbot/internal/stateMachine"
)

type RejectAction struct {
	state stateMachine.StateMachine
}

func NewRejectAction(state stateMachine.StateMachine) *RejectAction {
	return &RejectAction{state: state}
}

func (r RejectAction) CanHandle(ctx telebot.Context) bool {
	if ctx.Message().Text == "Отменить действие" {
		return true
	}
	return false
}
func (r RejectAction) Process(ctx telebot.Context) error {
	r.state.SetStatement(ctx.Sender().ID, stateMachine.Default)
	return ctx.Send(config.StartText, config.StartKeyboard)
}
