package sqlite

import (
	"fmt"
)

func (s *Sqlite) Notification(uid int64) (bool, error) {
	const op = "sqlite.Notification"

	q := `SELECT notification FROM main.users WHERE uid=?`
	pr, err := s.db.Prepare(q)
	if err != nil {
		return false, fmt.Errorf("%s error while preparing sql: %w", op, err)
	}

	var notification int

	row := pr.QueryRow(uid)
	if err := row.Scan(&notification); err != nil {
		return false, fmt.Errorf("%s error while scanning sql: %w", op, err)
	}

	return notification == 1, nil
}
