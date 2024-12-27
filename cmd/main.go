package main

import (
	"log"
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

	b.Handle(tele.OnText, func(c tele.Context) error {
		return c.Send("Hello!")
	})

	b.Start()
}
