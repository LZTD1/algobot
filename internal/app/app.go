package app

import (
	"algobot/internal/app/telegram"
	"algobot/internal/config"
	"algobot/internal/lib/backoffice"
	"algobot/internal/storage/sqlite"
	"log/slog"
)

type App struct {
	log         *slog.Logger
	cfg         *config.Config
	TelegramBot *telegram.App
}

func New(log *slog.Logger, cfg *config.Config) *App {

	storage, err := sqlite.NewDB(cfg)
	if err != nil {
		panic(err)
	}
	bo := backoffice.NewBackoffice(&cfg.Backoffice)

	botApplication := telegram.New(
		log,
		cfg,
		storage,
		bo,
	)

	return &App{log: log, cfg: cfg, TelegramBot: botApplication}
}
