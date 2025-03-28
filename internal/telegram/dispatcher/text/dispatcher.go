package text

import (
	"algobot/internal/lib/fsm"
	tele "gopkg.in/telebot.v4"
	"log/slog"
)

type Dispatcher struct {
	log     *slog.Logger
	handler Handlers
}

type Handler interface {
	Handle(c tele.Context) error
}

type Handlers map[fsm.State]Handler

func NewDispatcher(log *slog.Logger) *Dispatcher {

	return &Dispatcher{log: log, handler: make(Handlers)}
}

func (d *Dispatcher) Register(state fsm.State, handler Handler) {
	d.handler[state] = handler
}

func (d *Dispatcher) GetHandlers(state fsm.State) Handler {
	if val, ok := d.handler[state]; ok {
		return val
	}

	return d.handler[fsm.Default] // TODO : refactor default val, change to ret error
}
