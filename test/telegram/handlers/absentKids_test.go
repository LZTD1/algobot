package test

import (
	"algobot/internal/domain/models"
	"algobot/internal/services/groups"
	"algobot/internal/telegram/handlers/text"
	mocks3 "algobot/test/mocks"
	mocks2 "algobot/test/mocks/telegram"
	mocks "algobot/test/mocks/telegram/handlers"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	tele "gopkg.in/telebot.v4"
	"testing"
)

func TestNewAbsentKids(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mctx := mocks2.NewMockContext(ctrl)
	log := mocks3.NewMockLogger()
	agroup := mocks.NewMockActualGroup(ctrl)

	handler := text.NewAbsentKids(agroup, log)

	mctx.EXPECT().Get(gomock.Any()).Return("").AnyTimes()
	mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).AnyTimes()
	t.Run("Happy path", func(t *testing.T) {
		mctx.EXPECT().Message().Return(&tele.Message{Text: "/abs 2025-04-06 14:44"}).Times(1)
		agroup.EXPECT().CurrentGroup(int64(1), gomock.Any(), "").Return(grAsset, nil).Times(1)
		mctx.EXPECT().Reply("Группа: title\nЛекция: lesson\n\nОбщее число детей: 3\nОтсутствуют: 2\n\n```Отсутствующие\n1 (Уже 2 занятие)\n1\n```", tele.ModeMarkdown).Times(1)

		err := handler(mctx)
		assert.NoError(t, err)
	})
	t.Run("Wrong date", func(t *testing.T) {
		mctx.EXPECT().Message().Return(&tele.Message{Text: "/abs 212344"}).Times(1)
		mctx.EXPECT().Reply("Не удалось распарсить дату, пожалуйста, введите дату в формате YYYY-MM-DD HH:MM").Times(1)

		err := handler(mctx)
		assert.NoError(t, err)
	})
	t.Run("no groups found", func(t *testing.T) {
		gomock.InOrder(
			mctx.EXPECT().Message().Return(&tele.Message{Text: "/abs 2025-04-06 14:44"}).Times(1),

			agroup.EXPECT().CurrentGroup(int64(1), gomock.Any(), "").Return(models.CurrentGroup{}, groups.ErrNoGroups).Times(1),
			mctx.EXPECT().Send("В данный момент, никакой группы не найдено!").Return(nil).Times(1),
		)
		err := handler(mctx)
		assert.NoError(t, err)
	})
	t.Run("ErrNotValidCookie", func(t *testing.T) {
		gomock.InOrder(
			mctx.EXPECT().Message().Return(&tele.Message{Text: "/abs 2025-04-06 14:44"}).Times(1),
			agroup.EXPECT().CurrentGroup(int64(1), gomock.Any(), "").Return(models.CurrentGroup{}, groups.ErrNotValidCookie).Times(1),
			mctx.EXPECT().Send("Вам необходимо установить свои cookie!").Return(nil).Times(1),
		)
		err := handler(mctx)
		assert.NoError(t, err)
	})
	t.Run("CurrentGroup return err", func(t *testing.T) {
		errExp := errors.New("err")
		gomock.InOrder(
			mctx.EXPECT().Message().Return(&tele.Message{Text: "/abs 2025-04-06 14:44"}).Times(1),
			agroup.EXPECT().CurrentGroup(int64(1), gomock.Any(), "").Return(models.CurrentGroup{}, errExp).Times(1),
		)
		err := handler(mctx)
		assert.ErrorIs(t, err, errExp)
	})
}

var grAsset = models.CurrentGroup{
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
