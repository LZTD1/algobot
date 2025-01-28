package handlersHolders

import (
	"tgbot/internal/contextHandlers/defaultHandler"
	"tgbot/internal/contextHandlers/textHandlers/defaultState"
	"tgbot/internal/service"
	"tgbot/internal/stateMachine"
)

type DefaultHolders struct {
	service service.Service
	state   stateMachine.StateMachine
}

func NewDefaultHolders(service service.Service, state stateMachine.StateMachine) *DefaultHolders {
	return &DefaultHolders{service: service, state: state}
}

func (d DefaultHolders) HolderType() stateMachine.Statement {
	return stateMachine.Default
}

func (d DefaultHolders) GetHandlers() []defaultHandler.ContextHandler {
	return []defaultHandler.ContextHandler{
		&defaultState.Start{},
		defaultState.NewMissingKids(d.service),
		defaultState.NewSettings(d.service),
		defaultState.NewMyGroups(d.service),
		defaultState.NewAbsentKids(d.service),
	}
}
