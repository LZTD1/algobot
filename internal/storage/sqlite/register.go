package sqlite

import (
	"fmt"
)

func (s *Sqlite) Register(uid int64) error {
	const op = "sqlite.Register"

	sqle := `
		INSERT INTO users (uid, cookie, last_notification_msg, notification)
		VALUES (?, NULL, NULL, 0);
	`

	pr, err := s.db.Prepare(sqle)
	if err != nil {
		return fmt.Errorf("%s error while preparing sql: %w", op, err)
	}

	_, err = pr.Exec(uid)
	if err != nil {
		return fmt.Errorf("%s error while exec sql: %w", op, err)
	}

	return nil
}
