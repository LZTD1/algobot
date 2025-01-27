package services

import (
	"testing"
	"tgbot/internal/service"
	"tgbot/tests/mocks"
)

func TestDefaultService(t *testing.T) {
	t.Run("Get missing kids", func(t *testing.T) {
		defaultService := service.NewDefaultService(mocks.MockDomain{}, mocks.MockWebClient{})
		_ = defaultService
		// TODO доделать сервис
	})
}
