package handlersHolders

import (
	"algobot/internal_old/contextHandlers/defaultHandler"
	"algobot/internal_old/contextHandlers/textHandlers/chattingAi"
	"algobot/internal_old/service"
	"algobot/internal_old/stateMachine"
)

type ChattingAi struct {
	state   stateMachine.StateMachine
	service service.Service
	ai      service.AIService
}

func NewChattingAi(service service.Service, state stateMachine.StateMachine, ai service.AIService) *ChattingAi {
	return &ChattingAi{
		service: service,
		state:   state,
		ai:      ai,
	}
}

func (c *ChattingAi) HolderType() stateMachine.Statement {
	return stateMachine.ChattingAI
}

func (c *ChattingAi) GetHandlers() []defaultHandler.ContextHandler {
	return []defaultHandler.ContextHandler{
		chattingAi.NewBackAction(c.state),
		chattingAi.NewClearHistory(c.ai),
		chattingAi.NewAnyMessage(c.ai),
	}
}
