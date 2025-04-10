package backoffice

type GroupInfo struct {
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
