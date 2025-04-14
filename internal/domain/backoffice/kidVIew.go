package backoffice

import "time"

type KidView struct {
	Status string  `json:"status"`
	Data   Student `json:"data"`
}

type GroupKidInfo struct {
	ID             int       `json:"id"`
	GroupStudentID int       `json:"groupStudentId"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	Track          int       `json:"track"`
	Status         int       `json:"status"`
	StartTime      time.Time `json:"startTime"`
	EndTime        time.Time `json:"endTime"`
	CourseID       int       `json:"courseId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	DeletedAt      any       `json:"deletedAt"`
}
