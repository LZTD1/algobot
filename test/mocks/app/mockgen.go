package mocks

//go:generate mockgen -destination=./context_mock.go -package=mocks algobot/internal/app/scheduler Domain
//go:generate mockgen -destination=./api_mock.go -package=mocks algobot/internal/app/scheduler Backoffice
