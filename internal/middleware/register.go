package middleware

import (
	"gopkg.in/telebot.v4"
	"tgbot/internal/config"
	"tgbot/internal/service"
)

type Register struct {
	svc service.Service
}

func NewRegister(svc service.Service) *Register {
	return &Register{svc: svc}
}

func (r *Register) Middleware(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		uid := context.Sender().ID

		if r.svc.IsUserRegistered(uid) == false {
			r.svc.RegisterUser(uid)
			context.Send(config.HelloWorld, config.StartKeyboard)
		}

		return next(context)
	}
}
