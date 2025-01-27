package defaultState

import (
	"fmt"
	"gopkg.in/telebot.v4"
	"strings"
	"tgbot/internal/config"
	"tgbot/internal/domain"
	"tgbot/internal/helpers"
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
		return ctx.Send(config.UserDontHaveGroup, config.MyGroupsKeyboard)
	}
	sorted := helpers.GetSortedGroups(g)

	return ctx.Send(GetMyGroupsMessage(sorted), config.MyGroupsKeyboard)
}

func GetMyGroupsMessage(g []domain.Group) string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("%s%d\n", config.MyGroups, len(g)))
	before := g[0].Time.Weekday()
	c := 1
	for _, group := range g {
		if before != group.Time.Weekday() {
			c = 1
			before = group.Time.Weekday()
			s.WriteString("\n")
		}
		s.WriteString("\n")
		s.WriteString(fmt.Sprintf("%d. %s 🕐 %s %s", c, group.Name, getLocale(group.Time), group.Time.Format("15:04")))
		c += 1
	}

	return s.String()
}

func getLocale(t time.Time) string {
	return locales[t.Weekday()]
}
