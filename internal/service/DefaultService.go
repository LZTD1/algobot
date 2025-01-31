package service

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"tgbot/internal/clients"
	"tgbot/internal/domain"
	appError "tgbot/internal/error"
	"tgbot/internal/helpers"
	"time"
)

type DefaultService struct {
	domain    domain.Domain
	webClient clients.WebClient
}

func NewDefaultService(domain domain.Domain, webClient clients.WebClient) *DefaultService {
	return &DefaultService{domain: domain, webClient: webClient}
}

func (d DefaultService) CurrentGroup(uid int64, t time.Time) (domain.Group, error) {
	// TODO зарефачить метод
	// TODO выделить 2 разных обьекта для Domain в Service
	// TODO Зарефачить схему БД

	cookie, err := d.Cookie(uid)

	allGroups, err := d.domain.Groups(uid)
	if err != nil {
		return domain.Group{}, fmt.Errorf("DefaultService.CurrentGroup(%d, %v) : %w", uid, t, err)
	}

	group, err := helpers.GetCurrentGroup(t, allGroups)
	if err != nil {
		return domain.Group{}, fmt.Errorf("DefaultService.CurrentGroup(%d, %v) : %w", uid, t, err)
	}

	names, err := d.KidsNamesMap(cookie, group.Id)
	group.AllKids = len(names)

	stats, err := d.webClient.GetKidsStatsByGroup(cookie, strconv.FormatInt(int64(group.Id), 10))
	if err != nil {
		return domain.Group{}, fmt.Errorf("DefaultService.CurrentGroup(%d, %v) : %w", uid, t, err)
	}

	var absentKids []string
	var lession string
	for _, datum := range stats.Data {
		for _, attendance := range datum.Attendance {
			if attendance.Status == "absent" && matchDates(attendance.StartTimeFormatted, t) {
				absentKids = append(absentKids, names[datum.StudentID])
				lession = attendance.LessonTitle
				break
			}
		}
	}

	group.MissingKids = absentKids
	group.Lesson = lession

	return group, nil
}

func (d DefaultService) Groups(uid int64) ([]domain.Group, error) {
	data, err := d.domain.Groups(uid)
	if err != nil {
		return nil, fmt.Errorf("DefaultService.Groups(%d) : %w", uid, err)
	}
	return data, nil
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

func (d DefaultService) RegisterUser(uid int64) error {
	err := d.domain.RegisterUser(uid)
	if err != nil {
		return fmt.Errorf("DefaultService.RegisterUser(%d) : %w", uid, err)
	}
	return nil
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
			Id:          groupIdInt,
			Name:        group.Title,
			Lesson:      "",
			Time:        getTime(group.TimeLesson),
			AllKids:     0,
			MissingKids: nil,
		}
	}
	err = d.domain.SetGroups(uid, groupsFormatted)
	if err != nil {
		return fmt.Errorf("DefaultService.RefreshGroups(%d) : %w", uid, err)
	}

	return nil
}

func (d DefaultService) KidsNamesMap(cookie string, groupId int) (map[int]string, error) {
	names, err := d.webClient.GetKidsNamesByGroup(cookie, groupId)
	if err != nil {
		return nil, fmt.Errorf("DefaultService.KidsNamesMap(%s, %d)  : %w", cookie, groupId, err)
	}
	returnMap := make(map[int]string)
	for _, item := range names.Data.Items {
		if item.LastGroup.ID == groupId && item.LastGroup.Status == 0 {
			returnMap[item.ID] = item.FullName
		}
	}
	return returnMap, nil
}

func getTime(lesson string) time.Time {
	parse, err := time.Parse("02.01.2006 15:04", lesson)
	if err != nil {
		return time.Time{}
	}
	return parse
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
