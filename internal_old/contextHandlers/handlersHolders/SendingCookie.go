package handlersHolders

import (
	"algobot/internal_old/contextHandlers/defaultHandler"
	"algobot/internal_old/contextHandlers/textHandlers/sendingCookieState"
	"algobot/internal_old/service"
	"algobot/internal_old/stateMachine"
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
