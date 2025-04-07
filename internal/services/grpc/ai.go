package grpc

import (
	"algobot/internal/config"
	"algobot/internal/domain/models"
	"algobot/internal/lib/logger/handlers/slogpretty"
	"algobot/internal/lib/logger/sl"
	aiv1 "algobot/protos"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
)

var (
	ErrNotValidResponse = errors.New("not valid response")
)

type AiOption func(*AIService)

type AIService struct {
	grpc aiv1.AiClient
	log  *slog.Logger
	cfg  config.GRPC
}

func NewAIService(cfg config.GRPC, fn ...func(*AIService)) *AIService {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	svr := &AIService{
		grpc: aiv1.NewAiClient(conn),
		log:  slog.New(slogpretty.NewHandler(&slog.HandlerOptions{Level: slog.LevelDebug})),
		cfg:  cfg,
	}

	for _, o := range fn {
		o(svr)
	}

	return svr
}
func WithLogger(log *slog.Logger) func(*AIService) {
	return func(s *AIService) {
		s.log = log
	}
}
func WithClient(client aiv1.AiClient) func(*AIService) {
	return func(s *AIService) {
		s.grpc = client
	}
}

func (a *AIService) GetAIInfo(traceID interface{}) (models.AIInfo, error) {
	const op = "grpc.AIService.GetAIInfo"
	log := a.log.With(
		slog.String("op", op),
		slog.Any("trace_id", traceID),
	)

	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.Timeout)
	defer cancel()

	information, err := a.grpc.GetInformation(ctx, &aiv1.GetInformationRequest{})
	if err != nil {
		log.Warn("error while calling gRPC GetInformation", sl.Err(err))
		return models.AIInfo{}, fmt.Errorf("%s error while calling gRPC GetInformation: %w", op, err)
	}

	return models.AIInfo{
		TextModel:  information.GetChatModel(),
		ImageModel: information.GetImageModel(),
	}, nil
}
