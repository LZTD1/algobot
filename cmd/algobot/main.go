package main

import (
	"algobot/internal/app"
	"algobot/internal/config"
	"algobot/internal/lib/logger/handlers/slogpretty"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envProd  string = "prod"
	envLocal string = "local"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("starting application")

	application := app.New(log, cfg)

	go application.TelegramBot.Run()
	log.Info("started telegram bot")
	// TODO : start bot app
	// TODO : start message scheduler app

	// graceful shutdown
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-ch
	log.Info("shutting down application")
	application.TelegramBot.Stop()
	log.Info("application gracefully stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slogpretty.NewHandler(&slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
