package handlersHolders

import (
	"algobot/internal_old/contextHandlers/defaultHandler"
	"algobot/internal_old/contextHandlers/textHandlers/defaultState"
	"algobot/internal_old/service"
	"algobot/internal_old/stateMachine"
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
		defaultState.NewStartWithPayload(d.service),
		defaultState.NewAIChat(d.service, d.state),
	}
}
