package helpers

import (
	"fmt"
	"gopkg.in/telebot.v4"
	"log"
	"math/rand"
	"time"
)

func token() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	token := make([]byte, 7)
	for i := range token {
		token[i] = charset[rng.Intn(len(charset))]
	}

	return string(token)
}

func LogError(err error, ctx telebot.Context, reason string) error {
	token := token()

	log.Printf("ERR: %s | %s\n", token, err.Error())
	if reason == "" {
		reason = "Произошла ошибка, классифицировать не удалось :("
	}

	return ctx.Send(fmt.Sprintf("<b>[%s]</b> %s", token, reason), telebot.ModeHTML)
}
