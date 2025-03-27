package scheduler_test

import (
	"algobot/internal_old/schedulers"
	"algobot/tests/mocks"
	"testing"
)

func TestScheduler(t *testing.T) {
	mockService := mocks.NewMockService(nil)
	mockBot := mocks.MockBot{}

	scheduler := schedulers.NewMessage(&mockBot, mockService)
	scheduler.Schedule()

	if len(mockBot.Calls) != 1 {
		t.Errorf("Wanted 1, got %d", len(mockBot.Calls))
	}
	wanted := "Send(1) = 🔔 Новое сообщение\n\nОт: 1\nТема: 2\nСсылка: 3\n\n```Сообщение:\n4\n```"
	if mockBot.Calls[0] != wanted {
		t.Errorf("Wanted %s, got %s", wanted, mockBot.Calls[0])
	}
}
