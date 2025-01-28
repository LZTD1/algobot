package services

import (
	"reflect"
	"testing"
	"tgbot/internal/domain"
	"tgbot/internal/service"
	"tgbot/tests/mocks"
	"time"
)

func TestDefaultService(t *testing.T) {
	t.Run("Get missing kids", func(t *testing.T) {
		defaultService := service.NewDefaultService(&mocks.MockDomain{}, mocks.MockWebClient{})
		kids, err := defaultService.MissingKids(
			1,
			time.Date(2024, 10, 6, 14, 55, 55, 3, time.UTC),
			33,
		)
		if err != nil {
			t.Fatalf("Got error: %v", err)
		}

		if len(kids) != 1 {
			t.Errorf("kids length should be 1, but got %d", len(kids))
		}
		if kids[0] != "Мария Петрова" {
			t.Fatalf("Wanted Мария Петрова, got - %s", kids[0])
		}
	})
	t.Run("Get CurrentGroup", func(t *testing.T) {
		defaultService := service.NewDefaultService(&mocks.MockDomain{}, mocks.MockWebClient{})
		group, err := defaultService.CurrentGroup(
			1,
			time.Date(2024, 10, 6, 14, 55, 55, 3, time.UTC),
		)
		if err != nil {
			t.Fatalf("Got error: %v", err)
		}
		wanted := domain.Group{
			Id:          1,
			Name:        "test1",
			Time:        time.Date(2024, 10, 6, 14, 55, 55, 3, time.UTC),
			Lesson:      "",
			MissingKids: []string{"Мария Петрова"},
		}
		if !reflect.DeepEqual(wanted, group) {
			t.Fatalf("Wanted %v, got %v", wanted, group)
		}
	})
	t.Run("Refresh groups", func(t *testing.T) {
		d := mocks.MockDomain{}
		defaultService := service.NewDefaultService(&d, mocks.MockWebClient{})
		err := defaultService.RefreshGroups(1)
		if err != nil {
			t.Fatalf("Got error: %v", err)
		}
		wanted := domain.Group{
			Id:          1,
			Name:        "Title",
			Time:        time.Date(2025, 2, 1, 14, 0, 0, 0, time.UTC),
			Lesson:      "",
			MissingKids: nil,
		}
		if !reflect.DeepEqual(d.MockGroups[0], []domain.Group{wanted}[0]) {
			t.Fatalf("Wanted %#v, got %#v", []domain.Group{wanted}[0], d.MockGroups[0])
		}
	})
}
