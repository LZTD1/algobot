package test

import (
	"algobot/internal/domain"
	"algobot/internal/domain/models"
	"algobot/internal/domain/telegram/keyboards"
	"algobot/internal/telegram/handlers/text"
	"algobot/test/mocks"
	mocks3 "algobot/test/mocks/telegram"
	mocks2 "algobot/test/mocks/telegram/handlers"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	tele "gopkg.in/telebot.v4"
	"strconv"
	"testing"
	"time"
)

func TestMyGroups(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocks.NewMockLogger()
	grouper := mocks2.NewMockGrouper(ctrl)
	ser := mocks2.NewMockGroupSerializer(ctrl)
	botName := "name"
	mctx := mocks3.NewMockContext(ctrl)

	handler := text.NewMyGroup(log, grouper, ser, botName)

	mctx.EXPECT().Get(gomock.Any()).Return("trace_id").AnyTimes()
	t.Run("happyPath", func(t *testing.T) {

		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			grouper.EXPECT().Groups(int64(1), "trace_id").Return(mockGroups, nil).Times(1),
			ser.EXPECT().Serialize(domain.SerializeMessage{
				Type: domain.GroupType,
				Data: []string{strconv.Itoa(mockGroups[0].GroupID)},
			}).Return("ser-g1", nil).Times(1),
			ser.EXPECT().Serialize(domain.SerializeMessage{
				Type: domain.GroupType,
				Data: []string{strconv.Itoa(mockGroups[1].GroupID)},
			}).Return("ser-g2", nil).Times(1),
			ser.EXPECT().Serialize(domain.SerializeMessage{
				Type: domain.GroupType,
				Data: []string{strconv.Itoa(mockGroups[2].GroupID)},
			}).Return("", errors.New("ser")).Times(1),
			mctx.EXPECT().Send(mockStringRet, tele.ModeMarkdown, keyboards.RefreshGroups()).Return(nil).Times(1),
		)

		err := handler.ServeContext(mctx)
		assert.NoError(t, err)
	})
	t.Run("no one group", func(t *testing.T) {
		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			grouper.EXPECT().Groups(int64(1), "trace_id").Return([]models.Group{}, nil).Times(1),
			mctx.EXPECT().Send("Всего групп: 0\nПопробуйте обновить группы!", tele.ModeMarkdown, keyboards.RefreshGroups()).Return(nil).Times(1),
		)

		err := handler.ServeContext(mctx)
		assert.NoError(t, err)
	})
	t.Run("Groups return err", func(t *testing.T) {
		err := errors.New("groups err")

		gomock.InOrder(
			mctx.EXPECT().Sender().Return(&tele.User{ID: 1}).Times(1),
			grouper.EXPECT().Groups(int64(1), "trace_id").Return(nil, err).Times(1),
			mctx.EXPECT().Send("<b>[trace_id]</b> Ошибка при получении групп!", tele.ModeHTML).Return(nil).Times(1),
		)

		err = handler.ServeContext(mctx)
		assert.NoError(t, err)
	})
}

var mockStringRet = "Всего групп: 3\n\n1. [Name 1](t.me/name?start=ser-g1) 🕐 вт 12:00\n2. [Name 2](t.me/name?start=ser-g2) 🕐 вт 13:00\n\n1. Name 3 🕐 ср 12:00"

var mockGroups = []models.Group{
	{
		GroupID:    1,
		Title:      "Name 1",
		TimeLesson: time.Date(2020, time.April, 14, 12, 0, 0, 0, time.UTC),
	},
	{
		GroupID:    2,
		Title:      "Name 2",
		TimeLesson: time.Date(2020, time.April, 14, 13, 0, 0, 0, time.UTC),
	},
	{
		GroupID:    3,
		Title:      "Name 3",
		TimeLesson: time.Date(2020, time.April, 15, 12, 0, 0, 0, time.UTC),
	},
}
