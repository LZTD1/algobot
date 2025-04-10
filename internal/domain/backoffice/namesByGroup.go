package backoffice

import "time"

type NamesByGroup struct {
	Status string    `json:"status"`
	Data   GroupData `json:"data"`
}
type GroupData struct {
	Items []Student `json:"items"`
}

type Group struct {
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
}

type Links struct {
	Self SelfLink `json:"self"`
}

type SelfLink struct {
	Href string `json:"href"`
}

type KidsStats struct {
	Status string    `json:"status"`
	Data   []KidStat `json:"data"`
}

type KidStat struct {
	StudentID  int          `json:"student_id"`
	Attendance []Attendance `json:"attendance"`
}

type Attendance struct {
	LessonID           int    `json:"lesson_id"`
	LessonTitle        string `json:"lesson_title"`
	StartTimeFormatted string `json:"start_time_formatted"`
	Status             string `json:"status"`
}

type KidsMessages struct {
	Status string       `json:"status"`
	Data   MessagesData `json:"data"`
}

type MessagesData struct {
	Projects []Message `json:"projects"`
}

type Message struct {
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
}

type AllGroupsUser struct {
	Title       string
	GroupId     string
	TimeLesson  string
	RegularTime string
}
