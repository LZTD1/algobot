package service

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"tgbot/internal/clients"
	"tgbot/internal/domain"
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
		return domain.Group{}, err
	}

	group, err := helpers.GetCurrentGroup(t, allGroups)
	if err != nil {
		return domain.Group{}, err
	}

	fmt.Println(group)
	kids, err := d.MissingKids(uid, t, group.Id)
	if err != nil {
		return domain.Group{}, err
	}

	group.MissingKids = kids

	return group, nil
}

func (d DefaultService) Groups(uid int64) ([]domain.Group, error) {
	return d.domain.Groups(uid)
}

func (d DefaultService) MissingKids(uid int64, t time.Time, g int) ([]string, error) {
	cookie, err := d.domain.Cookie(uid)
	if err != nil {
		return nil, err
	}

	names, err := d.webClient.GetKidsNamesByGroup(cookie, strconv.FormatInt(int64(g), 10))
	if err != nil {
		return nil, err
	}

	fmt.Println(names)
	stats, err := d.webClient.GetKidsStatsByGroup(cookie, strconv.FormatInt(int64(g), 10))
	if err != nil {
		return nil, err
	}
	fmt.Println(stats)

	absentKids := make(map[int]string)
	for _, datum := range stats.Data {
		for _, attendance := range datum.Attendance {
			if attendance.Status == "absent" && matchDates(attendance.StartTimeFormatted, t) {
				absentKids[datum.StudentID] = ""
				break
			}
		}
	}

	fmt.Println(absentKids)
	var readyNames []string
	for _, item := range names.Data.Items {
		if _, exists := absentKids[item.ID]; exists {
			readyNames = append(readyNames, item.FullName)
		}
	}

	return readyNames, nil
}

func (d DefaultService) Cookie(uid int64) (string, error) {
	return d.domain.Cookie(uid)
}

func (d DefaultService) SetCookie(uid int64, cookie string) {
	d.domain.SetCookie(uid, cookie)
}

func (d DefaultService) Notification(uid int64) bool {
	return d.domain.Notification(uid)
}

func (d DefaultService) SetNotification(uid int64, notification bool) {
	d.domain.SetNotification(uid, notification)
}

func (d DefaultService) IsUserRegistered(uid int64) bool {
	_, err := d.domain.User(uid)
	if err != nil {
		return false
	}
	return true
}

func (d DefaultService) RegisterUser(uid int64) {
	d.domain.RegisterUser(uid)
}

func (d DefaultService) RefreshGroups(uid int64) error {
	cookie, err := d.domain.Cookie(uid)
	if err != nil {
		return err
	}
	groups, err := d.webClient.GetAllGroupsByUser(cookie)
	if err != nil {
		return err
	}

	groupsFormatted := make([]domain.Group, len(groups))
	for i, group := range groups {
		groupIdStr := group.GroupId
		groupIdInt, err := strconv.Atoi(groupIdStr)
		if err != nil {
			log.Printf("Error converting group id %s to int\n", groupIdStr)
			continue
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
	d.domain.SetGroups(uid, groupsFormatted)

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
