package main

import (
	middleware2 "gopkg.in/telebot.v4/middleware"
	"log"
	"tgbot/internal/contextHandlers"
	"tgbot/internal/middleware"
	"tgbot/tests/mocks"
	"time"

	tele "gopkg.in/telebot.v4"
)

const TOKEN = "6375608618:AAGtdaMkpj4SIJt495eNHOgw4oy5MZ_TIY4"

func main() {
	pref := tele.Settings{
		Token:  TOKEN,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer b.Stop()

	svc := mocks.NewMockService(make(map[int64]bool))

	regMid := middleware.NewRegister(svc)
	// TODO | Попробовать передавать синглетоны, а не инстансы обьектов

	msgHandler := contextHandlers.NewOnText(svc)

	b.Use(regMid.Middleware, middleware.MessageLogger, middleware2.AutoRespond())
	b.Handle(tele.OnText, msgHandler.Handle)

	b.Start()
}
