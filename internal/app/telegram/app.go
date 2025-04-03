package telegram

import (
	"algobot/internal/lib/fsm"
	"algobot/internal/lib/fsm/memory"
	"algobot/internal/lib/logger/sl"
	"algobot/internal/telegram/handlers/callback"
	"algobot/internal/telegram/handlers/text"
	"algobot/internal/telegram/middleware/auth"
	"algobot/internal/telegram/middleware/logger"
	"algobot/internal/telegram/middleware/stater"
	"algobot/internal/telegram/middleware/trace"
	"fmt"
	router "github.com/LZTD1/telebot-context-router"
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

func New(log *slog.Logger, token string, auther auth.Auther, set text.UserInformer) *App {
	const op = "telegram.New"

	nlog := log.With(
		slog.String("op", op),
	)

	pref := tele.Settings{
		Token: token,
		Poller: &tele.LongPoller{
			Timeout: 10 * time.Second,
		},
		OnError: func(e error, c tele.Context) { // TODO : refactor into handler
			traceID := c.Get("trace_id") // TODO : maybe send warnings to admin ?
			c.Send(fmt.Sprintf("<b>[%s]</b>\n\nУпс! Произошла какая-то непредвиденная ошибка!\nОбратитесь к администратору", traceID), tele.ModeHTML)
			log.Warn("cant handle error", sl.Err(e), slog.Any("trace_id", traceID))
		},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		nlog.Warn("error by creating telegram bot: ", sl.Err(err))
		os.Exit(1)
	}

	stateMachine := memory.New()

	// initialize routes
	b.Use(trace.New(log))
	b.Use(middleware.AutoRespond())
	b.Use(middleware.Recover())
	b.Use(logger.New(log))
	b.Use(auth.New(auther, log))

	r := router.NewRouter()
	//st := stater.New(stateMachine, fsm.Default)
	r.Group(func(r router.Router) { // Routes for default state
		r.Use(stater.New(stateMachine, fsm.Default))

		// message
		r.HandleFuncText("/start", text.NewStart(stateMachine))
		r.HandleFuncText("Настройки", text.NewSettings(set, log))

		// callbacks
		r.HandleFuncCallback("\fset_cookie", callback.NewChangeCookie(stateMachine))
		r.HandleFuncCallback("\fchange_notification", nil)
	})

	r.Group(func(r router.Router) { // Routes for SendingCookie state
		r.Use(stater.New(stateMachine, fsm.SendingCookie))

		// message
		r.HandleFuncText("⬅️ Назад", text.NewStart(stateMachine))
	})

	r.NotFound(text.NewStart(stateMachine))

	b.Handle(tele.OnText, r.ServeContext)
	b.Handle(tele.OnCallback, r.ServeContext)

	return &App{log: log, bot: b}
}

func (a *App) Run() {
	a.bot.Start()
}

func (a *App) Stop() {
	a.bot.Stop()
}
