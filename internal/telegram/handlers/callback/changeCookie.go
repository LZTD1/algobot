package callback

import (
	"algobot/internal/domain/telegram/keyboards"
	"algobot/internal/lib/fsm"
	"gopkg.in/telebot.v4"
)

type StateChanger interface {
	SetState(uid int64, state fsm.State)
}

func NewChangeCookie(stateChanger StateChanger) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		uid := ctx.Sender().ID

		stateChanger.SetState(uid, fsm.SendingCookie)

		return ctx.Send(
			"Отправьте мне свои cookie 🍪\nИнструкция: https://telegra.ph/Kak-dobavit-v-bota-svoi-Cookie-02-05",
			keyboards.RejectKeyboard(),
		)
	}
}
