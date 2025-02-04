package service

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"tgbot/internal/clients"
	"tgbot/internal/domain"
	appError "tgbot/internal/error"
	"tgbot/internal/helpers"
	"tgbot/internal/models"
	"time"
)

var dateMap map[string]string = map[string]string{
	"янв":  "01",
	"февр": "02",
	"мар":  "03",
	"апр":  "04",
	"мая":  "05",
	"июн":  "06",
	"июл":  "07",
	"авг":  "08",
	"сент": "09",
	"окт":  "10",
	"нояб": "11",
	"дек":  "12",
}

type DefaultService struct {
	domain    domain.Domain
	webClient clients.WebClient
}

func NewDefaultService(domain domain.Domain, webClient clients.WebClient) *DefaultService {
	return &DefaultService{domain: domain, webClient: webClient}
}

func (d DefaultService) UsersByNotif(status bool) ([]models.ScheduleData, error) {
	notif := 0
	if status {
		notif = 1
	}

	notification, err := d.domain.GetUsersByNotification(notif)
	if err != nil {
		return nil, fmt.Errorf("DefaultService.UserUidsByNotif(%v) : %w", status, err)
	}

	data := make([]models.ScheduleData, len(notification))
	for i, user := range notification {
		data[i] = models.ScheduleData{
			UID:    user.UID,
			Cookie: user.Cookie,
		}
	}
	return data, nil
}

func (d DefaultService) NewMessageByUID(uid int64) ([]models.Message, error) {
	cookie, err := d.Cookie(uid)
	if err != nil {
		return nil, fmt.Errorf("DefaultService.NewMessageByUID(%d) : %w", uid, err)
	}
	messages, err := d.webClient.GetKidsMessages(cookie)
	if err != nil {
		return nil, fmt.Errorf("DefaultService.NewMessageByUID(%d) : %w", uid, err)
	}

	var msgs []models.Message
	lastNotif, err := d.domain.LastNotificationDate(uid)
	var lastNotifData time.Time
	if err != nil {
		if !errors.Is(err, appError.ErrNotValid) {
			return nil, fmt.Errorf("DefaultService.NewMessageByUID(%d) : %w", uid, err)
		}
		lastNotifData = time.Time{}
	} else {
		lastNotifData = parseDate(lastNotif)
	}
	var lastNotifString string
	for i := len(messages.Data.Projects) - 1; i >= 0; i-- {
		if messages.Data.Projects[i].SenderScope == "student" {
			if lastNotifData.Before(parseDate(messages.Data.Projects[i].LastTime)) {
				m := models.Message{
					Id:      messages.Data.Projects[i].UID,
					Type:    messages.Data.Projects[i].Type,
					From:    messages.Data.Projects[i].Name,
					Theme:   messages.Data.Projects[i].Title,
					Link:    fmt.Sprintf("https://backoffice.algoritmika.org%s", messages.Data.Projects[i].Link),
					Content: messages.Data.Projects[i].Content,
				}
				if m.Type == "img" {
					m.Content = fmt.Sprintf("https://backoffice.algoritmika.org%s", m.Content)
				}
				msgs = append(msgs, m)
				lastNotifString = messages.Data.Projects[i].LastTime
			}
		}
	}
	if lastNotifString != "" {
		err := d.domain.SetLastNotificationDate(uid, lastNotifString)
		if err != nil {
			return nil, fmt.Errorf("DefaultService.NewMessageByUID(%d) : %w", uid, err)
		}
	}

	return msgs, nil
}

func parseDate(lastTime string) time.Time {
	parts := strings.Split(lastTime, " ")

	day := parts[0]
	month := dateMap[strings.Replace(parts[1], ".", "", -1)]

	if len(parts) == 3 {
		timeHour := parts[2]

		retTime, _ := time.Parse("2 01 2006 15:04", fmt.Sprintf("%s %s %d %s", day, month, time.Now().Year(), timeHour))
		return retTime
	}
	if len(parts) == 4 {
		year := strings.Replace(parts[2], "`", "", -1)
		year = strings.Replace(year, ",", "", -1)
		timeHour := parts[3]

		retTime, _ := time.Parse("2 01 06 15:04", fmt.Sprintf("%s %s %s %s", day, month, year, timeHour))
		return retTime
	}
	return time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
}

