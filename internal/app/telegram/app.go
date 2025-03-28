package telegram

import (
	"algobot/internal/lib/logger/sl"
	"algobot/internal/telegram/middleware/logger"
	"algobot/internal/telegram/middleware/trace"
	tele "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
	"log/slog"
	"os"
	"time"
)

type App struct {
	log *slog.Logger
	bot *tele.Bot
}

func New(log *slog.Logger, token string) *App {
	const op = "telegram.New"

	nlog := log.With(
		slog.String("op", op),
	)

	pref := tele.Settings{
		Token: token,
		Poller: &tele.LongPoller{
			Timeout: 10 * time.Second,
		},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		nlog.Warn("error by creating telegram bot: ", sl.Err(err))
		os.Exit(1)
	}

	// initialize routes
	b.Use(trace.New(log))
	b.Use(middleware.AutoRespond())
	b.Use(middleware.Recover())
	b.Use(logger.New(log))

	b.
		b.Handle(tele.OnText, func(c tele.Context) error {
		return c.Reply(c.Get("trace_id").(string))
	})
	return &App{log: log, bot: b}
}

func (a *App) Run() {
	a.bot.Start()
}

func (a *App) Stop() {
	a.bot.Stop()
}
