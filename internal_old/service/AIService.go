package service

import (
	pkg "algobot/protos"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type AIService interface {
	GetSuggestion(uid int64, text string) (string, error)
	ClearAllHistory(uid int64) error
}

type aiService struct {
	grpc pkg.AiClient
}

func NewAiService(port string) AIService {
	conn, err := grpc.NewClient("localhost:"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &aiService{
		grpc: pkg.NewAiClient(conn),
	}
}

func (a aiService) GetSuggestion(uid int64, text string) (string, error) {
	ctx := context.Background()
	suggest, err := a.grpc.GetSuggest(ctx, &pkg.SuggestRequest{
		Uid:     uid,
		Suggest: text,
	})
	if err != nil {
		return "", fmt.Errorf("aiService.GetSuggestion(%v, %v) : %w", uid, text, err)
	}

	return suggest.GetRequest(), nil
}

func (a aiService) ClearAllHistory(uid int64) error {
	ctx := context.Background()
	_, err := a.grpc.ClearHistory(ctx, &pkg.ClearHistoryRequest{
		Uid: uid,
	})
	if err != nil {
		return fmt.Errorf("aiService.ClearAllHistory(%v) : %w", uid, err)
	}

	return nil
}
