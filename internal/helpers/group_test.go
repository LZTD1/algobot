package helpers

import (
	"reflect"
	"testing"
	"tgbot/internal/domain"
	"time"
)

func Test_GetGroupsByDay(t *testing.T) {
	all := []domain.Group{
		{
			Id:   1,
			Name: "Lesson 1",
			Time: getDayByTime(21, 12, 30),
		},
		{
			Id:   2,
			Name: "Lesson 2",
			Time: getDayByTime(27, 14, 00),
		},
		{
			Id:   3,
			Name: "Lesson 3",
			Time: getDayByTime(26, 14, 00),
		},
		{
			Id:   4,
			Name: "Lesson 4",
			Time: getDayByTime(28, 15, 30),
		},
	}
	t.Run("If groups exists", func(t *testing.T) {
		g, e := GetGroupsByDay(getDayByTime(14, 9, 9), all)
		if e != nil {
			t.Fatalf("Unexpected error %v", e)
		}
		want := []domain.Group{
			all[0],
			all[3],
		}
		if reflect.DeepEqual(want, g) == false {
			t.Fatalf("Expected: %v, Got: %v", want, g)
		}
	})
	t.Run("If groups not exists", func(t *testing.T) {
		_, e := GetGroupsByDay(getDayByTime(18, 9, 9), all)
		if e == nil {
			t.Fatalf("Wanted error, got nil!")
		}
	})
}

func Test_GetCurrentGroup(t *testing.T) {
	all := []domain.Group{
		{
			Id:   1,
			Name: "Lesson 1", // вск
			Time: getDayByTime(21, 12, 30),
		},
		{
			Id:   2,
			Name: "Lesson 2", // суб
			Time: getDayByTime(27, 14, 00),
		},
		{
			Id:   3,
			Name: "Lesson 3", // птн
			Time: getDayByTime(26, 14, 00),
		},
		{
			Id:   4,
			Name: "Lesson 4", // вск
			Time: getDayByTime(28, 15, 30),
		},
	}
	t.Run("If group exists", func(t *testing.T) {
		date := getDayByTime(14, 12, 10)
		g, e := GetCurrentGroup(date, all)
		if e != nil {
			t.Fatalf("Unexpected error %v", e)
		}
		want := all[0]
		if reflect.DeepEqual(want, g) == false {
			t.Fatalf("Expected: %v, Got: %v", want, g)
		}
	})
	t.Run("If group end", func(t *testing.T) {
		date := getDayByTime(14, 21, 10)
		_, e := GetCurrentGroup(date, all)
		if e == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}

// getDayByTime 28, 21 вск ||  27, 20 сб
func getDayByTime(day, hour, min int) time.Time {
	return time.Date(2025, 9, day, hour, min, 0, 0, time.UTC)
}
