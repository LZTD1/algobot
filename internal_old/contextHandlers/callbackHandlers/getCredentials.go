package callbackHandlers

import (
	"algobot/internal_old/helpers"
	"algobot/internal_old/service"
	"fmt"
	"gopkg.in/telebot.v4"
	"strconv"
	"strings"
)

type GetCredentials struct {
	s service.Service
}

func NewGetCredentials(s service.Service) *GetCredentials {
	return &GetCredentials{s: s}
}

func (c GetCredentials) CanHandle(ctx telebot.Context) bool {
	if strings.HasPrefix(ctx.Callback().Data, "get_creds_") {
		return true
	}
	return false
}

func (c GetCredentials) Process(ctx telebot.Context) error {
	data := strings.TrimPrefix(ctx.Callback().Data, "get_creds_")

	groupID, err := strconv.Atoi(data)
	if err != nil {
		return helpers.LogError(err, ctx, "(1) Ошибка при анализе данных от кнопки!")
	}

	creds, err := c.s.AllCredentials(ctx.Callback().Sender.ID, groupID)
	if err != nil {
		return helpers.LogError(err, ctx, "(2) Ошибка при анализе данных от кнопки!")
	}

	return ctx.Send(getCreds(creds))
}

func getCreds(creds map[string]string) string {
	sb := strings.Builder{}
	sb.WriteString("Учетные записи детей:\n")
	for key, val := range creds {
		sb.WriteString(fmt.Sprintf("\n%v [%v]", key, val))
	}
	return sb.String()
}
