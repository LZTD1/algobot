package sqlite

import "fmt"

func (s *Sqlite) ChaneNotifDate(uid int64, lastnotif string) error {
	const op = "sqlite.ChaneNotifDate"

	sqlq := "UPDATE users SET last_notification_msg = ? WHERE uid = ?"
	pr, err := s.db.Prepare(sqlq)
	if err != nil {
		return fmt.Errorf("%s error while preparing statement: %w", op, err)
	}
	defer pr.Close()

	if _, err := pr.Exec(lastnotif, uid); err != nil {
		return fmt.Errorf("%s error while executing statement: %w", op, err)
	}

	return nil
}
