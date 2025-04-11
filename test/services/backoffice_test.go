package test

import (
	backoffice2 "algobot/internal/domain/backoffice"
	"algobot/internal/domain/models"
	backoffice3 "algobot/internal/lib/backoffice"
	"algobot/internal/services/backoffice"
	"algobot/test/mocks"
	mocks2 "algobot/test/mocks/services"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestBackoffice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocks.NewMockLogger()
	cookieGetter := mocks2.NewMockCookieGetter(ctrl)
	groupView := mocks2.NewMockGroupView(ctrl)
	kidViewer := mocks2.NewMockKidViewer(ctrl)

	sbo := backoffice.NewBackoffice(log, cookieGetter, groupView, kidViewer)
	t.Run("KidView", func(t *testing.T) {
		uid := int64(1)
		kidID := "1"
		groupID := "2"
		traceID := "trace"
		cookie := "cookie"
		errExp := errors.New("err exp")

		t.Run("happy path by KidView", func(t *testing.T) {
			gomock.InOrder(
				cookieGetter.EXPECT().Cookies(uid).Return(cookie, nil).Times(1),
				kidViewer.EXPECT().KidView(kidID, cookie).Return(KidViewBackoffice, nil).Times(1),
			)

			view, err := sbo.KidView(uid, kidID, groupID, traceID)
			assert.NoError(t, err)
			assert.Equal(t, kidExpected, view)
		})
		t.Run("happy path by KidsNamesByGroup", func(t *testing.T) {
			kidExpectedP := kidExpected
			kidExpectedP.Extra = models.NotAccessible

			gomock.InOrder(
				cookieGetter.EXPECT().Cookies(uid).Return(cookie, nil).Times(1),
				kidViewer.EXPECT().KidView(kidID, cookie).Return(backoffice2.KidView{}, backoffice3.ErrNotFound).Times(1),
				kidViewer.EXPECT().KidsNamesByGroup(groupID, cookie).Return(kidsNamesByGroupBackoffice, nil).Times(1),
			)

			view, err := sbo.KidView(uid, kidID, groupID, traceID)
			assert.NoError(t, err)
			assert.Equal(t, kidExpectedP, view)
		})
		t.Run("cookie returns err", func(t *testing.T) {
			gomock.InOrder(
				cookieGetter.EXPECT().Cookies(uid).Return("", errExp).Times(1),
			)

			_, err := sbo.KidView(uid, kidID, groupID, traceID)
			assert.ErrorIs(t, err, errExp)
		})
		t.Run("KidView returns err", func(t *testing.T) {
			gomock.InOrder(
				cookieGetter.EXPECT().Cookies(uid).Return(cookie, nil).Times(1),
				kidViewer.EXPECT().KidView(kidID, cookie).Return(backoffice2.KidView{}, errExp).Times(1),
			)

			_, err := sbo.KidView(uid, kidID, groupID, traceID)
			assert.ErrorIs(t, err, errExp)
		})
		t.Run("KidsNamesByGroup returns err", func(t *testing.T) {
			gomock.InOrder(
				cookieGetter.EXPECT().Cookies(uid).Return(cookie, nil).Times(1),
				kidViewer.EXPECT().KidView(kidID, cookie).Return(backoffice2.KidView{}, backoffice3.ErrNotFound).Times(1),
				kidViewer.EXPECT().KidsNamesByGroup(groupID, cookie).Return(backoffice2.NamesByGroup{}, errExp).Times(1),
			)

			_, err := sbo.KidView(uid, kidID, groupID, traceID)
			assert.ErrorIs(t, err, errExp)
		})
	})
	t.Run("Group view", func(t *testing.T) {
		uid := int64(1)
		kidID := "1"
		groupID := "2"
		traceID := "trace"
		cookie := "cookie"
		errExp := errors.New("err exp")

		_ = errExp
		_ = kidID
		t.Run("happy path", func(t *testing.T) {
			gomock.InOrder(
				cookieGetter.EXPECT().Cookies(uid).Return(cookie, nil).Times(1),
				groupView.EXPECT().GroupView(groupID, cookie).Return(groupInfoBackoffice, nil).Times(1),
				groupView.EXPECT().KidsNamesByGroup(groupID, cookie).Return(kidsNamesByGroupBackoffice, nil).Times(1),
			)

			view, err := sbo.GroupView(uid, groupID, traceID)
			assert.NoError(t, err)
			assert.Equal(t, expectedGroupView, view)
		})
		t.Run("cookie return err", func(t *testing.T) {
			gomock.InOrder(
				cookieGetter.EXPECT().Cookies(uid).Return("", errExp).Times(1),
			)

			_, err := sbo.GroupView(uid, groupID, traceID)
			assert.ErrorIs(t, err, errExp)
		})
		t.Run("GroupView return err", func(t *testing.T) {
			gomock.InOrder(
				cookieGetter.EXPECT().Cookies(uid).Return(cookie, nil).Times(1),
				groupView.EXPECT().GroupView(groupID, cookie).Return(groupInfoBackoffice, errExp).Times(1),
			)

			_, err := sbo.GroupView(uid, groupID, traceID)
			assert.ErrorIs(t, err, errExp)
		})
		t.Run("GroupView return err", func(t *testing.T) {
			gomock.InOrder(
				cookieGetter.EXPECT().Cookies(uid).Return(cookie, nil).Times(1),
				groupView.EXPECT().GroupView(groupID, cookie).Return(groupInfoBackoffice, nil).Times(1),
				groupView.EXPECT().KidsNamesByGroup(groupID, cookie).Return(kidsNamesByGroupBackoffice, errExp).Times(1),
			)

			_, err := sbo.GroupView(uid, groupID, traceID)
			assert.ErrorIs(t, err, errExp)
		})
	})

}

