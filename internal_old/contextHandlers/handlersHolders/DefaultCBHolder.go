package handlersHolders

import (
	"algobot/internal_old/contextHandlers/callbackHandlers"
	"algobot/internal_old/contextHandlers/defaultHandler"
	"algobot/internal_old/service"
	"algobot/internal_old/stateMachine"
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
		callbackHandlers.NewCloseLesson(d.service),
		callbackHandlers.NewOpenLesson(d.service),
		callbackHandlers.NewGetCredentials(d.service),
	}
}
