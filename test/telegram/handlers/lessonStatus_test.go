package test

import (
	"algobot/internal/services/backoffice"
	"algobot/internal/telegram/handlers/callback"
	"algobot/test/mocks"
	mocks3 "algobot/test/mocks/telegram"
	mocks2 "algobot/test/mocks/telegram/handlers"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	tele "gopkg.in/telebot.v4"
	"testing"
)

func TestLessonStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocks.NewMockLogger()
	ls := mocks2.NewMockLessonStatuser(ctrl)
	mctx := mocks3.NewMockContext(ctrl)

	mctx.EXPECT().Get(gomock.Any()).Return("").AnyTimes()
	mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).AnyTimes()

	t.Run("Happy path open", func(t *testing.T) {
		handler := callback.LessonStatus(ls, backoffice.OpenLesson, log)

		gomock.InOrder(
			mctx.EXPECT().Callback().Return(&tele.Callback{Data: "\fopen_lesson_1_1"}),
			ls.EXPECT().SetLessonStatus(int64(1), "1", "1", backoffice.OpenLesson, "").Return(nil).Times(1),
			mctx.EXPECT().Send("Статус переключен!").Return(nil).Times(1),
		)

		err := handler(mctx)
		assert.NoError(t, err)
	})
	t.Run("Happy path close", func(t *testing.T) {
		handler := callback.LessonStatus(ls, backoffice.CloseLesson, log)

		gomock.InOrder(
			mctx.EXPECT().Callback().Return(&tele.Callback{Data: "\fclose_lesson_1_1"}),
			ls.EXPECT().SetLessonStatus(int64(1), "1", "1", backoffice.CloseLesson, "").Return(nil).Times(1),
			mctx.EXPECT().Send("Статус переключен!").Return(nil).Times(1),
		)

		err := handler(mctx)
		assert.NoError(t, err)
	})
	t.Run("wrong data", func(t *testing.T) {
		handler := callback.LessonStatus(ls, backoffice.CloseLesson, log)

		gomock.InOrder(
			mctx.EXPECT().Callback().Return(&tele.Callback{Data: "close_lesson_1"}),
			mctx.EXPECT().Send("⚠️ Ошибка при анализе данных от кнопки").Return(nil).Times(1),
		)

		err := handler(mctx)
		assert.NoError(t, err)
	})
	t.Run("SetLessonStatus return err", func(t *testing.T) {
		errExp := errors.New("1")
		handler := callback.LessonStatus(ls, backoffice.CloseLesson, log)

		gomock.InOrder(
			mctx.EXPECT().Callback().Return(&tele.Callback{Data: "\fclose_lesson_1_1"}),
			ls.EXPECT().SetLessonStatus(int64(1), "1", "1", backoffice.CloseLesson, "").Return(errExp).Times(1),
		)

		err := handler(mctx)
		assert.ErrorIs(t, err, errExp)
	})
}
