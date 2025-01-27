package handlersHolders

import (
	"tgbot/internal/contextHandlers/callbackHandlers"
	"tgbot/internal/contextHandlers/defaultHandler"
	"tgbot/internal/service"
	"tgbot/internal/stateMachine"
)

type DefaultCBHolder struct {
	service service.Service
	state   stateMachine.StateMachine
}

func NewDefaultCBHolder(service service.Service, state stateMachine.StateMachine) *DefaultCBHolder {
	return &DefaultCBHolder{service: service, state: state}
}

func (d DefaultCBHolder) HolderType() stateMachine.Statement {
	return stateMachine.Default
}

func (d DefaultCBHolder) GetHandlers() []defaultHandler.ContextHandler {
	return []defaultHandler.ContextHandler{
		callbackHandlers.NewSetCookie(d.service, d.state),
		callbackHandlers.NewChangeNotification(d.service),
		callbackHandlers.NewRefreshGroups(d.service),
	}
}
