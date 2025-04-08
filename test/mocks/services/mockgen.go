package mocks

//go:generate mockgen -destination=./aiClient_mock.go -package=mocks algobot/protos AiClient

//go:generate mockgen -destination=./groupGetter_mock.go -package=mocks algobot/internal/services/groups GroupGetter
//go:generate mockgen -destination=./domainSetter_mock.go -package=mocks algobot/internal/services/groups DomainSetter
//go:generate mockgen -destination=./groupFetcher_mock.go -package=mocks algobot/internal/services/groups GroupFetcher
