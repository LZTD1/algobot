package test

import (
	"algobot/internal/config"
	"algobot/internal/domain/models"
	sqlite2 "algobot/internal/storage/sqlite"
	"database/sql"
	"fmt"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

const (
	testDBPath     = "test.db"
	migrationsPath = "./migrations-suite"
)

func TestSqlite(t *testing.T) {
	t.Cleanup(func() {
		err := os.Remove(testDBPath)
		if err != nil {
			panic(err)
		}
	})

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", testDBPath))
	defer db.Close()
	assert.NoError(t, err)

	err = goose.SetDialect("sqlite3")
	assert.NoError(t, err)

	err = goose.Up(db, migrationsPath)
	assert.NoError(t, err)

	sqlite, err := sqlite2.NewDB(&config.Config{
		StoragePath: testDBPath,
	})
	defer sqlite.MustClose()

	t.Run("IsRegistered", func(t *testing.T) {
		t.Run("user have reg", func(t *testing.T) {
			registered, err := sqlite.IsRegistered(1001)
			assert.NoError(t, err)
			assert.True(t, registered)
		})
		t.Run("user dont have reg", func(t *testing.T) {
			registered, err := sqlite.IsRegistered(0)
			assert.NoError(t, err)
			assert.False(t, registered)
		})
	})
	t.Run("Register", func(t *testing.T) {
		t.Run("successful registered", func(t *testing.T) {
			err := sqlite.Register(1002)
			assert.NoError(t, err)

			registered, err := sqlite.IsRegistered(1002)
			assert.NoError(t, err)
			assert.True(t, registered)
		})
		t.Run("already registered", func(t *testing.T) {
			err := sqlite.Register(1001)
			assert.Error(t, err)
		})
	})
	t.Run("Cookie", func(t *testing.T) {
		t.Run("Cookie set", func(t *testing.T) {
			cookie, err := sqlite.Cookies(1000)
			assert.NoError(t, err)
			assert.Equal(t, "", cookie)
		})
		t.Run("Cookie null", func(t *testing.T) {
			cookie, err := sqlite.Cookies(1001)
			assert.NoError(t, err)
			assert.Equal(t, "cookie", cookie)
		})
	})
	t.Run("Notification", func(t *testing.T) {
		notfication, err := sqlite.Notification(1000)
		assert.NoError(t, err)
		assert.Equal(t, false, notfication)
	})
	t.Run("SetCookie", func(t *testing.T) {
		cookie, err := sqlite.Cookies(999)
		assert.NoError(t, err)
		assert.Equal(t, "", cookie)

		err = sqlite.SetCookie(999, "a@a")
		assert.NoError(t, err)

		cookie, err = sqlite.Cookies(999)
		assert.NoError(t, err)
		assert.Equal(t, "a@a", cookie)
	})
	t.Run("SetNotification", func(t *testing.T) {
		notif, err := sqlite.Notification(999)
		assert.NoError(t, err)
		assert.False(t, notif)

		err = sqlite.SetNotification(999, true)
		assert.NoError(t, err)

		notif, err = sqlite.Notification(999)
		assert.NoError(t, err)
		assert.True(t, notif)

		err = sqlite.SetNotification(999, false)
		assert.NoError(t, err)

		notif, err = sqlite.Notification(999)
		assert.NoError(t, err)
		assert.False(t, notif)
	})
	t.Run("Groups", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			groups, err := sqlite.Groups(999)
			assert.NoError(t, err)
			assert.Len(t, groups, 3)
			assert.Equal(t, []models.Group{
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
				{
					GroupID:    999,
					Title:      "group 3",
					TimeLesson: time.Date(2025, time.March, 22, 14, 0, 0, 0, time.UTC),
				},
			}, groups)
		})
		t.Run("no one group", func(t *testing.T) {
			groups, err := sqlite.Groups(1000)
			assert.NoError(t, err)
			assert.Len(t, groups, 0)
		})
	})
	t.Run("SetGroups", func(t *testing.T) {
		gr := []models.Group{
			{
				GroupID:    1,
				Title:      "title",
				TimeLesson: time.Date(2025, time.March, 23, 14, 0, 0, 0, time.UTC),
			},
		}

		groups, err := sqlite.Groups(999)
		assert.NoError(t, err)
		assert.Len(t, groups, 3)

		err = sqlite.SetGroups(999, gr)
		assert.NoError(t, err)

		groups, err = sqlite.Groups(999)
		assert.NoError(t, err)
		assert.Len(t, groups, 1)
		assert.Equal(t, gr, groups)
	})
}
