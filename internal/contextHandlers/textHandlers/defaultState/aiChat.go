package defaultState

import (
	"gopkg.in/telebot.v4"
	"tgbot/internal/config"
	"tgbot/internal/service"
	"tgbot/internal/stateMachine"
)

type AIChat struct {
	svc   service.Service
	state stateMachine.StateMachine
}

func NewAIChat(s service.Service, state stateMachine.StateMachine) *AIChat {
	return &AIChat{svc: s, state: state}
}

func (a *AIChat) CanHandle(ctx telebot.Context) bool {
	if ctx.Message().Text == config.AIBtn.Text {
		return true
	}
	return false
}

func (a *AIChat) Process(ctx telebot.Context) error {
	a.state.SetStatement(ctx.Sender().ID, stateMachine.ChattingAI)
	return ctx.Send("Привет! Используй чат и клавиатуру для общения со мной!", config.AIKeyboard)
}
