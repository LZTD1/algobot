package grpc

import (
	"algobot/internal/lib/logger/sl"
	aiv1 "algobot/protos"
	"context"
	"fmt"
	"log/slog"
)

func (a *AIService) ChatAI(uid int64, message string, traceID interface{}) (string, error) {
	const op = "grpc.AIService.ChatAI"
	log := a.log.With(
		slog.String("op", op),
		slog.Any("trace_id", traceID),
	)

	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.Timeout)
	defer cancel()

	resp, err := a.grpc.GetSuggest(ctx, &aiv1.SuggestRequest{
		Uid:     uid,
		Suggest: message,
	})
	if err != nil {
		log.Warn("error while calling gRPC GetSuggest", sl.Err(err))
		return "", fmt.Errorf("%s error while calling gRPC GetSuggest: %w", op, err)
	}
	if !resp.GetOk() {
		log.Warn("error while checking status grpc, not ok", slog.Bool("ok", resp.GetOk()))
		return "", fmt.Errorf("%s error while checking status grpc, not ok: %w", op, ErrNotValidResponse)
	}

	return resp.GetRequest(), nil
}
