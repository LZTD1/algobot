package backoffice

import (
	"algobot/internal/domain/backoffice"
	"algobot/internal/domain/scheduler"
	"algobot/internal/lib/logger/sl"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

var dateMap = map[string]string{
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
var dateReverseMap = map[int]string{
	1:  "янв",
	2:  "февр",
	3:  "мар",
	4:  "апр",
	5:  "мая",
	6:  "июн",
	7:  "июл",
	8:  "авг",
	9:  "сент",
	10: "окт",
	11: "нояб",
	12: "дек",
}

type MessageFetcher interface {
	KidsMessages(cookie string) (backoffice.KidsMessages, error)
}

func (bo *Backoffice) MessagesUser(uid int64, lastTime string) ([]scheduler.Message, error) {
	const op = "services.backoffice.MessagesUser"
	log := bo.log.With(
		slog.String("op", op),
	)

	cookie, err := bo.cookieGetter.Cookies(uid)
	if err != nil {
		log.Warn("failed to get cookies", sl.Err(err))
		return nil, fmt.Errorf("%s failed to get cookies: %w", op, err)
	}

	messages, err := bo.msgFetcher.KidsMessages(cookie)
	if err != nil {
		return nil, fmt.Errorf("%s failed to get KidsMessages: %w", op, err)
	}

	var dateNotif time.Time
	if lastTime != "" {
		dateNotif = parseDate(lastTime)
	} else {
		timeNow := time.Now()
		return []scheduler.Message{{
			From:  "",
			Theme: "",
			Link:  "",
			Text:  "",
			Time:  timeNow.Format(fmt.Sprintf("2 %s. 15:04", dateReverseMap[int(timeNow.Month())])),
			To:    uid,
		}}, nil
	}

	var msgs []scheduler.Message
	for i := len(messages.Data.Projects) - 1; i >= 0; i-- {
		if messages.Data.Projects[i].SenderScope == "student" {
			if dateNotif.Before(parseDate(messages.Data.Projects[i].LastTime)) {
				m := scheduler.Message{
					From:  messages.Data.Projects[i].Name,
					Theme: messages.Data.Projects[i].Title,
					Link:  fmt.Sprintf("https://backoffice.algoritmika.org%s", messages.Data.Projects[i].Link),
					Text:  messages.Data.Projects[i].Content,
					Time:  messages.Data.Projects[i].LastTime,
					To:    uid,
				}

				if messages.Data.Projects[i].Type == "img" {
					m.LinkURL = fmt.Sprintf("https://backoffice.algoritmika.org%s", m.Text)
				}

				msgs = append(msgs, m)
			}
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
