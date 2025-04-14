package scheduler

import (
	"algobot/internal/config"
	"algobot/internal/domain/models"
	"algobot/internal/domain/scheduler"
	"algobot/internal/lib/logger/sl"
	"log/slog"
	"time"
)

type Domain interface {
	UsersByNotification(wantNotif int) ([]models.User, error)
	ChaneNotifDate(uid int64, lastnotif string) error
}
type Backoffice interface {
	MessagesUser(uid int64, lastTime string) ([]scheduler.Message, error)
}

type App struct {
	ch     chan scheduler.Message
	cfg    *config.Config
	ticker *time.Ticker
	log    *slog.Logger
	domain Domain
	bo     Backoffice
}

func New(log *slog.Logger, cfg *config.Config, domain Domain, bo Backoffice) *App {

	return &App{
		log:    log,
		cfg:    cfg,
		bo:     bo,
		domain: domain,
		ch:     make(chan scheduler.Message),
	}
}

func (a *App) Run() {
	const op = "scheduler.Run"
	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("start scheduling")

	ticker := time.NewTicker(a.cfg.Backoffice.MessageTimer)
	a.ticker = ticker

	go func() {
		for {
			select {
			case <-ticker.C:
				a.GetMessage()
			}
		}
	}()
}

func (a *App) GetMessage() {
	// TODO: maybe refactor into other struct
	const op = "scheduler.app.GetMessage"
	log := a.log.With(
		slog.String("op", op),
	)
	users, err := a.domain.UsersByNotification(1)
	if err != nil {
		log.Warn("error while get users by notif", sl.Err(err))
	}
	for _, user := range users {
		msgs, err := a.bo.MessagesUser(user.Uid, user.LastNotification)
		if err != nil {
			log.Warn("error while fetch MessagesUser", sl.Err(err))
			continue
		}
		for _, msg := range msgs {
			a.ch <- msg
			err := a.domain.ChaneNotifDate(msg.To, msg.Time)
			if err != nil {
				log.Warn("error while fetch ChaneNotifDate", sl.Err(err))
			}
		}
	}

}

func (a *App) Stop() {
	const op = "scheduler.Stop"
	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("stop scheduling")

	if a.ticker != nil {
		a.ticker.Stop()
	}
	close(a.ch)
}

func (a *App) Chan() chan scheduler.Message {
	return a.ch
}
