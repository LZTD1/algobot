package mocks

//go:generate mockgen -destination=./set_stater_mock.go -package=mocks algobot/internal/telegram/handlers/text SetStater
//go:generate mockgen -destination=./userInformer_mock.go -package=mocks algobot/internal/telegram/handlers/text UserInformer
//go:generate mockgen -destination=./stateChanger_mock.go -package=mocks algobot/internal/telegram/handlers/callback StateChanger
//go:generate mockgen -destination=./notificationChanger_mock.go -package=mocks algobot/internal/telegram/handlers/callback NotificationChanger
//go:generate mockgen -destination=./aiInformer_mock.go -package=mocks algobot/internal/telegram/handlers/text AIInformer
//go:generate mockgen -destination=./aiStater_mock.go -package=mocks algobot/internal/telegram/handlers/text AIStater

//go:generate mockgen -destination=./cookieSetter_mock.go -package=mocks algobot/internal/telegram/handlers/text CookieSetter
//go:generate mockgen -destination=./cookieStater_mock.go -package=mocks algobot/internal/telegram/handlers/text CookieStater

//go:generate mockgen -destination=./grouper_mock.go -package=mocks algobot/internal/telegram/handlers/text Grouper
//go:generate mockgen -destination=./groupSerializer_mock.go -package=mocks algobot/internal/telegram/handlers/text GroupSerializer

//go:generate mockgen -destination=./reseter_mock.go -package=mocks algobot/internal/telegram/handlers/text Reseter

//go:generate mockgen -destination=./generatorImage_mock.go -package=mocks algobot/internal/telegram/handlers/text GeneratorImage
