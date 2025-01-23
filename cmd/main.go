package main

import (
	"log"
	"tgbot/internal"
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

	svc := mocks.NewMockService(nil)
	msgHandler := internal.NewMessageHandler(svc)

	b.Handle(tele.OnText, func(c tele.Context) error {
		resp := msgHandler.Process(c)
		return c.Send(resp.Message, resp.Keyboard)
	})

	b.Start()
}
