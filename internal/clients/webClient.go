package clients

import (
	"fmt"
	"time"
)

type ClientError struct {
	Code    int
	Message string
}

func (c ClientError) Error() string {
	return fmt.Sprintf("%d: %s", c.Code, c.Message)
}
func GetError(code int, message string) *ClientError {
	return &ClientError{code, message}
}

type WebClient interface {
	// GetKidsNamesByGroup получить всех детей в группе
	GetKidsNamesByGroup(cookie, group string) (*GroupResponse, *ClientError)
	// GetKidsStatsByGroup получить статистику посещения детей в группе
	GetKidsStatsByGroup(cookie, group string) (*KidsStats, *ClientError)
	// OpenLession открыть лекцию с идентификатором {lession}
	OpenLession(cookie, group, lession string) *ClientError
	// CloseLession закрыть лекцию с идентификатором {lession}
	CloseLession(cookie, group, lession string) *ClientError
	// GetKidsMessages получить новые сообщения детей на платформе
	GetKidsMessages(cookie string) (*KidsMessages, *ClientError)
	// GetAllGroupsByUser получить все группы
	GetAllGroupsByUser(cookie string) ([]AllGroupsUser, *ClientError)
}

type GroupResponse struct {
	Status string `json:"status"`
	Data   struct {
		Items []struct {
			ID              int         `json:"id"`
			FirstName       string      `json:"firstName"`
			LastName        string      `json:"lastName"`
			FullName        string      `json:"fullName"`
			ParentName      string      `json:"parentName"`
			Email           string      `json:"email"`
			HasLaptop       int         `json:"hasLaptop"`
			Phone           string      `json:"phone"`
			Age             int         `json:"age"`
			BirthDate       time.Time   `json:"birthDate"`
			CreatedAt       time.Time   `json:"createdAt"`
			UpdatedAt       time.Time   `json:"updatedAt"`
			DeletedAt       interface{} `json:"deletedAt"`
			HasBranchAccess bool        `json:"hasBranchAccess"`
			Username        string      `json:"username"`
			Password        string      `json:"password"`
			LastGroup       struct {
				ID             int         `json:"id"`
				GroupStudentID int         `json:"groupStudentId"`
				Title          string      `json:"title"`
				Content        string      `json:"content"`
				Track          int         `json:"track"`
				Status         int         `json:"status"`
				StartTime      time.Time   `json:"startTime"`
				EndTime        time.Time   `json:"endTime"`
				CourseID       int         `json:"courseId"`
				CreatedAt      time.Time   `json:"createdAt"`
				UpdatedAt      time.Time   `json:"updatedAt"`
				DeletedAt      interface{} `json:"deletedAt"`
			} `json:"lastGroup"`
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"_links"`
		} `json:"items"`
	} `json:"data"`
}

type KidsStats struct {
	Status string `json:"status"`
	Data   []struct {
		StudentID  int `json:"student_id"`
		Attendance []struct {
			LessonID           int    `json:"lesson_id"`
			LessonTitle        string `json:"lesson_title"`
			StartTimeFormatted string `json:"start_time_formatted"`
			Status             string `json:"status"`
		} `json:"attendance"`
	} `json:"data"`
}

type KidsMessages struct {
	Status string `json:"status"`
	Data   struct {
		Projects []struct {
			UID         string `json:"uid"`
			New         bool   `json:"new"`
			SenderID    int    `json:"senderId"`
			SenderScope string `json:"senderScope"`
			Type        string `json:"type"`
			Content     string `json:"content"`
			Name        string `json:"name"`
			LastTime    string `json:"lastTime"`
			Title       string `json:"title"`
			Link        string `json:"link"`
		} `json:"projects"`
	} `json:"data"`
}

type AllGroupsUser struct {
	Title       string
	GroupId     string
	TimeLesson  string
	RegularTime string
}
