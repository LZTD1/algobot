package app

import (
	"algobot/internal/app/telegram"
	"algobot/internal/config"
	"log/slog"
)

type App struct {
	log         *slog.Logger
	cfg         *config.Config
	TelegramBot *telegram.App
}

func New(log *slog.Logger, cfg *config.Config) *App {

	botApplication := telegram.New(log, cfg.TelegramToken)

	return &App{log: log, cfg: cfg, TelegramBot: botApplication}
}
