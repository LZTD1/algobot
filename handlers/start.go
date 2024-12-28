package handlers

import (
	"gopkg.in/telebot.v4"
	"tgbot/config"
)

type StartHandler struct {
}

func (h *StartHandler) Message() string {
	return "/start"
}
func (h *StartHandler) Process(ctx telebot.Context) Response {
	return Response{
		Message: config.StartText,
	}
}
