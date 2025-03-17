package handlersHolders

import (
	"tgbot/internal/contextHandlers/defaultHandler"
	"tgbot/internal/contextHandlers/textHandlers/chattingAi"
	"tgbot/internal/service"
	"tgbot/internal/stateMachine"
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
