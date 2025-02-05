package defaultState

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/telebot.v4"
	"os"
	"strconv"
	"strings"
	"tgbot/internal/config"
	appError "tgbot/internal/error"
	"tgbot/internal/helpers"
	"tgbot/internal/models"
	"tgbot/internal/service"
	"time"
)

var locales = map[time.Weekday]string{
	time.Monday:    "пн",
	time.Tuesday:   "вт",
	time.Wednesday: "ср",
	time.Thursday:  "чт",
	time.Friday:    "пт",
	time.Saturday:  "сб",
	time.Sunday:    "вс",
}

type MyGroups struct {
	s service.Service
}

func NewMyGroups(s service.Service) *MyGroups {
	return &MyGroups{s}
}

func (m MyGroups) CanHandle(ctx telebot.Context) bool {
	if ctx.Message().Text == "Мои группы" {
		return true
	}
	return false
}
func (m MyGroups) Process(ctx telebot.Context) error {
	g, e := m.s.Groups(ctx.Message().Sender.ID)

	if e != nil {
		if errors.Is(e, appError.ErrHasNone) {
			return ctx.Send(config.UserDontHaveGroup, config.MyGroupsKeyboard)
		}
		return helpers.LogError(e, ctx, "Ошибка при попытке получить группы!")
	}
	sorted := helpers.GetSortedGroups(g)

	return ctx.Send(GetMyGroupsMessage(sorted), config.MyGroupsKeyboard, telebot.ModeMarkdown)
}

func GetMyGroupsMessage(g []models.Group) string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("%s%d\n", config.MyGroups, len(g)))

	before := g[0].TimeLesson.Weekday()
	c := 1
	for _, group := range g {
		if before != group.TimeLesson.Weekday() {
			c = 1
			before = group.TimeLesson.Weekday()
			s.WriteString("\n")
		}
		s.WriteString("\n")
		s.WriteString(fmt.Sprintf("%d. %s 🕐 %s %s", c, getGroupTitle(group), getLocale(group.TimeLesson), group.TimeLesson.Format("15:04")))
		c += 1
	}

	return s.String()
}

func getGroupTitle(group models.Group) string {
	marshal, err := json.Marshal(models.StartPayload{
		Action:  models.GetGroupInfo,
		Payload: []string{strconv.Itoa(group.GroupID)},
	})
	encodedStr := base64.StdEncoding.EncodeToString(marshal)

	if err != nil {
		fmt.Println(err)
		return group.Title
	}
	return fmt.Sprintf("[%s](t.me/%s?start=%s)", group.Title, os.Getenv("TELEGRAM_NAME"), encodedStr)
}

func getLocale(t time.Time) string {
	return locales[t.Weekday()]
}
