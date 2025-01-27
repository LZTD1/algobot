package mocks

import "gopkg.in/telebot.v4"

type SentMessage struct {
	What interface{}
	Opts []interface{}
}

type MockContext struct {
	userId      int64
	userMessage string
	unixTime    int64
	telebot.Context

	SentMessages []SentMessage
}

func (m *MockContext) Sender() *telebot.User {
	return &telebot.User{
		ID: m.userId,
	}
}

func (m *MockContext) Send(what interface{}, opts ...interface{}) error {
	m.SentMessages = append(m.SentMessages, SentMessage{
		What: what,
		Opts: opts,
	})
	return nil
}
func (m *MockContext) Edit(what interface{}, opts ...interface{}) error {
	m.SentMessages = append(m.SentMessages, SentMessage{
		What: what,
		Opts: opts,
	})
	return nil
}

func (m *MockContext) Message() *telebot.Message {
	return &telebot.Message{
		Sender: &telebot.User{
			ID: m.userId,
		},
		Text:     m.userMessage,
		Unixtime: m.unixTime,
	}
}

func (m *MockContext) SetUserMessage(uid int64, msg string) {
	m.userId = uid
	m.userMessage = msg
}

func (m *MockContext) SetUserMessageWithTime(uid int64, msg string, unix int64) {
	m.userId = uid
	m.userMessage = msg
	m.unixTime = unix
}
func (m *MockContext) Callback() *telebot.Callback {
	return &telebot.Callback{
		Sender: &telebot.User{
			ID: m.userId,
		},
		Data: m.userMessage,
	}
}
