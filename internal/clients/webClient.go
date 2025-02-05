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
	GetKidsNamesByGroup(cookie string, group int) (*GroupResponse, error)
	// GetKidsStatsByGroup получить статистику посещения детей в группе
	GetKidsStatsByGroup(cookie, group string) (*KidsStats, error)
	// OpenLession открыть лекцию с идентификатором {lession}
	OpenLession(cookie, group, lession string) error
	// CloseLession закрыть лекцию с идентификатором {lession}
	CloseLession(cookie, group, lession string) error
	// GetKidsMessages получить новые сообщения детей на платформе
	GetKidsMessages(cookie string) (*KidsMessages, error)
	// GetAllGroupsByUser получить все группы
	GetAllGroupsByUser(cookie string) ([]AllGroupsUser, error)
	// GetGroupInfo получить информацию о группе
	GetGroupInfo(cookie string, group string) (*FullGroupInfo, error)
	// GetKidInfo получить информацию о ребенке
	GetKidInfo(cookie string, kidID string) (*FullKidInfo, error)
}

type GroupResponse struct {
	Status string    `json:"status"`
	Data   GroupData `json:"data"`
}

type GroupData struct {
	Items []Student `json:"items"`
}

type Student struct {
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
	LastGroup       LastGroup   `json:"lastGroup"`
	Links           Links       `json:"_links"`
}

type LastGroup struct {
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

// ////////////////////
type StatusFull struct {
	Value int    `json:"value"`
	Label string `json:"label"`
	Tag   string `json:"tag"`
}

type TypeFull struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Tag   string `json:"tag"`
}

type ProfileFull struct {
	PhotoURL string `json:"photo_url"`
	Promo    string `json:"promo"`
}

type LinksFull struct {
	Self string `json:"self"`
}

type BranchFull struct {
	ID                            int    `json:"id"`
	Title                         string `json:"title"`
	Code                          string `json:"code"`
	Description                   string `json:"description"`
	Phone                         string `json:"phone"`
	Email                         string `json:"email"`
	SiteURL                       string `json:"site_url"`
	TemplateVersion               int    `json:"templateVersion"`
	UseAmo                        bool   `json:"use_amo"`
	AmoConfigID                   int    `json:"amoConfigId"`
	ShowFinanceInfo               bool   `json:"show_finance_info"`
	LmsDisplayStudentCredentials  bool   `json:"lms_display_student_credentials"`
	ShowOnlineRoomURLField        int    `json:"show_online_room_url_field"`
	UseSms                        bool   `json:"use_sms"`
	LanguageID                    int    `json:"language_id"`
	OrderName                     int    `json:"order_name"`
	UseFullyPaidLabel             int    `json:"use_fully_paid_label"`
	BrandName                     string `json:"brandName"`
	MaxCountStudentsForShowOnline int    `json:"max_count_students_for_show_online"`
	IsFillPaymentSystem           bool   `json:"isFillPaymentSystem"`
	FirstLessonNoRoyalty          int    `json:"firstLessonNoRoyalty"`
	RootBranchID                  int    `json:"root_branch_id"`
}

type VenueFull struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Address      string    `json:"address"`
	ContactName  string    `json:"contact_name"`
	ContactEmail string    `json:"contact_email"`
	ContactPhone string    `json:"contact_phone"`
	Links        LinksFull `json:"_links"`
}

type UserFull struct {
	ID       int         `json:"id"`
	Username string      `json:"username"`
	Phone    string      `json:"phone"`
	Email    string      `json:"email"`
	Name     string      `json:"name"`
	Profile  ProfileFull `json:"profile"`
	Status   int         `json:"status"`
	Links    LinksFull   `json:"_links"`
}

type TeacherFull struct {
	ID                 int                     `json:"id"`
	Username           string                  `json:"username"`
	Phone              string                  `json:"phone"`
	Email              string                  `json:"email"`
	Name               string                  `json:"name"`
	Profile            ProfileFull             `json:"profile"`
	AllowedUserCourses []AllowedUserCourseFull `json:"allowedUserCourses"`
	Status             int                     `json:"status"`
	Links              LinksFull               `json:"_links"`
}

type AllowedUserCourseFull struct {
	UserID    int `json:"userId"`
	CourseID  int `json:"courseId"`
	IsAllowed int `json:"isAllowed"`
}

