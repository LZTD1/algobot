package middleware

import (
	"gopkg.in/telebot.v4"
	"tgbot/internal/config"
	"tgbot/internal/helpers"
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

		reg, err := r.svc.IsUserRegistered(uid)
		if err != nil {
			t := helpers.LogWithRandomToken(err)
			context.Send(t + " | Произошла ошибка при проверки регистрации!")
		}
		if reg == false {
			r.svc.RegisterUser(uid)
			context.Send(config.HelloWorld, config.StartKeyboard)
		}

		return next(context)
	}
}
