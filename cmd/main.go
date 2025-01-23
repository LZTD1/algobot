package main

import (
	"log"
	"tgbot/internal/contextHandlers"
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

	// TODO | Попробовать передавать синглетоны, а не инстансы обьектов
	// TODO | Add middleware to check user registration
	// TODO | And logging every message

	msgHandler := contextHandlers.NewOnText(svc)

	b.Handle(tele.OnText, msgHandler.Handle)

	b.Start()
}
