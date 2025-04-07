package grpc

import (
	"algobot/internal/lib/logger/sl"
	aiv1 "algobot/protos"
	"context"
	"fmt"
	"log/slog"
)

func (a *AIService) ResetHistory(uid int64, traceID interface{}) error {
	const op = "grpc.AIService.ResetHistory"
	log := a.log.With(
		slog.String("op", op),
		slog.Any("trace_id", traceID),
	)

	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.Timeout)
	defer cancel()

	ok, err := a.grpc.ClearHistory(ctx, &aiv1.ClearHistoryRequest{
		Uid: uid,
	})
	if err != nil {
		log.Warn("error while calling gRPC ClearHistory", sl.Err(err))
		return fmt.Errorf("%s error while calling gRPC ClearHistory: %w", op, err)
	}
	if !ok.GetOk() {
		log.Warn("error while checking status grpc, not ok", slog.Bool("ok", ok.GetOk()))
		return fmt.Errorf("%s error while checking status grpc, not ok: %w", op, ErrNotValidResponse)
	}

	return nil
}