var KidViewBackoffice = backoffice2.KidView{
	Status: "Активен",
	Data: backoffice2.Student{
		ID:              101,
		FirstName:       "Иван",
		LastName:        "Иванов",
		FullName:        "Иван Иванов",
		ParentName:      "Мария Ивановна Иванова",
		Email:           "ivanov@example.com",
		HasLaptop:       1,
		Phone:           "+79171234567",
		Age:             10,
		BirthDate:       time.Date(2013, 5, 15, 0, 0, 0, 0, time.UTC),
		CreatedAt:       time.Date(2022, 9, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt:       time.Date(2023, 3, 15, 14, 30, 0, 0, time.UTC),
		DeletedAt:       nil,
		HasBranchAccess: true,
		Username:        "ivan101",
		Password:        "securepassword123",
		LastGroup: backoffice2.Group{
			ID:             201,
			GroupStudentID: 301,
			Title:          "Группа по математике",
			Content:        "Основы математики для младших школьников",
			Track:          1,
			Status:         1,
			StartTime:      time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
			EndTime:        time.Date(2023, 6, 30, 11, 0, 0, 0, time.UTC),
			CourseID:       401,
			CreatedAt:      time.Date(2022, 12, 1, 9, 0, 0, 0, time.UTC),
			UpdatedAt:      time.Date(2023, 2, 10, 10, 0, 0, 0, time.UTC),
			DeletedAt:      nil,
		},
		Groups: []backoffice2.Group{
			{
				ID:             201,
				GroupStudentID: 301,
				Title:          "Группа по математике",
				Content:        "Основы математики для младших школьников",
				Track:          1,
				Status:         1,
				StartTime:      time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
				EndTime:        time.Date(2023, 6, 30, 11, 0, 0, 0, time.UTC),
				CourseID:       401,
				CreatedAt:      time.Date(2022, 12, 1, 9, 0, 0, 0, time.UTC),
				UpdatedAt:      time.Date(2023, 2, 10, 10, 0, 0, 0, time.UTC),
				DeletedAt:      nil,
			},
			{
				ID:             202,
				GroupStudentID: 302,
				Title:          "Уравнения и геометрия",
				Content:        "Разбор уравнений и элементы геометрии",
				Track:          2,
				Status:         1,
				StartTime:      time.Date(2023, 7, 1, 10, 0, 0, 0, time.UTC),
				EndTime:        time.Date(2023, 12, 31, 11, 0, 0, 0, time.UTC),
				CourseID:       402,
				CreatedAt:      time.Date(2023, 6, 1, 9, 0, 0, 0, time.UTC),
				UpdatedAt:      time.Date(2023, 6, 10, 10, 0, 0, 0, time.UTC),
				DeletedAt:      nil,
			},
		},
	},
}
var kidExpected = models.KidView{
	Extra: "",
	Kid: models.Kid{
		FullName:   "Иван Иванов",
		ParentName: "Мария Ивановна Иванова",
		Email:      "ivanov@example.com",
		Phone:      "+79171234567",
		Age:        10,
		BirthDate:  time.Date(2013, 5, 15, 0, 0, 0, 0, time.UTC),
		Username:   "ivan101",
		Password:   "securepassword123",
		Groups: []models.KidViewGroup{
			{
				ID:        201,
				Title:     "Группа по математике",
				Content:   "Основы математики для младших школьников",
				Status:    1,
				StartTime: time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 6, 30, 11, 0, 0, 0, time.UTC),
			},
			{
				ID:        202,
				Title:     "Уравнения и геометрия",
				Content:   "Разбор уравнений и элементы геометрии",
				Status:    1,
				StartTime: time.Date(2023, 7, 1, 10, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2023, 12, 31, 11, 0, 0, 0, time.UTC),
			},
		},
	},
}

