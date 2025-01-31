package helpers

import (
	"log"
	"math/rand"
	"time"
)

func GenerateRandomToken() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	token := make([]byte, 7)
	for i := range token {
		token[i] = charset[rng.Intn(len(charset))]
	}

	return string(token)
}

func LogWithRandomToken(err error) string {
	token := GenerateRandomToken()
	log.Printf("ERR: %s | %v", token, err)
	return token
}
