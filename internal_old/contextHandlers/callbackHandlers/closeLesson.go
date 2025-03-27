package callbackHandlers

import (
	"algobot/internal_old/config"
	appError "algobot/internal_old/error"
	"algobot/internal_old/helpers"
	"algobot/internal_old/service"
	"fmt"
	"gopkg.in/telebot.v4"
	"strconv"
	"strings"
)

type CloseLesson struct {
	s service.Service
}

func NewCloseLesson(s service.Service) *CloseLesson {
	return &CloseLesson{s: s}
}

func (c CloseLesson) CanHandle(ctx telebot.Context) bool {
	if strings.HasPrefix(ctx.Callback().Data, "close_lesson_") {
		return true
	}
	return false
}

func (c CloseLesson) Process(ctx telebot.Context) error {
	data := strings.Split(strings.TrimPrefix(ctx.Callback().Data, "close_lesson_"), "_")
	if len(data) != 2 {
		return helpers.LogError(fmt.Errorf("%s : %w", ctx.Callback().Data, appError.NotEnoughArgs), ctx, "(1) Ошибка при анализе данных от кнопки!")
	}
	groupID, err := strconv.Atoi(data[0])
	if err != nil {
		return helpers.LogError(err, ctx, "(2) Ошибка при анализе данных от кнопки!")
	}
	lessionID, err := strconv.Atoi(data[1])
	if err != nil {
		return helpers.LogError(err, ctx, "(3) Ошибка при анализе данных от кнопки!")
	}
	err = c.s.CloseLesson(ctx.Callback().Sender.ID, groupID, lessionID)
	if err != nil {
		return helpers.LogError(err, ctx, "Ошибка при закрытии лекции")
	}

	return ctx.Send(config.SuccessfulChangeStatus)
}