var kidsNamesByGroupBackoffice = backoffice2.NamesByGroup{
	Status: "success",
	Data: backoffice2.GroupData{
		Items: []backoffice2.Student{
			{
				ID:              1,
				FirstName:       "Иван",
				LastName:        "Иванов",
				FullName:        "Иван Иванов",
				ParentName:      "Мария Ивановна Иванова",
				Email:           "ivanov@example.com",
				HasLaptop:       1,
				Phone:           "+79171234567",
				Age:             10,
				BirthDate:       time.Date(2013, 5, 15, 0, 0, 0, 0, time.UTC),
				CreatedAt:       time.Date(2022, time.September, 1, 12, 0, 0, 0, time.UTC),
				UpdatedAt:       time.Date(2022, time.September, 1, 12, 15, 0, 0, time.UTC),
				DeletedAt:       nil,
				HasBranchAccess: false,
				Username:        "ivan101",
				Password:        "securepassword123",
				LastGroup: backoffice2.Group{
					ID:             201,
					GroupStudentID: 301,
					Title:          "Группа по математике",
					Content:        "Основы математики для младших школьников",
					Track:          1,
					Status:         1,
					StartTime:      time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
					EndTime:        time.Date(2023, 6, 30, 11, 0, 0, 0, time.UTC),
					CourseID:       401,
					CreatedAt:      time.Date(2022, 12, 1, 9, 0, 0, 0, time.UTC),
					UpdatedAt:      time.Date(2023, 2, 10, 10, 0, 0, 0, time.UTC),
					DeletedAt:      nil,
				},
				Groups: []backoffice2.Group{
					{
						ID:             201,
						GroupStudentID: 301,
						Title:          "Группа по математике",
						Content:        "Основы математики для младших школьников",
						Track:          1,
						Status:         1,
						StartTime:      time.Date(2023, 2, 1, 10, 0, 0, 0, time.UTC),
						EndTime:        time.Date(2023, 6, 30, 11, 0, 0, 0, time.UTC),
						CourseID:       401,
						CreatedAt:      time.Date(2022, 12, 1, 9, 0, 0, 0, time.UTC),
						UpdatedAt:      time.Date(2023, 2, 10, 10, 0, 0, 0, time.UTC),
						DeletedAt:      nil,
					},
					{
						ID:             202,
						GroupStudentID: 302,
						Title:          "Уравнения и геометрия",
						Content:        "Разбор уравнений и элементы геометрии",
						Track:          2,
						Status:         1,
						StartTime:      time.Date(2023, 7, 1, 10, 0, 0, 0, time.UTC),
						EndTime:        time.Date(2023, 12, 31, 11, 0, 0, 0, time.UTC),
						CourseID:       402,
						CreatedAt:      time.Date(2023, 6, 1, 9, 0, 0, 0, time.UTC),
						UpdatedAt:      time.Date(2023, 6, 10, 10, 0, 0, 0, time.UTC),
						DeletedAt:      nil,
					},
				},
			},
		},
	},
}
var groupInfoBackoffice = backoffice2.GroupInfo{
	Status: "success",
	Data: backoffice2.GroupDataFull{
		ID:      101,
		Title:   "Основы программирования",
		Content: "Изучение основ программирования для детей",
		Type: backoffice2.TypeFull{
			Value: "programming",
			Label: "Программирование",
			Tag:   "coding",
		},
		Status: backoffice2.StatusFull{
			Value: 1,
			Label: "Активная",
			Tag:   "active",
		},
		StatusChangedAt: "2022-08-15T12:00:00Z",
		StartTime:       "2022-09-01T10:00:00Z",
		NextLessonTime:  "2022-09-08T10:00:00Z",
		LessonsTotal:    12,
		LessonsPassed:   3,
		HardwareNeeded:  1,
		Branch: backoffice2.BranchFull{
			ID:                            1,
			Title:                         "Отделение Центральный",
			Code:                          "CTR001",
			Description:                   "Центральное отделение",
			Phone:                         "+79998887766",
			Email:                         "info@branch.example.com",
			SiteURL:                       "https://branch.example.com",
			TemplateVersion:               1,
			UseAmo:                        true,
			AmoConfigID:                   1001,
			ShowFinanceInfo:               true,
			LmsDisplayStudentCredentials:  true,
			ShowOnlineRoomURLField:        1,
			UseSms:                        true,
			LanguageID:                    1,
			OrderName:                     1,
			UseFullyPaidLabel:             0,
			BrandName:                     "Учебный Центр",
			MaxCountStudentsForShowOnline: 30,
			IsFillPaymentSystem:           true,
			FirstLessonNoRoyalty:          0,
			RootBranchID:                  0,
		},
		Venue: backoffice2.VenueFull{
			ID:           10,
			Title:        "Салон №1",
			Address:      "ул. Большая Красная, д.1",
			ContactName:  "Иван Иванов",
			ContactEmail: "ivan.ivanov@venue.example.com",
			ContactPhone: "+79991112233",
		},
		Curator: backoffice2.UserFull{
			ID:       201,
			Username: "ivan_ivanov",
			Phone:    "+79991234567",
			Email:    "ivan.ivanov@curator.example.com",
			Name:     "Иван Иванов",
			Profile: backoffice2.ProfileFull{
				PhotoURL: "https://example.com/photos/ivan_ivanov.jpg",
				Promo:    "Популярный куратор",
			},
			Status: 1,
			Links: backoffice2.LinksFull{
				Self: "https://backoffice.example.com/api/users/201",
			},
		},
		Teacher: backoffice2.TeacherFull{
			ID:       301,
			Username: "anna_smirnova",
			Phone:    "+79995556677",
			Email:    "anna.smirnova@teacher.example.com",
			Name:     "Анна Смирнова",
			Profile: backoffice2.ProfileFull{
				PhotoURL: "https://example.com/photos/anna_smirnova.jpg",
				Promo:    "Стажированный преподаватель",
			},
			AllowedUserCourses: nil,
			Status:             1,
			Links: backoffice2.LinksFull{
				Self: "https://backoffice.example.com/api/teachers/301",
			},
		},
		Teachers:      nil,
		ClientManager: nil,
		Course: backoffice2.CourseFull{
			ID:          401,
			Name:        "Основы программирования",
			GUID:        "COURSE001",
			Description: "Изучение основ программирования для детей",
			ContentType: "interactive",
			CourseType: backoffice2.CourseTypeFull{
				ID:    1,
				Title: "Технологии",
				Code:  "tech",
			},
			LessonsCount:                12,
			GroupLessonsAmount:          4,
			LessonsCountFormatted:       "12",
			GroupLessonsAmountFormatted: "4",
			IsDeleted:                   0,
			Links: backoffice2.LinksFull{
				Self: "https://backoffice.example.com/api/courses/401",
			},
		},
		LanguageID:                     1,
		Journal:                        true,
		ShowJournal:                    true,
		ShowOnlineRoom:                 true,
		IsOnline:                       true,
		ActiveStudentCount:             25,
		OnlineRoomURL:                  "https://online-room.example.com/group/101",
		UseClientManager:               1,
		DisplayLessonDurationInMinutes: 90,
		DeletedAt:                      nil,
		DeletedBy:                      nil,
		PriorityLevel: backoffice2.PriorityLevelFull{
			Value: "high",
			Label: "Высокий",
			Tag:   "high",
		},
		IsFull:    false,
		CreatedAt: "2022-07-10T10:00:00Z",
		CreatedBy: backoffice2.UserFull{
			ID:       202,
			Username: "petr_petrov",
			Phone:    "+79992223344",
			Email:    "petr.petrov@teacher.example.com",
			Name:     "Пётр Петров",
			Profile: backoffice2.ProfileFull{
				PhotoURL: "https://example.com/photos/petr_petrov.jpg",
				Promo:    "Учитель с большим опытом",
			},
			Status: 1,
			Links: backoffice2.LinksFull{
				Self: "https://backoffice.example.com/api/users/202",
			},
		},
		Related: backoffice2.RelatedFull{
			Statuses:       nil,
			Types:          nil,
			PriorityLevels: nil,
		},
	},
}

var expectedGroupView = models.GroupView{
	GroupID:        101,
	GroupTitle:     "Основы программирования",
	GroupContent:   "Изучение основ программирования для детей",
	NextLessonTime: "2022-09-08T10:00:00Z",
	LessonsTotal:   12,
	LessonsPassed:  3,
	ActiveKids:     []models.GroupKid(nil),
	NotActiveKids: []models.GroupKid{{
		ID:       1,
		FullName: "Иван Иванов",
		LastGroup: models.KidGroup{
			ID:        201,
			StartTime: time.Date(2023, time.February, 1, 10, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2023, time.June, 30, 11, 0, 0, 0, time.UTC),
		},
	}},
}
