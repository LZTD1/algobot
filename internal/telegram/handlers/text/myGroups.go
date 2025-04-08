package text

import (
	"algobot/internal/domain/models"
	"algobot/internal/domain/telegram/keyboards"
	"algobot/internal/lib/logger/sl"
	"fmt"
	"gopkg.in/telebot.v4"
	"log/slog"
	"strings"
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

type Grouper interface {
	Groups(uid int64, traceID interface{}) ([]models.Group, error)
}

type GroupSerializer interface {
	Serialize(group models.Group, traceID interface{}) (string, error)
}
type MyGroup struct {
	log        *slog.Logger
	grouper    Grouper
	serializer GroupSerializer
	botName    string
}

func NewMyGroup(log *slog.Logger, grouper Grouper, serializer GroupSerializer, name string) *MyGroup {
	return &MyGroup{
		log:        log,
		grouper:    grouper,
		serializer: serializer,
		botName:    name,
	}
}

func (g *MyGroup) ServeContext(ctx telebot.Context) error {
	const op = "text.MyGroup.ServeContext"
	log := g.log.With(
		slog.String("op", op),
		slog.Any("trace_id", ctx.Get("trace_id")),
	)

	uid := ctx.Sender().ID
	groups, err := g.grouper.Groups(uid, ctx.Get("trace_id"))
	if err != nil {
		log.Warn("error while getting groups", sl.Err(err))
		return ctx.Send(fmt.Sprintf("<b>[%s]</b> Ошибка при получении групп!", ctx.Get("trace_id")), telebot.ModeHTML)
	}

	return ctx.Send(g.msgMyGroups(groups, ctx), telebot.ModeMarkdown, keyboards.RefreshGroups())
}

func (g *MyGroup) msgMyGroups(groups []models.Group, ctx telebot.Context) string {
	s := &strings.Builder{}
	s.WriteString(fmt.Sprintf("Всего групп: %d\n", len(groups)))

	if len(groups) == 0 {
		s.WriteString("Попробуйте обновить группы!")
		return s.String()
	}

	beforeDay := groups[0].TimeLesson.Weekday()
	c := 1
	for _, group := range groups {
		if beforeDay != group.TimeLesson.Weekday() {
			c = 1
			beforeDay = group.TimeLesson.Weekday()
			s.WriteString("\n")
		}
		s.WriteString("\n")
		s.WriteString(fmt.Sprintf(
			"%d. %s 🕐 %s %s",
			c,
			g.getFormattedTitle(group, ctx),
			g.getLocale(group.TimeLesson),
			group.TimeLesson.Format("15:04"),
		))
		c += 1
	}

	return s.String()
}

func (g *MyGroup) getFormattedTitle(group models.Group, ctx telebot.Context) string {
	const op = "text.MyGroup.getFormattedTitle"
	log := g.log.With(
		slog.String("op", op),
		slog.Any("trace_id", ctx.Get("trace_id")),
	)

	serialized, err := g.serializer.Serialize(group, ctx.Get("trace_id"))
	if err != nil {
		log.Warn("error while serializing group", sl.Err(err))
		return group.Title
	}
	return fmt.Sprintf("[%s](t.me/%s?start=%s)", group.Title, g.botName, serialized)
}

func (g *MyGroup) getLocale(t time.Time) string {
	return locales[t.Weekday()]
}
