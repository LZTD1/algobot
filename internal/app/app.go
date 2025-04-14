package app

import (
	"algobot/internal/app/scheduler"
	"algobot/internal/app/telegram"
	"algobot/internal/config"
	"algobot/internal/lib/backoffice"
	backoffice2 "algobot/internal/services/backoffice"
	"algobot/internal/storage/sqlite"
	"log/slog"
)

type App struct {
	log         *slog.Logger
	cfg         *config.Config
	TelegramBot *telegram.App
	Scheduler   *scheduler.App
}

func New(log *slog.Logger, cfg *config.Config) *App {

	storage, err := sqlite.NewDB(cfg)
	if err != nil {
		panic(err)
	}

	bo := backoffice.NewBackoffice(&cfg.Backoffice)
	boSvc := backoffice2.NewBackoffice(log, storage, bo, bo, bo, bo)

	sch := scheduler.New(log, cfg, storage, boSvc)
	botApplication := telegram.New(
		log,
		cfg,
		storage,
		bo,
		boSvc,
		sch.Chan(),
	)

	return &App{log: log, cfg: cfg, TelegramBot: botApplication, Scheduler: sch}
}
