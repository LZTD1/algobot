package handlersHolders

import (
	"tgbot/internal/contextHandlers/defaultHandler"
	"tgbot/internal/contextHandlers/textHandlers/sendingCookieState"
	"tgbot/internal/service"
	"tgbot/internal/stateMachine"
)

type SendingCookie struct {
	service service.Service
	state   stateMachine.StateMachine
}

func NewSendingCookie(service service.Service, state stateMachine.StateMachine) *SendingCookie {
	return &SendingCookie{service: service, state: state}
}

func (s SendingCookie) HolderType() stateMachine.Statement {
	return stateMachine.SendingCookie
}

func (s SendingCookie) GetHandlers() []defaultHandler.ContextHandler {
	return []defaultHandler.ContextHandler{
		sendingCookieState.NewRejectAction(s.state),
		sendingCookieState.NewSendingCookieAction(s.state, s.service),
	}
}
