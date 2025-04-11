package groups

import (
	"algobot/internal/domain/backoffice"
	"algobot/internal/domain/models"
	"algobot/internal/lib/logger/sl"
	"fmt"
	"log"
	"log/slog"
	"regexp"
	"strconv"
	"time"
)

type KidStats interface {
	KidsStats(cookie string, groupID int) (backoffice.KidsStats, error)
	KidsNamesByGroup(groupId string, cookie string) (backoffice.NamesByGroup, error)
}

func (g *Group) CurrentGroup(uid int64, time time.Time, traceID interface{}) (models.CurrentGroup, error) {
	const op = "services.groups.CurrentGroup"
	log := g.log.With(
		slog.String("op", op),
		slog.Any("trace_id", traceID),
	)

	cookie, err := g.domain.Cookies(uid)
	if err != nil {
		log.Warn("failed to get cookies", sl.Err(err))
		return models.CurrentGroup{}, fmt.Errorf("%s failed to get cookies: %w", op, err)
	}
	if cookie == "" {
		return models.CurrentGroup{}, fmt.Errorf("%s cookie is empty: %w", op, ErrNotValidCookie)
	}

	group, err := g.Groups(uid, traceID)
	if err != nil {
		log.Warn("failed to get groups", sl.Err(err))
		return models.CurrentGroup{}, fmt.Errorf("%s failed to get groups: %w", op, err)
	}

	m := models.CurrentGroup{}
	missingKids := make(map[int]models.MissingKid)

	actual, err := CurrentGroup(time, group)
	if err != nil {
		return models.CurrentGroup{}, fmt.Errorf("%s : %w", op, err)
	}
	m.Title = actual.Title
	m.GroupID = actual.GroupID

	stats, err := g.kidsStats.KidsStats(cookie, actual.GroupID)
	if err != nil {
		log.Warn("error while fetching kids stats", sl.Err(err))
		return models.CurrentGroup{}, fmt.Errorf("%s error while fetching kids stats: %w", op, err)
	}
	for _, datum := range stats.Data {
		studentID := datum.StudentID
		count := 0
		for _, attendance := range datum.Attendance {
			count++
			if attendance.Status != "absent" {
				count = 0
			}
			if matchDates(attendance.StartTimeFormatted, time) {
				m.Lesson = attendance.LessonTitle
				m.LessonID = attendance.LessonID

				if attendance.Status == "absent" {
					missingKids[studentID] = models.MissingKid{
						Fullname: "",
						Count:    count,
						KidID:    studentID,
					}
					break
				}
			}
		}
	}

	names, err := g.kidsStats.KidsNamesByGroup(strconv.Itoa(actual.GroupID), cookie)
	if err != nil {
		log.Warn("error while fetching KidsNamesByGroup", sl.Err(err))
		return models.CurrentGroup{}, fmt.Errorf("%s error while fetching KidsNamesByGroup: %w", op, err)
	}
	for _, datum := range names.Data.Items {
		if datum.LastGroup.ID == actual.GroupID && datum.LastGroup.Status == 0 {
			missingKids[datum.ID] = models.MissingKid{
				Fullname: datum.FullName,
				Count:    missingKids[datum.ID].Count,
				KidID:    missingKids[datum.ID].KidID,
			}
		}
	}

	split(&m, missingKids)

	return m, nil
}

func split(m *models.CurrentGroup, kids map[int]models.MissingKid) {
	m.MissingKids = make([]models.MissingKid, 0, len(kids))
	for _, kid := range kids {
		if kid.Count != 0 {
			m.MissingKids = append(m.MissingKids, kid)
		}
		m.Kids = append(m.Kids, kid.Fullname)
	}
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