type CourseTypeFull struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Code  string `json:"code"`
}

type CourseFull struct {
	ID                          int            `json:"id"`
	Name                        string         `json:"name"`
	GUID                        string         `json:"guid"`
	Description                 string         `json:"description"`
	ContentType                 string         `json:"contentType"`
	CourseType                  CourseTypeFull `json:"courseType"`
	LessonsCount                int            `json:"lessons_count"`
	GroupLessonsAmount          int            `json:"group_lessons_amount"`
	LessonsCountFormatted       string         `json:"lessons_count_formatted"`
	GroupLessonsAmountFormatted string         `json:"group_lessons_amount_formatted"`
	IsDeleted                   int            `json:"is_deleted"`
	Links                       LinksFull      `json:"_links"`
}

type PriorityLevelFull struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Tag   string `json:"tag"`
}

type RelatedFull struct {
	Statuses       []StatusFull        `json:"statuses"`
	Types          []TypeFull          `json:"types"`
	PriorityLevels []PriorityLevelFull `json:"priority_levels"`
}

type FullGroupInfo struct {
	Status string        `json:"status"`
	Data   GroupDataFull `json:"data"`
}

type GroupDataFull struct {
	ID                             int               `json:"id"`
	Title                          string            `json:"title"`
	Content                        string            `json:"content"`
	Type                           TypeFull          `json:"type"`
	Status                         StatusFull        `json:"status"`
	StatusChangedAt                string            `json:"status_changed_at"`
	StartTime                      string            `json:"start_time"`
	NextLessonTime                 string            `json:"next_lesson_time"`
	LessonsTotal                   int               `json:"lessons_total"`
	LessonsPassed                  int               `json:"lessons_passed"`
	HardwareNeeded                 int               `json:"hardware_needed"`
	Branch                         BranchFull        `json:"branch"`
	Venue                          VenueFull         `json:"venue"`
	Curator                        UserFull          `json:"curator"`
	Teacher                        TeacherFull       `json:"teacher"`
	Teachers                       []TeacherFull     `json:"teachers"`
	ClientManager                  interface{}       `json:"client_manager"`
	Course                         CourseFull        `json:"course"`
	LanguageID                     interface{}       `json:"language_id"`
	Journal                        bool              `json:"journal"`
	ShowJournal                    bool              `json:"show_journal"`
	ShowOnlineRoom                 bool              `json:"showOnlineRoom"`
	IsOnline                       bool              `json:"isOnline"`
	ActiveStudentCount             int               `json:"active_student_count"`
	OnlineRoomURL                  string            `json:"online_room_url"`
	UseClientManager               int               `json:"use_client_manager"`
	DisplayLessonDurationInMinutes int               `json:"display_lesson_duration_in_minutes"`
	DeletedAt                      interface{}       `json:"deleted_at"`
	DeletedBy                      interface{}       `json:"deleted_by"`
	PriorityLevel                  PriorityLevelFull `json:"priority_level"`
	IsFull                         bool              `json:"is_full"`
	CreatedAt                      string            `json:"created_at"`
	CreatedBy                      UserFull          `json:"created_by"`
	Related                        RelatedFull       `json:"_related"`
}

// FullKidInfo
type LinksKidInfo struct {
	Self struct {
		Href string `json:"href"`
	} `json:"self"`
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

type DataKidInfo struct {
	ID              int            `json:"id"`
	FirstName       string         `json:"firstName"`
	LastName        string         `json:"lastName"`
	FullName        string         `json:"fullName"`
	ParentName      string         `json:"parentName"`
	Email           string         `json:"email"`
	HasLaptop       int            `json:"hasLaptop"`
	Phone           string         `json:"phone"`
	Age             int            `json:"age"`
	BirthDate       time.Time      `json:"birthDate"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       any            `json:"deletedAt"`
	HasBranchAccess bool           `json:"hasBranchAccess"`
	Username        string         `json:"username"`
	Password        string         `json:"password"`
	Groups          []GroupKidInfo `json:"groups"`
	Links           LinksKidInfo   `json:"_links"`
}
type FullKidInfo struct {
	Status string      `json:"status"`
	Data   DataKidInfo `json:"data"`
}
