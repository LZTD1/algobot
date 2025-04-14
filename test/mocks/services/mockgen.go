package mocks

//go:generate mockgen -destination=./aiClient_mock.go -package=mocks algobot/protos AiClient

//go:generate mockgen -destination=./groupGetter_mock.go -package=mocks algobot/internal/services/groups GroupGetter
//go:generate mockgen -destination=./domainSetter_mock.go -package=mocks algobot/internal/services/groups DomainSetter
//go:generate mockgen -destination=./groupFetcher_mock.go -package=mocks algobot/internal/services/groups GroupFetcher

//go:generate mockgen -destination=./groupView_mock.go -package=mocks algobot/internal/services/backoffice GroupView
//go:generate mockgen -destination=./kidViewer_mock.go -package=mocks algobot/internal/services/backoffice KidViewer
//go:generate mockgen -destination=./cookieGetter_mock.go -package=mocks algobot/internal/services/backoffice CookieGetter
//go:generate mockgen -destination=./lessonStatuser_mock.go -package=mocks algobot/internal/services/backoffice LessonStatuser

//go:generate mockgen -destination=./kidStats_mock.go -package=mocks algobot/internal/services/groups KidStats
//go:generate mockgen -destination=./sender_mock.go -package=mocks algobot/internal/services/schedule Sender

//go:generate mockgen -destination=./messageFetcher_mock.go -package=mocks algobot/internal/services/backoffice MessageFetcher
