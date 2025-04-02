package mocks

//go:generate mockgen -destination=./set_stater_mock.go -package=mocks algobot/internal/telegram/handlers/text SetStater
//go:generate mockgen -destination=./userInformer_mock.go -package=mocks algobot/internal/telegram/handlers/text UserInformer
//go:generate mockgen -destination=./stateChanger_mock.go -package=mocks algobot/internal/telegram/handlers/callback StateChanger
//go:generate mockgen -destination=./notificationChanger_mock.go -package=mocks algobot/internal/telegram/handlers/callback NotificationChanger
//go:generate mockgen -destination=./aiInformer_mock.go -package=mocks algobot/internal/telegram/handlers/text AIInformer
//go:generate mockgen -destination=./aiStater_mock.go -package=mocks algobot/internal/telegram/handlers/text AIStater
