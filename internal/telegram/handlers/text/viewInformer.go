package text

import (
	"algobot/internal/domain"
	"algobot/internal/domain/models"
	"algobot/internal/lib/logger/sl"
	"fmt"
	"gopkg.in/telebot.v4"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
)

var statuses = map[int]string{
	0:  "🟢 Учится",
	20: "🔴 Выбыл",
	10: "🟡 Переведен",
}

type ViewFetcher interface {
	GroupView(uid int64, groupID string, traceID interface{}) (models.GroupView, error)
	KidView(uid int64, kidID string, groupId string, traceID interface{}) (models.KidView, error)
}

type Serializator interface {
	Serialize(msg domain.SerializeMessage) (string, error)
	Deserialize(decoded string) (*domain.SerializeMessage, error)
}

type ViewInformer struct {
	serdes      Serializator
	viewFetcher ViewFetcher
	log         *slog.Logger
	botName     string
}

func NewViewInformer(serdes Serializator, viewFetcher ViewFetcher, log *slog.Logger, botName string) *ViewInformer {
	return &ViewInformer{serdes: serdes, viewFetcher: viewFetcher, log: log, botName: botName}
}

func (v *ViewInformer) ServeContext(ctx telebot.Context) error {
	const op = "text.GenerateImage"

	uid := ctx.Sender().ID
	traceID := ctx.Get("trace_id")
	data := getData(ctx.Message().Text)

	log := v.log.With(
		slog.String("op", op),
		slog.Any("trace_id", traceID),
	)

	encodedMsg, err := v.serdes.Deserialize(data)
	if err != nil {
		log.Warn("can't get ser type", sl.Err(err))
		return ctx.Send("⚠️ Ошибка при расшифровке запроса!")
	}

	switch encodedMsg.Type {
	case domain.UserType:
		view, err := v.userInfo(encodedMsg, uid, traceID)
		if err != nil {
			log.Warn("can't get group info", sl.Err(err))
			return ctx.Send("⚠️ Невозможно получить данного ученика!")
		}
		return ctx.Send(view, telebot.ModeHTML, telebot.NoPreview)
	case domain.GroupType:
		view, err := v.groupInfo(encodedMsg, uid, traceID)
		if err != nil {
			log.Warn("can't get group info", sl.Err(err))
			return ctx.Send("⚠️ Невозможно получить данную группу!")
		}
		return ctx.Send(view, telebot.ModeHTML, telebot.NoPreview)
	default:
		return ctx.Send("⚠️ Не удалось определить обработчик!")
	}
}

func (v *ViewInformer) userInfo(data *domain.SerializeMessage, uid int64, traceID interface{}) (string, error) {
	const op = "viewInformer.groupInfo"

	if len(data.Data) != 2 {
		return "", fmt.Errorf("%s: kid view required 2 fields", op)
	}

	full, err := v.viewFetcher.KidView(uid, data.Data[0], data.Data[1], traceID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return v.GetKidInfoMessage(full), nil
}

func (v *ViewInformer) GetKidInfoMessage(full models.KidView) string {
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

func (v *ViewInformer) groupInfo(data *domain.SerializeMessage, uid int64, traceID interface{}) (string, error) {
	const op = "viewInformer.groupInfo"

	full, err := v.viewFetcher.GroupView(uid, data.Data[0], traceID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return v.GetGroupInfoMessage(full), nil
}

func (v *ViewInformer) GetGroupInfoMessage(full models.GroupView) string {
	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("<a href=\"https://backoffice.algoritmika.org/group/view/%d\">%s %s</a>\n", full.GroupID, full.GroupTitle, full.GroupContent))
	msg.WriteString(fmt.Sprintf("\n<b>Следующая лекция</b>: %s\n", full.NextLessonTime))
	msg.WriteString(fmt.Sprintf("<b>Всего пройдено</b> %d лекций из %d\n", full.LessonsPassed, full.LessonsTotal))
	msg.WriteString(fmt.Sprintf("\nАктивные дети: %d | Выбыло: %d | Всего: %d\n", len(full.ActiveKids), len(full.NotActiveKids), len(full.ActiveKids)+len(full.NotActiveKids)))
	msg.WriteString("<b>Активные дети</b>:\n")

	for i, kid := range full.ActiveKids {
		ser, err := v.serdes.Serialize(domain.SerializeMessage{
			Type: domain.UserType,
			Data: []string{strconv.Itoa(kid.ID), strconv.Itoa(full.GroupID)},
		})
		if err != nil {
			msg.WriteString(fmt.Sprintf("%d. %s\n", i+1, kid.FullName))
			continue
		}

		msg.WriteString(fmt.Sprintf("%d. <a href=\"https://t.me/%s?start=%s\">%s</a>\n", i+1, v.botName, ser, kid.FullName))
	}

	msg.WriteString("<b>Выбыли дети</b>:\n")
	for i, kid := range full.NotActiveKids {
		ser, err := v.serdes.Serialize(domain.SerializeMessage{
			Type: domain.UserType,
			Data: []string{strconv.Itoa(kid.ID), strconv.Itoa(full.GroupID)},
		})
		if err != nil {
			if kid.LastGroup.ID == full.GroupID {
				msg.WriteString(fmt.Sprintf("%d. %s (🔴 Выбыл: %s)\n", i+1, kid.FullName, kid.LastGroup.EndTime.Format("2006-01-02")))
			} else {
				msg.WriteString(fmt.Sprintf("%d. %s (🟡 Переведен: %s)\n", i+1, kid.FullName, kid.LastGroup.StartTime.Format("2006-01-02")))
			}
		}

		if kid.LastGroup.ID == full.GroupID {
			msg.WriteString(fmt.Sprintf("%d. <a href=\"https://t.me/%s?start=%s\">%s</a> (🔴 Выбыл: %s)\n", i+1, v.botName, ser, kid.FullName, kid.LastGroup.EndTime.Format("2006-01-02")))
		} else {
			msg.WriteString(fmt.Sprintf("%d. <a href=\"https://t.me/%s?start=%s\">%s</a> (🟡 Переведен: %s)\n", i+1, v.botName, ser, kid.FullName, kid.LastGroup.StartTime.Format("2006-01-02")))
		}
	}

	return msg.String()
}

func getData(text string) string {
	return strings.TrimSpace(strings.TrimLeft(text, "/start"))
}
