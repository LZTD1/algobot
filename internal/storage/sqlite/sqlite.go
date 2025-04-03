package sqlite

import (
	"algobot/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type Sqlite struct {
	db *sql.DB
}

func NewDB(cfg *config.Config) (*Sqlite, error) {
	const op = "sqlite.NewDB"

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", cfg.StoragePath))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Sqlite{db: db}, nil
}
