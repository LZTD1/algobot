package mocks

import (
	"tgbot/internal/clients"
	"time"
)

type MockWebClient struct {
}

func (m MockWebClient) GetKidsNamesByGroup(cookie, group string) (*clients.GroupResponse, error) {
	return &groupResponse, nil
}

func (m MockWebClient) GetKidsStatsByGroup(cookie, group string) (*clients.KidsStats, error) {
	return &kidsStats, nil
}

func (m MockWebClient) OpenLession(cookie, group, lession string) error {
	//TODO implement me
	panic("implement me")
}

func (m MockWebClient) CloseLession(cookie, group, lession string) error {
	//TODO implement me
	panic("implement me")
}

func (m MockWebClient) GetKidsMessages(cookie string) (*clients.KidsMessages, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockWebClient) GetAllGroupsByUser(cookie string) ([]clients.AllGroupsUser, error) {
	return []clients.AllGroupsUser{
		{
			Title:       "Title",
			GroupId:     "1",
			TimeLesson:  "01.02.2025 14:00",
			RegularTime: "4",
		},
	}, nil
}

var groupResponse = clients.GroupResponse{
	Status: "success",
	Data: clients.GroupData{
		Items: []clients.Student{
			{
				ID:              1,
				FirstName:       "Иван",
				LastName:        "Иванов",
				FullName:        "Иван Иванов",
				ParentName:      "Алексей Иванов",
				Email:           "ivanov@example.com",
				HasLaptop:       1,
				Phone:           "+79161234567",
				Age:             16,
				BirthDate:       time.Date(2008, time.Month(5), 10, 0, 0, 0, 0, time.UTC),
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				DeletedAt:       nil,
				HasBranchAccess: true,
				Username:        "ivan_ivanov",
				Password:        "secret_password_123",
				LastGroup: clients.LastGroup{
					ID:             101,
					GroupStudentID: 1,
					Title:          "Группа 1",
					Content:        "Основы программирования",
					Track:          1,
					Status:         1,
					StartTime:      time.Date(2025, time.Month(1), 15, 10, 0, 0, 0, time.UTC),
					EndTime:        time.Date(2025, time.Month(5), 15, 18, 0, 0, 0, time.UTC),
					CourseID:       201,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
					DeletedAt:      nil,
				},
				Links: clients.Links{
					Self: clients.SelfLink{
						Href: "http://example.com/students/1",
					},
				},
			},
			{
				ID:              2,
				FirstName:       "Мария",
				LastName:        "Петрова",
				FullName:        "Мария Петрова",
				ParentName:      "Елена Петрова",
				Email:           "petrova@example.com",
				HasLaptop:       0,
				Phone:           "+79261234567",
				Age:             15,
				BirthDate:       time.Date(2009, time.Month(7), 20, 0, 0, 0, 0, time.UTC),
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				DeletedAt:       nil,
				HasBranchAccess: false,
				Username:        "maria_petrov",
				Password:        "password_321",
				LastGroup: clients.LastGroup{
					ID:             102,
					GroupStudentID: 2,
					Title:          "Группа 2",
					Content:        "Алгоритмы и структуры данных",
					Track:          2,
					Status:         1,
					StartTime:      time.Date(2025, time.Month(2), 1, 10, 0, 0, 0, time.UTC),
					EndTime:        time.Date(2025, time.Month(6), 1, 18, 0, 0, 0, time.UTC),
					CourseID:       202,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
					DeletedAt:      nil,
				},
				Links: clients.Links{
					Self: clients.SelfLink{
						Href: "http://example.com/students/2",
					},
				},
			},
		},
	},
}

var kidsStats = clients.KidsStats{
	Status: "success",
	Data: []clients.KidStat{
		{
			StudentID: 1,
			Attendance: []clients.Attendance{
				{
					LessonID:           1,
					LessonTitle:        "Алгоритмы и структуры данных",
					StartTimeFormatted: "вс 22.09.24 14:00",
					Status:             "present",
				},
				{
					LessonID:           2,
					LessonTitle:        "Алгоритмы и структуры данных",
					StartTimeFormatted: "вс 29.09.24 14:00",
					Status:             "present",
				},
				{
					LessonID:           3,
					LessonTitle:        "Алгоритмы и структуры данных",
					StartTimeFormatted: "вс 06.10.24 14:00",
					Status:             "present",
				},
			},
		},
		{
			StudentID: 2,
			Attendance: []clients.Attendance{
				{
					LessonID:           1,
					LessonTitle:        "Алгоритмы и структуры данных",
					StartTimeFormatted: "вс 22.09.24 14:00",
					Status:             "present",
				},
				{
					LessonID:           2,
					LessonTitle:        "Алгоритмы и структуры данных",
					StartTimeFormatted: "вс 29.09.24 14:00",
					Status:             "present",
				},
				{
					LessonID:           3,
					LessonTitle:        "Алгоритмы и структуры данных",
					StartTimeFormatted: "вс 06.10.24 14:00",
					Status:             "absent",
				},
			},
		},
	},
}
