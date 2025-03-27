package defaultState

import (
	"algobot/internal_old/config"
	"algobot/internal_old/service"
	"algobot/internal_old/stateMachine"
	"gopkg.in/telebot.v4"
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
