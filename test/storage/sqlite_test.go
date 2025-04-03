package test

import (
	"algobot/internal/config"
	sqlite2 "algobot/internal/storage/sqlite"
	"database/sql"
	"fmt"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
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
}
