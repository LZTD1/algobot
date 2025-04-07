package test

import (
	"algobot/internal/domain/models"
	"algobot/internal/services/groups"
	"algobot/test/mocks"
	mocks2 "algobot/test/mocks/services"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocks.NewMockLogger()
	gGetter := mocks2.NewMockGroupGetter(ctrl)

	service := groups.NewGroup(log, gGetter)

	t.Run("happy path", func(t *testing.T) {
		gGetter.EXPECT().Groups(int64(1)).Return(assets, nil).Times(1)
		gr, err := service.Groups(1, "trace_id")
		assert.NoError(t, err)
		assert.Equal(t, []models.Group{
			{
				GroupID:    999,
				Title:      "group 3",
				TimeLesson: time.Date(2025, time.March, 22, 14, 0, 0, 0, time.UTC),
			},
			{
				GroupID:    1001,
				Title:      "group 1",
				TimeLesson: time.Date(2025, time.March, 23, 14, 0, 0, 0, time.UTC),
			},
			{
				GroupID:    1000,
				Title:      "group 2",
				TimeLesson: time.Date(2025, time.March, 23, 16, 0, 0, 0, time.UTC),
			},
		}, gr)
	})
	t.Run("Groups return err", func(t *testing.T) {
		errExp := errors.New("some error")
		gGetter.EXPECT().Groups(int64(1)).Return(nil, errExp).Times(1)
		_, err := service.Groups(1, "trace_id")
		assert.ErrorIs(t, err, errExp)
	})

}

var assets = []models.Group{
	{
		GroupID:    1000,
		Title:      "group 2",
		TimeLesson: time.Date(2025, time.March, 23, 16, 0, 0, 0, time.UTC),
	},
	{
		GroupID:    1001,
		Title:      "group 1",
		TimeLesson: time.Date(2025, time.March, 23, 14, 0, 0, 0, time.UTC),
	},

	{
		GroupID:    999,
		Title:      "group 3",
		TimeLesson: time.Date(2025, time.March, 22, 14, 0, 0, 0, time.UTC),
	},
}
