package test

import (
	"algobot/internal/config"
	"algobot/internal/domain/models"
	"algobot/internal/services/grpc"
	aiv1 "algobot/protos"
	"algobot/test/mocks"
	mocks2 "algobot/test/mocks/services"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestGRPCAI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mocks.NewMockLogger()
	aiClient := mocks2.NewMockAiClient(ctrl)

	svc := grpc.NewAIService(config.GRPC{
		Host:    "0.0.0.0",
		Port:    "1111",
		Timeout: 100 * time.Millisecond,
	}, grpc.WithLogger(log), grpc.WithClient(aiClient))

	t.Run("GetAIInfo", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			aiClient.EXPECT().GetInformation(gomock.Any(), gomock.Any()).Return(
				&aiv1.GetInformationResponse{
					ChatModel:  "chat",
					ImageModel: "image",
				}, nil,
			).Times(1)

			info, err := svc.GetAIInfo("")
			assert.NoError(t, err)
			assert.Equal(t, models.AIInfo{
				TextModel:  "chat",
				ImageModel: "image",
			}, info)
		})
		t.Run("GetInformation return err", func(t *testing.T) {
			errExp := errors.New("GetInformation err")
			aiClient.EXPECT().GetInformation(gomock.Any(), gomock.Any()).Return(nil, errExp).Times(1)

			_, err := svc.GetAIInfo("")
			assert.Error(t, err)
			assert.ErrorIs(t, err, errExp)
		})
	})
	t.Run("ResetHistory", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			aiClient.EXPECT().ClearHistory(gomock.Any(), &aiv1.ClearHistoryRequest{
				Uid: int64(1),
			}).Return(&aiv1.ClearHistoryResponse{Ok: true}, nil).Times(1)

			err := svc.ResetHistory(1, "")
			assert.NoError(t, err)
		})
		t.Run("ClearHistory return err", func(t *testing.T) {
			errExp := errors.New("ClearHistory err")

			aiClient.EXPECT().ClearHistory(gomock.Any(), &aiv1.ClearHistoryRequest{
				Uid: int64(1),
			}).Return(&aiv1.ClearHistoryResponse{Ok: false}, errExp).Times(1)

			err := svc.ResetHistory(1, "")
			assert.Error(t, err)
			assert.ErrorIs(t, err, errExp)
		})
		t.Run("ClearHistory return false", func(t *testing.T) {
			aiClient.EXPECT().ClearHistory(gomock.Any(), &aiv1.ClearHistoryRequest{
				Uid: int64(1),
			}).Return(&aiv1.ClearHistoryResponse{Ok: false}, nil).Times(1)

			err := svc.ResetHistory(1, "")
			assert.Error(t, err)
			assert.ErrorIs(t, err, grpc.ErrNotValidResponse)
		})
	})
	t.Run("GenerateImage", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			aiClient.EXPECT().GenerateImage(gomock.Any(), &aiv1.GenerateImageRequest{
				Uid:   int64(1),
				Promt: "firefox",
			}).Return(&aiv1.GenerateImageResponse{
				Url: "https",
			}, nil).Times(1)

			url, err := svc.GenerateImage(int64(1), "firefox", "")
			assert.NoError(t, err)
			assert.Equal(t, "https", url)
		})
		t.Run("GenerateImage return err", func(t *testing.T) {
			errExp := errors.New("GenerateImage err")

			aiClient.EXPECT().GenerateImage(gomock.Any(), &aiv1.GenerateImageRequest{
				Uid:   int64(1),
				Promt: "firefox",
			}).Return(nil, errExp).Times(1)

			_, err := svc.GenerateImage(int64(1), "firefox", "")
			assert.Error(t, err)
			assert.ErrorIs(t, err, errExp)
		})
	})
	t.Run("GetSuggest", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			aiClient.EXPECT().GetSuggest(gomock.Any(), &aiv1.SuggestRequest{
				Uid:     int64(1),
				Suggest: "Suggest",
			}).Return(&aiv1.SuggestResponse{
				Ok:      true,
				Request: "Request",
			}, nil).Times(1)

			msg, err := svc.ChatAI(int64(1), "Suggest", "")
			assert.NoError(t, err)
			assert.Equal(t, "Request", msg)
		})
		t.Run("GetSuggest return err", func(t *testing.T) {
			errExp := errors.New("GetSuggest err")

			aiClient.EXPECT().GetSuggest(gomock.Any(), &aiv1.SuggestRequest{
				Uid:     int64(1),
				Suggest: "Suggest",
			}).Return(&aiv1.SuggestResponse{
				Ok:      false,
				Request: "",
			}, errExp).Times(1)

			_, err := svc.ChatAI(int64(1), "Suggest", "")
			assert.Error(t, err)
			assert.ErrorIs(t, err, errExp)
		})
		t.Run("ClearHistory return false", func(t *testing.T) {
			aiClient.EXPECT().GetSuggest(gomock.Any(), &aiv1.SuggestRequest{
				Uid:     int64(1),
				Suggest: "Suggest",
			}).Return(&aiv1.SuggestResponse{
				Ok:      false,
				Request: "",
			}, nil).Times(1)

			_, err := svc.ChatAI(int64(1), "Suggest", "")
			assert.Error(t, err)
			assert.ErrorIs(t, err, grpc.ErrNotValidResponse)
		})
	})

}
