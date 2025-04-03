package sqlite

import "fmt"

func (s *Sqlite) IsRegistered(uid int64) (bool, error) {
	const op = "sqlite.IsRegistered"

	sql := `SELECT COUNT(*) FROM users WHERE uid = ?`

	pr, err := s.db.Prepare(sql)
	if err != nil {
		return false, fmt.Errorf("%s error while preparing sql: %w", op, err)
	}

	var count int

	row := pr.QueryRow(uid)
	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("%s error while scanning sql: %w", op, err)
	}

	return count > 0, nil
}
