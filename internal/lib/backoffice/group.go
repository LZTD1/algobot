package backoffice

import (
	"algobot/internal/domain/models"
	"algobot/internal/lib/logger/sl"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

func (bo *Backoffice) Group(cookie string) ([]models.Group, error) {
	const op = "backoffice.Group"
	log := bo.log.With(
		slog.String("op", op),
	)

	req, err := bo.createReq("GET", "/group", cookie, map[string]string{
		"GroupSearch[status][]": "active",
		"presetType":            "all",
		"_pjax":                 "#group-grid-pjax",
	}, nil)

	if err != nil {
		log.Warn("failed to create request", sl.Err(err))
		return nil, fmt.Errorf("%s failed to create request: %w", op, err)
	}
	data, err := bo.doReq(req)
	if err != nil {
		return nil, fmt.Errorf("%s failed to doReq: %w", op, err)
	}

	res, err := parseHTML(data.Body, bo.log)
	if err != nil {
		return nil, fmt.Errorf("%s failed to parse HTML: %w", op, err)
	}

	return res, nil
}

func parseHTML(body io.ReadCloser, log *slog.Logger) ([]models.Group, error) {
	const op = "backoffice.parseHTML"
	log = log.With(slog.String("op", op))

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("%s failed parse doc: %w", op, err)
	}

	var groups []models.Group

	doc.Find("tr.group-grid").Each(func(i int, row *goquery.Selection) {
		groupId := row.Find("td[data-col-seq='id']").First().Text()
		groupId = strings.TrimSpace(groupId)

		titleCell := row.Find("td[data-col-seq='title']").First()

		groupTitle := titleCell.Find("p").First().Text()
		groupTitle = strings.TrimSpace(groupTitle)

		groupTime := titleCell.Find("a").First().Text()
		groupTime = strings.TrimSpace(groupTime)

		nextLessonTime := row.Find("td[data-col-seq='nextLessonTime']").First().Text()
		nextLessonTime = strings.TrimSpace(nextLessonTime)

		if nextLessonTime != "" {
			groupID, err := strconv.Atoi(strings.ReplaceAll(groupId, "\u00A0", " "))
			if err != nil {
				log.Warn("failed to convert group id to int", sl.Err(err))
				return
			}
			timeLession, err := time.Parse("02.01.2006 15:04", strings.ReplaceAll(nextLessonTime, "\u00A0", " "))
			if err != nil {
				log.Warn("failed to convert timeLession to time", sl.Err(err))
				return
			}

			groups = append(groups, models.Group{
				GroupID:    groupID,
				Title:      strings.ReplaceAll(groupTitle, "\u00A0", " "),
				TimeLesson: timeLession,
			})
		}
	})

	return groups, nil
}
