package helpers

import (
	"algobot/internal_old/models"
	"reflect"
	"testing"
	"time"
)

func Test_GetGroupsByDay(t *testing.T) {
	all := []models.Group{
		{
			GroupID:    1,
			Title:      "Lesson 1",
			TimeLesson: getDayByTime(21, 12, 30),
		},
		{
			GroupID:    2,
			Title:      "Lesson 2",
			TimeLesson: getDayByTime(27, 14, 00),
		},
		{
			GroupID:    3,
			Title:      "Lesson 3",
			TimeLesson: getDayByTime(26, 14, 00),
		},
		{
			GroupID:    4,
			Title:      "Lesson 4",
			TimeLesson: getDayByTime(28, 15, 30),
		},
	}
	t.Run("If groups exists", func(t *testing.T) {
		g, e := GetGroupsByDay(getDayByTime(14, 9, 9), all)
		if e != nil {
			t.Fatalf("Unexpected error %v", e)
		}
		want := []models.Group{
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
	all := []models.Group{
		{
			GroupID:    1,
			Title:      "Lesson 1", // вск
			TimeLesson: getDayByTime(21, 12, 30),
		},
		{
			GroupID:    2,
			Title:      "Lesson 2", // суб
			TimeLesson: getDayByTime(27, 14, 00),
		},
		{
			GroupID:    3,
			Title:      "Lesson 3", // птн
			TimeLesson: getDayByTime(26, 14, 00),
		},
		{
			GroupID:    4,
			Title:      "Lesson 4", // вск
			TimeLesson: getDayByTime(28, 15, 30),
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
