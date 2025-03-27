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

type OpenLesson struct {
	s service.Service
}

func NewOpenLesson(s service.Service) *OpenLesson {
	return &OpenLesson{s: s}
}

func (c OpenLesson) CanHandle(ctx telebot.Context) bool {
	if strings.HasPrefix(ctx.Callback().Data, "open_lesson_") {
		return true
	}
	return false
}

func (c OpenLesson) Process(ctx telebot.Context) error {
	data := strings.Split(strings.TrimPrefix(ctx.Callback().Data, "open_lesson_"), "_")
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
	err = c.s.OpenLesson(ctx.Callback().Sender.ID, groupID, lessionID)
	if err != nil {
		return helpers.LogError(err, ctx, "Ошибка при открытии лекции")
	}

	return ctx.Send(config.SuccessfulChangeStatus)
}
