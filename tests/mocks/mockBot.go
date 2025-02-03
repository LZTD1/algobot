package mocks

import (
	"fmt"
	"gopkg.in/telebot.v4"
)

type MockBot struct {
	telebot.API
	Calls []string
}

func (b *MockBot) Send(to telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error) {
	b.Calls = append(b.Calls, fmt.Sprintf("Send(%s) = %s", to.Recipient(), what))
	return nil, nil
}
