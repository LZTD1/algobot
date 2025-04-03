package app

import (
	"algobot/internal/app/telegram"
	"algobot/internal/config"
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
	botApplication := telegram.New(log, cfg.TelegramToken, storage, storage)

	return &App{log: log, cfg: cfg, TelegramBot: botApplication}
}
