package defaultState

import (
	"fmt"
	"gopkg.in/telebot.v4"
	"strings"
	"tgbot/internal/service"
	"time"
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
	parsedTime, err := time.Parse("2006-01-02 15:04", ctx.Message().Payload)
	if err != nil {
		return err
	}
	fmt.Println(parsedTime)
	uid := ctx.Sender().ID
	group, err := a.s.CurrentGroup(uid, parsedTime)
	if err != nil {
		return err
	}
	fmt.Println("")
	fmt.Println(group)

	return ctx.Send(strings.Join(group.MissingKids, "\n"))
}
