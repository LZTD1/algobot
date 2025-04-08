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
	setter := mocks2.NewMockDomainSetter(ctrl)
	fetcher := mocks2.NewMockGroupFetcher(ctrl)

	service := groups.NewGroup(
		log,
		gGetter,
		setter,
		fetcher,
	)

	t.Run("Groups", func(t *testing.T) {
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
	})
	t.Run("RefreshGroup", func(t *testing.T) {
		stubGroups := []models.Group{
			{
				GroupID:    1,
				Title:      "title1",
				TimeLesson: time.Date(2025, time.March, 22, 14, 0, 0, 0, time.UTC),
			},
			{
				GroupID:    2,
				Title:      "title2",
				TimeLesson: time.Date(2025, time.March, 22, 16, 0, 0, 0, time.UTC),
			},
		}

		t.Run("happy path", func(t *testing.T) {
			gomock.InOrder(
				setter.EXPECT().Cookies(int64(1)).Return("cookie", nil).Times(1),
				fetcher.EXPECT().Group("cookie").Return(stubGroups, nil).Times(1),
				setter.EXPECT().SetGroups(int64(1), stubGroups).Return(nil).Times(1),
			)

			err := service.RefreshGroup(1, "trace_id")
			assert.NoError(t, err)
		})
		t.Run("return empty group", func(t *testing.T) {
			gomock.InOrder(
				setter.EXPECT().Cookies(int64(1)).Return("cookie", nil).Times(1),
				fetcher.EXPECT().Group("cookie").Return([]models.Group{}, nil).Times(1),
			)

			err := service.RefreshGroup(1, "trace_id")
			assert.ErrorIs(t, err, groups.ErrNoGroups)
		})
		t.Run("Cookies return err", func(t *testing.T) {
			errExp := errors.New("some error")

			setter.EXPECT().Cookies(int64(1)).Return("cookie", errExp).Times(1)
			err := service.RefreshGroup(1, "trace_id")
			assert.ErrorIs(t, err, errExp)
		})
		t.Run("Cookies return empty cookie", func(t *testing.T) {
			setter.EXPECT().Cookies(int64(1)).Return("", nil).Times(1)
			err := service.RefreshGroup(1, "trace_id")
			assert.ErrorIs(t, err, groups.ErrNotValidCookie)
		})
		t.Run("fetcher Group return err", func(t *testing.T) {
			errExp := errors.New("some error")

			gomock.InOrder(
				setter.EXPECT().Cookies(int64(1)).Return("cookie", nil).Times(1),
				fetcher.EXPECT().Group("cookie").Return(stubGroups, errExp).Times(1),
			)

			err := service.RefreshGroup(1, "trace_id")
			assert.ErrorIs(t, err, errExp)
		})
		t.Run("fetcher SetGroups return err", func(t *testing.T) {
			errExp2 := errors.New("some error")

			gomock.InOrder(
				setter.EXPECT().Cookies(int64(1)).Return("cookie", nil).Times(1),
				fetcher.EXPECT().Group("cookie").Return(stubGroups, nil).Times(1),
				setter.EXPECT().SetGroups(int64(1), stubGroups).Return(errExp2).Times(1),
			)

			err := service.RefreshGroup(1, "trace_id")
			assert.ErrorIs(t, err, errExp2)
		})
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
