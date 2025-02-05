package defaultState

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gopkg.in/telebot.v4"
	"os"
	"regexp"
	"strconv"
	"strings"
	"tgbot/internal/helpers"
	"tgbot/internal/models"
	"tgbot/internal/service"
)

var statuses = map[int]string{
	0:  "🟢 Учиться",
	20: "🔴 Выбыл",
	10: "🟡 Переведен",
}

type StartWithPayload struct {
	svc service.Service
}

func NewStartWithPayload(svc service.Service) *StartWithPayload {
	return &StartWithPayload{svc: svc}
}

func (s StartWithPayload) CanHandle(ctx telebot.Context) bool {
	if strings.HasPrefix(ctx.Message().Text, "/start") {
		return true
	}

	return false
}

func (s StartWithPayload) Process(ctx telebot.Context) error {
	decodedBytes, err := base64.StdEncoding.DecodeString(ctx.Message().Payload)
	if err != nil {
		return ctx.Send("Ошибка декодирования")
	}
	var payload models.StartPayload

	err = json.Unmarshal(decodedBytes, &payload)
	if err != nil {
		return ctx.Send("(1) Ошибка декодирования")
	}

	switch payload.Action {
	case models.GetGroupInfo:
		return s.getGroupInfo(ctx, payload)
	case models.GetKidInfo:
		return s.getKidInfo(ctx, payload)
	default:
		return ctx.Send("Not supported")
	}
}

func (s StartWithPayload) getGroupInfo(ctx telebot.Context, payload models.StartPayload) error {
	g, _ := strconv.Atoi(payload.Payload[0])
	full, err := s.svc.FullGroupInfo(ctx.Sender().ID, g)
	if err != nil {
		return helpers.LogError(err, ctx, "Ошибка при получении данной группы!")
	}

	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("[%s %s](https://backoffice.algoritmika.org/group/view/%d)\n", full.GroupTitle, full.GroupContent, full.GroupID))
	msg.WriteString(fmt.Sprintf("\n***Следующая лекция***: %s\n", full.NextLessonTime))
	msg.WriteString(fmt.Sprintf("***Всего пройдено*** %d лекций из %d\n", full.LessonsPassed, full.LessonsTotal))
	msg.WriteString(fmt.Sprintf("\nАктивные дети: %d | Выбыло: %d | Всего: %d\n", len(full.ActiveKids), len(full.NotActiveKids), len(full.ActiveKids)+len(full.NotActiveKids)))
	msg.WriteString("***Активные дети***:\n")
	for i, kid := range full.ActiveKids {
		marshal, _ := json.Marshal(models.StartPayload{
			Action:  models.GetKidInfo,
			Payload: []string{strconv.Itoa(kid.ID)},
		})
		encodedStr := base64.StdEncoding.EncodeToString(marshal)

		msg.WriteString(fmt.Sprintf("%d. [%s](t.me/%s?start=%s)\n", i+1, kid.FullName, os.Getenv("TELEGRAM_NAME"), encodedStr))
	}
	msg.WriteString("***Выбыли дети***:\n")
	for i, kid := range full.NotActiveKids {
		marshal, _ := json.Marshal(models.StartPayload{
			Action:  models.GetKidInfo,
			Payload: []string{strconv.Itoa(kid.ID)},
		})
		encodedStr := base64.StdEncoding.EncodeToString(marshal)

		if kid.LastGroup.ID == g {
			msg.WriteString(fmt.Sprintf("%d. [%s](t.me/%s?start=%s) (Выбыл: %s)\n", i+1, kid.FullName, os.Getenv("TELEGRAM_NAME"), encodedStr, kid.LastGroup.EndTime.Format("2006-01-02")))
		} else {
			msg.WriteString(fmt.Sprintf("%d. [%s](t.me/%s?start=%s) (Переведен: %s)\n", i+1, kid.FullName, os.Getenv("TELEGRAM_NAME"), encodedStr, kid.LastGroup.StartTime.Format("2006-01-02")))
		}
	}
	return ctx.Send(msg.String(), telebot.ModeMarkdown)
}

func (s StartWithPayload) getKidInfo(ctx telebot.Context, payload models.StartPayload) error {
	id, _ := strconv.Atoi(payload.Payload[0])
	full, err := s.svc.FullKidInfo(ctx.Sender().ID, id)
	if err != nil {
		return helpers.LogError(err, ctx, "Ошибка при получении данного ученика!")
	}

	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("***%s***\n", full.Kid.Data.FullName))
	msg.WriteString(fmt.Sprintf("Возраст: %d\n", full.Kid.Data.Age))
	msg.WriteString(fmt.Sprintf("День рождения: %s\n", full.Kid.Data.BirthDate.Format("2006-01-02")))
	msg.WriteString("\n***Данные от аккаунта:***\n")
	msg.WriteString(fmt.Sprintf("Логин: _%s_\n", full.Kid.Data.Username))
	msg.WriteString(fmt.Sprintf("Пароль: _%s_\n", full.Kid.Data.Password))
	msg.WriteString("\n***Родитель:***\n")
	msg.WriteString(fmt.Sprintf("Имя: %s\n", full.Kid.Data.ParentName))

	msg.WriteString(fmt.Sprintf("Телефон: %s\n", regexp.MustCompile(`[^0-9+]`).ReplaceAllString(full.Kid.Data.Phone, "")))
	msg.WriteString(fmt.Sprintf("Почта: %s\n", full.Kid.Data.Email))
	msg.WriteString("\n***Группы***\n")

	groups := full.Kid.Data.Groups
	for i := len(groups) - 1; i >= 0; i-- {
		msg.WriteString(fmt.Sprintf("%d . [%s %s](https://backoffice.algoritmika.org/group/view/%d)\n", len(groups)-i, groups[i].Title, groups[i].Content, groups[i].ID))
		v, ok := statuses[groups[i].Status]
		if !ok {
			v = fmt.Sprintf("Статус [%d]", groups[i].Status)
		}
		msg.WriteString(fmt.Sprintf("%s (%s - %s)\n\n", v, groups[i].StartTime.Format("2006-01-02"), groups[i].EndTime.Format("2006-01-02")))
	}

	return ctx.Send(msg.String(), telebot.ModeMarkdown)
}
