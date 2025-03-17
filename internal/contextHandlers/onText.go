package contextHandlers

import (
	"errors"
	"fmt"
	"gopkg.in/telebot.v4"
	"tgbot/internal/config"
	"tgbot/internal/contextHandlers/handlersHolders"
	"tgbot/internal/helpers"
	"tgbot/internal/service"
	"tgbot/internal/stateMachine"
)

type OnText struct {
	holders []handlersHolders.HandlersHolder
	state   stateMachine.StateMachine
	ai      service.AIService
}

func NewOnText(service service.Service, state stateMachine.StateMachine, ai service.AIService) *OnText {

	h := []handlersHolders.HandlersHolder{
		handlersHolders.NewDefaultHolders(service, state),
		handlersHolders.NewSendingCookie(service, state),
		handlersHolders.NewChattingAi(service, state, ai),
	}

	return &OnText{holders: h, state: state, ai: ai}
}

func (m *OnText) Handle(ctx telebot.Context) error {
	st := m.state.GetStatement(ctx.Sender().ID)

	holder := m.getHolder(st)
	if holder != nil {
		for _, h := range holder.GetHandlers() {
			if h.CanHandle(ctx) {
				return h.Process(ctx)
			}
		}
		m.state.SetStatement(ctx.Sender().ID, stateMachine.Default)
		return ctx.Send(config.Incorrect, config.StartKeyboard)
	}

	m.state.SetStatement(ctx.Sender().ID, stateMachine.Default)
	return helpers.LogError(errors.New(fmt.Sprintf("HolderNotFound(%v,%v)", st, ctx)), ctx, "Ошибка обработки запроса!\nНажмите - /start , для возвращения в меню")
}

func (m *OnText) getHolder(st stateMachine.Statement) handlersHolders.HandlersHolder {
	for _, holder := range m.holders {
		if holder.HolderType() == st {
			return holder
		}
	}
	return nil
}
