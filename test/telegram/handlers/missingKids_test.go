package test

import (
	"algobot/internal/domain/models"
	"algobot/internal/services/groups"
	"algobot/internal/telegram/handlers/text"
	mocks2 "algobot/test/mocks"
	mocks3 "algobot/test/mocks/telegram"
	mocks "algobot/test/mocks/telegram/handlers"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gopkg.in/telebot.v4"
	"testing"
)

func TestMissingKids(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	actualGroup := mocks.NewMockActualGroup(ctrl)
	log := mocks2.NewMockLogger()
	mctx := mocks3.NewMockContext(ctrl)

	handler := text.NewMissingKids(log, actualGroup)

	mctx.EXPECT().Get(gomock.Any()).Return("").AnyTimes()
	mctx.EXPECT().Sender().Return(&telebot.User{ID: 1}).AnyTimes()
	t.Run("happy path", func(t *testing.T) {
		gomock.InOrder(
			actualGroup.EXPECT().CurrentGroup(int64(1), gomock.Any(), "").Return(acGroupasset, nil).Times(1),
			mctx.EXPECT().Send("Группа: title\nЛекция: lesson\n\nОбщее число детей: 3\nОтсутствуют: 2\n\n```Отсутствующие\n1 (Уже 2 занятие)\n1\n```", telebot.ModeMarkdown).Return(nil).Times(1),
		)
		err := handler(mctx)
		assert.NoError(t, err)
	})
	t.Run("no groups found", func(t *testing.T) {
		gomock.InOrder(
			actualGroup.EXPECT().CurrentGroup(int64(1), gomock.Any(), "").Return(models.CurrentGroup{}, groups.ErrNoGroups).Times(1),
			mctx.EXPECT().Send("В данный момент, никакой группы не найдено!").Return(nil).Times(1),
		)
		err := handler(mctx)
		assert.NoError(t, err)
	})
	t.Run("ErrNotValidCookie", func(t *testing.T) {
		gomock.InOrder(
			actualGroup.EXPECT().CurrentGroup(int64(1), gomock.Any(), "").Return(models.CurrentGroup{}, groups.ErrNotValidCookie).Times(1),
			mctx.EXPECT().Send("Вам необходимо установить свои cookie!").Return(nil).Times(1),
		)
		err := handler(mctx)
		assert.NoError(t, err)
	})
	t.Run("CurrentGroup return err", func(t *testing.T) {
		errExp := errors.New("err")
		gomock.InOrder(
			actualGroup.EXPECT().CurrentGroup(int64(1), gomock.Any(), "").Return(models.CurrentGroup{}, errExp).Times(1),
		)
		err := handler(mctx)
		assert.ErrorIs(t, err, errExp)
	})
}

var acGroupasset = models.CurrentGroup{
	GroupID:  1,
	Title:    "title",
	Lesson:   "lesson",
	LessonID: 1,
	Kids: []string{
		"1",
		"2",
		"3",
	},
	MissingKids: []models.MissingKid{
		{
			Fullname: "1",
			KidID:    1,
			Count:    2,
		},
		{
			Fullname: "1",
			KidID:    2,
			Count:    1,
		},
	},
}
