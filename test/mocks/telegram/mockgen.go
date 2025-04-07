package mocks

//go:generate mockgen -destination=./context_mock.go -package=mocks gopkg.in/telebot.v4 Context
//go:generate mockgen -destination=./api_mock.go -package=mocks gopkg.in/telebot.v4 API
