package handlersHolders

import (
	"algobot/internal_old/contextHandlers/defaultHandler"
	"algobot/internal_old/stateMachine"
)

type HandlersHolder interface {
	HolderType() stateMachine.Statement
	GetHandlers() []defaultHandler.ContextHandler
}
