package contextHandlers

import (
	"gopkg.in/telebot.v4"
	"tgbot/internal/config"
	"tgbot/internal/contextHandlers/defaultHandler"
	"tgbot/internal/contextHandlers/textHandlers"
	"tgbot/internal/service"
)

type OnText struct {
	h []defaultHandler.ContextHandler
	s service.Service
}

func NewOnText(s service.Service) *OnText {
	h := []defaultHandler.ContextHandler{
		&textHandlers.Start{},
		textHandlers.NewMissingKids(s),
		textHandlers.NewSettings(s),
		textHandlers.NewMyGroups(s),
	}

	return &OnText{h: h, s: s}
}

func (m *OnText) Handle(ctx telebot.Context) error {
	for _, h := range m.h {
		if h.CanHandle(ctx) {
			return h.Process(ctx)
		}
	}
	return ctx.Send(config.Incorrect, config.StartKeyboard)
}
