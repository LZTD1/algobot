package defaultState

import (
	"gopkg.in/telebot.v4"
	"strings"
	"tgbot/internal/service"
)

type AbsentKids struct {
	s service.Service
}

func NewAbsentKids(s service.Service) *AbsentKids {
	return &AbsentKids{s: s}
}

func (a AbsentKids) CanHandle(ctx telebot.Context) bool {
	if strings.HasPrefix(ctx.Message().Text, "/abs") {
		return true
	}
	return false
}

func (a AbsentKids) Process(ctx telebot.Context) error {
	// TODO
	//parsedTime, err := time.Parse("2006-01-02 15:04", ctx.Message().Payload)
	//if err != nil {
	//	return err
	//}
	//uid := ctx.Sender().ID
	//_, err = a.s.CurrentGroup(uid, parsedTime)
	//if err != nil {
	//	return err
	//}

	return ctx.Send("В разработке... ")
}
