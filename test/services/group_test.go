package test

import (
	"algobot/internal/domain/backoffice"
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
	kidStats := mocks2.NewMockKidStats(ctrl)

	service := groups.NewGroup(
		log,
		gGetter,
		fetcher,
		setter,
		kidStats,
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
	t.Run("CurrentGroup", func(t *testing.T) {
		userID := int64(1)
		timeStub := time.Date(2025, time.March, 23, 14, 0, 0, 0, time.UTC)
		traceID := ""
		cookie := "cookie"
		errExp := errors.New("")

		t.Run("happy path", func(t *testing.T) {
			gomock.InOrder(
				setter.EXPECT().Cookies(userID).Return(cookie, nil).Times(1),
				gGetter.EXPECT().Groups(userID).Return(assets, nil).Times(1),
				kidStats.EXPECT().KidsStats(cookie, 1001).Return(kidsStats, nil).Times(1),
				kidStats.EXPECT().KidsNamesByGroup("1001", cookie).Return(KidsNames, nil).Times(1),
			)

			group, err := service.CurrentGroup(userID, timeStub, traceID)
			assert.NoError(t, err)
			assert.Equal(
				t,
				models.CurrentGroup{
					GroupID:  1001,
					Title:    "group 1",
					Lesson:   "LessonTitle1",
					LessonID: 1,
					Kids:     []string{"FullName1", "FullName0"},
					MissingKids: []models.MissingKid{{
						"FullName1",
						1,
						2,
					}},
				},
				group,
			)
		})
		t.Run("Cookies return err", func(t *testing.T) {
			t.Run("Empty cookie", func(t *testing.T) {
				setter.EXPECT().Cookies(userID).Return("", nil).Times(1)
				_, err := service.CurrentGroup(userID, timeStub, traceID)
				assert.ErrorIs(t, err, groups.ErrNotValidCookie)
			})
			t.Run("Cookies return error", func(t *testing.T) {
				setter.EXPECT().Cookies(userID).Return("", errExp).Times(1)
				_, err := service.CurrentGroup(userID, timeStub, traceID)
				assert.ErrorIs(t, err, errExp)
			})
		})
		t.Run("Groups return err", func(t *testing.T) {

			gomock.InOrder(
				setter.EXPECT().Cookies(userID).Return(cookie, nil).Times(1),
				gGetter.EXPECT().Groups(userID).Return(nil, errExp).Times(1),
			)

			_, err := service.CurrentGroup(userID, timeStub, traceID)
			assert.ErrorIs(t, err, errExp)
		})
		t.Run("returns ErrNoGroups ", func(t *testing.T) {
			gomock.InOrder(
				setter.EXPECT().Cookies(userID).Return(cookie, nil).Times(1),
				gGetter.EXPECT().Groups(userID).Return(assets, nil).Times(1),
			)

			_, err := service.CurrentGroup(userID, time.Date(2025, time.April, 2025, 22, 0, 0, 0, time.UTC), traceID)
			assert.ErrorIs(t, err, groups.ErrNoGroups)
		})
		t.Run("KidsStats return err", func(t *testing.T) {
			gomock.InOrder(
				setter.EXPECT().Cookies(userID).Return(cookie, nil).Times(1),
				gGetter.EXPECT().Groups(userID).Return(assets, nil).Times(1),
				kidStats.EXPECT().KidsStats(cookie, 1001).Return(backoffice.KidsStats{}, errExp).Times(1),
			)

			_, err := service.CurrentGroup(userID, timeStub, traceID)
			assert.ErrorIs(t, err, errExp)
		})
		t.Run("KidsNamesByGroup return err", func(t *testing.T) {
			gomock.InOrder(
				setter.EXPECT().Cookies(userID).Return(cookie, nil).Times(1),
				gGetter.EXPECT().Groups(userID).Return(assets, nil).Times(1),
				kidStats.EXPECT().KidsStats(cookie, 1001).Return(kidsStats, nil).Times(1),
				kidStats.EXPECT().KidsNamesByGroup("1001", cookie).Return(backoffice.NamesByGroup{}, errExp).Times(1),
			)

			_, err := service.CurrentGroup(userID, timeStub, traceID)
			assert.ErrorIs(t, err, errExp)
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
var KidsNames = backoffice.NamesByGroup{
	Status: "",
	Data: backoffice.GroupData{
		Items: []backoffice.Student{
			{
				ID:       0,
				FullName: "FullName0",
				LastGroup: backoffice.Group{
					Status: 0,
					ID:     1001,
				},
			},
			{
				ID:       1,
				FullName: "FullName1",
				LastGroup: backoffice.Group{
					Status: 0,
					ID:     1001,
				},
			},
		},
	},
}

var (
	kidsStats = backoffice.KidsStats{
		Status: "",
		Data: []backoffice.KidStat{
			{
				StudentID: 0,
				Attendance: []backoffice.Attendance{
					{
						LessonID:           0,
						LessonTitle:        "LessonTitle0",
						StartTimeFormatted: "вс 12.03.25 14:00",
						Status:             "present",
					},
					{
						LessonID:           1,
						LessonTitle:        "LessonTitle1",
						StartTimeFormatted: "вс 23.03.25 14:00",
						Status:             "present",
					},
					{
						LessonID:           2,
						LessonTitle:        "LessonTitle2",
						StartTimeFormatted: "вс 23.04.25 14:00",
						Status:             "future",
					},
				},
			},
			{
				StudentID: 1,
				Attendance: []backoffice.Attendance{
					{
						LessonID:           0,
						LessonTitle:        "LessonTitle0",
						StartTimeFormatted: "вс 12.03.25 14:00",
						Status:             "absent",
					},
					{
						LessonID:           1,
						LessonTitle:        "LessonTitle1",
						StartTimeFormatted: "вс 23.03.25 14:00",
						Status:             "absent",
					},
					{
						LessonID:           2,
						LessonTitle:        "LessonTitle2",
						StartTimeFormatted: "вс 23.04.25 14:00",
						Status:             "future",
					},
				},
			},
		},
	}
)
