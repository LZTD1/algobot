package rate

import (
	"algobot/internal/config"
	"golang.org/x/time/rate"
	"gopkg.in/telebot.v4"
	"log/slog"
	"sync"
)

var userLimits = make(map[int64]*rate.Limiter)
var mu sync.Mutex

func New(log *slog.Logger, rateCfg config.RateLimit) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		log = log.With(
			slog.String("component", "middleware/rate"),
		)

		return func(ctx telebot.Context) error {
			traceID := ctx.Get("trace_id")
			log = log.With("trace_id", traceID)
			uid := ctx.Sender().ID

			limiter := getUserLimiter(uid, rateCfg)
			if !limiter.Allow() {
				log.Warn("user limit exceeded")
				return ctx.Send("🙈 Слишком много запросов, давай помедленнее!")
			}

			return next(ctx)
		}
	}
}

func getUserLimiter(uid int64, cfg config.RateLimit) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if limiter, ok := userLimits[uid]; ok {
		return limiter
	}

	limiter := rate.NewLimiter(rate.Every(cfg.FillPeriod), cfg.BucketLimit)
	userLimits[uid] = limiter
	return limiter
}
