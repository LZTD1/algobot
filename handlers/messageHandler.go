package handlers

import (
	"gopkg.in/telebot.v4"
	"tgbot"
	"tgbot/storage"
)

type MessageHandler struct {
	storage storage.Storage
}

type HandlerResponse struct {
	Text string
}

func NewMessageHandler(s storage.Storage) *MessageHandler {
	return &MessageHandler{s}
}

func (m *MessageHandler) Process(ctx telebot.Context) HandlerResponse {
	return HandlerResponse{Text: tgbot.HelloWorld}
}