func (d DefaultService) OpenLesson(uid int64, groupId int, lessonId int) error {
	cookie, err := d.Cookie(uid)
	if err != nil {
		return fmt.Errorf("DefaultService.OpenLesson(%d, %d, %d) : %w", uid, lessonId, groupId, err)
	}
	err = d.webClient.OpenLession(cookie, strconv.Itoa(groupId), strconv.Itoa(lessonId))
	if err != nil {
		return fmt.Errorf("DefaultService.OpenLesson(%d, %d, %d) : %w", uid, lessonId, groupId, err)
	}

	return nil
}
func (d DefaultService) CloseLesson(uid int64, groupId int, lessonId int) error {
	cookie, err := d.Cookie(uid)
	if err != nil {
		return fmt.Errorf("DefaultService.CloseLesson(%d, %d, %d) : %w", uid, lessonId, groupId, err)
	}
	err = d.webClient.CloseLession(cookie, strconv.Itoa(groupId), strconv.Itoa(lessonId))
	if err != nil {
		return fmt.Errorf("DefaultService.CloseLesson(%d, %d, %d) : %w", uid, lessonId, groupId, err)
	}

	return nil
}

func (d DefaultService) AllCredentials(uid int64, groupId int) (map[string]string, error) {
	names, err := d.AllKidsNames(uid, groupId)
	if err != nil {
		return nil, fmt.Errorf("DefaultService.AllCredentials(%d, %d) : %w", uid, groupId, err)
	}

	creds := make(map[string]string, len(names))
	for _, kid := range names {
		creds[kid.FullName] = fmt.Sprintf("%s:%s", kid.Login, kid.Password)
	}

	return creds, nil
}

func (d DefaultService) Groups(uid int64) ([]models.Group, error) {
	data, err := d.domain.Groups(uid)
	if err != nil {
		return nil, fmt.Errorf("DefaultService.Groups(%d) : %w", uid, err)
	}
	return models.GroupMap(data), nil
}

func (d DefaultService) Cookie(uid int64) (string, error) {
	c, err := d.domain.Cookie(uid)
	if err != nil {
		if errors.Is(err, appError.ErrNotValid) {
			return "", nil
		}
		return "", fmt.Errorf("DefaultService.Cookie(%d) : %w", uid, err)
	}
	return c, nil
}
func (d DefaultService) SetCookie(uid int64, cookie string) error {
	err := d.domain.SetCookie(uid, cookie)
	if err != nil {
		return fmt.Errorf("DefaultService.SetCookie(%d, %s) : %w", uid, cookie, err)
	}
	return nil
}

func (d DefaultService) Notification(uid int64) (bool, error) {
	n, err := d.domain.Notification(uid)
	if err != nil {
		if errors.Is(err, appError.ErrNotValid) {
			return false, nil
		}
		return false, fmt.Errorf("DefaultService.Notification(%d) : %w", uid, err)
	}
	return n, nil
}
func (d DefaultService) SetNotification(uid int64, notification bool) error {
	err := d.domain.SetNotification(uid, notification)
	if err != nil {
		return fmt.Errorf("DefaultService.SetNotification(%d, %v) : %w", uid, notification, err)
	}
	return nil
}

