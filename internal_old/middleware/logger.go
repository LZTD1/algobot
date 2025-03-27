package middleware

import (
	tele "gopkg.in/telebot.v4"
	"log"
)

func MessageLogger(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if c.Callback() != nil {
			log.Printf("[CLB] %s, запрос \"%s\"\n", c.Callback().Sender.FirstName+c.Callback().Sender.LastName, c.Callback().Data)
		} else {
			log.Printf("[MSG] %s, запрос \"%s\"\n", c.Message().Sender.FirstName+c.Message().Sender.LastName, c.Message().Text)
		}

		return next(c)
	}
}
