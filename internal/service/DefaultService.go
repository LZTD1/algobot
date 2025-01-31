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
	allGroups, err := d.domain.Groups(uid)
	if err != nil {
		return domain.Group{}, fmt.Errorf("DefaultService.CurrentGroup(%d, %v) : %w", uid, t, err)
	}

	group, err := helpers.GetCurrentGroup(t, allGroups)
	if err != nil {
		return domain.Group{}, fmt.Errorf("DefaultService.CurrentGroup(%d, %v) : %w", uid, t, err)
	}

	kids, err := d.MissingKids(uid, t, group.Id)
	if err != nil {
		return domain.Group{}, fmt.Errorf("DefaultService.CurrentGroup(%d, %v) : %w", uid, t, err)
	}

	group.MissingKids = kids

	return group, nil
}

func (d DefaultService) Groups(uid int64) ([]domain.Group, error) {
	data, err := d.domain.Groups(uid)
	if err != nil {
		return nil, fmt.Errorf("DefaultService.Groups(%d) : %w", uid, err)
	}
	return data, nil
}

func (d DefaultService) MissingKids(uid int64, t time.Time, g int) ([]string, error) {
	cookie, err := d.Cookie(uid)
	if err != nil {
		return nil, fmt.Errorf("DefaultService.MissingKids(%d, %v, %d) : %w", uid, t, g, err)
	}

	names, err := d.webClient.GetKidsNamesByGroup(cookie, strconv.FormatInt(int64(g), 10))
	if err != nil {
		return nil, fmt.Errorf("DefaultService.MissingKids(%d, %v, %d) : %w", uid, t, g, err)
	}

	stats, err := d.webClient.GetKidsStatsByGroup(cookie, strconv.FormatInt(int64(g), 10))
	if err != nil {
		return nil, fmt.Errorf("DefaultService.MissingKids(%d, %v, %d) : %w", uid, t, g, err)
	}

	absentKids := make(map[int]string)
	for _, datum := range stats.Data {
		for _, attendance := range datum.Attendance {
			if attendance.Status == "absent" && matchDates(attendance.StartTimeFormatted, t) {
				absentKids[datum.StudentID] = ""
				break
			}
		}
	}

	var readyNames []string
	for _, item := range names.Data.Items {
		if _, exists := absentKids[item.ID]; exists {
			readyNames = append(readyNames, item.FullName)
		}
	}

	return readyNames, nil
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

	fmt.Println(timeStr)
	fmt.Println(t)
	fmt.Println()
	fmt.Println()
	fmt.Println()
	if t.YearDay() == timeFormatted.YearDay() {
		return true
	}
	return false
}
