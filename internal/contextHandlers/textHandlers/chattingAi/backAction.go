package chattingAi

import (
	"gopkg.in/telebot.v4"
	"tgbot/internal/config"
	"tgbot/internal/stateMachine"
)

type BackAction struct {
	state stateMachine.StateMachine
}

func NewBackAction(s stateMachine.StateMachine) *BackAction {
	return &BackAction{state: s}
}

func (b *BackAction) CanHandle(ctx telebot.Context) bool {
	if ctx.Message().Text == config.BackBtn.Text {
		return true
	}
	return false
}

func (b *BackAction) Process(ctx telebot.Context) error {
	b.state.SetStatement(ctx.Sender().ID, stateMachine.Default)
	return ctx.Send(config.StartText, config.StartKeyboard)
}
