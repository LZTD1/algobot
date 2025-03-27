package middleware

import (
	"algobot/internal_old/config"
	appError "algobot/internal_old/error"
	"algobot/internal_old/helpers"
	"algobot/internal_old/service"
	"errors"
	"gopkg.in/telebot.v4"
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
