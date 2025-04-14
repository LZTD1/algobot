package text

import (
	"algobot/internal/domain/telegram/keyboards"
	"algobot/internal/lib/fsm"
	"gopkg.in/telebot.v4"
)

type SetStater interface {
	SetState(uid int64, state fsm.State)
}

func NewStart(setStater SetStater) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		setStater.SetState(ctx.Sender().ID, fsm.Default)

		return ctx.Send("Открыто главное меню:", keyboards.Start())
	}
}
