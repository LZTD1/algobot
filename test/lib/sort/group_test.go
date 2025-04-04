package sort

import (
	"algobot/internal/domain/models"
	"algobot/internal/lib/sort"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSortGroup(t *testing.T) {
	sort.GroupsByDate(assets)
	assert.Equal(t, []models.Group{
		{
			GroupID:    999,
			Title:      "group 3",
			TimeLesson: time.Date(2025, time.March, 22, 14, 0, 0, 0, time.UTC),
		},
		{
			GroupID:    1001,
			Title:      "group 1",
			TimeLesson: time.Date(2025, time.March, 23, 14, 0, 0, 0, time.UTC),
		},
		{
			GroupID:    1000,
			Title:      "group 2",
			TimeLesson: time.Date(2025, time.March, 23, 16, 0, 0, 0, time.UTC),
		},
	}, assets)
}

var assets = []models.Group{
	{
		GroupID:    1000,
		Title:      "group 2",
		TimeLesson: time.Date(2025, time.March, 23, 16, 0, 0, 0, time.UTC),
	},
	{
		GroupID:    1001,
		Title:      "group 1",
		TimeLesson: time.Date(2025, time.March, 23, 14, 0, 0, 0, time.UTC),
	},

	{
		GroupID:    999,
		Title:      "group 3",
		TimeLesson: time.Date(2025, time.March, 22, 14, 0, 0, 0, time.UTC),
	},
}
