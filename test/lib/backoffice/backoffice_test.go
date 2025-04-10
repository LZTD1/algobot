package backoffice

import (
	"algobot/internal/config"
	backoffice2 "algobot/internal/domain/backoffice"
	"algobot/internal/domain/models"
	"algobot/internal/lib/backoffice"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestBackoffice(t *testing.T) {
	t.Run("Group", func(t *testing.T) {
		groupsExp := []models.Group{
			{GroupID: 98637162, Title: "Группа по курсу КГ", TimeLesson: time.Date(2025, time.April, 13, 14, 0, 0, 0, time.UTC)},
			{GroupID: 98623404, Title: "Группа по курсу ОЛиП МП", TimeLesson: time.Date(2025, time.April, 13, 12, 0, 0, 0, time.UTC)},
			{GroupID: 98621252, Title: "Группа по курсу Пст", TimeLesson: time.Date(2025, time.April, 13, 18, 0, 0, 0, time.UTC)},
			{GroupID: 98619913, Title: "Группа по курсу КГ", TimeLesson: time.Date(2025, time.April, 12, 10, 0, 0, 0, time.UTC)},
			{GroupID: 98619873, Title: "Группа по курсу ГД", TimeLesson: time.Date(2025, time.April, 12, 14, 0, 0, 0, time.UTC)},
			{GroupID: 98619867, Title: "Группа по курсу ВП", TimeLesson: time.Date(2025, time.April, 12, 12, 0, 0, 0, time.UTC)},
			{GroupID: 98589447, Title: "Группа по курсу ВП", TimeLesson: time.Date(2025, time.April, 13, 10, 0, 0, 0, time.UTC)},
			{GroupID: 985504, Title: "Группа по курсу Пст 2", TimeLesson: time.Date(2025, time.April, 12, 18, 0, 0, 0, time.UTC)},
			{GroupID: 978298, Title: "Группа по курсу Пст 2", TimeLesson: time.Date(2025, time.April, 13, 16, 0, 0, 0, time.UTC)},
		}

		responseHTML := readFile("group_example")

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte(responseHTML))
			rw.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		bo := backoffice.NewBackoffice(&config.Backoffice{
			Retries:         5,
			RetriesTimeout:  time.Second,
			ResponseTimeout: time.Second,
		}, backoffice.WithURL(server.URL))

		group, err := bo.Group("cookie")
		assert.NoError(t, err)
		assert.Equal(t, groupsExp, group)
	})
	t.Run("GroupView", func(t *testing.T) {
		responseHTML := readFile("GroupView_example")
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte(responseHTML))
			rw.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		bo := backoffice.NewBackoffice(&config.Backoffice{
			Retries:         5,
			RetriesTimeout:  time.Second,
			ResponseTimeout: time.Second,
		}, backoffice.WithURL(server.URL))

		group, err := bo.GroupView("", "")
		assert.NoError(t, err)
		assert.Equal(t, backofficeGroupViewExpected, group)
	})
	t.Run("KidsNamesByGroup", func(t *testing.T) {
		responseHTML := readFile("KidsNamesByGroup_example")
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte(responseHTML))
			rw.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		bo := backoffice.NewBackoffice(&config.Backoffice{
			Retries:         5,
			RetriesTimeout:  time.Second,
			ResponseTimeout: time.Second,
		}, backoffice.WithURL(server.URL))

		group, err := bo.KidsNamesByGroup("", "")
		assert.NoError(t, err)
		assert.Equal(t, backofficeKidsByGroupExpected, group)
	})
	t.Run("KidView", func(t *testing.T) {
		responseHTML := readFile("KidView_example")
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte(responseHTML))
			rw.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		bo := backoffice.NewBackoffice(&config.Backoffice{
			Retries:         5,
			RetriesTimeout:  time.Second,
			ResponseTimeout: time.Second,
		}, backoffice.WithURL(server.URL))

		group, err := bo.KidView("", "")
		assert.NoError(t, err)
		assert.Equal(t, backofficeKidViewExpected, group)
	})
}

func readFile(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	responseHTML := string(b)
	return responseHTML
}

