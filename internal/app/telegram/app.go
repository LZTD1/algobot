package telegram

import (
	"algobot/internal/config"
	backoffice3 "algobot/internal/lib/backoffice"
	"algobot/internal/lib/fsm"
	"algobot/internal/lib/fsm/memory"
	"algobot/internal/lib/logger/sl"
	"algobot/internal/lib/serdes/base62"
	"algobot/internal/services/backoffice"
	"algobot/internal/services/groups"
	grpc2 "algobot/internal/services/grpc"
	"algobot/internal/storage/sqlite"
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
	cfg *config.Config,
	storage *sqlite.Sqlite,
	bo *backoffice3.Backoffice,
) *App {
	const op = "telegram.New"

	nlog := log.With(
		slog.String("op", op),
	)

	pref := tele.Settings{
		Token: cfg.TelegramToken,
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
	groupServ := groups.NewGroup(log, storage, bo, storage, bo)
	stateMachine := memory.New()
	serdes := base62.NewSerdes(log)
	grpc := grpc2.NewAIService(
		cfg.GRPC,
		grpc2.WithLogger(log),
	)
	boSvc := backoffice.NewBackoffice(log, storage, bo, bo, bo)

	// initialize routes
	b.Use(trace.New(log))
	b.Use(middleware.AutoRespond())
	b.Use(middleware.Recover())
	b.Use(logger.New(log))
	b.Use(auth.New(storage, log))
	b.Use(rate.New(log, cfg.RateLimit))

	// create routing
	r := router.NewRouter()
	r.Group(func(r router.Router) { // Routes for default state
		r.Use(stater.New(stateMachine, fsm.Default))

		// message
		r.HandleFuncText("/start", text.NewStart(stateMachine))
		r.HandleFuncText("Настройки", text.NewSettings(storage, log))
		r.HandleFuncText("AI 🔹", text.NewAI(grpc, log, stateMachine))
		r.HandleText("Мои группы", text.NewMyGroup(log, groupServ, serdes, b.Me.Username))
		r.HandleFuncText("Получить отсутсвующих", text.NewMissingKids(log, groupServ))
		r.HandleFuncRegexpText(regexp.MustCompile(`^(?m)\/abs(.*)$`), text.NewAbsentKids(groupServ, log))

		r.HandleRegexpText(regexp.MustCompile(`^(?m)\/start\s(.+)$`), text.NewViewInformer(serdes, boSvc, log, b.Me.Username))

		// callbacks
		r.HandleFuncCallback("\fset_cookie", callback.NewChangeCookie(stateMachine))
		r.HandleFuncCallback("\fchange_notification", callback.NewChangeNotification(storage, log))
		r.HandleFuncCallback("\frefresh_groups", callback.RefreshGroup(groupServ, log))

		r.HandleFuncRegexpCallback(regexp.MustCompile(`^\fclose_lesson_(.+)$`), callback.LessonStatus(boSvc, backoffice.CloseLesson, log))
		r.HandleFuncRegexpCallback(regexp.MustCompile(`^\fopen_lesson_(.+)$`), callback.LessonStatus(boSvc, backoffice.OpenLesson, log))
	})

	r.Group(func(r router.Router) { // Routes for SendingCookie state
		r.Use(stater.New(stateMachine, fsm.SendingCookie))

		// message
		r.HandleFuncText("⬅️ Назад", text.NewStart(stateMachine))
		r.HandleFuncRegexpText(regexp.MustCompile(".+"), text.NewSendingCookie(log, storage, stateMachine))
	})

	r.Group(func(r router.Router) { // Routes for ChattingAI state
		r.Use(stater.New(stateMachine, fsm.ChattingAI))

		// message
		r.HandleFuncText("⬅️ Назад", text.NewStart(stateMachine))
		r.HandleFuncText("/reset", text.NewReset(grpc, log))
		r.HandleFuncRegexpText(regexp.MustCompile(`^(?m)\/image\s(.+)$`), text.GenerateImage(grpc, log))
		r.HandleFuncRegexpText(regexp.MustCompile(`^[^/].*$`), text.ChatAI(grpc, log))
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
