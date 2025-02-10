package defaultState

import (
	"fmt"
	"gopkg.in/telebot.v4"
	"os"
	"regexp"
	"strconv"
	"strings"
	"tgbot/internal/helpers"
	"tgbot/internal/models"
	"tgbot/internal/serdes"
	"tgbot/internal/service"
)

var statuses = map[int]string{
	0:  "🟢 Учится",
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
	payload, err := serdes.Deserialize(ctx.Message().Payload)
	if err != nil {
		return helpers.LogError(err, ctx, "Ошибка десериализации")
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

	msg := GetGroupInfoMessage(full)
	return ctx.Send(msg, telebot.ModeHTML, telebot.NoPreview)
}

func GetGroupInfoMessage(full models.FullGroupInfo) string {
	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("<a href=\"https://backoffice.algoritmika.org/group/view/%d\">%s %s</a>\n", full.GroupID, full.GroupTitle, full.GroupContent))
	msg.WriteString(fmt.Sprintf("\n<b>Следующая лекция</b>: %s\n", full.NextLessonTime))
	msg.WriteString(fmt.Sprintf("<b>Всего пройдено</b> %d лекций из %d\n", full.LessonsPassed, full.LessonsTotal))
	msg.WriteString(fmt.Sprintf("\nАктивные дети: %d | Выбыло: %d | Всего: %d\n", len(full.ActiveKids), len(full.NotActiveKids), len(full.ActiveKids)+len(full.NotActiveKids)))
	msg.WriteString("<b>Активные дети</b>:\n")
	for i, kid := range full.ActiveKids {
		ser := serdes.Serialize(models.StartPayload{
			Action:  models.GetKidInfo,
			Payload: []string{strconv.Itoa(kid.ID), strconv.Itoa(full.GroupID)},
		})

		msg.WriteString(fmt.Sprintf("%d. <a href=\"https://t.me/%s?start=%s\">%s</a>\n", i+1, os.Getenv("TELEGRAM_NAME"), ser, kid.FullName))
	}
	msg.WriteString("<b>Выбыли дети</b>:\n")
	for i, kid := range full.NotActiveKids {
		ser := serdes.Serialize(models.StartPayload{
			Action:  models.GetKidInfo,
			Payload: []string{strconv.Itoa(kid.ID), strconv.Itoa(full.GroupID)},
		})

		if kid.LastGroup.ID == full.GroupID {
			msg.WriteString(fmt.Sprintf("%d. <a href=\"https://t.me/%s?start=%s\">%s</a> (🔴 Выбыл: %s)\n", i+1, os.Getenv("TELEGRAM_NAME"), ser, kid.FullName, kid.LastGroup.EndTime.Format("2006-01-02")))
		} else {
			msg.WriteString(fmt.Sprintf("%d. <a href=\"https://t.me/%s?start=%s\">%s</a> (🟡 Переведен: %s)\n", i+1, os.Getenv("TELEGRAM_NAME"), ser, kid.FullName, kid.LastGroup.StartTime.Format("2006-01-02")))
		}
	}
	return msg.String()
}

func (s StartWithPayload) getKidInfo(ctx telebot.Context, payload models.StartPayload) error {

	id, _ := strconv.Atoi(payload.Payload[0])
	groupId, _ := strconv.Atoi(payload.Payload[1])
	full, err := s.svc.FullKidInfo(ctx.Sender().ID, id, groupId)
	if err != nil {
		return helpers.LogError(err, ctx, "Ошибка при получении данного ученика!")
	}

	m := GetKidInfoMessage(full)
	return ctx.Send(m, telebot.ModeHTML, telebot.NoPreview)
}

func GetKidInfoMessage(full models.FullKidInfo) string {
	parentPhone := regexp.MustCompile(`[^0-9+]`).ReplaceAllString(full.Kid.Phone, "")

	msg := strings.Builder{}
	if full.Extra == models.NotAccessible {
		msg.WriteString(fmt.Sprintf("⚠️ У вас больше нету доступа к ребенку\n"))
	}
	msg.WriteString(fmt.Sprintf("<b>%s</b>\n", full.Kid.FullName))
	msg.WriteString(fmt.Sprintf("Возраст: %d\n", full.Kid.Age))
	msg.WriteString(fmt.Sprintf("День рождения: %s\n", full.Kid.BirthDate.Format("2006-01-02")))
	msg.WriteString("\n<b>Данные от аккаунта:</b>\n")
	msg.WriteString(fmt.Sprintf("Логин: <i>%s</i>\n", full.Kid.Username))
	msg.WriteString(fmt.Sprintf("Пароль: <i>%s</i>\n", full.Kid.Password))
	msg.WriteString("\n<b>Родитель:</b>\n")
	msg.WriteString(fmt.Sprintf("Имя: %s\n", full.Kid.ParentName))

	msg.WriteString(fmt.Sprintf("Телефон: %s <a href=\"https://wa.me/%s\">🟩 Whatsapp</a>\n", parentPhone, strings.TrimPrefix(parentPhone, "+")))
	msg.WriteString(fmt.Sprintf("Почта: %s\n", full.Kid.Email))
	msg.WriteString("\n<b>Группы</b>\n")

	groups := full.Kid.Groups
	for i := len(groups) - 1; i >= 0; i-- {
		msg.WriteString(fmt.Sprintf("%d . <a href=\"https://backoffice.algoritmika.org/group/view/%d\">%s %s</a>\n", len(groups)-i, groups[i].ID, groups[i].Title, groups[i].Content))
		v, ok := statuses[groups[i].Status]
		if !ok {
			v = fmt.Sprintf("Статус [%d]", groups[i].Status)
		}
		msg.WriteString(fmt.Sprintf("%s (%s - %s)\n\n", v, groups[i].StartTime.Format("2006-01-02"), groups[i].EndTime.Format("2006-01-02")))
	}
	m := msg.String()
	return m
}
