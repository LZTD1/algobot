package mocks

//go:generate mockgen -destination=./groupGetter_mock.go -package=mocks algobot/internal/services GroupGetter

//go:generate mockgen -destination=./aiClient_mock.go -package=mocks algobot/protos AiClient
