package mocks

//go:generate mockgen -destination=./auther_mock.go -package=mocks algobot/internal/telegram/middleware/auth Auther
