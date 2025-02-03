package middleware

import (
	"errors"
	"gopkg.in/telebot.v4"
	"tgbot/internal/config"
	appError "tgbot/internal/error"
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
		if err != nil && !errors.Is(err, appError.ErrNotFound) {
			return helpers.LogError(err, context, "Произошла ошибка при проверки регистрации!")
		}
		if reg == false {
			err := r.svc.RegisterUser(uid)
			if err != nil {
				return helpers.LogError(err, context, "Произошла ошибка при регистрации!")
			}
			context.Send(config.HelloWorld, config.StartKeyboard)
		}

		return next(context)
	}
}
