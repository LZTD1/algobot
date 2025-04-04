package telegram

import (
	"algobot/internal/config"
	"algobot/internal/lib/fsm"
	"algobot/internal/lib/fsm/memory"
	"algobot/internal/lib/logger/sl"
	"algobot/internal/lib/serdes/base62"
	"algobot/internal/services"
	grpc2 "algobot/internal/services/grpc"
	"algobot/internal/telegram/handlers/callback"
	"algobot/internal/telegram/handlers/text"
	"algobot/internal/telegram/middleware/auth"
	"algobot/internal/telegram/middleware/logger"
	"algobot/internal/telegram/middleware/rate"
	"algobot/internal/telegram/middleware/stater"
	"algobot/internal/telegram/middleware/trace"
	"fmt"
	router "github.com/LZTD1/telebot-context-router"
	tele "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
	"log/slog"
	"os"
	"regexp"
	"time"
)

type App struct {
	log *slog.Logger
	bot *tele.Bot
}

func New(
	log *slog.Logger,
	token string,
	grGetter services.GroupGetter,
	auther auth.Auther,
	set text.UserInformer,
	cookieSetter text.CookieSetter,
	notifChanger callback.NotificationChanger,
	rateCfg config.RateLimit,
	grpcCfg config.GRPC,

) *App {
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

	// dependencies
	groupServ := services.NewGroup(log, grGetter)
	stateMachine := memory.New()
	serdes := base62.NewSerdes(log)
	grpc := grpc2.NewAIService(grpcCfg, log)

	// initialize routes
	b.Use(trace.New(log))
	b.Use(middleware.AutoRespond())
	b.Use(middleware.Recover())
	b.Use(logger.New(log))
	b.Use(auth.New(auther, log))
	b.Use(rate.New(log, rateCfg))

	// create routing
	r := router.NewRouter()
	r.Group(func(r router.Router) { // Routes for default state
		r.Use(stater.New(stateMachine, fsm.Default))

		// message
		r.HandleFuncText("/start", text.NewStart(stateMachine))
		r.HandleFuncText("Настройки", text.NewSettings(set, log))
		r.HandleFuncText("AI 🔹", text.NewAI(grpc, log, stateMachine))
		r.HandleText("Мои группы", text.NewMyGroup(log, groupServ, serdes, b.Me.Username))

		// callbacks
		r.HandleFuncCallback("\fset_cookie", callback.NewChangeCookie(stateMachine))
		r.HandleFuncCallback("\fchange_notification", callback.NewChangeNotification(notifChanger, log))
	})

	r.Group(func(r router.Router) { // Routes for SendingCookie state
		r.Use(stater.New(stateMachine, fsm.SendingCookie))

		// message
		r.HandleFuncText("⬅️ Назад", text.NewStart(stateMachine))
		r.HandleFuncRegexpText(regexp.MustCompile(".+"), text.NewSendingCookie(log, cookieSetter, stateMachine))
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
