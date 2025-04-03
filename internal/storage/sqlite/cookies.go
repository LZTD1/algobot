package sqlite

import (
	"database/sql"
	"fmt"
)

func (s *Sqlite) Cookies(uid int64) (string, error) {
	const op = "sqlite.Cookies"

	q := `SELECT cookie FROM main.users WHERE uid=?`
	pr, err := s.db.Prepare(q)
	if err != nil {
		return "", fmt.Errorf("%s error while preparing sql: %w", op, err)
	}

	var cookie sql.NullString

	row := pr.QueryRow(uid)
	if err := row.Scan(&cookie); err != nil {
		return "", fmt.Errorf("%s error while scanning sql: %w", op, err)
	}

	if !cookie.Valid {
		return "", nil
	}

	return cookie.String, nil
}
