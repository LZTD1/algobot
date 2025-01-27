package contextHandlers

import (
	"fmt"
	"gopkg.in/telebot.v4"
	"tgbot/internal/config"
	"tgbot/internal/contextHandlers/handlersHolders"
	"tgbot/internal/service"
	"tgbot/internal/stateMachine"
)

type OnText struct {
	holders []handlersHolders.HandlersHolder
	state   stateMachine.StateMachine
}

func NewOnText(service service.Service, state stateMachine.StateMachine) *OnText {

	h := []handlersHolders.HandlersHolder{
		handlersHolders.NewDefaultHolders(service, state),
		handlersHolders.NewSendingCookie(service, state),
	}

	return &OnText{holders: h, state: state}
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
	return fmt.Errorf("не найден холдер для данного state")
}

func (m *OnText) getHolder(st stateMachine.Statement) handlersHolders.HandlersHolder {
	for _, holder := range m.holders {
		if holder.HolderType() == st {
			return holder
		}
	}
	return nil
}