var backofficeGroupViewExpected = backoffice2.GroupInfo{
	Status: "success",
	Data: backoffice2.GroupDataFull{ID: 12345678, Title: "Занятие в библиотеке 5 в 14.00", Content: "Группа по курсу Python", Type: backoffice2.TypeFull{Value: "regular", Label: "Группа", Tag: "default"}, Status: backoffice2.StatusFull{Value: 10, Label: "Активная", Tag: "success"}, StatusChangedAt: "20.09.2024 12:04", StartTime: "22.09.2024 14:00", NextLessonTime: "13.04.2025 14:00", LessonsTotal: 35, LessonsPassed: 28, HardwareNeeded: 0, Branch: backoffice2.BranchFull{
		ID:                            987,
		Title:                         "Город Х",
		Code:                          "city_x",
		Description:                   "",
		Phone:                         "+7 (491) 555-22-33",
		Email:                         "cityx@algoritmika.org",
		SiteURL:                       "https://cityx.algoritmika.org",
		TemplateVersion:               2,
		UseAmo:                        true,
		AmoConfigID:                   321,
		ShowFinanceInfo:               true,
		LmsDisplayStudentCredentials:  true,
		ShowOnlineRoomURLField:        0,
		UseSms:                        false,
		LanguageID:                    2,
		OrderName:                     1,
		UseFullyPaidLabel:             0,
		BrandName:                     "",
		MaxCountStudentsForShowOnline: 10,
		IsFillPaymentSystem:           false,
		FirstLessonNoRoyalty:          0,
		RootBranchID:                  123,
	}, Venue: backoffice2.VenueFull{ID: 4321, Title: "Библиотека №5", Address: "390000, Город Х, ул Ленина, д 10", ContactName: "", ContactEmail: "", ContactPhone: "", Links: backoffice2.LinksFull{Self: "/venue/view/4321"}}, Curator: backoffice2.UserFull{ID: 1234, Username: "random_user", Phone: "+7 (900) 555-12-34", Email: "example123@mail.com", Name: "Иван Иванов", Profile: backoffice2.ProfileFull{PhotoURL: "/uploads/avatar/avatar/avatar_1234_1603094230-96x96.jpg", Promo: ""}, Status: 10, Links: backoffice2.LinksFull{Self: "/user/update/1234"}}, Teacher: backoffice2.TeacherFull{ID: 5678, Username: "random_user2", Phone: "+7 (922) 555-00-11", Email: "randomuser2@mail.com", Name: "Алексей Смирнов", Profile: backoffice2.ProfileFull{PhotoURL: "", Promo: ""}, AllowedUserCourses: []backoffice2.AllowedUserCourseFull{{UserID: 42407, CourseID: 84, IsAllowed: 1}, {UserID: 42407, CourseID: 305, IsAllowed: 1}, {UserID: 42407, CourseID: 389, IsAllowed: 1}, {UserID: 42407, CourseID: 405, IsAllowed: 1}, {UserID: 42407, CourseID: 406, IsAllowed: 1}, {UserID: 42407, CourseID: 407, IsAllowed: 1}, {UserID: 42407, CourseID: 408, IsAllowed: 1}, {UserID: 42407, CourseID: 416, IsAllowed: 1}, {UserID: 42407, CourseID: 417, IsAllowed: 1}, {UserID: 42407, CourseID: 448, IsAllowed: 1}, {UserID: 42407, CourseID: 465, IsAllowed: 1}, {UserID: 42407, CourseID: 606, IsAllowed: 1}, {UserID: 42407, CourseID: 640, IsAllowed: 1}, {UserID: 42407, CourseID: 641, IsAllowed: 1}, {UserID: 42407, CourseID: 661, IsAllowed: 1}, {UserID: 42407, CourseID: 662, IsAllowed: 1}, {UserID: 42407, CourseID: 665, IsAllowed: 1}, {UserID: 42407, CourseID: 666, IsAllowed: 1}, {UserID: 42407, CourseID: 686, IsAllowed: 1}, {UserID: 42407, CourseID: 706, IsAllowed: 1}, {UserID: 42407, CourseID: 707, IsAllowed: 1}, {UserID: 42407, CourseID: 716, IsAllowed: 1}, {UserID: 42407, CourseID: 727, IsAllowed: 1}, {UserID: 42407, CourseID: 729, IsAllowed: 1}, {UserID: 42407, CourseID: 734, IsAllowed: 1}, {UserID: 42407, CourseID: 735, IsAllowed: 1}, {UserID: 42407, CourseID: 777, IsAllowed: 1}, {UserID: 42407, CourseID: 783, IsAllowed: 1}, {UserID: 42407, CourseID: 797, IsAllowed: 1}, {UserID: 42407, CourseID: 799, IsAllowed: 1}, {UserID: 42407, CourseID: 809, IsAllowed: 1}, {UserID: 42407, CourseID: 810, IsAllowed: 1}, {UserID: 42407, CourseID: 823, IsAllowed: 1}, {UserID: 42407, CourseID: 831, IsAllowed: 1}, {UserID: 42407, CourseID: 852, IsAllowed: 1}, {UserID: 42407, CourseID: 853, IsAllowed: 1}, {UserID: 42407, CourseID: 857, IsAllowed: 1}, {UserID: 42407, CourseID: 858, IsAllowed: 1}, {UserID: 42407, CourseID: 859, IsAllowed: 1}, {UserID: 42407, CourseID: 860, IsAllowed: 1}, {UserID: 42407, CourseID: 861, IsAllowed: 1}, {UserID: 42407, CourseID: 862, IsAllowed: 1}, {UserID: 42407, CourseID: 864, IsAllowed: 1}, {UserID: 42407, CourseID: 1338, IsAllowed: 1}, {UserID: 42407, CourseID: 1339, IsAllowed: 1}, {UserID: 42407, CourseID: 1346, IsAllowed: 1}, {UserID: 42407, CourseID: 1347, IsAllowed: 1}, {UserID: 42407, CourseID: 1387, IsAllowed: 1}, {UserID: 42407, CourseID: 1459, IsAllowed: 1}, {UserID: 42407, CourseID: 1484, IsAllowed: 1}, {UserID: 42407, CourseID: 1543, IsAllowed: 1}, {UserID: 42407, CourseID: 1554, IsAllowed: 1}, {UserID: 42407, CourseID: 1613, IsAllowed: 1}, {UserID: 42407, CourseID: 1614, IsAllowed: 1}, {UserID: 42407, CourseID: 1615, IsAllowed: 1}, {UserID: 42407, CourseID: 1616, IsAllowed: 1}, {UserID: 42407, CourseID: 1653, IsAllowed: 1}, {UserID: 42407, CourseID: 1654, IsAllowed: 1}, {UserID: 42407, CourseID: 1664, IsAllowed: 1}, {UserID: 42407, CourseID: 1665, IsAllowed: 1}, {UserID: 42407, CourseID: 1668, IsAllowed: 1}, {UserID: 42407, CourseID: 1686, IsAllowed: 1}, {UserID: 42407, CourseID: 1688, IsAllowed: 1}, {UserID: 42407, CourseID: 1692, IsAllowed: 1}, {UserID: 42407, CourseID: 1710, IsAllowed: 1}, {UserID: 42407, CourseID: 1748, IsAllowed: 1}, {UserID: 42407, CourseID: 1767, IsAllowed: 1}, {UserID: 42407, CourseID: 1810, IsAllowed: 1}, {UserID: 42407, CourseID: 1904, IsAllowed: 1}, {UserID: 42407, CourseID: 2007, IsAllowed: 1}, {UserID: 42407, CourseID: 2033, IsAllowed: 1}, {UserID: 42407, CourseID: 2125, IsAllowed: 1}, {UserID: 42407, CourseID: 2126, IsAllowed: 1}, {UserID: 42407, CourseID: 2231, IsAllowed: 1}, {UserID: 42407, CourseID: 2259, IsAllowed: 1}}, Status: 10, Links: backoffice2.LinksFull{Self: "/user/update/42407"}}, Teachers: []backoffice2.TeacherFull{{
		ID:                 5678,
		Username:           "random_user2",
		Phone:              "+7 (922) 555-00-11",
		Email:              "randomuser2@mail.com",
		Name:               "Алексей Смирнов",
		Profile:            backoffice2.ProfileFull{PhotoURL: "", Promo: ""},
		AllowedUserCourses: []backoffice2.AllowedUserCourseFull{{UserID: 42407, CourseID: 84, IsAllowed: 1}, {UserID: 42407, CourseID: 305, IsAllowed: 1}, {UserID: 42407, CourseID: 389, IsAllowed: 1}, {UserID: 42407, CourseID: 405, IsAllowed: 1}, {UserID: 42407, CourseID: 406, IsAllowed: 1}, {UserID: 42407, CourseID: 407, IsAllowed: 1}, {UserID: 42407, CourseID: 408, IsAllowed: 1}, {UserID: 42407, CourseID: 416, IsAllowed: 1}, {UserID: 42407, CourseID: 417, IsAllowed: 1}, {UserID: 42407, CourseID: 448, IsAllowed: 1}, {UserID: 42407, CourseID: 465, IsAllowed: 1}, {UserID: 42407, CourseID: 606, IsAllowed: 1}, {UserID: 42407, CourseID: 640, IsAllowed: 1}, {UserID: 42407, CourseID: 641, IsAllowed: 1}, {UserID: 42407, CourseID: 661, IsAllowed: 1}, {UserID: 42407, CourseID: 662, IsAllowed: 1}, {UserID: 42407, CourseID: 665, IsAllowed: 1}, {UserID: 42407, CourseID: 666, IsAllowed: 1}, {UserID: 42407, CourseID: 686, IsAllowed: 1}, {UserID: 42407, CourseID: 706, IsAllowed: 1}, {UserID: 42407, CourseID: 707, IsAllowed: 1}, {UserID: 42407, CourseID: 716, IsAllowed: 1}, {UserID: 42407, CourseID: 727, IsAllowed: 1}, {UserID: 42407, CourseID: 729, IsAllowed: 1}, {UserID: 42407, CourseID: 734, IsAllowed: 1}, {UserID: 42407, CourseID: 735, IsAllowed: 1}, {UserID: 42407, CourseID: 777, IsAllowed: 1}, {UserID: 42407, CourseID: 783, IsAllowed: 1}, {UserID: 42407, CourseID: 797, IsAllowed: 1}, {UserID: 42407, CourseID: 799, IsAllowed: 1}, {UserID: 42407, CourseID: 809, IsAllowed: 1}, {UserID: 42407, CourseID: 810, IsAllowed: 1}, {UserID: 42407, CourseID: 823, IsAllowed: 1}, {UserID: 42407, CourseID: 831, IsAllowed: 1}, {UserID: 42407, CourseID: 852, IsAllowed: 1}, {UserID: 42407, CourseID: 853, IsAllowed: 1}, {UserID: 42407, CourseID: 857, IsAllowed: 1}, {UserID: 42407, CourseID: 858, IsAllowed: 1}, {UserID: 42407, CourseID: 859, IsAllowed: 1}, {UserID: 42407, CourseID: 860, IsAllowed: 1}, {UserID: 42407, CourseID: 861, IsAllowed: 1}, {UserID: 42407, CourseID: 862, IsAllowed: 1}, {UserID: 42407, CourseID: 864, IsAllowed: 1}, {UserID: 42407, CourseID: 1338, IsAllowed: 1}, {UserID: 42407, CourseID: 1339, IsAllowed: 1}, {UserID: 42407, CourseID: 1346, IsAllowed: 1}, {UserID: 42407, CourseID: 1347, IsAllowed: 1}, {UserID: 42407, CourseID: 1387, IsAllowed: 1}, {UserID: 42407, CourseID: 1459, IsAllowed: 1}, {UserID: 42407, CourseID: 1484, IsAllowed: 1}, {UserID: 42407, CourseID: 1543, IsAllowed: 1}, {UserID: 42407, CourseID: 1554, IsAllowed: 1}, {UserID: 42407, CourseID: 1613, IsAllowed: 1}, {UserID: 42407, CourseID: 1614, IsAllowed: 1}, {UserID: 42407, CourseID: 1615, IsAllowed: 1}, {UserID: 42407, CourseID: 1616, IsAllowed: 1}, {UserID: 42407, CourseID: 1653, IsAllowed: 1}, {UserID: 42407, CourseID: 1654, IsAllowed: 1}, {UserID: 42407, CourseID: 1664, IsAllowed: 1}, {UserID: 42407, CourseID: 1665, IsAllowed: 1}, {UserID: 42407, CourseID: 1668, IsAllowed: 1}, {UserID: 42407, CourseID: 1686, IsAllowed: 1}, {UserID: 42407, CourseID: 1688, IsAllowed: 1}, {UserID: 42407, CourseID: 1692, IsAllowed: 1}, {UserID: 42407, CourseID: 1710, IsAllowed: 1}, {UserID: 42407, CourseID: 1748, IsAllowed: 1}, {UserID: 42407, CourseID: 1767, IsAllowed: 1}, {UserID: 42407, CourseID: 1810, IsAllowed: 1}, {UserID: 42407, CourseID: 1904, IsAllowed: 1}, {UserID: 42407, CourseID: 2007, IsAllowed: 1}, {UserID: 42407, CourseID: 2033, IsAllowed: 1}, {UserID: 42407, CourseID: 2125, IsAllowed: 1}, {UserID: 42407, CourseID: 2126, IsAllowed: 1}, {UserID: 42407, CourseID: 2231, IsAllowed: 1}, {UserID: 42407, CourseID: 2259, IsAllowed: 1}},
		Status:             10,
		Links:              backoffice2.LinksFull{Self: "/user/update/42407"},
	}}, ClientManager: interface{}(nil), Course: backoffice2.CourseFull{
		ID:                          729,
		Name:                        "Компьютерная грамотность",
		GUID:                        "a131029b-cde5-11eb-a724-6cb31107bf10",
		Description:                 "Курс для внеурочных занятий (дополнительное образование) с детьми в возрасте 7-9 лет, на русском языке, версия 2021/2022",
		ContentType:                 "course",
		CourseType:                  backoffice2.CourseTypeFull{ID: 19, Title: "компьютерная грамотность", Code: "comp"},
		LessonsCount:                0,
		GroupLessonsAmount:          0,
		LessonsCountFormatted:       "нет модулей",
		GroupLessonsAmountFormatted: "нет уроков",
		IsDeleted:                   0,
		Links:                       backoffice2.LinksFull{Self: "/course/view/a131029b-cde5-11eb-a724-6cb31107bf10"},
	}, LanguageID: interface{}(nil), Journal: true, ShowJournal: true, ShowOnlineRoom: true, IsOnline: false, ActiveStudentCount: 8, OnlineRoomURL: "", UseClientManager: 0, DisplayLessonDurationInMinutes: 60, DeletedAt: interface{}(nil), DeletedBy: interface{}(nil), PriorityLevel: backoffice2.PriorityLevelFull{Value: "normal", Label: "Обычный приоритет", Tag: "default"}, IsFull: false, CreatedAt: "17.09.2024 10:41", CreatedBy: backoffice2.UserFull{
		ID:       1234,
		Username: "random_user",
		Phone:    "+7 (900) 555-12-34",
		Email:    "example123@mail.com",
		Name:     "Иван Иванов",
		Profile:  backoffice2.ProfileFull{PhotoURL: "/uploads/avatar/avatar/avatar_1234_1603094230-96x96.jpg", Promo: ""},
		Status:   10,
		Links:    backoffice2.LinksFull{Self: "/user/update/1234"},
	}, Related: backoffice2.RelatedFull{
		Statuses:       []backoffice2.StatusFull{{Value: 10, Label: "Активная", Tag: "success"}, {Value: 1, Label: "Не стартовала", Tag: "warning"}, {Value: 20, Label: "Идет набор", Tag: "warning"}, {Value: 30, Label: "Приостановлена", Tag: "warning"}, {Value: 0, Label: "Окончена", Tag: "default"}, {Value: 2, Label: "Развалилась", Tag: "default"}},
		Types:          []backoffice2.TypeFull{{Value: "regular", Label: "Группа", Tag: "default"}, {Value: "masterclass", Label: "Мастер-класс", Tag: "info"}, {Value: "intensive", Label: "Интенсив", Tag: "warning"}, {Value: "demo", Label: "Обучение сотрудников", Tag: "inactive"}, {Value: "individual", Label: "Индивидуальная", Tag: "default"}},
		PriorityLevels: []backoffice2.PriorityLevelFull{{Value: "normal", Label: "Обычный", Tag: "default"}, {Value: "high", Label: "Высокий", Tag: "warning"}},
	}},
}
var backofficeKidsByGroupExpected = backoffice2.NamesByGroup{
	Status: "success",
	Data: backoffice2.GroupData{Items: []backoffice2.Student{{
		ID:              70245813,
		FirstName:       "Иван",
		LastName:        "Петров",
		FullName:        "Иван Петров",
		ParentName:      "Ольга",
		Email:           "petrov_ivan@mail.ru",
		HasLaptop:       -1,
		Phone:           "+7 (915) 123-45-67",
		Age:             10,
		BirthDate:       time.Date(2014, time.December, 16, 0, 0, 0, 0, time.Local),
		HasBranchAccess: true,
		Username:        "petrov_i",
		Password:        "7605",
		LastGroup: backoffice2.Group{
			ID:             98637162,
			GroupStudentID: 6553709,
			Title:          "Библиотека 7 вс 14.00",
			Content:        "Группа по курсу КГ",
			Track:          2,
			Status:         0,
			StartTime:      time.Date(2024, time.November, 25, 10, 54, 55, 0, time.Local),
			EndTime:        time.Date(9999, time.December, 31, 0, 0, 0, 0, time.Local),
			CourseID:       729,
		},
		Groups: []backoffice2.Group(nil),
		Links:  backoffice2.Links{Self: backoffice2.SelfLink{Href: "/student/update/70245813"}},
	}}},
}
var backofficeKidViewExpected = backoffice2.KidView{
	Status: "success",
	Data: backoffice2.Student{
		ID:              70245813,
		FirstName:       "Иван",
		LastName:        "Петров",
		FullName:        "Иван Петров",
		ParentName:      "Ольга",
		Email:           "student123@example.com",
		HasLaptop:       -1,
		Phone:           "+7 (900) 123-45-67",
		Age:             10,
		BirthDate:       time.Date(2014, time.December, 16, 0, 0, 0, 0, time.Local),
		DeletedAt:       interface{}(nil),
		HasBranchAccess: true,
		Username:        "student_01",
		Password:        "7605",
		LastGroup: backoffice2.Group{
			ID:             0,
			GroupStudentID: 0,
			Title:          "",
			Content:        "",
			Track:          0,
			Status:         0,
			StartTime:      time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			EndTime:        time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			CourseID:       0,
			DeletedAt:      interface{}(nil),
		},
		Groups: []backoffice2.Group{{
			ID:             98637162,
			GroupStudentID: 6543284,
			Title:          "Библиотека 7 вс 14.00",
			Content:        "Группа по курсу КГ",
			Track:          1,
			Status:         20,
			StartTime:      time.Date(2024, time.November, 18, 11, 18, 45, 0, time.Local),
			EndTime:        time.Date(2024, time.November, 23, 16, 56, 45, 0, time.Local),
			CourseID:       729,
			DeletedAt:      interface{}(nil),
		}, {ID: 98637162, GroupStudentID: 6553709, Title: "Библиотека 7 вс 14.00", Content: "Группа по курсу КГ", Track: 2, Status: 0, StartTime: time.Date(2024, time.November, 25, 10, 54, 55, 0, time.Local), EndTime: time.Date(9999, time.December, 31, 0, 0, 0, 0, time.Local), CourseID: 729, CreatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), UpdatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), DeletedAt: interface{}(nil)}},
		Links: backoffice2.Links{Self: backoffice2.SelfLink{Href: "/student/update/70245813"}},
	},
}
