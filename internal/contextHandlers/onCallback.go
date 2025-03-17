package contextHandlers

import (
	"errors"
	"fmt"
	"gopkg.in/telebot.v4"
	"strings"
	"tgbot/internal/config"
	"tgbot/internal/contextHandlers/handlersHolders"
	"tgbot/internal/helpers"
	"tgbot/internal/service"
	"tgbot/internal/stateMachine"
)

type OnCallback struct {
	holders []handlersHolders.HandlersHolder
	state   stateMachine.StateMachine
}

func NewOnCallback(service service.Service, state stateMachine.StateMachine) *OnCallback {
	h := []handlersHolders.HandlersHolder{
		handlersHolders.NewDefaultCBHolder(service, state),
	}

	return &OnCallback{holders: h, state: state}
}

func (h *OnCallback) Handle(ctx telebot.Context) error {
	st := h.state.GetStatement(ctx.Sender().ID)
	holder := h.getHolder(st)

	if holder != nil {

		if strings.HasPrefix(ctx.Callback().Data, "\f") {
			ctx.Callback().Data = strings.TrimPrefix(ctx.Callback().Data, "\f")
		}

		for _, hand := range holder.GetHandlers() {
			if hand.CanHandle(ctx) {
				return hand.Process(ctx)
			}
		}
		h.state.SetStatement(ctx.Sender().ID, stateMachine.Default)
		return ctx.Send(config.Incorrect, config.StartKeyboard)
	}

	h.state.SetStatement(ctx.Sender().ID, stateMachine.Default)
	return helpers.LogError(errors.New(fmt.Sprintf("HolderNotFound(%v,%v)", st, ctx)), ctx, "Ошибка обработки запроса!\nНажмите - /start , для возвращения в меню")
}

func (h *OnCallback) getHolder(st stateMachine.Statement) handlersHolders.HandlersHolder {
	for _, holder := range h.holders {
		if holder.HolderType() == st {
			return holder
		}
	}
	return nil
}
