package grpc

import (
	"algobot/internal/lib/logger/sl"
	aiv1 "algobot/protos"
	"context"
	"fmt"
	"log/slog"
)

func (a *AIService) GenerateImage(uid int64, promt string, traceID interface{}) (string, error) {
	const op = "grpc.AIService.GenerateImage"
	log := a.log.With(
		slog.String("op", op),
		slog.Any("trace_id", traceID),
	)

	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.Timeout)
	defer cancel()

	resp, err := a.grpc.GenerateImage(ctx, &aiv1.GenerateImageRequest{
		Uid:   uid,
		Promt: promt,
	})

	if err != nil {
		log.Warn("error while calling gRPC GenerateImage", sl.Err(err))
		return "", fmt.Errorf("%s error while calling gRPC GenerateImage: %w", op, err)
	}

	return resp.GetUrl(), nil
}
