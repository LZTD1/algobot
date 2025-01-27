package handlersHolders

import (
	"tgbot/internal/contextHandlers/defaultHandler"
	"tgbot/internal/stateMachine"
)

type HandlersHolder interface {
	HolderType() stateMachine.Statement
	GetHandlers() []defaultHandler.ContextHandler
}
