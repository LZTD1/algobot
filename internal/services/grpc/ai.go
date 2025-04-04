package grpc

import (
	"algobot/internal/config"
	"algobot/internal/domain/models"
	"algobot/internal/lib/logger/sl"
	aiv1 "algobot/protos"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
)

type AIService struct {
	grpc    aiv1.AiClient
	timeout time.Duration
	log     *slog.Logger
}

func NewAIService(grpcCfg config.GRPC, log *slog.Logger) *AIService {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", grpcCfg.Host, grpcCfg.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	return &AIService{
		grpc:    aiv1.NewAiClient(conn),
		timeout: grpcCfg.Timeout,
		log:     log,
	}
}

func (a *AIService) GetAIInfo(traceID interface{}) (models.AIInfo, error) {
	const op = "grpc.AIService.GetAIInfo"
	log := a.log.With(
		slog.String("op", op),
		slog.Any("trace_id", traceID),
	)

	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
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