func (d DefaultService) CurrentGroup(uid int64, t time.Time) (models.Group, error) {
	allGroups, err := d.Groups(uid)
	if err != nil {
		return models.Group{}, fmt.Errorf("DefaultService.CurrentGroup(%d, %v) : %w", uid, t, err)
	}

	group, err := helpers.GetCurrentGroup(t, allGroups)
	if err != nil {
		return models.Group{}, fmt.Errorf("DefaultService.CurrentGroup(%d, %v) : %w", uid, t, err)
	}

	return group, nil
}
func (d DefaultService) ActualInformation(uid int64, t time.Time, groupId int) (models.ActualInformation, error) {
	cookie, err := d.Cookie(uid)
	if err != nil {
		return models.ActualInformation{}, fmt.Errorf("DefaultService.ActualInformation(%d, %v, %d) : %w", uid, t, groupId, err)
	}
	stats, err := d.webClient.GetKidsStatsByGroup(cookie, strconv.Itoa(groupId))
	if err != nil {
		return models.ActualInformation{}, fmt.Errorf("DefaultService.ActualInformation(%d, %v, %d) : %w", uid, t, groupId, err)
	}

	actual := models.ActualInformation{}
	for _, datum := range stats.Data {
		studentID := datum.StudentID
		for _, attendance := range datum.Attendance {
			if matchDates(attendance.StartTimeFormatted, t) {
				actual.LessonTitle = attendance.LessonTitle
				actual.LessonId = attendance.LessonID

				if attendance.Status == "absent" {
					actual.MissingKids = append(actual.MissingKids, studentID)
					break
				}
			}
		}
	}

	return actual, nil
}
func (d DefaultService) AllKidsNames(uid int64, groupId int) (models.AllKids, error) {
	cookie, err := d.Cookie(uid)
	if err != nil {
		return nil, fmt.Errorf("DefaultService.AllKidsNames(%d, %d) : %w", uid, groupId, err)
	}
	group, err := d.webClient.GetKidsNamesByGroup(cookie, groupId)
	if err != nil {
		return nil, fmt.Errorf("DefaultService.AllKidsNames(%d, %d) : %w", uid, groupId, err)
	}

	names := make(map[int]models.KidData, len(group.Data.Items))
	for _, datum := range group.Data.Items {
		if datum.LastGroup.ID == groupId && datum.LastGroup.Status == 0 {
			names[datum.ID] = models.KidData{
				FullName: datum.FullName,
				Login:    datum.Username,
				Password: datum.Password,
			}
		}
	}

	return names, nil
}

func (d DefaultService) RefreshGroups(uid int64) error {
	cookie, err := d.Cookie(uid)
	if err != nil {
		return fmt.Errorf("DefaultService.RefreshGroups(%d) : %w", uid, err)
	}
	groups, err := d.webClient.GetAllGroupsByUser(cookie)
	if err != nil {
		return fmt.Errorf("DefaultService.RefreshGroups(%d) : %w", uid, err)
	}

	groupsFormatted := make([]domain.Group, len(groups))
	for i, group := range groups {
		groupIdStr := group.GroupId
		groupIdInt, err := strconv.Atoi(groupIdStr)

		if err != nil {
			return fmt.Errorf("DefaultService.RefreshGroups(%d) : %w", uid, err)
		}

		groupsFormatted[i] = domain.Group{
			GroupID:    groupIdInt,
			Title:      group.Title,
			TimeLesson: getTime(group.TimeLesson),
		}
	}

	err = d.domain.SetGroups(uid, groupsFormatted)
	if err != nil {
		return fmt.Errorf("DefaultService.RefreshGroups(%d) : %w", uid, err)
	}

	return nil
}

func (d DefaultService) RegisterUser(uid int64) error {
	err := d.domain.RegisterUser(uid)
	if err != nil {
		return fmt.Errorf("DefaultService.RegisterUser(%d) : %w", uid, err)
	}
	return nil
}
func (d DefaultService) IsUserRegistered(uid int64) (bool, error) {
	_, err := d.domain.User(uid)
	if err != nil {
		if errors.Is(err, appError.ErrNotValid) {
			return false, nil
		}
		return false, fmt.Errorf("DefaultService.IsUserRegistered(%d) : %w", uid, err)
	}
	return true, nil
}

func matchDates(timeStr string, t time.Time) bool {
	timeStr = regexp.MustCompile(`^[а-яА-Я]+(\s+)?`).ReplaceAllString(timeStr, "")
	timeFormatted, err := time.Parse("02.01.06 15:04", timeStr)
	if err != nil {
		log.Printf("Cant convert date str to Time - '%s'\n", timeStr)
		return false
	}

	if t.YearDay() == timeFormatted.YearDay() {
		return true
	}
	return false
}
func getTime(lesson string) time.Time {
	parse, err := time.Parse("02.01.2006 15:04", lesson)
	if err != nil {
		return time.Time{}
	}
	return parse
}
